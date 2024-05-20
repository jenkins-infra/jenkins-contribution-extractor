/*
Copyright © 2024 Jean-Marc Meessen jean-marc@meessen-web.org

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var honoredDataDir string
var honoredOutput string
var honoredMonth string

// honoredCmd represents the honored command
var honoredCmd = &cobra.Command{
	Use:   "honored",
	Short: "Gets a contributor to honor",
	Long: `A command to get a random submitter from a given month and
format his data in such a way that it can be used to format an honoring
message at the bottom of the https://contributors.jenkins.io/ page`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performHonoredContributorSelection(honoredDataDir, honoredOutput, honoredMonth)
	},
}

// Initialize command parameters and defaults
func init() {
	rootCmd.AddCommand(honoredCmd)
	honoredCmd.Flags().StringVarP(&honoredDataDir, "data_dir", "", "consolidated_data", "Directory containing the data to be read")
	honoredCmd.Flags().StringVarP(&honoredOutput, "output", "", "", "File to output the data to (default: \"[data_dir]/honored_contributor.csv\")")
	honoredCmd.Flags().StringVarP(&honoredMonth, "month", "", "", "the month to select the submitter from (format \"YYYY-MM\")")
}

// Command processing entry point
func performHonoredContributorSelection(dataDir string, outputFileName string, monthToSelectFrom string) error {
	//does the dataDir exist ?
	if !isValidDir(dataDir) {
		return fmt.Errorf("Supplied DataDir \"%s\" does not exist.", dataDir)
	}
	//TODO: if output is not defined, build it
	//TODO: if month is not defined, try to guess it. If defined does it have the right format?
	return nil
}
