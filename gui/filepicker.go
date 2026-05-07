package gui

import (
	"github.com/maicek/go-mib-browser/smi"
	"github.com/sqweek/dialog"
)

var (
	isChoosingFile bool
	lastLoadedFile string
)

func openFilePicker() {
	if isChoosingFile {
		return
	}

	go func() {
		defer func() { isChoosingFile = false }()

		filePath, err := dialog.File().
			Title("Select MIB file").
			Filter("MIB", "mib", "txt", "my").
			Filter("All files", "*").
			Load()

		if err != nil {
			return
		}

		smi.LoadFromFile(filePath, false)
	}()
}
