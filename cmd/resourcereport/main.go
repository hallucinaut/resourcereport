package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ResourceMetrics struct {
	Service    string
	CPUPercent float64
	MemoryUsed float64
	MemoryTotal float64
	DiskUsed   float64
	DiskTotal  float64
	NetworkRX  float64
	NetworkTX  float64
	Status     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(color.CyanString("resourcereport - Resource Usage Report Generator"))
		fmt.Println()
		fmt.Println("Usage: resourcereport [--json] [--html] [service1] [service2] ...")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  resourcereport")
		fmt.Println("  resourcereport --json")
		fmt.Println("  resourcereport --html > report.html")
		os.Exit(1)
	}

	var format string
	services := []string{}

	for _, arg := range os.Args[1:] {
		if arg == "--json" {
			format = "json"
		} else if arg == "--html" {
			format = "html"
		} else if arg != "" {
			services = append(services, arg)
		}
	}

	if format == "" {
		format = "text"
	}

	metrics := collectMetrics(services)
	generateReport(metrics, format)
}

func collectMetrics(services []string) []ResourceMetrics {
	var metrics []ResourceMetrics

	// Get all running containers
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.Container}}|{{.CPUPerc}}|{{.MemUsage}}|{{.NetIO}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(color.YellowString("Warning: Could not get Docker stats"))
		return metrics
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			service := parts[0]
			
			// Filter services if specified
			if len(services) > 0 {
				found := false
				for _, s := range services {
					if strings.Contains(service, s) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			metric := parseMetrics(parts)
			metrics = append(metrics, metric)
		}
	}

	return metrics
}

func parseMetrics(parts []string) ResourceMetrics {
	metric := ResourceMetrics{}
	metric.Service = parts[0]

	// Parse CPU percentage
	cpuRegex := regexp.MustCompile(`([\d.]+)%`)
	cpuMatch := cpuRegex.FindString(parts[1])
	if cpuMatch != "" {
		cpu, _ := strconv.ParseFloat(strings.TrimSuffix(cpuMatch, "%"), 64)
		metric.CPUPercent = cpu
	}

	// Parse memory usage
	memRegex := regexp.MustCompile(`([\d.]+) ([^/]+)/([\d.]+) ([^/]+)`)
	memMatch := memRegex.FindStringSubmatch(parts[2])
	if len(memMatch) >= 5 {
		metric.MemoryUsed, _ = strconv.ParseFloat(memMatch[1], 64)
		metric.MemoryTotal, _ = strconv.ParseFloat(memMatch[3], 64)
	}

	// Parse disk usage (if available)
	if len(parts) > 3 {
		netRegex := regexp.MustCompile(`([\d.]+) ([^/]+)/([\d.]+) ([^/]+)`)
		netMatch := netRegex.FindStringSubmatch(parts[3])
		if len(netMatch) >= 5 {
			metric.NetworkRX, _ = strconv.ParseFloat(netMatch[1], 64)
			metric.NetworkTX, _ = strconv.ParseFloat(netMatch[3], 64)
		}
	}

	// Set status based on CPU usage
	if metric.CPUPercent > 80 {
		metric.Status = "HIGH"
	} else if metric.CPUPercent > 50 {
		metric.Status = "NORMAL"
	} else {
		metric.Status = "LOW"
	}

	return metric
}

func generateReport(metrics []ResourceMetrics, format string) {
	if format == "json" {
		generateJSONReport(metrics)
	} else if format == "html" {
		generateHTMLReport(metrics)
	} else {
		generateTextReport(metrics)
	}
}

func generateTextReport(metrics []ResourceMetrics) {
	fmt.Println(color.CyanString("\n=== RESOURCE USAGE REPORT ===\n"))

	fmt.Printf("%-20s %8s %12s %12s %8s\n",
		"SERVICE", "CPU%", "MEM USED", "MEM TOTAL", "STATUS")
	fmt.Println(strings.Repeat("-", 70))

	for _, m := range metrics {
		statusColor := color.GreenString
		if m.Status == "HIGH" {
			statusColor = color.RedString
		} else if m.Status == "NORMAL" {
			statusColor = color.YellowString
		}

		fmt.Printf("%-20s %8.1f%% %10.1fMB %10.1fMB %8s\n",
			m.Service,
			m.CPUPercent,
			m.MemoryUsed,
			m.MemoryTotal,
			statusColor(m.Status),
		)
	}

	fmt.Println()
	fmt.Printf("Report generated at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func generateJSONReport(metrics []ResourceMetrics) {
	jsonData, _ := json.MarshalIndent(metrics, "", "  ")
	fmt.Println(string(jsonData))
}

func generateHTMLReport(metrics []ResourceMetrics) {
	fmt.Println(`<!DOCTYPE html>
<html>
<head>
    <title>Resource Usage Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #4CAF50; color: white; }
        tr:nth-child(even) { background-color: #f2f2f2; }
        .high { color: red; }
        .normal { color: orange; }
        .low { color: green; }
    </style>
</head>
<body>
    <h1>Resource Usage Report</h1>
    <p>Generated: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
    <table>
        <tr>
            <th>Service</th>
            <th>CPU%</th>
            <th>MEM Used</th>
            <th>MEM Total</th>
            <th>Status</th>
        </tr>`)

	for _, m := range metrics {
		fmt.Printf(`        <tr>
            <td>%s</td>
            <td>%.1f%%</td>
            <td>%.1fMB</td>
            <td>%.1fMB</td>
            <td class="%s">%s</td>
        </tr>`, m.Service, m.CPUPercent, m.MemoryUsed, m.MemoryTotal, strings.ToLower(m.Status), m.Status)
	}

	fmt.Println(`    </table>
</body>
</html>`)
}