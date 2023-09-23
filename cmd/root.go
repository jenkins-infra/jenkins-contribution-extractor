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
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
	Use:   "jenkins-get-commenters",
	Short: "Retrieve the commenters from a PR list",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "out", "o", "jenkins_commenters_data.csv", "Output file name.")
	rootCmd.PersistentFlags().StringVarP(&ghTokenVar, "token_var", "t", "GITHUB_TOKEN", "The environment variable containing the GitHub token.")
	rootCmd.PersistentFlags().BoolVarP(&globalIsAppend, "append", "a", false, "Appends data to existing output file.")
	rootCmd.PersistentFlags().BoolVarP(&globalIsNoHeader, "no_header", "", false, "Doesn't add a header to file (implied when appending to existing file).")
	rootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Displays useful info during the extraction.")

	//Disable the Cobra completion options
	rootCmd.CompletionOptions.DisableDefaultCmd = true

}

// Load the GitHub token from the specified environment variable
func loadGitHubToken(envVariableName string) string {
	token := os.Getenv(envVariableName)
	if token == "" {
		fmt.Println("Unauthorized: No token present")
		//This is a major error: we crash out of the program
		os.Exit(0)
	}
	return token
}
