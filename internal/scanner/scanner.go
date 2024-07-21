package scanner

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"port-scanner/internal/models"
	"port-scanner/internal/utils"

	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
)

var (
	timeout    time.Duration
	maxWorkers int
	retries    int
	envFile    string
)

func init() {
	utils.ClearConsole()
	// Get the directory of the current file
	
	projRoot, err := utils.GetProjectRoot()
	if err != nil {
		pterm.Warning.PrintOnError(err)
	}

	envFile = fmt.Sprintf("%s\\.env", projRoot)

	if err := godotenv.Load(envFile); err != nil {
		pterm.Warning.Printfln("No .env file found")
	}

	timeOutMs, err := getEnvInt("PORT_TIMEOUT", 1000)
    if err != nil {
        pterm.Warning.Printfln("Invalid PORT_TIMEOUT value: %v, using default: %v", err, timeOutMs)
    }
    timeout = time.Duration(timeOutMs) * time.Millisecond

    maxWorkers, err = getEnvInt("SCAN_WORKERS", 5000)
    if err != nil {
        pterm.Warning.Printfln("Invalid SCAN_WORKERS value: %v, using default: %v", err, maxWorkers)
    }

    retries, err = getEnvInt("SCAN_RETRIES", 3)
    if err != nil {
        pterm.Warning.Printfln("Invalid SCAN_RETRIES value: %v, using default: %v", err, retries)
    }
}

func getEnv(key, fallback string) string {
    value := os.Getenv(key)
    if value == "" {
        pterm.Warning.Printfln("Environment variable %s not set, using default value: %s", key, fallback)
        return fallback
    }
    return value
}

func getEnvInt(key string, fallback int) (int, error) {
    strValue := getEnv(key, strconv.Itoa(fallback))
    intValue, err := strconv.Atoi(strValue)
    if err != nil {
        return fallback, fmt.Errorf("invalid value for %s: %v", key, err)
    }
    return intValue, nil
}

func scanPort(ctx context.Context, ip string, port int) (bool, string) {
	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return false, ""
		default:
			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err == nil {
				defer conn.Close()

				// Try to upgrade to TLS
				tlsConn := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
				err = tlsConn.SetDeadline(time.Now().Add(timeout))
				if err == nil {
					err = tlsConn.Handshake()
					if err == nil {
						return true, "TLS"
					}
				}

				// If not TLS, try to get service banner
				conn.SetDeadline(time.Now().Add(timeout))
				banner := make([]byte, 1024)
				n, _ := conn.Read(banner)
				if n > 0 {
					return true, string(banner[:n])
				}

				return true, "Open"
			}
			time.Sleep(timeout / 2)
		}
	}
	return false, ""
}

func worker(ctx context.Context, ip string, ports <-chan int, results chan<- models.PortResult, wg *sync.WaitGroup, p *pterm.ProgressbarPrinter) {
	defer wg.Done()
	for port := range ports {
		open, info := scanPort(ctx, ip, port)
		if open {
			results <- models.PortResult{Port: port, Info: info}
		}
		p.Increment()
	}
}

func ScanPorts(target string, startPort, endPort int) models.ScanResult {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ip, hostname, err := resolveTarget(target)
	if err != nil {
		return models.NewScanResult(target, "", startPort, endPort, nil, fmt.Sprintf("Resolution failed: %v", err))
	}

	totalPorts := endPort - startPort + 1
	
	// scanning text
	pterm.Println(pterm.Sprintf("Scanning: %s (%s)", hostname, ip))
    // Create and start the progress bar on the next line
    p, _ := pterm.DefaultProgressbar.
        WithTotal(totalPorts).
        WithTitle("Scanning Port").
        Start()
	defer p.Stop()

	var wg sync.WaitGroup
	ports := make(chan int, maxWorkers)
	results := make(chan models.PortResult, totalPorts)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ctx, ip.String(), ports, results, &wg, p)
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			select {
			case ports <- port:
			case <-ctx.Done():
				return
			}
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var openPorts []models.PortResult
	for result := range results {
		openPorts = append(openPorts, models.PortResult{Port: result.Port, Info: result.Info})
	}

	sort.Slice(openPorts, func(i, j int) bool {
		return openPorts[i].Port < openPorts[j].Port
	})

	return models.NewScanResult(hostname, ip.String(), startPort, endPort, openPorts, "")
}

func resolveTarget(target string) (net.IP, string, error) {
	ip := net.ParseIP(target)
	if ip == nil {
		ips, err := net.LookupIP(target)
		if err != nil {
			return nil, "", err
		}
		ip = ips[0]
		return ip, target, nil
	}

	names, err := net.LookupAddr(ip.String())
	if err == nil && len(names) > 0 {
		return ip, names[0], nil
	}
	return ip, ip.String(), nil
}
