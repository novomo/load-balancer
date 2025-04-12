# Golang Load Balancer

This project is a simple load balancer written in Go. It uses Redis to manage state and supports round-robin load balancing for multiple backend services.

## Features

- Round-robin load balancing for backend services.
- Configuration via a YAML file.
- REST API to balance requests and list available services.
- Nginx integration for proxying HTTP requests.

## Prerequisites

- Go 1.20 or later
- Redis server
- Nginx (optional, for HTTP proxying)

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/load-balancer.git
   cd load-balancer
   ```

## Run the setup script to install dependencies and initialize the project:

./setup.sh

## Ensure Redis is running:

sudo systemctl start redis-server

## Start the load balancer:

go run main.go

## (Optional) Configure Nginx using the provided nginx.template.conf file.

## Configuration

The load balancer is configured using a YAML file (config.yaml). Define your services and their backend servers in this file. Example:

````services:

- name: ServiceA
  servers:
  - "http://localhost:8081"
  - "http://localhost:8082"
- name: ServiceB
  servers:
  - "http://localhost:8083"
  - "http://localhost:8084"```

## API Endpoints

POST /balance
Redirects to the next backend server for the specified service.
Request body:

```{
"serviceName": "ServiceA"
}```

GET /services
Lists all configured services.
````
