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

// RenderGui renders all panels as free-floating windows (lib mode).
func RenderGui() {
	RenderMibTree()
	RenderLoadedMibs()
	RenderNodeDetails()
	RenderResultsWindow()
	RenderToolbox()
}

// ─── MIB Tree ────────────────────────────────────────────────────────────────

func RenderMibTree() {
	imgui.Begin("MIB Tree##treeview")
	if smi.RootNode != nil {
		RenderMibNode(smi.RootNode)
	} else {
		imgui.TextDisabled("No MIB loaded. Use the Toolbox to load a file.")
	}
	imgui.End()
}

// ─── Loaded MIBs ─────────────────────────────────────────────────────────────

func RenderLoadedMibs() {
	imgui.Begin("Loaded MIB's##mibloader")
	renderLoadedMibsContent()
	imgui.End()
}

func renderLoadedMibsContent() {
	if imgui.Button("  Load MIB file  ") {
		openFilePicker()
	}
	imgui.Spacing()
	imgui.Separator()
	imgui.Spacing()

	mods := gosmi.GetLoadedModules()
	if len(mods) == 0 {
		imgui.TextDisabled("No modules loaded.")
		return
	}

	if imgui.TreeNodeExStrV("Loaded modules", imgui.TreeNodeFlagsDefaultOpen) {
		for _, mod := range mods {
			imgui.BulletText(mod.Name)
		}
		imgui.TreePop()
	}
}

// ─── Node Details ─────────────────────────────────────────────────────────────

func RenderNodeDetails() {
	imgui.Begin("Node Details##nodedetails")
	renderNodeDetailsContent()
	imgui.End()
}

func renderNodeDetailsContent() {
	if SelectedNode == nil {
		imgui.TextDisabled("Select a node from the tree.")
		return
	}

	n := SelectedNode

	imgui.PushStyleColorVec4(imgui.ColText, imgui.Vec4{X: 0.46, Y: 0.70, Z: 1.00, W: 1.00})
	imgui.Text(n.Name)
	imgui.PopStyleColor()

	imgui.Spacing()

	imgui.TextDisabled("OID")
	imgui.SameLine()
	imgui.Text(n.OID)

	if n.Type != "" {
		imgui.TextDisabled("Type")
		imgui.SameLine()
		imgui.Text(n.Type)
	}

	if n.Description != "" {
		imgui.Spacing()
		imgui.Separator()
		imgui.Spacing()
		imgui.TextWrapped(n.Description)
	}
}
