package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const timeout = 500 * time.Millisecond

func ScanPort(host string, port int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	conn.Close()
	results <- port
}

func ScanPorts(host string, startPort, endPort int) []int {
	var wg sync.WaitGroup
	results := make(chan int)
	var openPorts []int

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go ScanPort(host, port, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		openPorts = append(openPorts, port)
	}

	return openPorts
}