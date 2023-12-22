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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Flag indicating whether a backup of the file is required.
var remove_requireBackup bool

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <user> <filename>",
	Short: "Removes given user's data in CSV",
	Long: `This command will remove, for a given user, every data line from the data CSV.
A backup of the treated file can be requested.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		//call requires two parameters (org and month)
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// First argument is the GitHub user to remove
		// Second argument is the filename where to remove the user
		err := performRemove(args[0], args[1], remove_requireBackup)
		if err != nil {
			return err
		}
		return nil
	},
}

// Initialises COBRA for this command
func init() {
	rootCmd.AddCommand(removeCmd)

	// Add local flag
	removeCmd.Flags().BoolVarP(&remove_requireBackup, "backup", "b", true, "Make a backup of the original file")
}

// Main function of the REMOVE command
func performRemove(githubUser string, fileToClean_name string, isBackup bool) error {

	//test whether it is a valid GitHub user
	if !isValidOrgFormat(githubUser) {
		return fmt.Errorf("ERROR: %s is not a valid GitHub user.\n", githubUser)
	}

	//Do we have an existing file to clean ?
	if !fileExist(fileToClean_name) {
		return fmt.Errorf("ERROR: %s is not an existing file.\n", githubUser)
	}

	//Load input file
	if isVerbose {
		fmt.Printf("Loading the file to clean (%s) \n", fileToClean_name)
	}
	err, csvToClean_List := loadCSVtoClean(fileToClean_name)
	if err != nil {
		return err
	}

	// Try to clean the file
	if isVerbose {
		fmt.Printf("Removing entries for user \"%s\" \n", githubUser)
	}
	cleanedCsv_List := cleanCsvList(csvToClean_List, githubUser)

	//Was it useful ?
	// cleaned file should be shorter than the initial file
	if len(cleanedCsv_List) < len(csvToClean_List) {
		//if backup
		//  compute backup filename
		//  write in as the backup file
		//endif
		//write out (cleaned file)
	} else {
		fmt.Printf("Didn't find an entry for user \"%s\" in file \"%s\" \n", githubUser, fileToClean_name)
	}

	//If the cleaned file is larger than the original file something went horribly wrong....
	if len(cleanedCsv_List) > len(csvToClean_List) {
		return fmt.Errorf("[ERROR] Something went horribly wrong: the cleaned file increased in size !!!!???\n")
	}

	return nil
}

// load input file
func loadCSVtoClean(fileName string) (error, []string) {

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Unable to read input file %s: %v\n", fileName, err), nil
	}
	defer f.Close()

	var loadedFile []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		loadedFile = append(loadedFile, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error loading \"%s\": %v", fileName, err), nil
	}

	if len(loadedFile) <= 1 {
		return fmt.Errorf("Error: \"%s\" seems empty. Retrieved %d lines.", fileName, len(loadedFile)), nil
	}

	return nil, loadedFile
}

// Removes every list item where the gitHub user is present
func cleanCsvList(csvToClean_List []string, githubUser string) []string {
	var cleanedList []string

	for _, line := range csvToClean_List {
		if !strings.Contains(line, githubUser) {
			cleanedList = append(cleanedList, line)
		}
	}

	return cleanedList
}
