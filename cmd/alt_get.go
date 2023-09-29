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

type comment struct {
	Author struct {
		Login string
	}
	CreatedAt githubv4.DateTime
}

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
						Login githubv4.String
					}
				}
			} `graphql:"comments(first: 100)"`
			Reviews struct {
				Nodes []struct {
					CreatedAt githubv4.DateTime
					BodyText  string
					Author    struct {
						Login githubv4.String
					}
					Comments struct {
						Nodes []struct {
							CreatedAt githubv4.DateTime
							Body      string
							Author    struct {
								Login githubv4.String
							}
						}
					} `graphql:"comments(first: 100)"`
				}
			} `graphql:"reviews(first: 100)"`
		} `graphql:"pullRequest(number: $pr)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func fetchComments_alt(org string, prj string, pr int) {
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
		log.Panic(err)
	}
	fmt.Println(prQuery2.Repository.Description)
	fmt.Println(prQuery2.Repository.PullRequest.Title)

	totalComments := 0

	for i, comment := range prQuery2.Repository.PullRequest.Comments.Nodes {
		fmt.Printf("%d. %s %s %s\n", i+1, comment.Author.Login, comment.CreatedAt, comment.Body)
		totalComments++
	}
	fmt.Printf("Nbr PR Comments: %d\n", len(prQuery2.Repository.PullRequest.Comments.Nodes))

	for i, comment := range prQuery2.Repository.PullRequest.Reviews.Nodes {
		fmt.Printf("%d. %s %s\n", i+1, comment.Author.Login, comment.CreatedAt)
		if len(comment.Comments.Nodes) > 1 {
			totalComments++
		}
		for ii, comment := range comment.Comments.Nodes {
			// if ii == 0 {
			// 	continue
			// }
			fmt.Printf("  %d. %s %s\n", ii+1, comment.Author.Login, comment.CreatedAt)
			totalComments++
		}
	}
	fmt.Printf("Nbr PR Reviews: %d\n", len(prQuery2.Repository.PullRequest.Reviews.Nodes))
	fmt.Printf("Grand total de reviews: %d\n", totalComments)
}
