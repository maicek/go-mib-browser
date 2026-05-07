package gui

import (
	"fmt"
	"sync"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/gosnmp/gosnmp"
	"github.com/maicek/go-mib-browser/smi"
	"github.com/sleepinggenius2/gosmi"
)

type ResultItem struct {
	MibNode *gosmi.SmiNode
	PDU     *gosnmp.SnmpPDU
}

type MainResultTable struct {
	lock    *sync.Mutex
	Results []ResultItem
}

var mainResultTable *MainResultTable = &MainResultTable{
	lock:    &sync.Mutex{},
	Results: make([]ResultItem, 0),
}

func (m *MainResultTable) AddResult(p *gosnmp.SnmpPDU) {
	m.lock.Lock()
	defer m.lock.Unlock()

	node, err := smi.GetOidInfo(p.Name)
	if err != nil {
		return
	}

	m.Results = append(m.Results, ResultItem{
		MibNode: node,
		PDU:     p,
	})
}

func (m *MainResultTable) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Results = make([]ResultItem, 0)
}

func RenderResultsWindow() {
	flags := imgui.WindowFlagsMenuBar

	imgui.BeginV("Results", nil, flags)

	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Tools") {
			if imgui.MenuItemBool("Clear") {
				mainResultTable.Clear()
			}
			imgui.EndMenu()
		}

		imgui.EndMenuBar()
	}

	if imgui.BeginTableV("Result table", 2, imgui.TableFlagsSortable|imgui.TableFlagsResizable|imgui.TableFlagsBorders, imgui.Vec2{X: 0, Y: 0}, 0.0) {
		imgui.TableSetupColumn("Name")
		imgui.TableSetupColumn("Value")

		mainResultTable.lock.Lock()
		defer mainResultTable.lock.Unlock()

		for _, pdu := range mainResultTable.Results {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.Text(fmt.Sprintf("%s (%s)", pdu.MibNode.Name, pdu.PDU.Name))
			imgui.TableSetColumnIndex(1)

			if pdu.MibNode != nil {
				if pdu.MibNode.Type != nil {
					val := pdu.MibNode.FormatValue(pdu.PDU.Value)
					imgui.Text(fmt.Sprintf("%v", val))
				} else {
					imgui.Text(fmt.Sprintf("%v", pdu.PDU.Value))
				}
			}
		}

		imgui.EndTable()
	}

	imgui.End()
}
