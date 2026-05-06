package devices

import (
	"encoding/json"
	"os"

	"github.com/adrg/xdg"
)

func LoadDevices() ([]*Device, error) {
	datafilePath, err := xdg.DataFile("maicek_mib_browser/devices.json")
	if err != nil {
		return []*Device{}, err
	}

	fileData, err := os.ReadFile(datafilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Device{}, nil
		}
		return []*Device{}, err
	}

	var devices []*Device
	err = json.Unmarshal(fileData, &devices)
	if err != nil {
		return []*Device{}, err
	}

	return devices, nil
}

func SaveDevices() error {
	datafilePath, err := xdg.DataFile("maicek_mib_browser/devices.json")
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(Devices, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(datafilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
