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
	"regexp"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var isSkipClosed bool

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

		// We probably have a file with users to exclude
		if excludeFileName != "" {
			var err error
			err, excludedGithubUsers = load_exclusions(excludeFileName)
			if err != nil {
				return fmt.Errorf("invalid excluded user list => %v\n", err)
			}
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
	prCmd.PersistentFlags().BoolVarP(&isSkipClosed, "skip_closed", "", false, "Skip PR marked as closed.")
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

	// Check whether we will not get too many items, forcing us to split
	nbrOfItems, errGetTotal := getTotalNumberOfItems(searchedOrg, searchedMonth)
	if errGetTotal != nil {
		return errGetTotal
	}
	if isRootDebug {
		loggers.debug.Printf("Total number of items in month: %d\n", nbrOfItems)
	}

	var output_data_list []string

	if nbrOfItems > 1000 {
		hasMore := true
		i := 0
		startDate := ""
		endDate := ""
		loadedItems := 0
		for hasMore {
			startDate, endDate, hasMore = splitPeriodForMaxQueryItem(nbrOfItems, searchedMonth, i)
			output_list, itemsInIterations, err := getData(searchedOrg, startDate, endDate)
			if err != nil {
				return err
			}
			output_data_list = append(output_data_list, output_list...)
			loadedItems = loadedItems + itemsInIterations
			i++
		}
		if isRootDebug {
			loggers.debug.Printf("expected nbr of items (%d) vs. retrieved nbr of items (%d)\n", nbrOfItems, loadedItems)
		}
		if nbrOfItems != loadedItems {
			return fmt.Errorf("Expected nbr of items (%d) does not match retrieved nbr of items (%d)", nbrOfItems, loadedItems)
		}

	} else {
		startDate, endDate := getStartAndEndOfMonth(searchedMonth)
		// A value of 0001-01-01 and 0001-01-31 indicates a rubbish input. Input is validated higher, so we don't check this here

		//get the data from GitHub
		var err error
		output_data_list, _, err = getData(searchedOrg, startDate, endDate)
		if err != nil {
			return err
		}
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

	// We make no difference  whether data was found or not

	// Creates, overwrites, or opens for append depending on the combination
	out, newIsNoHeader := openOutputCSV(outputFileName, isAppend, globalIsNoHeader)
	defer out.Close()

	//TODO: Refactor
	header := "org,repository,number,url,state,created_at,merged_at,user.login,month_year,title"
	writeCSVtoFile(out, isAppend, newIsNoHeader, header, output_data_list)
	out.Close()

	return nil
}

// Gets the data from GitHub for all PRs created in the given month
func getData(searchedOrg string, startDate string, endDate string) ([]string, int, error) {
	// initLoggers()

	//note: parameters are checked at Cobra API level

	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var prList []string
	issueCount := 0

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
								Login        string
								ResourcePath string
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
		barDescription := fmt.Sprintf("%s %s->%s    ", searchedOrg, startDate, endDate)
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
			//TODO: treat error
			_ = bar.Add(1)
		}

		i := 0
		for {
			err := client.Query(context.Background(), &prQuery, variables)
			if err != nil {
				if isRootDebug {
					loggers.debug.Printf("Error performing query: %v\n", err)
				}
				var emptyList []string
				return emptyList, 0, err
			}

			if isRootDebug {
				loggers.debug.Printf("GitHub query successful: retrieved %d PRs", len(prQuery.Search.Edges))
			}

			// We update the progress bar with the total size we get with the first call
			totalIssues := prQuery.Search.IssueCount
			if i == 0 && !isVerbose {
				if isRootDebug {
					loggers.debug.Printf("Expecting to treat %d items. Resetting progress bar\n", totalIssues)
				}
				issueCount = totalIssues
				// +1 to compensate the initial add() we used to display the bar
				bar.ChangeMax(totalIssues + 1)
			}

			for ii, singlePr := range prQuery.Search.Edges {

				if !isVerbose {
					//TODO: treat error
					_ = bar.Add(1)
				}

				createdAtStr := ""
				if !singlePr.Node.PullRequest.CreatedAt.IsZero() {
					createdAtStr = singlePr.Node.PullRequest.CreatedAt.Format(time.RFC3339) //created At
				}

				mergedAtStr := ""
				if !singlePr.Node.PullRequest.MergedAt.IsZero() {
					mergedAtStr = singlePr.Node.PullRequest.MergedAt.Format(time.RFC3339) //mergedAt, if available
				}

				author := ""
				// Applications have a RessourcePath that starts with "/apps" and we don't count them
				regexpApp := regexp.MustCompile(`^\/apps\/`)
				if regexpApp.MatchString(singlePr.Node.PullRequest.Author.ResourcePath) {
					if isRootDebug {
						loggers.debug.Printf("   %d-%d (%d/%d)  Skipping %s because user %s is an application.\n",
							i, ii, (i*100)+ii, totalIssues,
							singlePr.Node.PullRequest.Url,
							singlePr.Node.PullRequest.Author.ResourcePath)
					}
					continue
				} else {
					// Is it an author that we don't want to track ?
					authorToCheck := singlePr.Node.PullRequest.Author.Login
					if !isExcludedAuthor(excludedGithubUsers,authorToCheck){
						author = authorToCheck
					}else {
						continue
					}
					
				}

				// Skip PR if the status is CLOSED (Same behavior as the bash extraction)
				if isSkipClosed {
					if singlePr.Node.PullRequest.State == "CLOSED" {
						if isRootDebug {
							loggers.debug.Printf("   %d-%d (%d/%d)  Skipping %s because it is CLOSED\n",
								i, ii, (i*100)+ii, totalIssues, singlePr.Node.PullRequest.Url)
						}
						continue
					}
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
					author,                                                // PR's author
					singlePr.Node.PullRequest.CreatedAt.Format("2006-01"), // Creation month-year
					cleanedTitle,                                          // PR's description
				)

				if isRootDebug {
					loggers.debug.Printf("   %d-%d (%d/%d)  %s\n", i, ii, (i*100)+ii, totalIssues, dataLine)
				}
				prList = append(prList, dataLine)

				if isVerbose {
					fmt.Printf("%d-%d (%d/%d)  %s    %s\n", i, ii, (i*100)+ii, totalIssues, singlePr.Node.PullRequest.Author.Login, singlePr.Node.PullRequest.Url)
				}
			}

			if !prQuery.Search.PageInfo.HasNextPage {
				if isRootDebug {
					loggers.debug.Printf("HasNextPage is set to false. Exiting loop...\n")
				}
				break
			}
			variables["pullRequestCursor"] = githubv4.NewString(prQuery.Search.PageInfo.EndCursor)
			i++

			// Function has its own debug trace
			checkIfSufficientQuota_2(2,
				prQuery.RateLimit.Remaining,
				prQuery.RateLimit.Limit,
				prQuery.RateLimit.ResetAt)
		}
	}
	// as the progress exist doesn't do it
	fmt.Printf("\n")
	return prList, issueCount, nil
}

// Makes a call to GitHub to get the total number of items. We can handle only 1K items in one
// series of call. If above 1K we will have to split by decreasing the date range.
func getTotalNumberOfItems(searchedOrg string, searchedMonth string) (int, error) {
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

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
						Url    string
						Number int
						Title  string
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

	// Make the call
	err := client.Query(context.Background(), &prQuery, variables)
	if err != nil {
		if isRootDebug {
			loggers.debug.Printf("Error performing query: %v\n", err)
		}
		return 0, err
	}

	if isRootDebug {
		loggers.debug.Printf("GitHub query successful: retrieved %d PRs", len(prQuery.Search.Edges))
	}

	// We update the progress bar with the total size we get with the first call
	totalIssues := prQuery.Search.IssueCount

	return totalIssues, nil
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
			resourcePath
          }
          title
        }
      }
    }
  }
}
*/
