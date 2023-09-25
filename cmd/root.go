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
		// check if the file is a correctly formatted CSV
		loadPrListFile(args[0], isVerbose)

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
func loadPrListFile(fileName string, isVerbose bool) bool {

	var isValidTable = true

	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Unable to read input file "+fileName+"\n", err)
		return false
	}
	defer f.Close()

	r := csv.NewReader(f)

	headerLine, err1 := r.Read()
	if err1 != nil {
		log.Printf("Unexpected error loading"+fileName+"\n", err)
		return false
	}

	if isVerbose {
		fmt.Println("Checking input file format")
		fmt.Printf("  - Number of columns defined in header: %d\n", len(headerLine))
	}

	//TODO: check result
	validateHeader(headerLine, referenceCSVheader, isVerbose)

	// // first column should be empty
	// if firstLine[0] != "" {
	// 	fmt.Println("Not the expected first column name (should be empty)")
	// 	return false
	// }
	// if isVerboseCheck {
	// 	fmt.Println("  - File's header start with empty column name.")
	// }

	// //loop through columns to check headings
	// month_regexp, _ := regexp.Compile("20[0-9]{2}-[0-9]{2}")
	// for i, s := range firstLine {
	// 	if i != 0 {
	// 		if !month_regexp.MatchString(s) {
	// 			fmt.Printf("Column header %s is not of the expected format (YYYY-MM)\n", s)
	// 			return false
	// 		}
	// 	}
	// }
	// if isVerboseCheck {
	// 	endMonth := firstLine[len(firstLine)-1]
	// 	fmt.Printf("  - File's header data column format (\"20YY-MM\"). Most recent data is \"%s\"\n", endMonth)
	// }

	// nbrOfColumns := len(firstLine)
	// if nbrOfColumns < 3 {
	// 	fmt.Printf("Not enough monthly data available\n")
	// 	return false
	// }
	// if isVerboseCheck {
	// 	fmt.Printf("  - More than one month data available\n")
	// }

	// records, err := r.ReadAll()
	// if err != nil {
	// 	log.Printf("Unexpected error loading"+fileName+"\n", err)
	// 	return false
	// }

	// if len(records) < 2 {
	// 	fmt.Printf("No data available after the header\n")
	// 	return false
	// }
	// if isVerboseCheck {
	// 	fmt.Println("  - At least one submitter's data available")
	// }

	// //The GitHub user validation regexp (see https://stackoverflow.com/questions/58726546/github-username-convention-using-regex)
	// // should be regexp.Compile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`). But the dataset contains "invalid" data: username ending with a "-" or
	// // a double "-" in the name.
	// name_exp, _ := regexp.Compile(`^[a-zA-Z0-9\-]+$`)

	// //Check the loaded data
	// for i, dataLine := range records {
	// 	//Skip header line as it has already been checked
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	for ii, column := range dataLine {
	// 		//check the GitHub user (first columns)
	// 		if ii == 0 {
	// 			if !(len(column) < 40 && len(column) > 0 && name_exp.MatchString(column)) {
	// 				fmt.Printf("Submitter \"%s\" at line %d does not follow GitHub rules\n", column, i)
	// 				return false
	// 			}
	// 		} else {
	// 			// check the other columns is an integer (we don't check the sign)
	// 			if data_value, err := strconv.Atoi(column); err != nil {
	// 				fmt.Printf("Value \"%s\" at line %d (column %d) isn't an integer\n", column, i, ii)
	// 				return false
	// 			} else {
	// 				if data_value < 0 {
	// 					fmt.Printf("Value \"%s\" at line %d (column %d) is negative\n", column, i, ii)
	// 					return false
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// if isVerboseCheck {
	// 	fmt.Println("  - Number of data columns match header columns.")
	// 	fmt.Printf("  - Records have a valid GitHub username and number of submitted PRs. (%d data records)\n", len(records)-1)
	// }

	// if !isSilent {
	// 	fmt.Printf("\nSuccessfully checked \"%s\"\n   It is a valid Jenkins Submitter Pivot Table and can be processes\n\n", fileName)
	// }

	return isValidTable
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
