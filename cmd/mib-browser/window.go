package main

import (
	"os"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/glfwbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/maicek/go-mib-browser/gui"
	"github.com/maicek/go-mib-browser/smi"
)

// Handle OS window
func CreateOsWindow() {
	bnd, err := backend.CreateBackend(glfwbackend.NewGLFWBackend())
	if err != nil {
		panic(err)
	}

	bnd.SetBgColor(imgui.NewVec4(0.09, 0.10, 0.13, 1.0))
	bnd.CreateWindow("MIB Browser", 1280, 720)

	io := imgui.CurrentIO()

	// Enable docking
	flags := io.ConfigFlags() | imgui.ConfigFlagsDockingEnable
	// Disable viewports (single window mode)
	flags = flags &^ imgui.ConfigFlagsViewportsEnable
	io.SetConfigFlags(flags)

	io.SetIniFilename("imgui.ini")

	bnd.SetCloseCallback(func() {
		os.Exit(1)
	})

	smi.Init()
	gui.Init()

	gui.InitFont()
	gui.ApplyTheme()

	bnd.Run(func() {
		gui.RenderStandaloneLayout()
	})
}
