package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"golang.org/x/net/context"
)

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name    string   `yaml:"name"`
	Servers []string `yaml:"servers"`
}

type LoadBalanceRequest struct {
	ServiceName string `json:"serviceName"`
}

var (
	ctx         = context.Background()
	config      Config
	redisClient *redis.Client
	mu          sync.Mutex
)

func init() {
	// Load YAML configuration
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse YAML config: %v", err)
	}

	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust if Redis is not running locally
	})
}

func getNextServer(serviceName string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	key := fmt.Sprintf("service:%s:next", serviceName)
	nextIndex, err := redisClient.Get(ctx, key).Int()
	if err != nil {
		// If no value exists in Redis, start from the first server
		nextIndex = 0
	}
	service := getServiceByName(serviceName)
	if service == nil {
		return "", fmt.Errorf("service not found: %s", serviceName)
	}
	// Get the next server
	server := service.Servers[nextIndex%len(service.Servers)]
	// Update Redis to point to the next server
	redisClient.Set(ctx, key, (nextIndex+1)%len(service.Servers), 0)
	return server, nil
}

func getServiceByName(name string) *Service {
	for _, service := range config.Services {
		if service.Name == name {
			return &service
		}
	}
	return nil
}

func loadBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Only POST is allowed.", http.StatusMethodNotAllowed)
		return
	}

	var req LoadBalanceRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.ServiceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	server, err := getNextServer(req.ServiceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, server, http.StatusTemporaryRedirect)
}

func listServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.Services)
}

func main() {
	http.HandleFunc("/balance", loadBalanceHandler)
	http.HandleFunc("/services", listServicesHandler)

	fmt.Println("Load balancer running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}