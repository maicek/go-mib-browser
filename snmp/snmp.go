package snmp

import (
	"context"
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/maicek/go-mib-browser/devices"
)

type SnmpClient struct {
	Device *devices.Device
	Snmp   *gosnmp.GoSNMP
}

func SetupSnmp(device *devices.Device) (*SnmpClient, error) {

	fmt.Printf("Device: %+v\n", device)
	snmp := &gosnmp.GoSNMP{
		Target:                  device.IpAddr,
		Port:                    uint16(device.SnmpPort),
		Community:               device.Community,
		Transport:               "udp",
		UseUnconnectedUDPSocket: true,
		Retries:                 1,
		Version:                 device.SnmpVersion,
		Timeout:                 10 * time.Second,
		MaxOids:                 10,
		MaxRepetitions:          100,
	}

	client := &SnmpClient{
		Device: device,
		Snmp:   snmp,
	}

	err := client.Snmp.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *SnmpClient) Get(oid string) (*gosnmp.SnmpPDU, error) {
	result, err := s.Snmp.Get([]string{oid})
	if err != nil {
		return nil, err
	}

	if len(result.Variables) > 0 {
		return &result.Variables[0], nil
	}

	return nil, fmt.Errorf("OID %s not found", oid)
}

func (s *SnmpClient) Walk(oid string) (chan *gosnmp.SnmpPDU, chan error, context.CancelFunc) {
	results := make(chan *gosnmp.SnmpPDU, 100)
	errorChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		err := s.Snmp.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
			fmt.Printf("Result: %v\n", pdu)
			select {
			case results <- &pdu:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})

		close(results)
		errorChan <- err
		close(errorChan)
	}()

	return results, errorChan, cancel
}
