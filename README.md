An SNMP MIB Browser written in Go with:

- [cimgui-go](https://github.com/AllenDang/cimgui-go) GUI library
- [gosnmp](https://github.com/gosnmp/gosnmp) for SNMP communication
- [gosmi](https://github.com/sleepinggenius2/gosmi) for MIB parsing

## Motivation

Project was primarily made as an internal devtool because other tools lacked snmpv3 support and also has limited number of loaded MIBs at once.

## Key features

- SNMP v1, v2c, v3 support
- Unlimited loaded MIBs
- Tree view of loaded MIBs
- Lightweight and responsive GUI
- Cross-platform (Windows, Linux, macOS)
- Persistent device list
- Works as standalone app or could be embedded into existing imgui devtool. (Lacks proper API for now but I plan to add it.)

## Building

```bash
make build
```

## Development

Requires [air](https://github.com/air-verse/air) installed.

```bash
air
```

## Todo

- [ ] Walk table.
- [ ] SNMP Set.
- [ ] Tree view.
- [ ] Polling and graph of values.
- [ ] Improved results table.
  - [ ] Clear values.
- [ ] Add SNMP Trap receiver.

## License

The project is licensed under the [MIT License](LICENSE).
