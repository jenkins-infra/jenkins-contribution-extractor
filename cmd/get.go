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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [commenters|pr]",
	Short: "Retrieves data from GitHub (PRs or Commenters)",
	Long: `Long description
`,
}

// Cobra initialize
func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVarP(&outputFileName, "out", "o", "jenkins_commenters_data.csv", "Output file name.")
	getCmd.PersistentFlags().StringVarP(&excludeFileName, "excludeFile", "x", "", "Name of the file containing the github handles to exclude from the data collection.")
	getCmd.PersistentFlags().BoolVarP(&globalIsAppend, "append", "a", false, "Appends data to existing output file.")
	getCmd.PersistentFlags().BoolVarP(&globalIsNoHeader, "no_header", "", false, "Doesn't add a header to file (implied when appending to existing file).")

	rootCmd.PersistentFlags().BoolVarP(&isRootDebug, "debug", "", false, "Display debug information (super verbose mode)")

}
