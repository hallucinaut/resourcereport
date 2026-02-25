# resourcereport

Resource Usage Report Generator - generates reports on Docker container resource usage.

## Purpose

Collect and report CPU, memory, and network usage statistics from Docker containers in multiple formats.

## Installation

```bash
go build -o resourcereport ./cmd/resourcereport
```

## Usage

```bash
resourcereport [--json|--html] [service1] [service2] ...
```

### Examples

```bash
# Generate text report for all containers
resourcereport

# Generate JSON report
resourcereport --json

# Generate HTML report
resourcereport --html > report.html

# Filter by service name
resourcereport api worker
```

## Output

### Text Format

```
=== RESOURCE USAGE REPORT ===

SERVICE             CPU%   MEM USED   MEM TOTAL   STATUS
----------------------------------------------------------------------
api-server           45.2%     512.0MB    1024.0MB   NORMAL
worker-1             78.5%     768.0MB    1024.0MB   HIGH
database             12.3%     2048.0MB   4096.0MB   LOW

Report generated at: 2026-02-25 14:30:00
```

### JSON Format

```json
[
  {
    "Service": "api-server",
    "CPUPercent": 45.2,
    "MemoryUsed": 512.0,
    "MemoryTotal": 1024.0,
    "Status": "NORMAL"
  }
]
```

### HTML Format

Generates a complete HTML report with CSS styling suitable for email distribution.

## Status Levels

- LOW: CPU usage <50%
- NORMAL: CPU usage 50-80%
- HIGH: CPU usage >80%

## Dependencies

- Go 1.21+
- github.com/fatih/color
- Docker installed

## Build and Run

```bash
# Build
go build -o resourcereport ./cmd/resourcereport

# Run
go run ./cmd/resourcereport

# Generate JSON
go run ./cmd/resourcereport --json

# Generate HTML
go run ./cmd/resourcereport --html > report.html
```

## License

MIT