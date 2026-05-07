package main

import (
	"os"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/glfwbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/maicek/go-mib-browser/gui"
	"github.com/maicek/go-mib-browser/smi"
)

var (
	showMetricsWindow bool = false
	showDemoWindow    bool = false
)

// Handle OS window
func CreateOsWindow() {
	bnd, err := backend.CreateBackend(glfwbackend.NewGLFWBackend())
	if err != nil {
		panic(err)
	}

	bnd.SetBgColor(imgui.NewVec4(0.09, 0.10, 0.13, 1.0))
	bnd.CreateWindow("MIB Browser", 1280, 720)
	bnd.SetTargetFPS(120)
	bnd.SetSwapInterval(glfwbackend.GLFWSwapIntervalVsync)

	io := imgui.CurrentIO()

	// Enable docking
	flags := io.ConfigFlags() | imgui.ConfigFlagsDockingEnable
	// Disable viewports (single window mode)
	flags = flags &^ imgui.ConfigFlagsViewportsEnable
	io.SetConfigFlags(flags)

	io.SetIniFilename("imgui.ini")

	bnd.SetCloseCallback(func() {
		os.Exit(0)
	})

	smi.Init()
	gui.Init()

	gui.InitFont()
	gui.ApplyTheme()

	bnd.Run(func() {
		gui.RenderStandaloneLayout()

		if showMetricsWindow {
			imgui.ShowMetricsWindowV(&showMetricsWindow)
		}

		if showDemoWindow {
			imgui.ShowDemoWindowV(&showDemoWindow)
		}

		if imgui.BeginMainMenuBar() {
			if imgui.BeginMenu("View") {
				if imgui.MenuItemBoolV("ImGui Metrics (DevTools)", "", showMetricsWindow, true) {
					showMetricsWindow = !showMetricsWindow
				}
				if imgui.MenuItemBoolV("ImGui Demo Window", "", showDemoWindow, true) {
					showDemoWindow = !showDemoWindow
				}
				imgui.EndMenu()
			}
			imgui.EndMainMenuBar()
		}
	})
}
