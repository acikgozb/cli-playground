/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"io"
	"os"

	"github.com/acikgozb/cli-playground/pscan/scan"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add <host1> ... <hostn>", // let consumer know that he/she can provide any number of hosts.
	Aliases: []string{"a"},
	Short:   "Add new host(s) to the list",
	Args: cobra.MinimumNArgs(
		1,
	), // a small validation which Cobra can provide, in this case we want at least one host to be added with this command.
	SilenceUsage: true, // prevent showing command usage when an error occurs to not cause confusion. The user can still see the usage with -h
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return addAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(addCmd)
}

// addAction runs when users use add command to add a host to the host list.
// PS: The current design is not suitable for concurrent usage, beware of it.
func addAction(out io.Writer, hostsFile string, args []string) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, host := range args {
		if err := hl.Add(host); err != nil {
			return err
		}

		fmt.Fprintln(out, "Added host:", host)
	}

	return hl.Save(hostsFile)
}
