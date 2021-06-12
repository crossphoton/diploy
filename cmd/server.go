/*
Copyright © 2021 Aditya Agrawal adiag1200@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/crossphoton/diploy/src"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var serviceFile = `[Unit]
Description=diploy server
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
User=root
ExecStart=/usr/bin/diploy server --addr %s >> %s/diploy-log.txt

[Install]
WantedBy=multi-user.target
`

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start diploy server",
	Run:   server,
}

var serverSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup diploy server as a service",
	RunE:  serveSetup,
}

func serveSetup(cmd *cobra.Command, args []string) error {
	fmt.Println(strings.ToUpper("diploy server setup"))

	// Create directories
	fmt.Println("using the below path for logs and database (build manually to change):")
	fmt.Println("Logs:", src.LOG_PATH)
	fmt.Println("Config Database:", src.DB_PATH)

	err := os.MkdirAll(src.DB_PATH, 0700)
	err = os.MkdirAll(src.LOG_PATH, 0700)

	if err != nil {
		return fmt.Errorf("couldn't create directories: %v", err)
	}

	executablePath, err := exec.LookPath("diploy")
	if err != nil {
		dir, _ := os.Getwd()
		executablePath = dir + "/diploy"
	}

	binaryPath, _ := exec.LookPath("echo")
	tempF := strings.Split(binaryPath, "/")
	tempF = tempF[:len(tempF)-1]
	binaryPath = strings.Join(tempF, "/") + "/diploy"

	fmt.Printf("bin path [to store binary] (%s): ", binaryPath)

	temp := ""
	if count, _ := fmt.Scanf("%s", &temp); count > 0 {
		binaryPath = temp
	}

	err = exec.Command("cp", executablePath, binaryPath).Run()
	if err != nil {
		return fmt.Errorf("couldn't copy binary to %s: %s", binaryPath, err)
	}

	fmt.Print("setup systemd service file (in /etc/systemd/system)? (y/N)")
	if count, _ := fmt.Scanf("%s", &temp); count > 0 {
		if strings.ToLower(temp) == ("y") {
			serverAddress := ""
			fmt.Print("server address (0.0.0.0:80): ")
			if count, _ = fmt.Scanf("%s", &serverAddress); count == 0 {
				serverAddress = "0.0.0.0:80"
			}
			servicefile := []byte(fmt.Sprintf(serviceFile, serverAddress))
			file, err := os.Create("/etc/systemd/system/diploy.service")
			if err != nil {
				return fmt.Errorf("couldn't create file /etc/systemd/system/diploy.service: %s", err)
			}
			_, err = file.Write([]byte(servicefile))
			if err != nil {
				return fmt.Errorf("couldn't write to file: %s", err)
			}
			file.Close()
		}
	}

	return nil
}

var server_address string

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&server_address, "addr", "0.0.0.0:80", "specify address for server [ip:port]")
	serverCmd.AddCommand(serverSetupCmd)
}

func httpHandler() (handler *mux.Router) {
	handler = mux.NewRouter()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is diploy home")
	})
	handler.HandleFunc("/start/{mode}/{name}", startWithName).Methods("POST")
	handler.HandleFunc("/stop/{name}", stopWithName).Methods("POST")
	handler.HandleFunc("/restart/{name}", restartWithName).Methods("POST")
	return
}

func server(cmd *cobra.Command, args []string) {

	fmt.Printf("Initializing server at http://%s\n", server_address)

	server := http.Server{
		Addr:         server_address,
		Handler:      httpHandler(),
		WriteTimeout: time.Second * 3,
	}

	log.Fatal(server.ListenAndServe())
}

func stopWithName(w http.ResponseWriter, r *http.Request) {
	config, err := repetitive(w, r)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "{\"message\": \"started stopping processes\", \"error\": null}")
	go config.Stop()
}

func startWithName(w http.ResponseWriter, r *http.Request) {
	var err error
	mode := mux.Vars(r)["mode"]
	if !src.MODES[mode] {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"failed\", \"error\": \"%s\"}",
			fmt.Sprintf("mode '%s' not supported", mode))))
		return
	}
	config, err := repetitive(w, r)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "{\"message\": \"started in background\", \"error\": null}")

	go config.Start(mode)
}

func restartWithName(w http.ResponseWriter, r *http.Request) {
	config, err := repetitive(w, r)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "{\"message\": \"started in background\", \"error\": null}")
	go config.Restart()
}

func repetitive(w http.ResponseWriter, r *http.Request) (config src.Config, err error) {
	name := mux.Vars(r)["name"]
	config, err = src.SearchConfig(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"failed\", \"error\": \"%s\"}", err)))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	return
}
