#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Starting setup for Golang Load Balancer..."

# Check if Go is installed
if ! [ -x "$(command -v go)" ]; then
  echo "Go is not installed. Installing Go..."
  wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
  rm go1.20.5.linux-amd64.tar.gz
  echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
  source ~/.bashrc
else
  echo "Go is already installed."
fi

# Check if Redis is installed
if ! [ -x "$(command -v redis-server)" ]; then
  echo "Redis is not installed. Installing Redis..."
  sudo apt update
  sudo apt install -y redis-server
  sudo systemctl enable redis-server.service
  sudo systemctl start redis-server.service
else
  echo "Redis is already installed."
fi

# Initialize Go module if not already initialized
if [ ! -f "go.mod" ]; then
  echo "Initializing Go module..."
  go mod init loadbalancer
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go get github.com/go-redis/redis/v8
go get gopkg.in/yaml.v3

echo "Setup completed successfully."
echo "Next steps:"
echo "1. Ensure Redis is running: sudo systemctl start redis-server"
echo "2. Run the application: go run main.go"
echo "3. Access the load balancer: http://localhost:8080"