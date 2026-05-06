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
	"github.com/sleepinggenius2/gosmi/types"
)

//go:embed native/*
var fs embed.FS

type NodeDetails struct {
	Syntax  string
	DefVal  string
	Indexes string
}

type MibNode struct {
	Name        string
	OID         string
	Label       string // Pre-calculated for performance
	Description string
	SourceMib   string
	Access      string
	Syntax      string
	Type        string
	Details     *NodeDetails
	Children    []*MibNode
}

var (
	RootNode   *MibNode
	smiMutex   sync.Mutex
	TreeMutex  sync.RWMutex
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
			// Pass true for cached so it doesn't trigger side effects, 
			// and we should add a non-rebuilding version of LoadFromFile if needed,
			// but for now, let's just make sure we only rebuild at the end.
			loadFromFileInternal(path)
		}
	}

	rebuildTreeLocked()
}

func loadFromFileInternal(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	moduleName, err := extractModuleName(absPath)
	if err != nil {
		return fmt.Errorf("cannot read module name: %w", err)
	}

	fmt.Printf("Loading file: %s (%s)\n", filePath, moduleName)

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
	os.WriteFile(idealPath, content, 0644)

	gosmi.LoadModule(moduleName)
	return nil
}

func LoadFromFile(filePath string, cached bool) error {
	smiMutex.Lock()
	defer smiMutex.Unlock()

	err := loadFromFileInternal(filePath)
	if err != nil {
		return err
	}

	if !cached {
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

func ExtractNodeDetails(node gosmi.SmiNode) NodeDetails {
	var details NodeDetails

	if node.Type != nil {
		details.Syntax = node.Type.Name

		if len(node.Type.Ranges) > 0 {
			var ranges []string
			for _, r := range node.Type.Ranges {
				ranges = append(ranges, fmt.Sprintf("%v..%v", r.MinValue, r.MaxValue))
			}
			details.Syntax += fmt.Sprintf(" { %s }", strings.Join(ranges, " | "))
		}

		if node.Type.Enum != nil && len(node.Type.Enum.Values) > 0 {
			var enums []string
			for _, e := range node.Type.Enum.Values {
				enums = append(enums, fmt.Sprintf("%s(%v)", e.Name, e.Value))
			}
			details.Syntax += fmt.Sprintf(" { %s }", strings.Join(enums, ", "))
		}
	}

	raw := node.GetRaw()
	if raw != nil {
		if raw.Value.Value != nil {
			details.DefVal = fmt.Sprintf("%v", raw.Value.Value)
		}

		if raw.IndexKind == types.IndexUnknown || raw.IndexKind == types.IndexAugment {
			var idxNames []string
			for _, idxNode := range node.GetIndex() {
				idxNames = append(idxNames, idxNode.Name)
			}
			details.Indexes = strings.Join(idxNames, ", ")
		}
	}

	return details
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

				details := ExtractNodeDetails(n)

				localNodeMap[oidStr] = &MibNode{
					Name:   n.Name,
					OID:    oidStr,
					Label:  fmt.Sprintf("%s##%s", n.Name, oidStr),
					Access: n.Access.String(),

					SourceMib:   mod.Name,
					Description: n.Description,
					Type:        typeStr,
					Details:     &details,
				}
			}
		}
	}

	newRoot := &MibNode{Name: "Root", OID: "", Label: "Root##root"}

	for oidStr, node := range localNodeMap {
		parentOid := getParentOID(oidStr)
		if parentNode, exists := localNodeMap[parentOid]; exists {
			parentNode.Children = append(parentNode.Children, node)
		} else {
			newRoot.Children = append(newRoot.Children, node)
		}
	}

	for _, child := range newRoot.Children {
		flattenBaseOIDs(child)
	}

	sortTree(newRoot)

	TreeMutex.Lock()
	RootNode = newRoot
	TreeMutex.Unlock()
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
		node.Label = fmt.Sprintf("%s##%s", node.Name, node.OID)
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
