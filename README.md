[![CodeFactor](https://www.codefactor.io/repository/github/ondrovic/port-scanner/badge)](https://www.codefactor.io/repository/github/ondrovic/port-scanner)
# Port Scanner

A concurrent port scanner

## Prerequisites

- Go 1.22.5 or later

## Local Development

### Setup

1. Clone the repository:
    ```
    git clone https://github.com/yourusername/port-scanner.git
    cd port-scanner
    ```
2. Install dependencies:
    ```
    go mod tidy
    ```
3. Run the application:
    ```
    go run .\cmd\scanner\port-scanner.go   
    ```
4. Build the application
    ```
    go build .\cmd\scanner\port-scanner.go
    ```
### Response
```json
{
  "host": "localhost",
  "ip": "127.0.0.1",
  "start_port": 1,
  "end_port": 9999,
  "num_of_ports": 18,
  "open_ports": [
    80
  ]
}
```

# Security Note
### Ensure that you have permission to scan the target hosts. Unauthorized port scanning may be illegal in some jurisdictions.
