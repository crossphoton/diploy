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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a service with name",
	RunE:  start,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVar(&server_address, "addr", "0.0.0.0:80", "specify address for server [ip:port]")
}

func start(cmd *cobra.Command, args []string) error {
	var allFailed = true
	if len(args) < 2 {
		return fmt.Errorf("Usage: start [method] [config_name]")
	}

	mode := args[0]
	args = args[1:]
	for _, name := range args {
		err := startUtil(name, mode)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
			continue
		}
		allFailed = false
	}
	if allFailed {
		os.Exit(1)
	}

	return nil
}

func startUtil(name, mode string) (err error) {
	url := fmt.Sprintf("http://%s/start/%s/%s", server_address, mode, name)

	response, err := http.Post(url, "", bytes.NewBuffer([]byte("")))
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		message, _ := ioutil.ReadAll(response.Body)
		err = fmt.Errorf("%s", string(message))
		return
	}

	return
}
