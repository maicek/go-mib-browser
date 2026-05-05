package gui

import (
	"embed"
	"unsafe"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/maicek/go-mib-browser/smi"
	"github.com/sleepinggenius2/gosmi"
)

//go:embed Roboto.ttf
var fontFs embed.FS

func InitFont() {
	io := imgui.CurrentIO()

	loadedFontBytes, err := fontFs.ReadFile("Roboto.ttf")
	if err != nil {
		panic("Nie udało się odczytać fontu z embed: " + err.Error())
	}

	fontConfig := imgui.NewFontConfig()

	fontConfig.SetFontDataOwnedByAtlas(false)

	ranges := io.Fonts().GlyphRangesDefault()
	fontPtr := uintptr(unsafe.Pointer(&loadedFontBytes[0]))
	fontDataSize := int32(len(loadedFontBytes))

	io.Fonts().AddFontFromMemoryTTFV(
		fontPtr,
		fontDataSize,
		16.0,
		fontConfig,
		ranges,
	)
}

// Handle gui rendering
func RenderGui() {

	RenderMibTree()
	RenderLoadedMibs()
	RenderNodeDetails()
	RenderResultsWindow()
	RenderToolbox()
}

func RenderLoadedMibs() {
	imgui.Begin("Loaded MIB's##mibloader")

	if imgui.Button("Load MIB's") {
		openFilePicker()
	}

	if imgui.TreeNodeExStr("Loaded modules") {
		for _, mod := range gosmi.GetLoadedModules() {
			imgui.BulletText(mod.Name)
		}

		imgui.TreePop()
	}

	imgui.End()
}

func RenderMibTree() {
	imgui.Begin("MIB Tree##treeview")

	if smi.RootNode != nil {
		RenderMibNode(smi.RootNode)
	} else {
		imgui.Text("Load any mib file first")
	}

	imgui.End()
}

func RenderNodeDetails() {
	imgui.Begin("Node Details##nodedetails")

	if SelectedNode != nil {
		imgui.Text("Name: " + SelectedNode.Name)
		imgui.Text("OID: " + SelectedNode.OID)
		if SelectedNode.Type != "" {
			imgui.Text("Type: " + SelectedNode.Type)
		}
		imgui.Separator()
		imgui.TextWrapped(SelectedNode.Description)
	} else {
		imgui.Text("Select a node from the tree.")
	}

	imgui.End()
}
