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
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var honorDataDir string
var honorOutput string

// honorCmd represents the honor command
var honorCmd = &cobra.Command{
	Use:   "honor <month>",
	Short: "Gets a contributor to honor",
	Long: `A command to get a random submitter from a given month and
format his data in such a way that it can be used to format an honoring
message at the bottom of the https://contributors.jenkins.io/ page.

\"month\" is a required parameter. It is in YYYY-MM format.`,
	Args: func(cmd *cobra.Command, args []string) error {
		//call requires two parameters (org and month)
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			if err.Error() == "requires at least 1 arg(s), only received 0" {
				return fmt.Errorf("\"month\" argument is missing.")
			} else {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return performHonorContributorSelection(honorDataDir, honorOutput, args[0])
	},
}

// Initialize command parameters and defaults
func init() {
	rootCmd.AddCommand(honorCmd)
	honorCmd.Flags().StringVarP(&honorDataDir, "data_dir", "", "data", "Directory containing the data to be read")
	honorCmd.Flags().StringVarP(&honorOutput, "output", "", "", "File to output the data to (default: \"[data_dir]/honored_contributor.csv\")")
}

// Command processing entry point
func performHonorContributorSelection(dataDir string, outputFileName string, monthToSelectFrom string) error {
	// validate the month
	if !isValidMonthFormat(monthToSelectFrom) {
		return fmt.Errorf("\"%s\" is not a valid month.", monthToSelectFrom)
	}

	// does the dataDir exist ?
	if !isValidDir(dataDir) {
		return fmt.Errorf("Supplied DataDir \"%s\" does not exist.", dataDir)
	}

	// if output is not defined, build it
	if outputFileName == "" {
		outputFileName = filepath.Join(dataDir, "honored_contributor.csv")
	}

	//compute the correct input filename (pr_per_submitter-YYYY-MM.csv)
	inputFileName := filepath.Join(dataDir, "pr_per_submitter-"+monthToSelectFrom+".csv")

	// fail if the file does not exist else open the file
	f, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("Unable to read input file "+inputFileName+"\n", err)
	}
	defer f.Close()

	// validate that it has the correct format (CSV and column names)
	r := csv.NewReader(f)

	headerLine, err1 := r.Read()
	if err1 != nil {
		return fmt.Errorf("Unexpected error loading"+inputFileName+"\n", err)
	}

	if isVerbose {
		fmt.Println("Checking input file")
	}

	referencePrPerSubmitterHeader := []string{"user", "PR"}
	if !validateHeader(headerLine, referencePrPerSubmitterHeader, isVerbose) {
		return fmt.Errorf(" Error: header is incorrect.")
	} else {
		if isVerbose {
			fmt.Printf("  - Header is correct\n")
		}
	}

	// load the file in memory
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("Unexpected error loading \""+inputFileName+"\"\n", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("Error: No data available after the header\n")
	}
	if isVerbose {
		fmt.Println("  - At least one Submitter data available")
	}

	// TODO: pick a data line randomly
	// TODO: make a GitHub query to retrieve the contributors information (URL, avatar)
	// TODO: for the given user, retrieve all the PRs of that user in the given month
	// TODO: pick the required data and assemble it so that it can be outputed
	// TODO: output the file

	return nil
}
