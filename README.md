# Port Scanner

A concurrent port scanner with a REST API, written in Go.

## Prerequisites

- Go 1.21 or later
- Docker (for containerized deployment)

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
3. Create a `.env` file in the root directory:
    ```
    PORT=8080
    ```
4. Run the application:
    ```
    ./port-scanner
    ```

The server will start on the port specified in your `.env` file (default is 8080).

### Usage

Use curl or any HTTP client to make a GET request to the `/scan` endpoint:
```
curl "http://localhost:8080/scan?host=example.com&start=1&end=100"
```

## Docker Deployment

### Building the Docker Image

1. Build the image with default settings:
    ``` 
    docker build --build-arg PORT=9000 -t port-scanner .
    ```
### Running the Docker Container

1. Run with default settings:
    ```
    docker run -p 8080:8080 port-scanner
    ```
    Or specify a custom port:
    ```
    docker run -p 9000:9000 -e PORT=9000 port-scanner
    ```
### Usage with Docker

The usage is the same as local deployment, just ensure you're using the correct port:
```
curl "http://localhost:<port-num>/scan?host=example.com&start=1&end=100"
```

## API Endpoints

### GET /scan

Scans the specified host for open ports.

Query Parameters:
- `host`: The target host to scan (required)
- `start`: The starting port number (optional, default: 1)
- `end`: The ending port number (optional, default: 1024)

Response:
```json
{
  "host": "example.com",
  "start_port": 1,
  "end_port": 100,
  "open_ports": [80, 443]
}