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
	"regexp"
	"unicode"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

/*
{
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
}

func fetchComments_alt(org string, prj string, pr int) (nbrComment int, output [][]string) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"owner": githubv4.String(org),
		"name":  githubv4.String(prj),
		"pr":    githubv4.Int(pr),
	}

	err := client.Query(context.Background(), &prQuery2, variables)
	if err != nil {
		//FIXME: Better error handling
		log.Panic(err)
	}

	prSpec := fmt.Sprintf("%s/%s/%d", org, prj, pr)
	totalComments := 0
	dbgDateFormat := "2006-01-02 15:04:05"

	var output_slice [][]string

	for i, comment := range prQuery2.Repository.PullRequest.Comments.Nodes {
		//When there is no info about the user, it means it has been deleted
		author := comment.Author.Login
		if author == "" {
			author = "deleted_user"
		}

		output_slice = append(output_slice, createRecord(prSpec, author, comment.CreatedAt))
		fmt.Printf("%d. %s, %s, \"%s\"\n", i+1, author, comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.Body))
		totalComments++
	}
	fmt.Printf("Nbr PR Comments: %d\n", len(prQuery2.Repository.PullRequest.Comments.Nodes))

	for i, comment := range prQuery2.Repository.PullRequest.Reviews.Nodes {
		//When there is no info about the user, it means it has been deleted
		author := comment.Author.Login
		if author == "" {
			author = "deleted_user"
		}

		fmt.Printf("%d. %s, %s, \"%s\"\n", i+1, author, comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.BodyText))
		//Just guessing
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
			fmt.Printf("  %d. %s %s \"%s\"\n", ii+1, author, comment.CreatedAt.Format(dbgDateFormat), cleanBody(comment.Body))
			totalComments++
		}
	}
	fmt.Printf("Nbr PR Reviews: %d\n", len(prQuery2.Repository.PullRequest.Reviews.Nodes))
	fmt.Printf("Grand total de reviews: %d\n", totalComments)
	fmt.Printf("\n%v/n", output_slice)
	return totalComments, output_slice
}

// Removes and truncates a Body or BodyText element
func cleanBody(input string) (output string) {
	re := regexp.MustCompile(`\r?\n`)
	temp := re.ReplaceAllString(input, " ")

	output = truncateString(temp, 40)
	return output
}

func truncateString(input string, max int) (otput string) {
	lastSpaceIx := -1
	len := 0
	for i, r := range input {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len >= max {
			if lastSpaceIx != -1 {
				return input[:lastSpaceIx] + "..."
			}
			// If here, string is longer than max, but has no spaces
		}
	}
	// If here, string is shorter than max
	return input
}

func createRecord(prSpec string, user string, date githubv4.DateTime) []string {
	var output_record []string
	monthFormat := "2006-01"
	output_record = append(output_record, prSpec, user, date.Format(monthFormat))
	return output_record
}
