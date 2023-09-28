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
	// "log"
	"strconv"

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
retrieved from an environment variable (default is "GITHUB_TOKEN" but can be overridden with a flag)
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
		initLoggers()
		if isDebug {
			loggers.debug.Println("******** New debug session ********")
		}

		if isDebug {
			fmt.Println("*** Debug mode enabled ***\nSee \"debug.log\" for the trace")
		}

		if isDebug {
			limit, remaining := get_quota_data()
			loggers.debug.Printf("Start quota: %d/%d\n", remaining, limit)
		}

		getCommenters(args[0], globalIsAppend, globalIsNoHeader, outputFileName)

		if isDebug {
			limit, remaining := get_quota_data()
			loggers.debug.Printf("End quota: %d/%d\n", remaining, limit)
		}

	},
}

// Cobra initialize
func init() {
	rootCmd.AddCommand(getCmd)

}

//TODO: Logging
// per query, total and type of comments

//TODO: handle secondary quota error

//**********
// This is where it starts and the magic happens
//**********

// Get the requested commenter data, extract it, and write it to CSV
func getCommenters(prSpec string, isAppend bool, isNoHeader bool, outputFileName string) int {

	org, prj, pr, err := validatePRspec(prSpec)
	if err != nil {
		fmt.Printf("Unexpected error in PR specification (%v)\n Skipping %s\n", err, prSpec)
		return 0
	}

	if isVerbose {
		fmt.Printf("Fetching comments for %s\n", prSpec)
	}
	// ------
	// Retrieving all comments for the given PR from GitHub
	comments, err := fetchComments(org, prj, pr)
	if err != nil {
		if !isVerbose {
			fmt.Printf("Fetching comments for %s\n", prSpec)
		}
		fmt.Printf("Error: %v\n   Skipping....\n", err)
		return 0
	}
	// Load the collected comment data in the output data structure
	output_comment_list := load_issueComments(org, prj, strconv.Itoa(pr), comments)

	// ------
	// Retrieving all review comments for the given PR from GitHub
	review_comments, err := fetchReviews(org, prj, pr)
	if err != nil {
		if !isVerbose {
			fmt.Printf("Fetching review comments for %s\n", prSpec)
		}
		fmt.Printf("Error: %v\n   Skipping....\n", err)
		return 0
	}
	// Load the collected comment data in the output data structure
	output_review_list := load_reviewComments(org, prj, strconv.Itoa(pr), review_comments)

	// Assemble the two lists
	output_data_list := append(output_comment_list, output_review_list...)

	// Only process if data was found
	nbrOfComments := len(output_data_list)
	if nbrOfComments > 0 {

		if isVerbose {
			fmt.Printf("   Found %d comments (%d review comments and %d general comments).\n",
				nbrOfComments, len(output_review_list), len(output_comment_list))
		}

		if isDebug {
			loggers.debug.Printf("For \"%s\" found %d comments (%d review comments and %d general comments).\n",
			prSpec, nbrOfComments, len(output_review_list), len(output_comment_list))
		}

		// Creates, overwrites, or opens for append depending on the combination
		out, newIsNoHeader := openOutputCSV(outputFileName, isAppend, isNoHeader)
		defer out.Close()

		writeCSVtoFile(out, isAppend, newIsNoHeader, output_data_list)
		out.Close()
	} else {
		if isVerbose {
			fmt.Println("   No comments found for PR, skipping...")
		}
	}
	return nbrOfComments
}

// Get the comment data from GitHub.
func fetchComments(org string, project string, pr_nbr int) ([]*github.IssueComment, error) {

	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)

	client := github.NewClient(nil).WithAuthToken(ghToken)

	var allComments []*github.IssueComment
	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		comments, resp, err := client.Issues.ListComments(context.Background(), org, project, pr_nbr, opt)
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
// TODO: create a test
func load_issueComments(org string, prj string, pr_number string, comments []*github.IssueComment) [][]string {
	var output_slice [][]string
	for _, comment := range comments {
		var output_record []string

		pr_ref := fmt.Sprintf("%s/%s/%s", org, prj, pr_number)
		commenter := *comment.GetUser().Login
		timestamp := comment.GetCreatedAt().String()
		month := timestamp[0:7]

		// create record
		output_record = append(output_record, pr_ref, commenter, month)

		//append the record to the list we are building
		output_slice = append(output_slice, output_record)
	}

	return output_slice
}

// Get the reviews data from GitHub.
func fetchReviews(org string, project string, pr_nbr int) ([]*github.PullRequestReview, error) {

	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)

	client := github.NewClient(nil).WithAuthToken(ghToken)

	var allComments []*github.PullRequestReview
	opt := &github.ListOptions{PerPage: 10}

	for {
		comments, resp, err := client.PullRequests.ListReviews(context.Background(), org, project, pr_nbr, opt)
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
// TODO: create a test
func load_reviewComments(org string, prj string, pr_number string, reviews []*github.PullRequestReview) [][]string {
	var output_slice [][]string
	for _, comment := range reviews {
		var output_record []string

		pr_ref := fmt.Sprintf("%s/%s/%s", org, prj, pr_number)
		commenter := *comment.GetUser().Login
		timestamp := comment.GetSubmittedAt().String()
		month := timestamp[0:7]

		// create record
		output_record = append(output_record, pr_ref, commenter, month)

		//append the record to the list we are building
		output_slice = append(output_slice, output_record)
	}

	return output_slice
}
