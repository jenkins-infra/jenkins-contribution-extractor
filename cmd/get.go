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

	"github.com/google/go-github/v55/github"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves the commenter data, given org, project, and PR",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		getCommenters()

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//TODO: fix verb documentation
//TODO: add parameters (to verb and function)

//TODO: use authentication
// token := os.Getenv("GITHUB_AUTH_TOKEN")
// if token == "" {
// 	log.Fatal("Unauthorized: No token present")
// }



// Get the requested commenter data, extract it
func getCommenters() {

	var org string = "on4kjm"
	var prj string = "FLEcli"
	var pr int = 1

	fmt.Println("Fetching comments")
	comments, err := fetchComments(org, prj, pr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Load the collected comment data in the output data structure
	output_data_list := load_data(org,prj,strconv.Itoa(pr), comments)

	fmt.Printf("%v\n", output_data_list)


	//TODO: write slice to CSV and save it
}

// Get the comment data from GitHub.
func fetchComments(org string, project string, pr_nbr int) ([]*github.PullRequestComment, error) {

	client := github.NewClient(nil)

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
func load_data(org string, prj string, pr_number string, comments []*github.PullRequestComment) ([][]string){
	var output_slice [][]string
	for i, comment := range comments {
		var output_record []string
		// "PR_ref","commenter", "month"

		pr_ref := fmt.Sprintf("%s/%s/%s", org, prj, pr_number)
		commenter := *comment.GetUser().Login
		timestamp := comment.GetCreatedAt().String()
		month := timestamp[0:7]

		// create record
		output_record = append(output_record,pr_ref,commenter,month)

		fmt.Printf("%v. %s, %s, %s\n",i+1, pr_ref,commenter,month)
		//append the record to the list we are building
		output_slice = append(output_slice, output_record)
	}

	return output_slice
}