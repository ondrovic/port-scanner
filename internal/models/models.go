package models

type ScanResult struct {
	Host      string `json:"host"`
	StartPort int    `json:"start_port"`
	EndPort   int    `json:"end_port"`
	OpenPorts []int  `json:"open_ports"`
}