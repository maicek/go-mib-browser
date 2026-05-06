package devices

import "github.com/gosnmp/gosnmp"

type Device struct {
	Name        string             `json:"name"`
	IpAddr      string             `json:"ip_addr"`
	SnmpVersion gosnmp.SnmpVersion `json:"snmp_version"`
	SnmpPort    int32              `json:"snmp_port"`

	// SNMP v1, v2c
	Community      string `json:"community"`
	WriteCommunity string `json:"write_community"`

	// SNMP v3
	SecurityName     string                    `json:"security_name"`
	AuthKey          string                    `json:"auth_key"`
	AuthProtocol     gosnmp.SnmpV3AuthProtocol `json:"auth_protocol"`
	PrivProtocol     gosnmp.SnmpV3PrivProtocol `json:"priv_protocol"`
	PrivKey          string                    `json:"priv_key"`
	SecurityEngineID string                    `json:"security_engine_id"`
}

var Devices []*Device = make([]*Device, 0)
var SelectedDevice int

func init() {
	Devices, _ = LoadDevices()
}

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
