package smi

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sleepinggenius2/gosmi"
)

//go:embed native/*
var fs embed.FS

type MibNode struct {
	Name        string
	OID         string
	Description string
	SourceMib   string
	Access      string
	Type        string
	Children    []*MibNode
}

var (
	RootNode   *MibNode
	smiMutex   sync.Mutex
	tempMibDir string
)

func Init() {
	smiMutex.Lock()

	gosmi.Init()

	tempMibDir, _ = os.MkdirTemp("", "go-mibs-*")

	entries, err := fs.ReadDir("native")
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				content, _ := fs.ReadFile("native/" + entry.Name())
				os.WriteFile(filepath.Join(tempMibDir, entry.Name()), content, 0644)
			}
		}
	}

	gosmi.AppendPath(tempMibDir)

	gosmi.LoadModule("SNMPv2-SMI")
	gosmi.LoadModule("SNMPv2-TC")
	gosmi.LoadModule("SNMPv2-CONF")
	gosmi.LoadModule("SNMPv2-MIB")
	gosmi.LoadModule("IF-MIB")

	smiMutex.Unlock()

	mibPaths, err := GetCustomMibs()
	if err == nil {
		for _, path := range mibPaths {
			fmt.Printf("Load %s\n", path)
			LoadFromFile(path, true)
		}
	}

	rebuildTreeLocked()
}

func LoadFromFile(filePath string, cached bool) error {
	smiMutex.Lock()
	defer smiMutex.Unlock()

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	moduleName, err := extractModuleName(absPath)
	if err != nil {
		return fmt.Errorf("cannot read module name: %w", err)
	}

	fmt.Printf("Loading file from dist: %s (detected module: %s)\n", filePath, moduleName)

	if tempMibDir == "" {
		tempMibDir, _ = os.MkdirTemp("", "go-mibs-*")
	}

	gosmi.AppendPath(tempMibDir)
	gosmi.AppendPath(filepath.Dir(absPath))

	content, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	idealPath := filepath.Join(tempMibDir, moduleName)
	err = os.WriteFile(idealPath, content, 0644)
	if err != nil {
		return fmt.Errorf("cannot parse temp file: %w", err)
	}

	_, err = gosmi.LoadModule(moduleName)
	if err != nil {
		return fmt.Errorf("gosmi cannot load module %s: %w", moduleName, err)
	}

	if !cached {
		fmt.Printf("Cached...")
		PushCustomMib(filePath)
	}

	rebuildTreeLocked()
	return nil
}

func extractModuleName(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(b), "\n")
	var words []string

	for _, line := range lines {
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}
		words = append(words, strings.Fields(line)...)
	}

	for i, w := range words {
		if w == "DEFINITIONS" && i > 0 {
			return words[i-1], nil
		}
	}

	return "", fmt.Errorf("cannot find DEFINITIONS keyword")
}

func rebuildTreeLocked() {
	localNodeMap := make(map[string]*MibNode)

	for _, mod := range gosmi.GetLoadedModules() {
		for _, n := range mod.GetNodes() {
			oidStr := n.Oid.String()
			if _, exists := localNodeMap[oidStr]; !exists {
				typeStr := ""
				if n.Type != nil {
					typeStr = n.Type.Name
				}
				localNodeMap[oidStr] = &MibNode{
					Name:        n.Name,
					OID:         oidStr,
					Access:      n.Access.String(),
					SourceMib:   mod.Name,
					Description: n.Description,
					Type:        typeStr,
				}
			}
		}
	}

	RootNode = &MibNode{Name: "Root", OID: ""}

	for oidStr, node := range localNodeMap {
		parentOid := getParentOID(oidStr)
		if parentNode, exists := localNodeMap[parentOid]; exists {
			parentNode.Children = append(parentNode.Children, node)
		} else {
			RootNode.Children = append(RootNode.Children, node)
		}
	}

	for _, child := range RootNode.Children {
		flattenBaseOIDs(child)
	}

	sortTree(RootNode)
}

func flattenBaseOIDs(node *MibNode) {
	if node == nil {
		return
	}

	for len(node.Children) == 1 {
		single := node.Children[0]
		cleanOID := strings.TrimPrefix(single.OID, ".")

		if cleanOID != "1" && cleanOID != "1.3" && cleanOID != "1.3.6" && cleanOID != "1.3.6.1" {
			break
		}

		node.Name = node.Name + "." + single.Name
		node.OID = single.OID
		node.Description = single.Description
		node.Type = single.Type
		node.Children = single.Children
	}

	for _, child := range node.Children {
		flattenBaseOIDs(child)
	}
}

func getParentOID(oid string) string {
	lastDot := strings.LastIndex(oid, ".")
	if lastDot == -1 {
		return ""
	}
	return oid[:lastDot]
}

func compareOID(oid1, oid2 string) bool {
	parts1 := strings.Split(strings.TrimPrefix(oid1, "."), ".")
	parts2 := strings.Split(strings.TrimPrefix(oid2, "."), ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		n1, err1 := strconv.Atoi(parts1[i])
		n2, err2 := strconv.Atoi(parts2[i])
		if err1 == nil && err2 == nil {
			if n1 != n2 {
				return n1 < n2
			}
		} else {
			if parts1[i] != parts2[i] {
				return parts1[i] < parts2[i]
			}
		}
	}
	return len(parts1) < len(parts2)
}

func sortTree(node *MibNode) {
	if node == nil {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		return compareOID(node.Children[i].OID, node.Children[j].OID)
	})

	for _, child := range node.Children {
		sortTree(child)
	}
}
