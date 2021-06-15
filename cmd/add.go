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
	"io/ioutil"
	"os"

	"github.com/crossphoton/diploy/src"
	"github.com/spf13/cobra"
)

var cfgFile string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a configuration",
	RunE:  addConfiguration,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringVar(&cfgFile, "file", "diploy.yml", "config file (default is ./diploy.yml)")
}

func addConfiguration(cmd *cobra.Command, args []string) (err error) {
	fileFlag := cmd.Flag("file").Value
	file, err := ioutil.ReadFile(fileFlag.String())
	if err != nil {
		return
	}
	workdir, err := os.Getwd()
	if err != nil {
		return
	}
	consent("diploy logs location", &src.LOG_PATH)
	err = src.AddFromFile(file, workdir)
	if err != nil {
		return
	}

	return err
}
