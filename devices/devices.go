package devices

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
)

type Device struct {
	Name        string             `json:"name"`
	IpAddr      string             `json:"ip_addr"`
	SnmpVersion gosnmp.SnmpVersion `json:"snmp_version"`
	SnmpPort    int32              `json:"snmp_port"`

	// SNMP v1, v2c
	Community      string `json:"community"`
	WriteCommunity string `json:"write_community"`

	// SNMP v3
	UserName     string                    `json:"user_name"`
	AuthPassword string                    `json:"auth_password"`
	PrivPassword string                    `json:"priv_password"`
	AuthProtocol gosnmp.SnmpV3AuthProtocol `json:"auth_protocol"`
	PrivProtocol gosnmp.SnmpV3PrivProtocol `json:"priv_protocol"`
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

func (d *Device) Remove() error {
	Devices = append(Devices[:SelectedDevice], Devices[SelectedDevice+1:]...)
	if SelectedDevice >= len(Devices) {
		SelectedDevice = 0
	}
	return nil
}

func SelectDevice(id int) error {
	if id < 0 || id >= len(Devices) {
		return fmt.Errorf("Invalid device index")
	}
	SelectedDevice = id
	return nil
}

func GetSelected() *Device {
	if SelectedDevice >= len(Devices) {
		return nil
	}
	return Devices[SelectedDevice]
}
