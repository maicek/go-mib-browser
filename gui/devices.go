package gui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/gosnmp/gosnmp"
	"github.com/maicek/go-mib-browser/devices"
)

var managingDevices bool

func RenderDeviceToolbox() {
	if len(devices.Devices) != 0 {
		if imgui.BeginCombo("Selected Device", devices.Devices[devices.SelectedDevice].Name) {
			for i, device := range devices.Devices {
				isActive := devices.SelectedDevice == i

				if imgui.SelectableBool(fmt.Sprintf("%s##%d", device.Name, i)) {
					devices.SelectDevice(i)
				}

				if isActive {
					imgui.SetItemDefaultFocus()
				}
			}

			imgui.EndCombo()
		}

	} else {
		imgui.TextDisabled("No devices configured")
	}

	imgui.SameLine()
	if imgui.Button("Manage devices") {
		managingDevices = true
	}

	RenderDevicesConfig()
}

func RenderDevicesConfig() {
	if !managingDevices {
		return
	}

	flag := imgui.WindowFlags(0)
	flag |= imgui.WindowFlagsNoDocking
	flag |= imgui.WindowFlagsNoResize
	flag |= imgui.WindowFlagsNoCollapse

	imgui.SetNextWindowSize(imgui.Vec2{X: 500, Y: 600})

	imgui.BeginV("Devices configuration##devicesconfig", &managingDevices, flag)

	for i, device := range devices.Devices {
		flags := imgui.TreeNodeFlagsDefaultOpen

		if imgui.TreeNodeExStrV(fmt.Sprintf("Device: %s ####deviceconfig_%d", device.Name, i), flags) {

			imgui.InputTextWithHint("Name", "Provide device name", &device.Name, imgui.InputTextFlagsNone, nil)
			imgui.InputTextWithHint("IP", "Provide IP address", &device.IpAddr, imgui.InputTextFlagsNone, nil)

			portFlags := imgui.InputTextFlags(0)
			portFlags |= imgui.InputTextFlagsCharsDecimal

			imgui.InputIntV("Port", &device.SnmpPort, 0, 10, portFlags)

			if imgui.BeginCombo("SNMP Version", getSnmpVersionString(device.SnmpVersion)) {
				if imgui.SelectableBool("SNMP V1") {
					device.SnmpVersion = gosnmp.Version1
				}

				if imgui.SelectableBool("SNMP V2c") {
					device.SnmpVersion = gosnmp.Version2c
				}

				if imgui.SelectableBool("SNMP V3") {
					device.SnmpVersion = gosnmp.Version3
				}

				imgui.EndCombo()
			}

			switch device.SnmpVersion {
			case gosnmp.Version2c, gosnmp.Version1:
				imgui.InputTextWithHint("Community (read)", "Provide community", &device.Community, imgui.InputTextFlagsNone, nil)
				imgui.InputTextWithHint("Community (write)", "Provide community", &device.WriteCommunity, imgui.InputTextFlagsNone, nil)
			case gosnmp.Version3:
				//SNMP V3 security configuration
			}

			btnColor := imgui.Vec4{X: 0.5, Y: 0.1, Z: 0.1, W: 1}
			imgui.PushStyleColorVec4(imgui.ColButton, btnColor)
			if imgui.ButtonV(fmt.Sprintf("Remove device##remove_device_%d", i), imgui.Vec2{X: -1, Y: 0}) {
				device.Remove()
			}

			imgui.PopStyleColor()

			imgui.TreePop()
		}

		imgui.Separator()
	}

	if imgui.ButtonV("Add new device", imgui.Vec2{X: -1, Y: 0}) {
		devices.AddNewDevice()
	}

	availY := imgui.ContentRegionAvail().Y

	itemHeight := imgui.FrameHeightWithSpacing()
	if availY > itemHeight {
		imgui.SetCursorPos(imgui.CursorPos().Add(imgui.Vec2{X: 0, Y: availY}).Add(imgui.Vec2{X: 0, Y: itemHeight}).Add(imgui.Vec2{X: 0, Y: -74.0}))
	}

	imgui.TextDisabled("Settings needs to be saved to keep changes on next startup.")

	btnColor := imgui.Vec4{X: 0.1, Y: 0.4, Z: 0.1, W: 1}
	imgui.PushStyleColorVec4(imgui.ColButton, btnColor)

	if imgui.ButtonV("Save", imgui.Vec2{X: -1, Y: 0}) {
		devices.SaveDevices()
		managingDevices = false
	}
	imgui.PopStyleColor()

	if imgui.BeginItemTooltip() {
		imgui.Text("This will save the list of devices in app data file.")

		imgui.EndTooltip()

	}

	imgui.End()
}

func getSnmpVersionString(v gosnmp.SnmpVersion) string {
	switch v {
	case gosnmp.Version1:
		return "SNMP V1"
	case gosnmp.Version2c:
		return "SNMP V2c"
	case gosnmp.Version3:
		return "SNMP V3"
	default:
		return "Unknown"
	}
}
