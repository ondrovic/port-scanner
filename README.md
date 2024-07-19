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

### Pulling from Docker hub
#### Run
```
docker run --rm -p 5002:5002 -e PORT=5002 ondrovic/port-scanner:latest
```
#### Test
```
curl "http://localhost:5002/scan?host=127.0.0.1&start=1&end=80"
```
Response
```json
{
    "host": "127.0.0.1",
    "start_port": 1,
    "end_port": 80,
    "open_ports": null
}
```

### Usage with Docker Compose
1. Create a new file named `docker-compose.yml` in your desired directory.

2. Copy the following content into the `docker-compose.yml` file:

   ```yaml
   services:
     port-scanner:
       image: ondrovic/port-scanner:latest
       ports:
         - "9000:8080"
       environment:
         - PORT=8080
    ```
3. Save the file
4. Open a terminal and navigate to the directory containing your `docker-compose.yml` file.
5. Run the following command to start the service:
    ```
    docker compose up -d
    ```
    This command pulls the image from Docker Hub (if not already present locally) and starts the container in detached mode.
6. Verify the service is up and running"
    ```
    docker compose ps
    ```
    You should see the port-scanner service listed as `Up`.

### Usage
Once deployed, you can access the port scanner API at:
```
http://localhost:9000/scan?host=example.com&start=1&end=100
```
Replace localhost with your server's IP address if accessing remotely.

### Stopping the Service
To stop the service, run:
```
docker compose down
```
This command stops and removes the containers defined in the docker-compose.yml file.

### Updating the image
1. Pull the image:
    ```
    docker compose pull
    ```
2. Restart the service with the new image:
    ```
    docker compose up -d
    ```
### Troubleshooting
* If you encounter port conflicts, modify the host port in the docker-compose.yml file.
* Check logs using docker-compose logs port-scanner if the service isn't working as expected.
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
```

# Security Note
### Ensure that you have permission to scan the target hosts. Unauthorized port scanning may be illegal in some jurisdictions.
