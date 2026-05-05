package gui

import "github.com/AllenDang/cimgui-go/imgui"

var (
	dockInitialized bool
	dockspaceID     imgui.ID
)

// RenderStandaloneLayout renders all panels in a fixed, docked layout.
// Call this from the standalone binary's render loop.
func RenderStandaloneLayout() {
	viewport := imgui.MainViewport()
	pos := viewport.Pos()
	size := viewport.Size()

	imgui.SetNextWindowPos(pos)
	imgui.SetNextWindowSize(size)
	imgui.SetNextWindowViewport(viewport.ID())

	hostFlags := imgui.WindowFlagsNoTitleBar |
		imgui.WindowFlagsNoCollapse |
		imgui.WindowFlagsNoResize |
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoBringToFrontOnFocus |
		imgui.WindowFlagsNoDocking |
		imgui.WindowFlagsNoScrollbar |
		imgui.WindowFlagsNoScrollWithMouse

	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 0, Y: 0})

	imgui.BeginV("##DockHost", nil, hostFlags)

	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleVar()

	dockspaceID = imgui.IDStr("MainDockSpace")
	imgui.DockSpace(dockspaceID)

	if !dockInitialized {
		initDockLayout(dockspaceID, size)
		dockInitialized = true
	}

	imgui.End()

	// Render individual panels
	RenderMibTree()
	RenderBottomPanel()
	RenderToolbox()
	RenderResultsWindow()
}

func initDockLayout(id imgui.ID, size imgui.Vec2) {
	imgui.InternalDockBuilderRemoveNode(id)
	imgui.InternalDockBuilderAddNodeV(id, imgui.DockNodeFlagsNone)
	imgui.InternalDockBuilderSetNodeSize(id, size)

	// Split: left (tree ~37%) | right remainder
	var rightID imgui.ID
	var leftID imgui.ID
	imgui.InternalDockBuilderSplitNode(id, imgui.DirLeft, 0.37, &leftID, &rightID)

	// Split left vertically: top (tree ~75%) | bottom (details + loaded mibs)
	var leftTopID imgui.ID
	var leftBotID imgui.ID
	imgui.InternalDockBuilderSplitNode(leftID, imgui.DirUp, 0.75, &leftTopID, &leftBotID)

	// Split right: top (toolbox) | bottom (results)
	var rightTopID imgui.ID
	var rightBotID imgui.ID
	imgui.InternalDockBuilderSplitNode(rightID, imgui.DirUp, 0.15, &rightTopID, &rightBotID)

	imgui.InternalDockBuilderDockWindow("MIB Tree##treeview", leftTopID)
	imgui.InternalDockBuilderDockWindow("##BottomPanel", leftBotID)
	imgui.InternalDockBuilderDockWindow("Toolbox##toolbox", rightTopID)
	imgui.InternalDockBuilderDockWindow("Results", rightBotID)

	imgui.InternalDockBuilderFinish(id)
}

// RenderBottomPanel renders Node Details + Loaded MIBs as tabs in the bottom-left panel.
func RenderBottomPanel() {
	imgui.Begin("##BottomPanel")

	if imgui.BeginTabBar("##bottomtabs") {
		if imgui.BeginTabItem("Node Details") {
			renderNodeDetailsContent()
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Loaded MIB's") {
			renderLoadedMibsContent()
			imgui.EndTabItem()
		}
		imgui.EndTabBar()
	}

	imgui.End()
}
