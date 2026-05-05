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

	bnd.SetBgColor(imgui.NewVec4(0.35, 0.45, 0.5, 0.5))
	bnd.CreateWindow("MIB Browser", 800, 600)

	io := imgui.CurrentIO()
	flags := io.ConfigFlags()
	flags = flags &^ imgui.ConfigFlagsViewportsEnable
	io.SetConfigFlags(flags)

	io.SetIniFilename("imgui.ini")

	bnd.SetCloseCallback(func() {
		os.Exit(1)
	})

	smi.Init()

	gui.InitFont()

	bnd.Run(func() {
		gui.RenderGui()
	})
}
