package gui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/maicek/go-mib-browser/smi"
	"golang.design/x/clipboard"
)

var SelectedNode *smi.MibNode

func RenderMibNode(node *smi.MibNode) {
	flags := imgui.TreeNodeFlagsNone

	if len(node.Children) == 0 {
		flags |= imgui.TreeNodeFlagsLeaf
	}

	if SelectedNode != nil && SelectedNode.OID == node.OID {
		flags |= imgui.TreeNodeFlagsSelected
	}

	isOpen := imgui.TreeNodeExStrV(fmt.Sprintf("%s##%s", node.Name, node.OID), flags)

	if imgui.BeginPopupContextItem() {
		if imgui.MenuItemBool("Get") {
			go mainResultTable.Get(node.OID)
		}

		if imgui.MenuItemBool("Walk") {
			go mainResultTable.Walk(node.OID)
		}

		if imgui.MenuItemBool("Copy OID") {
			clipboard.Write(clipboard.FmtText, []byte(node.OID))
		}

		imgui.EndPopup()
	}

	if imgui.IsItemClicked() {
		SelectedNode = node
	}

	if isOpen {
		if len(node.Children) > 0 {

			drawList := imgui.WindowDrawList()
			cursorPos := imgui.CursorScreenPos()

			lineX := cursorPos.X - 12.0
			lineStartY := cursorPos.Y - 4.0
			lineEndY := lineStartY

			lineColor := imgui.ColorConvertFloat4ToU32(imgui.Vec4{X: 0.6, Y: 0.6, Z: 0.6, W: 0.6})

			halfHeight := imgui.TextLineHeight() / 2.0

			for _, child := range node.Children {

				childPos := imgui.CursorScreenPos()

				p1 := imgui.Vec2{X: lineX, Y: childPos.Y + halfHeight}
				p2 := imgui.Vec2{X: childPos.X - 4.0, Y: childPos.Y + halfHeight}

				drawList.AddLine(p1, p2, lineColor)

				lineEndY = p1.Y

				RenderMibNode(child)
			}

			pStart := imgui.Vec2{X: lineX, Y: lineStartY}
			pEnd := imgui.Vec2{X: lineX, Y: lineEndY}
			drawList.AddLine(pStart, pEnd, lineColor)
		}

		imgui.TreePop()
	}
}
