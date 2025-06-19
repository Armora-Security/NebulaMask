package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
    "time"

    "gopkg.in/yaml.v2"
)

type Backend struct {
    Name              string `yaml:"name"`
    URL               string `yaml:"url"`
    HealthCheckPath   string `yaml:"health_check_path"`
    IntervalSec       int    `yaml:"interval_sec"`
    Healthy           bool
}

type Config struct {
    Server struct {
        HTTPPort  int `yaml:"http_port"`
        HTTPSPort int `yaml:"https_port"`
    } `yaml:"server"`

    SSL struct {
        Enabled     bool   `yaml:"enabled"`
        CertFile    string `yaml:"cert_file"`
        KeyFile     string `yaml:"key_file"`
    } `yaml:"ssl"`

    Backends []Backend `yaml:"backends"`
}

var backends []Backend
var mutex = &sync.Mutex{}
var currentIdx int

func loadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    return &config, err
}

func startHealthChecks() {
    for {
        for i := range backends {
            resp, err := http.Get(backends[i].URL + backends[i].HealthCheckPath)
            mutex.Lock()
            if err == nil && resp.StatusCode == 200 {
                backends[i].Healthy = true
            } else {
                backends[i].Healthy = false
            }
            mutex.Unlock()
            time.Sleep(time.Second * time.Duration(backends[i].IntervalSec))
        }
    }
}

func getNextHealthyBackend() (*Backend, error) {
    mutex.Lock()
    defer mutex.Unlock()

    for i := 0; i < len(backends); i++ {
        idx := (currentIdx + i) % len(backends)
        if backends[idx].Healthy {
            currentIdx = (idx + 1) % len(backends)
            return &backends[idx], nil
        }
    }
    return nil, fmt.Errorf("no healthy backend available")
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
    backend, err := getNextHealthyBackend()
    if err != nil {
        http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        return
    }

    proxyReq, _ := http.NewRequest(r.Method, backend.URL+r.RequestURI, r.Body)
    proxyReq.Header = r.Header

    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Upstream Error", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    for k, v := range resp.Header {
        w.Header()[k] = v
    }
    w.WriteHeader(resp.StatusCode)
    body, _ := ioutil.ReadAll(resp.Body)
    w.Write(body)
}

func redirectHTTP(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "https://"+r.Host+r.RequestURI,  http.StatusMovedPermanently)
}

func main() {
    config, err := loadConfig("config.yaml")
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    backends = config.Backends
    go startHealthCheck()

    // Setup handlers
    mux := http.NewServeMux()
    mux.HandleFunc("/", proxyHandler)

    // Start HTTPS server
    go func() {
        port := fmt.Sprintf(":%d", config.Server.HTTPSPort)
        fmt.Printf("🔐 HTTPS Server running on port %s\n", port)
        log.Fatal(http.ListenAndServeTLS(port, config.SSL.CertFile, config.SSL.KeyFile, mux))
    }()

    // Start HTTP redirect server
    go func() {
        port := fmt.Sprintf(":%d", config.Server.HTTPPort)
        fmt.Printf("🌐 HTTP Redirect Server running on port %s\n", port)
        http.ListenAndServe(port, http.HandlerFunc(redirectHTTP))
    }()

    // Wait forever
    select {}
}
