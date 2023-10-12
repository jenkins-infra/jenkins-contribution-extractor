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
	"log"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// TODO: better variable name
var isDebugGet bool

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

	getCmd.PersistentFlags().BoolVarP(&isDebugGet, "debugGet", "", false, "Display debug information (super verbose mode) for the GET command")

	err := getCmd.PersistentFlags().MarkHidden("debugGet")
	if err != nil {
		log.Printf("Error hiding debug flag: %v\n", err)
	}

}

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

	if isDebugGet {
		loggers.debug.Printf("Fetching comments for %s\n", prSpec)
	}

	_, output_data_list := fetchComments_v4(org, prj, pr)

	// Only process if data was found
	nbrOfComments := len(output_data_list)
	if nbrOfComments > 0 {

		// Creates, overwrites, or opens for append depending on the combination
		out, newIsNoHeader := openOutputCSV(outputFileName, isAppend, isNoHeader)
		defer out.Close()

		header := "PR_ref,commenter,month"
		writeCSVtoFile(out, isAppend, newIsNoHeader, header, output_data_list)
		out.Close()
	} else {
		if isVerbose {
			fmt.Println("   No comments found for PR, skipping...")
		}
	}
	return nbrOfComments
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
  repository(name: "flecli", owner: "on4kjm") {
    pullRequest(number: 1) {
      reviews(first: 100) {
        nodes {
          bodyText
          createdAt
          author {
            login
          }
          comments(first: 100) {
            nodes {
              author {
                login
              }
              body
            }
          }
        }
      }
      comments(first: 100) {
        nodes {
          author {
            login
          }
          createdAt
          body
        }
        totalCount
      }
    }
  }
}
*/

var prQuery2 struct {
	Repository struct {
		Description string
		PullRequest struct {
			Title    string
			Comments struct {
				Nodes []struct {
					CreatedAt githubv4.DateTime
					Body      string
					Author    struct {
						Login string
					}
				}
			} `graphql:"comments(first: 100)"`
			Reviews struct {
				Nodes []struct {
					CreatedAt githubv4.DateTime
					BodyText  string
					Author    struct {
						Login string
					}
					Comments struct {
						Nodes []struct {
							CreatedAt githubv4.DateTime
							Body      string
							Author    struct {
								Login string
							}
						}
					} `graphql:"comments(first: 100)"`
				}
			} `graphql:"reviews(first: 100)"`
		} `graphql:"pullRequest(number: $pr)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
	RateLimit struct {
		Limit     int
		Cost      int
		Remaining int
		ResetAt   time.Time
	}
}

func fetchComments_v4(org string, prj string, pr int) (nbrComment int, output []string) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	//Check whether we have enough  quota left and wait if necessary
	// only if we are dealing with a PR list file bigger than the quota limit.
	if globalIsBigFile {
		checkIfSufficientQuota(5)
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(org),
		"name":  githubv4.String(prj),
		"pr":    githubv4.Int(pr),
	}

	err := client.Query(context.Background(), &prQuery2, variables)
	if err != nil {
		log.Printf("ERROR: Unexpected error getting comments: %v\n", err)
		return 0, nil
	}

	prSpec := fmt.Sprintf("%s/%s/%d", org, prj, pr)
	totalComments := 0
	dbgDateFormat := "2006-01-02 15:04:05"

	var output_slice []string

	for i, comment := range prQuery2.Repository.PullRequest.Comments.Nodes {
		//When there is no info about the user, it means it has been deleted
		author := comment.Author.Login
		if author == "" {
			author = "deleted_user"
		}

		output_slice = append(output_slice, createRecord(prSpec, author, comment.CreatedAt))
		if isDebugGet {
			loggers.debug.Printf("%d. %s, %s, \"%s\"\n", i+1, author, comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.Body))
		}
		totalComments++
	}
	if isDebugGet {
		loggers.debug.Printf("Nbr PR Comments: %d\n", len(prQuery2.Repository.PullRequest.Comments.Nodes))
	}

	for i, comment := range prQuery2.Repository.PullRequest.Reviews.Nodes {
		//When there is no info about the user, it means it has been deleted
		author := comment.Author.Login
		if author == "" {
			author = "deleted_user"
		}

		if isDebugGet {
			loggers.debug.Printf("%d. %s, %s, \"%s\"\n", i+1, author, comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.BodyText))
		}
		//Just guessing correct counting
		if comment.BodyText != "" {
			output_slice = append(output_slice, createRecord(prSpec, author, comment.CreatedAt))
			totalComments++
		}
		for ii, comment := range comment.Comments.Nodes {
			//When there is no info about the user, it means it has been deleted
			author := comment.Author.Login
			if author == "" {
				author = "deleted_user"
			}

			output_slice = append(output_slice, createRecord(prSpec, author, comment.CreatedAt))
			if isDebugGet {
				loggers.debug.Printf("  %d. %s %s \"%s\"\n", ii+1, author,
					comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.Body))
			}
			totalComments++
		}
	}

	prettyPrinted_prSpec := "\"" + prSpec + "\""
	if isRootDebug {
		if totalComments == 0 {
			loggers.debug.Printf("For %-40s no comment found. (quota cost: %d, remaining: %d)\n",
				prettyPrinted_prSpec, prQuery2.RateLimit.Cost, prQuery2.RateLimit.Remaining)
		} else {
			loggers.debug.Printf("For %-40s found %d comments. (quota cost: %d, remaining: %d)\n",
				prettyPrinted_prSpec, totalComments, prQuery2.RateLimit.Cost, prQuery2.RateLimit.Remaining)
		}
	}
	if isDebugGet {
		loggers.debug.Printf("Nbr PR Reviews: %d\n", len(prQuery2.Repository.PullRequest.Reviews.Nodes))
		loggers.debug.Printf("Grand total de reviews: %d\n", totalComments)
	}

	checkIfSufficientQuota_2(2, prQuery2.RateLimit.Remaining, prQuery2.RateLimit.Limit, prQuery2.RateLimit.ResetAt)

	return totalComments, output_slice
}

func createRecord(prSpec string, user string, date githubv4.DateTime) string {
	monthFormat := "2006-01"
	output_record := fmt.Sprintf("\"%s\",\"%s\",\"%s\"", prSpec, user, date.Format(monthFormat))
	return output_record
}
