/*
Copyright © 2023 Jean-Marc Meessen jean-marc@meessen-web.org

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
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	// "github.com/google/go-github/v55/github"
)

// var cfgFile string
var outputFileName string
var ghTokenVar string
var isVerbose bool
var globalIsAppend bool
var globalIsNoHeader bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jenkins-get-commenters [PR list CSV]",
	Short: "Retrieve the commenters from a PR list",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if !fileExist(args[0]) {
			return fmt.Errorf("Invalid file\n")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Processing \"%s\"\n", args[0])

		// read the relevant data from the file (and checking it)
		prList, result := loadPrListFile(args[0], isVerbose)
		if !result {
			fmt.Printf("Could not load \"%s\"\n",args[0])
			os.Exit(1)
		}

		for i, pr_line := range prList {
			if i ==0 {
				// handle the first create
			}
			getCommenters(pr_line,globalIsAppend,globalIsNoHeader,outputFileName)
		}
		
		// loop though the file
		//   call the get command for each PR

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Cobra initialization
func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "out", "o", "jenkins_commenters_data.csv", "Output file name.")
	rootCmd.PersistentFlags().StringVarP(&ghTokenVar, "token_var", "t", "GITHUB_TOKEN", "The environment variable containing the GitHub token.")
	rootCmd.PersistentFlags().BoolVarP(&globalIsAppend, "append", "a", false, "Appends data to existing output file.")
	rootCmd.PersistentFlags().BoolVarP(&globalIsNoHeader, "no_header", "", false, "Doesn't add a header to file (implied when appending to existing file).")
	rootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Displays useful info during the extraction.")

	//Disable the Cobra completion options
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Don't sort flags in alphabetical order
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

}

var referenceCSVheader = []string{"org", "repository", "number", "url", "state", "created_at", "merged_at", "user.login", "month_year", "title"}

// Loads the data from a file and try to parse it as a CSV
func loadPrListFile(fileName string, isVerbose bool) ([]string, bool) {

	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Unable to read input file "+fileName+"\n", err)
		return nil, false
	}
	defer f.Close()

	r := csv.NewReader(f)

	headerLine, err1 := r.Read()
	if err1 != nil {
		log.Printf("Unexpected error loading"+fileName+"\n", err)
		return nil, false
	}

	if isVerbose {
		fmt.Println("Checking input file")
	}

	if !validateHeader(headerLine, referenceCSVheader, isVerbose) {
		fmt.Println(" Error: header is incorrect.")
		return nil, false
	} else {
		if isVerbose {
			fmt.Printf("  - Header is correct\n")
		}
	}

	records, err := r.ReadAll()
	if err != nil {
		log.Printf("Unexpected error loading \""+fileName+"\"\n", err)
		return nil, false
	}

	if len(records) < 2 {
		fmt.Printf("Error: No data available after the header\n")
		return nil, false
	}
	if isVerbose {
		fmt.Println("  - At least one Pull Request data available")
	}

	var prList []string
	org_regexp, _ := regexp.Compile(`^(jenkinsci|jenkins-infra)$`)
	prj_regexp, _ := regexp.Compile(`^[\w-\.]+$`) // see https://stackoverflow.com/questions/59081778/rules-for-special-characters-in-github-repository-name
	pr_regexp, _ := regexp.Compile(`^\d+$`)

	// Check the loaded data
	for i, dataLine := range records {
		//Skip header line as it has already been checked
		if i == 0 {
			continue
		}

		// Org must be within list (jenkinsci and jenkins_infra)
		org := dataLine[0]
		if !org_regexp.MatchString(strings.ToLower(org)) {
			if isVerbose {
				fmt.Printf(" Error: ORG field \"%s\" is not the expected value (\"jenkinsci\" or \"jenkins-infra\")", org)
			}
			return nil, false
		}

		// project name must be "^[\w-\.]+$"
		prj := dataLine[1]
		if !prj_regexp.MatchString(strings.ToLower(prj)) {
			if isVerbose {
				fmt.Printf(" Error: PRJ field \"%s\" is not of the expected format", prj)
			}
			return nil, false
		}

		// PR number must be a number
		prNbr := dataLine[2]
		if !pr_regexp.MatchString(prNbr) {
			if isVerbose {
				fmt.Printf(" Error: PR field \"%s\" is not a (positive) number", prNbr)
			}
			return nil, false
		}

		prInfo := fmt.Sprintf("%s/%s/%s", org, prj, prNbr)
		prList = append(prList, prInfo)

	}

	if isVerbose {
		fmt.Printf("Successfully loaded \"%s\" (%d Pull Request to analyze)\n", fileName, len(prList))
	}

	return prList, true
}

// Checks whether the retrieved header is equivalent to the reference header
func validateHeader(header []string, referenceHeader []string, isVerbose bool) bool {
	if len(header) != len(referenceHeader) {
		if isVerbose {
			fmt.Printf(" Error: field number mismatch (found %d, wanted %d)\n", len(header), len(referenceHeader))
		}
		return false
	}
	for i, v := range header {
		if v != referenceHeader[i] {
			if isVerbose {
				fmt.Printf(" Error: not the expected header field at column %d (found \"%v\", wanted \"%v\")\n", i+1, v, referenceHeader[i])
			}
			return false
		}
	}
	return true

}
