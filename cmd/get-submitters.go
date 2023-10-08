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
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
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
			return fmt.Errorf("ERROR: %s is not a valid month (should be \"YYYY-MM\").\n", args[1])
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

	//TODO: separate output default: https://github.com/spf13/cobra/issues/553 and https://travis.media/how-to-use-subcommands-in-cobra-go-cobra-tutorial/

}

// *************************
// *************************

// Main function: it searches GitHub for all PRs created in the given month and writes it to a CSV
func performSearch(searchedOrg string, searchedMonth string) error {
	initLoggers()
	if isRootDebug {
		loggers.debug.Println("******** New \"Get Submitters\" debug session ********")
	}

	if isRootDebug {
		fmt.Print("*** Debug mode enabled ***\nSee \"debug.log\" for the trace\n\n")

		limit, remaining, _, _ := get_quota_data_v4()
		loggers.debug.Printf("Start quota: %d/%d\n", remaining, limit)
	}

	//get the data from GitHub
	output_data_list, err := getData(searchedOrg, searchedMonth)
	if err != nil {
		return err
	}

	// Write to CSV
	isAppend := globalIsAppend
	if !globalIsAppend {
		// Meaning that we need to create a new file
		if fileExist(outputFileName) {
			os.Remove(outputFileName)
		}
		isAppend = true
	}

	nbrOfPRs := len(output_data_list)
	if nbrOfPRs > 0 {

		// Creates, overwrites, or opens for append depending on the combination
		out, newIsNoHeader := openOutputCSV(outputFileName, isAppend, globalIsNoHeader)
		defer out.Close()

		header := "org,repository,number,url,state,created_at,merged_at,user.login,month_year,title"
		writeCSVtoFile(out, isAppend, newIsNoHeader, header, output_data_list)
		out.Close()
	} else {
		if isVerbose {
			fmt.Println("   No comments found for PR, skipping...")
		}
	}

	return nil
}

// Gets the data from GitHub for all PRs created in the given month
func getData(searchedOrg string, searchedMonth string) ([]string, error) {
	initLoggers()

	//note: parameters are checked at Cobra API level

	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var prList []string

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
							Repository struct {
								Name  string
								Owner struct {
									Login string
								}
							}
							Author struct {
								Login string
							}
							CreatedAt time.Time
							MergedAt  time.Time
							State     string
							Url       string
							Number    int
							Title     string
						} `graphql:"... on PullRequest"`
					}
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"search(first: $count, after: $pullRequestCursor, query: $searchQuery, type: ISSUE)"`
		}

		startDate, endDate := getStartAndEndOfMonth(searchedMonth)
		// A value of 0001-01-01 and 0001-01-31 indicates a rubbish input. Input is validated higher, so we don't check this here

		variables := map[string]interface{}{
			"searchQuery": githubv4.String(
				fmt.Sprintf(`org:%s is:pr -author:app/dependabot -author:app/renovate -author:app/github-actions -author:jenkins-infra-bot created:%s..%s`,
					githubv4.String(searchedOrg),
					githubv4.String(startDate),
					githubv4.String(endDate),
				),
			),
			"count":             githubv4.Int(100),
			"pullRequestCursor": (*githubv4.String)(nil), // Null after argument to get first page.
		}

		//TODO: solve issue of different default output file for this command
		//TODO: handle quota wait

		var bar *progressbar.ProgressBar
		barDescription := fmt.Sprintf("%s %s    ", searchedOrg, searchedMonth)
		if !isVerbose {
			bar = progressbar.NewOptions(
				1000,
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetDescription(barDescription),
				progressbar.OptionSetPredictTime(false),
				progressbar.OptionShowBytes(false),
				progressbar.OptionFullWidth(),
				progressbar.OptionShowCount(),
			)
			bar.Add(1)
		}

		i := 0
		for {
			err := client.Query(context.Background(), &prQuery, variables)
			if err != nil {
				var emptyList []string
				return emptyList, err
			}

			// We update the progress bar with the total size we get with the first call
			totalIssues := prQuery.Search.IssueCount
			if i == 0 && !isVerbose {
				// +1 to compensate the initial add() we used to display the bar
				bar.ChangeMax(totalIssues + 1)
			}

			for ii, singlePr := range prQuery.Search.Edges {

				createdAtStr := ""
				if !singlePr.Node.PullRequest.CreatedAt.IsZero() {
					createdAtStr = singlePr.Node.PullRequest.CreatedAt.Format(time.RFC3339) //created At
				}

				mergedAtStr := ""
				if !singlePr.Node.PullRequest.MergedAt.IsZero() {
					mergedAtStr = singlePr.Node.PullRequest.MergedAt.Format(time.RFC3339) //mergedAt, if available
				}

				// clean and shorten the title
				cleanedTitle := truncateString(cleanBody(singlePr.Node.PullRequest.Title), 30)

				// data format: "org,repository,number,url,state,created_at,merged_at,user.login,month_year,title"

				dataLine := fmt.Sprintf("\"%s\",\"%s\",%d,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"",
					singlePr.Node.PullRequest.Repository.Owner.Login,      // Org
					singlePr.Node.PullRequest.Repository.Name,             //repository
					singlePr.Node.PullRequest.Number,                      // PR number
					singlePr.Node.PullRequest.Url,                         // PR's URL
					singlePr.Node.PullRequest.State,                       // PR's state
					createdAtStr,                                          // Creation date&time
					mergedAtStr,                                           // Merged date&time
					singlePr.Node.PullRequest.Author.Login,                // PR's author
					singlePr.Node.PullRequest.CreatedAt.Format("2006-01"), // Creation month-year
					cleanedTitle,                                          // PR's description
				)

				prList = append(prList, dataLine)
				if !isVerbose {
					bar.Add(1)
				}

				//TODO: show this only if in verbose mode
				if isVerbose {
					fmt.Printf("%d-%d (%d/%d)  %s    %s\n", i, ii, (i*100)+ii, totalIssues, singlePr.Node.PullRequest.Author.Login, singlePr.Node.PullRequest.Url)
				}
			}

			if !prQuery.Search.PageInfo.HasNextPage {
				break
			}
			variables["pullRequestCursor"] = githubv4.NewString(prQuery.Search.PageInfo.EndCursor)
			i++

			checkIfSufficientQuota_2(2,
				prQuery.RateLimit.Remaining,
				prQuery.RateLimit.Limit,
				prQuery.RateLimit.ResetAt)
		}
	}
	return prList, nil
}

//GitHub Graphql query. Test at https://docs.github.com/en/graphql/overview/explorer
/*
{
  rateLimit {
    limit
    cost
    remaining
    resetAt
  }
  search(
    query: "org:jenkinsci is:pr -author:app/dependabot -author:app/renovate -author:jenkins-infra-bot created:2023-09-01..2023-09-30"
    type: ISSUE
    first: 100
  ) {
    issueCount
    pageInfo {
      endCursor
      hasNextPage
    }
    edges {
      node {
        ... on PullRequest {
          repository {
            name
            owner {
              login
            }
          }
          number
          url
          state
          createdAt
          closedAt
          author {
            login
          }
          title
        }
      }
    }
  }
}
*/

func testBar() {
	var bar *progressbar.ProgressBar
	bar = progressbar.Default(int64(30))

	for i := 0; i < 100; i++ {
		bar.Add(1)
		if i == 0 {
			bar.ChangeMax(100)
		}
		time.Sleep(40 * time.Millisecond)
	}
}
