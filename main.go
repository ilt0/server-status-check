package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
)

type StatusData struct {
	SERVER []ServerStatus
}

type ServerStatus struct {
	Title   string
	Content string
}

func getServerStatus() []ServerStatus {
	cmd := exec.Command("bash", "-c", "uptime")
	uptimeOutput, _ := cmd.CombinedOutput()

	cmd = exec.Command("bash", "-c", "cat /proc/loadavg")
	loadAvgOutput, _ := cmd.CombinedOutput()

	cmd = exec.Command("bash", "-c", "df -h /")
	diskUsageOutput, _ := cmd.CombinedOutput()

	cmd = exec.Command("bash", "-c", "free -h")
	memoryUsageOutput, _ := cmd.CombinedOutput()

	cmd = exec.Command("bash", "-c", "ifconfig eth0")
	networkStatusOutput, _ := cmd.CombinedOutput()

	cmd = exec.Command("cmd", "/c", "systeminfo")
	systemStatusOutput, _ := cmd.CombinedOutput()

	statusData := []ServerStatus{
		{Title: "System Uptime", Content: string(uptimeOutput)},
		{Title: "Load Average", Content: string(loadAvgOutput)},
		{Title: "Disk Usage for /", Content: string(diskUsageOutput)},
		{Title: "Memory Usage", Content: string(memoryUsageOutput)},
		{Title: "Network Status", Content: string(networkStatusOutput)},
		{Title: "System Info", Content: string(systemStatusOutput)},
	}

	return statusData
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	statusData := getServerStatus()

	tmpl, err := template.ParseFiles("server-status-check.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, StatusData{SERVER: statusData})
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

func main() {

	http.HandleFunc("/Server-status-check", serverStatusHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})


	fmt.Println("Starting HTTP server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP server error:", err)
	}
	

	
}
