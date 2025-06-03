package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	types "github.com/megh16123/loadbs/Types"
)

var State types.BS_State

func ReadAndLoadConfiguration(St *types.BS_State) {
	conf_bytes, err := os.ReadFile("configuration.json")
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		St.InitState(conf_bytes)
	}
}

func run_user_server(port string, serverName string, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	for _, service := range State.Services {
		for _, ep := range service.ServiceEndpoints {
			mux.HandleFunc("/"+ep, func(w http.ResponseWriter, r *http.Request) {
				res, err := State.GetAvailabe(service.ServiceName, ep)
				if err != nil {
					w.WriteHeader(404)
					fmt.Fprintf(w, "Not a valid endpoint")
				} else {

					// Construct the new request
					targetURL := res.ToString() + r.URL.Path // Replace with your target server address

					proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
					if err != nil {
						http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
						return
					}
					proxyReq.URL.RawQuery = r.URL.Query().Encode()
					// Copy all headers from the original request to the new request
					for name, values := range r.Header {
						for _, value := range values {
							proxyReq.Header.Add(name, value)
						}
					}

					// Optionally, add/modify headers like X-Forwarded-For
					proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

					// Execute the new request
					client := &http.Client{}
					resp, err := client.Do(proxyReq)
					if err != nil {
						http.Error(w, "Error forwarding request", http.StatusBadGateway)
						return
					}
					defer resp.Body.Close()

					// Copy headers from the proxy response to the original response writer
					for name, values := range resp.Header {
						for _, value := range values {
							w.Header().Add(name, value)
						}
					}

					// Set the status code
					w.WriteHeader(resp.StatusCode)

					// Copy the response body
					_, err = io.Copy(w, resp.Body)
					if err != nil {
						log.Printf("Error copying response body: %v", err)
					}

				}
			})
		}
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Starting %s on port %s\n", serverName, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting %s: %v\n", serverName, err)
	}
}

func run_app_server(port string, serverName string, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	var app_info types.IP_PORT
	// Give 204 to leader and 200 with ip and port of leader
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(503)
		} else {
			defer r.Body.Close()
			err := json.Unmarshal(bodyBytes, &app_info)
			if err != nil {
				w.WriteHeader(503)
				w.Write([]byte("Something went wrong"))
			} else {
				isLeader := State.CheckLeader(app_info)
				if isLeader {
					w.WriteHeader(204)
					w.Write([]byte("done"))
				} else {
					w.WriteHeader(200)
					w.Write(State.GetLeader())
				}
			}
		}
	})
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Starting %s on port %s\n", serverName, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error starting %s: %v\n", serverName, err)
	}
}
func main() {
	var wg sync.WaitGroup
	ReadAndLoadConfiguration(&State)
	wg.Add(2)
	go run_app_server("8080", "APP SERVER", &wg)
	go run_user_server("8090", "USER SERVER", &wg)
	wg.Wait()
}
