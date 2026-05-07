package gui

import (
	"fmt"
	"time"

	"github.com/maicek/go-mib-browser/devices"
	"github.com/maicek/go-mib-browser/snmp"
)

func (m *MainResultTable) Get(oid string) {
	client, err := snmp.SetupSnmp(devices.GetSelected())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer client.Snmp.Conn.Close()

	fmt.Printf("Getting %s from %s\n", oid, devices.GetSelected().IpAddr)
	result, err := client.Get(oid)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	m.AddResult(result)
}

func (m *MainResultTable) Walk(oid string) {
	client, err := snmp.SetupSnmp(devices.GetSelected())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer client.Snmp.Conn.Close()

	fmt.Printf("Walking %s from %s\n", oid, devices.GetSelected().IpAddr)
	results, errChan, _ := client.Walk(oid)

	fmt.Printf("Walk this shit")

	// handle channel messages, pop one message after another
	for {
		select {
		case result, ok := <-results:
			if !ok {
				return
			}
			m.AddResult(result)
		case err := <-errChan:
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}
