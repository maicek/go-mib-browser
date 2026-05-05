package smi

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrg/xdg"
)

// Loaded mibs are stored just like semicolon-separated values of file paths.

func GetCustomMibs() ([]string, error) {
	datafilePath, err := xdg.DataFile("maicek_mib_browser/custom_mibs.txt")
	if err != nil {
		return []string{}, err
	}

	fileData, err := os.ReadFile(datafilePath)
	if err != nil {
		return []string{}, err
	}

	paths := strings.Split(string(fileData), ";")

	return paths, nil
}

func deduplicateCacheEntries(entries *[]string) {
	seen := make(map[string]struct{}, len(*entries))

	j := 0
	for _, entry := range *entries {
		if _, ok := seen[entry]; !ok {
			seen[entry] = struct{}{}
			(*entries)[j] = entry
			j++
		}
	}

	*entries = (*entries)[:j]
}

func PushCustomMib(mibPath string) {
	currentData, err := GetCustomMibs()
	if err != nil {
		fmt.Printf("Error happened when loading current file: %s\n", err)
	}
	currentData = append(currentData, mibPath)

	deduplicateCacheEntries(&currentData)

	raw := strings.Join(currentData, ";")

	datafilePath, err := xdg.DataFile("maicek_mib_browser/custom_mibs.txt")
	if err != nil {
		return
	}

	file, err := os.Create(datafilePath)
	if err != nil {
		return
	}

	fmt.Printf("Wrote data file: %s with \"%s\"\n", datafilePath, raw)

	file.Write([]byte(raw))
}
