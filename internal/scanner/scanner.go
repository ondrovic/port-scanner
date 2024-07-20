package scanner

import (
	"fmt"
	"log"
	"net"

	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"port-scanner/internal/models"

	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
)

var (
	timeout    time.Duration
	maxWorkers int
	retries    int
	envMap     map[string]string
)

func init() {
	var err error
	// Parse the embedded .env file
	envMap, err = godotenv.Parse(strings.NewReader(envFile))
	if err != nil {
		log.Println("Error parsing embedded .env file:", err)
	}

	timeoutMs, err := strconv.Atoi(getEnv("PORT_TIMEOUT", "500"))
	if err != nil {
		log.Printf("Invalid PORT_TIMEOUT value: %v, using default", err)
		timeoutMs = 1000
	}
	timeout = time.Duration(timeoutMs) * time.Millisecond

	maxWorkers, err = strconv.Atoi(getEnv("SCAN_WORKERS", "5000"))
	if err != nil {
		log.Printf("Invalid SCAN_WORKERS value: %v, using default", err)
		maxWorkers = 100
	}

	retries, err = strconv.Atoi(getEnv("SCAN_RETRIES", "2"))
	if err != nil {
		log.Printf("Invalid SCAN_RETRIES value: %v, using default", err)
		retries = 3
	}
}

func getEnv(key, fallback string) string {
	if value, ok := envMap[key]; ok {
		return value
	}
	return fallback
}

func scanPort(ip string, port int) bool {
	for i := 0; i < retries; i++ {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(timeout / 2) // Wait a bit before retrying
	}
	return false
}

func worker(ip string, ports <-chan int, results chan<- int, wg *sync.WaitGroup, p *pterm.ProgressbarPrinter) {
	defer wg.Done()
	for port := range ports {
		if scanPort(ip, port) {
			results <- port
		}
		p.Increment() // Increment the progress bar
	}
}

func ScanPorts(host string, startPort, endPort int) models.ScanResult {
	ips, err := net.LookupIP(host)
	if err != nil {
		return models.NewScanResult(host, "", startPort, endPort, nil, fmt.Sprintf("DNS resolution failed: %v", err))
	}

	ip := ips[0].String()

	totalPorts := endPort - startPort + 1

	p, _ := pterm.DefaultProgressbar.WithTotal(totalPorts).WithTitle(fmt.Sprintf("Scanning ports on: %s (%s)", host, ip)).WithMaxWidth(100).Start()

	defer p.Stop()

	var wg sync.WaitGroup
	ports := make(chan int, maxWorkers)
	results := make(chan int, totalPorts)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ip, ports, results, &wg, p)
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			ports <- port
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var openPorts []int
	for port := range results {
		openPorts = append(openPorts, port)
	}

	sort.Ints(openPorts)
	return models.NewScanResult(host, ip, startPort, endPort, openPorts, "")
}
