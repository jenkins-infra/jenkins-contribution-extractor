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
)

// forPrCmd represents the forPr command
var forPrCmd = &cobra.Command{
	Use:   "forPr [PR Spec]",
	Short: "Retrieves the commenters of a given PR",
	Long: `This command will retrieve the commenters of the specified PR from GitHub. 

The PR is specified as "organization/project/PR number".

The output is a CVS file, specified with the "-o"/"--out" parameter. If not
defined it will take the default output filename.
Each record of the output contains the following information:
- PR specification
- Commenter's login name
- The month the comment was created (YYYY-MM)

The behavior can be controlled with various flags, such as appending to an existing
output file or overwriting it, header of no-header.

This query requires authenticated API call. The GitHub Token (Personal Access Token) is
retrieved from an environment variable (default is "GITHUB_TOKEN" but can be overridden with a flag)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if _, _, _, validateErr := validatePRspec(args[0]); validateErr != nil {
			return validateErr
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		initLoggers()
		if isRootDebug || isDebugGet {
			loggers.debug.Println("******** New debug session ********")
		}

		if isRootDebug || isDebugGet {
			fmt.Println("*** Debug mode enabled ***\nSee \"debug.log\" for the trace")
		}

		getCommenters(args[0], globalIsAppend, globalIsNoHeader, outputFileName)

	},
}

func init() {
	commentersCmd.AddCommand(forPrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// forPrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// forPrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
