package gui

import "github.com/AllenDang/cimgui-go/imgui"

func ApplyTheme() {
	imgui.StyleColorsDark()

	style := imgui.CurrentStyle()

	// Geometry
	style.SetWindowRounding(6)
	style.SetChildRounding(4)
	style.SetFrameRounding(4)
	style.SetPopupRounding(4)
	style.SetScrollbarRounding(4)
	style.SetGrabRounding(4)
	style.SetTabRounding(4)
	style.SetWindowPadding(imgui.Vec2{X: 10, Y: 10})
	style.SetFramePadding(imgui.Vec2{X: 6, Y: 4})
	style.SetItemSpacing(imgui.Vec2{X: 8, Y: 5})
	style.SetItemInnerSpacing(imgui.Vec2{X: 6, Y: 4})
	style.SetIndentSpacing(16)
	style.SetScrollbarSize(10)
	style.SetWindowBorderSize(1)
	style.SetFrameBorderSize(0)
	style.SetChildBorderSize(1)

	// Build color palette
	colors := style.Colors()

	c := func(r, g, b, a float32) imgui.Vec4 { return imgui.Vec4{X: r, Y: g, Z: b, W: a} }

	colors[imgui.ColText]             = c(0.90, 0.92, 0.96, 1.00)
	colors[imgui.ColTextDisabled]     = c(0.45, 0.47, 0.54, 1.00)
	colors[imgui.ColWindowBg]         = c(0.09, 0.10, 0.13, 1.00)
	colors[imgui.ColChildBg]          = c(0.10, 0.11, 0.14, 1.00)
	colors[imgui.ColPopupBg]          = c(0.10, 0.11, 0.15, 0.97)
	colors[imgui.ColBorder]           = c(0.20, 0.22, 0.29, 1.00)
	colors[imgui.ColBorderShadow]     = c(0, 0, 0, 0)
	colors[imgui.ColFrameBg]          = c(0.14, 0.15, 0.19, 1.00)
	colors[imgui.ColFrameBgHovered]   = c(0.18, 0.20, 0.26, 1.00)
	colors[imgui.ColFrameBgActive]    = c(0.22, 0.24, 0.32, 1.00)
	colors[imgui.ColTitleBg]          = c(0.07, 0.08, 0.11, 1.00)
	colors[imgui.ColTitleBgActive]    = c(0.12, 0.14, 0.20, 1.00)
	colors[imgui.ColTitleBgCollapsed] = c(0.07, 0.08, 0.11, 0.75)
	colors[imgui.ColMenuBarBg]        = c(0.07, 0.08, 0.11, 1.00)
	colors[imgui.ColScrollbarBg]         = c(0.07, 0.08, 0.11, 1.00)
	colors[imgui.ColScrollbarGrab]       = c(0.24, 0.28, 0.38, 1.00)
	colors[imgui.ColScrollbarGrabHovered]= c(0.32, 0.38, 0.52, 1.00)
	colors[imgui.ColScrollbarGrabActive] = c(0.38, 0.46, 0.64, 1.00)
	colors[imgui.ColCheckMark]      = c(0.46, 0.70, 1.00, 1.00)
	colors[imgui.ColSliderGrab]     = c(0.38, 0.58, 0.94, 1.00)
	colors[imgui.ColSliderGrabActive]= c(0.46, 0.68, 1.00, 1.00)
	colors[imgui.ColButton]         = c(0.22, 0.36, 0.70, 1.00)
	colors[imgui.ColButtonHovered]  = c(0.30, 0.46, 0.86, 1.00)
	colors[imgui.ColButtonActive]   = c(0.18, 0.30, 0.58, 1.00)
	colors[imgui.ColHeader]         = c(0.22, 0.34, 0.62, 0.55)
	colors[imgui.ColHeaderHovered]  = c(0.28, 0.42, 0.76, 0.80)
	colors[imgui.ColHeaderActive]   = c(0.32, 0.48, 0.88, 1.00)
	colors[imgui.ColSeparator]         = c(0.20, 0.22, 0.30, 1.00)
	colors[imgui.ColSeparatorHovered]  = c(0.32, 0.46, 0.78, 0.80)
	colors[imgui.ColSeparatorActive]   = c(0.38, 0.54, 0.92, 1.00)
	colors[imgui.ColResizeGrip]        = c(0.28, 0.42, 0.78, 0.25)
	colors[imgui.ColResizeGripHovered] = c(0.36, 0.54, 0.90, 0.67)
	colors[imgui.ColResizeGripActive]  = c(0.40, 0.60, 1.00, 0.95)
	colors[imgui.ColTab]               = c(0.11, 0.13, 0.18, 1.00)
	colors[imgui.ColTabHovered]        = c(0.28, 0.42, 0.76, 0.80)
	colors[imgui.ColTabSelected]       = c(0.22, 0.34, 0.62, 1.00)
	colors[imgui.ColTabDimmed]         = c(0.07, 0.08, 0.11, 1.00)
	colors[imgui.ColTabDimmedSelected] = c(0.16, 0.22, 0.38, 1.00)
	colors[imgui.ColDockingPreview]    = c(0.38, 0.58, 0.94, 0.70)
	colors[imgui.ColDockingEmptyBg]    = c(0.07, 0.08, 0.11, 1.00)

	style.SetColors(&colors)
}
