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
	"log"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string
var outputFileName string
var ghTokenVar string
var isVerbose bool
var isRootDebug bool
var globalIsAppend bool
var globalIsNoHeader bool
var globalIsBigFile bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	// Use:   "jenkins-stats [PR list CSV]",
	Short: "Retrieve Jenkins related usage stats from GitHub",
	Long: `Retrieve data from GitHub that can be useful to evaluate the health and activity of a community.
It currently gets data about Pull Request submitters and commenters on those Pull Requests.
`,
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

	rootCmd.PersistentFlags().StringVarP(&ghTokenVar, "token_var", "t", "GITHUB_TOKEN", "The environment variable containing the GitHub token.")
	rootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Displays useful info during the extraction.")

	//Disable the Cobra completion options
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Don't sort flags in alphabetical order
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	err := rootCmd.PersistentFlags().MarkHidden("debug")
	if err != nil {
		log.Printf("Error hiding debug flag: %v\n", err)
	}

}
