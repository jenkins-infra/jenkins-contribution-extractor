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
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "submitters [org] [YYYY-MM]",
	Short: "Get all PRs (and their submitters) for a given month and org.",
	Long:  `Get all PRs (and their submitters) for a given month and org.`,
	Args: func(cmd *cobra.Command, args []string) error {
		//call requires two parameters (org and month)
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return err
		}
		if !isValidOrgFormat(args[0]) {
			return fmt.Errorf("ERROR: %s is not a valid GitHub user or Org name.\n", args[0])
		}

		if !isValidMonthFormat(args[1]) {
			return fmt.Errorf("ERROR: %s is not a valid month (should be \"YYYY-MM\").\n", args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := performSearch(args[0], args[1])
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(prCmd)

}

// Main function: it searches GitHub for all PRs created in the given month
func performSearch(searchedOrg string, searchedMonth string) error {
	initLoggers()

	//note: parameters are checked at Cobra API level

	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	{
		var prQuery struct {
			Viewer struct {
				Login string
			}
			RateLimit struct {
				Limit     int
				Cost      int
				Remaining int
				ResetAt   time.Time
			}
			Search struct {
				IssueCount int
				Edges      []struct {
					Node struct {
						PullRequest struct {
							Author struct {
								Login string
							}
							CreatedAt time.Time
							ClosedAt  time.Time
							Url       string
							Number    int
						} `graphql:"... on PullRequest"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"search(first: $count, after: $pullRequestCursor, query: $searchQuery, type: ISSUE)"`
		}

		//TODO: pass date parameters to GraphQL
		//TODO: get beginning and end date of month

		// now := time.Now()
		// currentYear, currentMonth, _ := now.Date()
		// currentLocation := now.Location()

		// firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		// lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		variables := map[string]interface{}{
			"searchQuery":       githubv4.String(fmt.Sprintf(`org:%s is:pr -author:app/dependabot -author:app/renovate -author:app/github-actions -author:jenkins-infra-bot created:2023-09-01..2023-09-30`, githubv4.String(searchedOrg))),
			"count":             githubv4.Int(100),
			"pullRequestCursor": (*githubv4.String)(nil), // Null after argument to get first page.
		}

		//TODO: solve issue of different default output file for this command
		//TODO: write header
		//TODO: write CSV

		i := 0
		for {
			err := client.Query(context.Background(), &prQuery, variables)
			if err != nil {
				return (err)
			}

			totalIssues := prQuery.Search.IssueCount
			for ii, singlePr := range prQuery.Search.Edges {
				fmt.Printf("%d-%d (%d/%d)  %s %s\n", i, ii, (i*100)+ii, totalIssues, singlePr.Node.PullRequest.Author.Login, singlePr.Node.PullRequest.Url)
			}

			if !prQuery.Search.PageInfo.HasNextPage {
				break
			}
			variables["pullRequestCursor"] = githubv4.NewString(prQuery.Search.PageInfo.EndCursor)
			i++
		}
		// printJSON(prQuery)
	}
	return nil
}

//GitHub Graphql query
// {
// 	rateLimit {
// 	  limit
// 	  cost
// 	  remaining
// 	  resetAt
// 	}
// 	search(
// 	  query: "org:jenkinsci is:pr -author:app/dependabot -author:app/renovate -author:jenkins-infra-bot created:2023-09-01..2023-09-30"
// 	  type: ISSUE
// 	  first: 100
// 	) {
// 	  issueCount
// 	  pageInfo {
// 		endCursor
// 		hasNextPage
// 	  }
// 	  edges {
// 		node {
// 		  ... on PullRequest {
// 			author {
// 			  login
// 			}
// 			createdAt
// 			closedAt
// 			url
// 			number
// 		  }
// 		}
// 	  }
// 	}
//   }
