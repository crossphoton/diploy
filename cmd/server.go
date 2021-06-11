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
	"time"

	"github.com/crossphoton/diploy/src"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start diploy server",
	Run:   server,
}

var serverSetupCmd = &cobra.Command{
	Use:   "server",
	Short: "Start diploy server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server setup called")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverSetupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	server_address := cmd.Flag("addr").Value
	if server_address.String() == ":80" {
		server_address.Set("localhost:80")
	}

	fmt.Printf("Initializing server at http://%s\n", server_address)

	server := http.Server{
		Addr:         server_address.String(),
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
	mode := mux.Vars(r)["mode"]
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
