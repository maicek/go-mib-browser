package main

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/glfwbackend"
	"github.com/AllenDang/cimgui-go/imgui"
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
		fmt.Println("window is closing")
	})

	bnd.Run(func() {
		loop()
	})
}

func loop() {
	imgui.Begin("MIB Browser")

	imgui.End()
}
