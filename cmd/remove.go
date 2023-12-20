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

	"github.com/spf13/cobra"
)

//Flag indicating whether a backup of the file is required.
var remove_requireBackup bool

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [user] [filename]",
	Short: "Removes given user's data in CSV",
	Long: `This command will remove every data line for the user passed as an argument from the data CSV.
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

func init() {
	rootCmd.AddCommand(removeCmd)

	// Add local flag
	removeCmd.Flags().BoolVarP(&remove_requireBackup,"backup", "b", true, "Make a backup of the original file")
}

// isValidOrgFormat()

// Main function of the command
func performRemove(githubUser string, fileToClean_name string, isBackup bool) error {

	//TODO: test whether it is a valid GitHub user

	//TODO: Do we have an existing file to clean ?

	//TODO: what type is it.
	
	return nil
}
