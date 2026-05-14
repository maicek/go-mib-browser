# go-mib-browser

Fast, lightweight SNMP MIB browser built with Go and ImGui.

A high-performance diagnostic tool designed for network engineers who need a responsive and reliable way to explore SNMP-enabled devices. Unlike traditional MIB browsers, this tool focuses on speed, modern protocol support, and the ability to handle massive MIB collections without performance degradation.

## Motivation

This project was originally developed as an internal devtool to address the limitations of existing software—specifically the lack of robust SNMPv3 support and performance bottlenecks when loading large numbers of MIB files simultaneously.

## Key Features

- **Full SNMP Support**: Native support for SNMP v1, v2c, and v3 (including USM).
- **GPU-Accelerated GUI**: Built with [cimgui-go](https://github.com/AllenDang/cimgui-go) for a 60FPS, reactive user interface that remains smooth even with thousands of OIDs loaded.
- **Scalable MIB Handling**: Utilizes [gosmi](https://github.com/sleepinggenius2/gosmi) for robust MIB parsing with no arbitrary limits on loaded modules.
- **Persistent Management**: Built-in device list persistence for quick access to frequently monitored hardware.
- **Standalone or Embedded**: Designed as a self-contained application, but architected to be integrated into larger diagnostic toolsets.

## Technical Stack

- **Language**: Go 1.21+
- **SNMP Engine**: [gosnmp](https://github.com/gosnmp/gosnmp)
- **UI Framework**: [cimgui-go](https://github.com/AllenDang/cimgui-go) (Dear ImGui bindings)
- **MIB Parser**: [gosmi](https://github.com/sleepinggenius2/gosmi)

## Installation

Download binary from releases :D.

## Getting Started (Development)

### Prerequisites

Ensure you have Go installed on your system.

### Installation

```bash
git clone https://github.com/maicek/go-mib-browser.git
cd go-mib-browser
make build
```

### Development

For hot-reloading during development, use [air](https://github.com/air-verse/air):

```bash
air
```

## Roadmap

- [ ] **Data Visualization**: Real-time graphing of polled values.
- [ ] **Tabular Views**: Comprehensive SNMP Walk and Table data rendering.
- [ ] **Active Monitoring**: Integrated SNMP Trap receiver.
- [ ] **Advanced Interactions**: Support for SNMP Set operations.
- [ ] **Enhanced Filtering**: Search and filter results by OID, name, or value type.

## License

This project is licensed under the [MIT License](LICENSE).
