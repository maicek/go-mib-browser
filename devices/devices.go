package devices

import "github.com/gosnmp/gosnmp"

type Device struct {
	Name        string
	IpAddr      string
	SnmpVersion gosnmp.SnmpVersion
	SnmpPort    int32

	// SNMP v1, v2c
	Community      string
	WriteCommunity string

	// SNMP v3
	SecurityName     string
	AuthKey          string
	AuthProtocol     gosnmp.SnmpV3AuthProtocol
	PrivProtocol     gosnmp.SnmpV3PrivProtocol
	PrivKey          string
	SecurityEngineID string
}

var Devices []*Device = make([]*Device, 0)
var SelectedDevice int

func AddNewDevice() {
	Devices = append(Devices, &Device{
		Name:           "",
		IpAddr:         "",
		SnmpVersion:    gosnmp.Version2c,
		SnmpPort:       161,
		Community:      "public",
		WriteCommunity: "private",
	})
}

func (d *Device) Remove() {
	Devices = append(Devices[:SelectedDevice], Devices[SelectedDevice+1:]...)
	if SelectedDevice >= len(Devices) {
		SelectedDevice = 0
	}
}
