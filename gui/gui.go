package gui

import "github.com/AllenDang/cimgui-go/imgui"

// Handle gui rendering
func RenderGui() {

	imgui.Begin("MIB Browser")

	imgui.Text("Test")

	imgui.End()
}
