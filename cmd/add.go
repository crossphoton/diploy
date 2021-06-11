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
	"io/ioutil"
	"log"
	"os"

	"github.com/crossphoton/diploy/src"
	"github.com/spf13/cobra"
)

var cfgFile string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a configuration",
	Run:   addConfiguration,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringVar(&cfgFile, "file", "diploy.yml", "config file (default is ./diploy.yml)")
}

func addConfiguration(cmd *cobra.Command, args []string) {
	fileFlag := cmd.Flag("file").Value
	file, err := ioutil.ReadFile(fileFlag.String())
	if err != nil {
		fmt.Println("couldn't load file - ", fileFlag.String(), " : ", err)
		os.Exit(1)
	}
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal("couldn't get current working directory: ", err)
		os.Exit(137)
	}
	err = src.AddFromFile(file, workdir)
	if err != nil {
		fmt.Println("can't use the supplied file: ", err)
		os.Exit(1)
	}
}
