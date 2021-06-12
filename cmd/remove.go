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
	"os"

	"github.com/crossphoton/diploy/src"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a configuaration by name",
	RunE:  delete,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func delete(cmd *cobra.Command, args []string) error {
	var allFailed = true
	if len(args) < 2 {
		return fmt.Errorf("Usage: remove config_name...")
	}
	for _, name := range args {
		config, err := src.SearchConfig(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "no such config: ", name)
			continue
		}
		config.Delete()
		allFailed = false
	}

	if allFailed {
		os.Exit(1)
	}

	return nil
}
