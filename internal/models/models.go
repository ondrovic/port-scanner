package models

// CommandFlags structure of the command flags
type CommandFlags struct {
	Host      string
	StartPort int
	EndPort   int
}

// ScanResult the structure of teh scan results
type ScanResult struct {
	Host       string `json:"host"`
	IP         string `json:"ip,omitempty"`
	StartPort  int    `json:"start_port"`
	EndPort    int    `json:"end_port"`
	NumOfPorts int    `json:"num_of_ports"`
	OpenPorts  []int  `json:"open_ports"`
	Error      string `json:"error,omitempty"`
}

// NewScanResult creates and returns a new ScanResult
func NewScanResult(host, ip string, startPort, endPort int, openPorts []int, err string) ScanResult {
	return ScanResult{
		Host:       host,
		IP:         ip,
		StartPort:  startPort,
		EndPort:    endPort,
		NumOfPorts: len(openPorts),
		OpenPorts:  openPorts,
		Error:      err,
	}
}
