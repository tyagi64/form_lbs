package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "USER SERVER")
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
