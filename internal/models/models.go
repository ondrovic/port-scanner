package models

// CommandFlags structure of the command flags
type CommandFlags struct {
	Host      string
	StartPort int
	EndPort   int
}

// PortResult the structure of the port
type PortResult struct {
	Port int
	Info string
}

// ScanResult the structure of teh scan results
type ScanResult struct {
	Hostname   string
	IP         string
	StartPort  int
	EndPort    int
	NumOfPorts int
	OpenPorts  []PortResult
	Error      string
}

// NewScanResult creates and returns a new ScanResult
func NewScanResult(hostname, ip string, startPort, endPort int, openPorts []PortResult, err string) ScanResult {
	return ScanResult{
		Hostname:   hostname,
		IP:         ip,
		StartPort:  startPort,
		EndPort:    endPort,
		NumOfPorts: len(openPorts),
		OpenPorts:  openPorts,
		Error:      err,
	}
}
