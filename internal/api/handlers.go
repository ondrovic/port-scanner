package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"port-scanner/internal/models"
	"port-scanner/internal/scanner"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/scan", handleScan)
	return r
}

func handleScan(c *gin.Context) {
	host := c.Query("host")
	startPort, endPort := getPorts(c)

	if host == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Host parameter is required"})
		return
	}

	openPorts := scanner.ScanPorts(host, startPort, endPort)

	result := models.ScanResult{
		Host:      host,
		StartPort: startPort,
		EndPort:   endPort,
		OpenPorts: openPorts,
	}

	c.IndentedJSON(http.StatusOK, result)
}

func getPorts(c *gin.Context) (int, int) {
	startPort, _ := strconv.Atoi(c.DefaultQuery("start", "1"))
	endPort, _ := strconv.Atoi(c.DefaultQuery("end", "1024"))
	return startPort, endPort
}