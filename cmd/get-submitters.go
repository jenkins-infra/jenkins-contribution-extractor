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

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "submitters org year month",
	Short: "A brief description of your command",
	Long: `This command will get the list of comments
	(author and month) from GitHub for a given PR . 
	
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pr called")
	},
}

func init() {
	getCmd.AddCommand(prCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
