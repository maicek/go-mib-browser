package gui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

func RenderToolbox() {
	flags := imgui.WindowFlags(0)
	flags |= imgui.WindowFlagsNoResize
	flags |= imgui.WindowFlagsNoCollapse
	flags |= imgui.WindowFlagsAlwaysAutoResize

	imgui.BeginV("Toolbox##toolbox", nil, flags)

	RenderDeviceToolbox()

	imgui.End()
}
