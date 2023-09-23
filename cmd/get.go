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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [PR Specification]",
	Short: "Retrieves the commenter data, given a PR specification",
	Long: `This command will get from GitHub for a given PR the list of comments
(author and month). 

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
retrieved from an environment variable (default is "GITHUB_TOKEN" but can be overriden with a flag)
`,
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

		getCommenters(args[0], globalIsAppend, globalIsNoHeader, outputFileName)

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//TODO: add parameters (to verb and function)
//TODO: handle secondary quota error

// Get the requested commenter data, extract it, and write it to CSV
func getCommenters(prSpec string, isAppend bool, isNoHeader bool, outputFileName string) {

	org, prj, pr, err := validatePRspec(prSpec)
	if err != nil {
		fmt.Printf("Unexpected error in PR specification (%v)\n Skipping %s\n", err, prSpec)
		return
	}

	if isVerbose {
		fmt.Printf("Fetching comments for %s\n", prSpec)
	}
	comments, err := fetchComments(org, prj, pr)
	if err != nil {
		if !isVerbose {
			fmt.Printf("Fetching comments for %s\n", prSpec)
		}
		fmt.Printf("Error: %v\n   Skipping....\n", err)
		return
	}

	// Only process if data was found
	nbrOfComments := len(comments)
	if nbrOfComments > 0 {

		if isVerbose {
			fmt.Printf("   Found %d comments.\n", nbrOfComments)
		}
		// Load the collected comment data in the output data structure
		output_data_list := load_data(org, prj, strconv.Itoa(pr), comments)

		fmt.Printf("%v\n", output_data_list)

		//TODO: open output
		//TODO: write slice to CSV and save it
	} else {
		if isVerbose {
			fmt.Println("   No comments found for PR, skipping...")
		}
	}
}

// Get the comment data from GitHub.
func fetchComments(org string, project string, pr_nbr int) ([]*github.PullRequestComment, error) {

	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)

	client := github.NewClient(nil).WithAuthToken(ghToken)

	var allComments []*github.PullRequestComment
	opt := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		comments, resp, err := client.PullRequests.ListComments(context.Background(), org, project, pr_nbr, opt)
		if err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allComments, nil
}

// Load the collected comment data in the output data structure
func load_data(org string, prj string, pr_number string, comments []*github.PullRequestComment) [][]string {
	var output_slice [][]string
	for i, comment := range comments {
		var output_record []string
		// "PR_ref","commenter", "month"

		pr_ref := fmt.Sprintf("%s/%s/%s", org, prj, pr_number)
		commenter := *comment.GetUser().Login
		timestamp := comment.GetCreatedAt().String()
		month := timestamp[0:7]

		// create record
		output_record = append(output_record, pr_ref, commenter, month)

		fmt.Printf("%v. %s, %s, %s\n", i+1, pr_ref, commenter, month)
		//append the record to the list we are building
		output_slice = append(output_slice, output_record)
	}

	return output_slice
}

// validates that the supplied string is a valid PR specification
// in the form of "org/project/pr_nbr"
func validatePRspec(prSpec string) (org string, project string, prNbr int, err error) {
	splittedString := strings.Split(strings.TrimSpace(prSpec), "/")

	if len(splittedString) != 3 {
		return "", "", -1, fmt.Errorf("Invalid number of elements in prSpec. (expecting 3, found %v)\n", len(splittedString))
	}

	work_Org := splittedString[0]
	work_Project := splittedString[1]
	prString := splittedString[2]

	if strings.TrimSpace(work_Org) == "" {
		return "", "", -1, fmt.Errorf("Organization element in prSpec is empty\n")
	}
	if strings.TrimSpace(work_Project) == "" {
		return "", "", -1, fmt.Errorf("Project element in prSpec is empty\n")
	}
	if strings.TrimSpace(prString) == "" {
		return "", "", -1, fmt.Errorf("PR element in prSpec is empty\n")
	}

	work_prNbr, err := strconv.Atoi(strings.TrimSpace(prString))
	if err != nil {
		return "", "", -1, fmt.Errorf("PR part of psSpec is not numerical (%v)\n", err)
	}
	return work_Org, work_Project, work_prNbr, nil
}
