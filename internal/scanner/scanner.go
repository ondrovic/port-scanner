package scanner

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
	EnvFile    string
)

func init() {
    utils.ClearConsole()

	if os.Getenv("PORT_TIMEOUT") == "" && os.Getenv("SCAN_WORKERS") == "" && os.Getenv("SCAN_RETRIES") == "" {
		projRoot, err := utils.GetProjectRoot()
        if err != nil {
            pterm.Warning.PrintOnError(err)
        }

        envLocations := []string{
            filepath.Join(projRoot, ".env"),
            "/app/.env",
        }

        envLoaded := false
        for _, loc := range envLocations {
            if err := godotenv.Load(loc); err == nil {
                EnvFile = loc
                envLoaded = true
                pterm.Success.Printfln("Loaded .env file from: %s", loc)
                break
            }
        }

        if !envLoaded {
            pterm.Info.Printfln("No .env file found, using default or provided environment variables")
        }
    }

    // Load environment variables or use defaults
    timeOutMs, err := getEnvInt("PORT_TIMEOUT", 1000)
    if err != nil {
        pterm.Info.Printfln("Using default PORT_TIMEOUT: %v", timeOutMs)
    }
    timeout = time.Duration(timeOutMs) * time.Millisecond

    maxWorkers, err = getEnvInt("SCAN_WORKERS", 5000)
    if err != nil {
        pterm.Info.Printfln("Using default SCAN_WORKERS: %v", maxWorkers)
    }

    retries, err = getEnvInt("SCAN_RETRIES", 3)
    if err != nil {
        pterm.Info.Printfln("Using default SCAN_RETRIES: %v", retries)
    }
}

func getEnvInt(key string, fallback int) (int, error) {
    if value, exists := os.LookupEnv(key); exists {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return fallback, fmt.Errorf("invalid %s value: %v", key, err)
        }
        return intValue, nil
    }
    return fallback, nil
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
