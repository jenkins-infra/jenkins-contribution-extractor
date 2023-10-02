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

	"github.com/google/go-github/v55/github"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// quotaCmd represents the quota command
var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Gets the current GitHub API quota status",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		get_quota()
	},
}

func init() {
	rootCmd.AddCommand(quotaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quotaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quotaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ---
// Retrieves the GitHub API Quota
func get_quota() {
	limit, remaining := get_quota_data()
	fmt.Printf("V3 Limit: %d \nV3 Remaining %d \n\n", limit, remaining)

	limit_v4, remaining_v4, reset_time := get_quota_data_v4()
	resetTimeString := reset_time.Format(time.RFC1123)
	now := time.Now()
	diff := reset_time.Sub(now)
	secondsToGo := diff.Seconds()
	

	fmt.Printf("V4 Limit: %d \nV4 Remaining: %d \nV4 Reset time: %s (in %.0f secs)\n", limit_v4, remaining_v4,resetTimeString,secondsToGo)
}

// Retrieves the GitHub Quota.
func get_quota_data() (limit int, remaining int) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)

	client := github.NewClient(nil).WithAuthToken(ghToken)

	limitsData, _, err := client.RateLimits(context.Background())
	if err != nil {
		log.Printf("Error getting limit: %v", err)
		return 0, 0
	}
	return limitsData.Core.Limit, limitsData.Core.Remaining
}

/*
query {
  viewer {
    login
  }
  rateLimit {
    limit
    cost
    remaining
    resetAt
  }
}
*/

var quotaQuery struct {
	Viewer struct {
		Login string
	}
	RateLimit struct {
		Limit     int
		Cost      int
		Remaining int
		ResetAt   time.Time
	}
}

func get_quota_data_v4() (limit int, remaining int, resetAt time.Time) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	err := client.Query(context.Background(), &quotaQuery, nil)
	if err != nil {
		//FIXME: Better error handling
		log.Panic(err)
	}

	// fmt.Printf("User: %s\n", quotaQuery.Viewer.Login)
	// fmt.Printf("limit:     %d\n", quotaQuery.RateLimit.Limit)
	// fmt.Printf("cost:      %d\n", quotaQuery.RateLimit.Cost)
	// fmt.Printf("remaining: %d\n", quotaQuery.RateLimit.Remaining)
	// fmt.Printf("resetAt:   %s\n", quotaQuery.RateLimit.ResetAt)
	return quotaQuery.RateLimit.Limit, quotaQuery.RateLimit.Remaining, quotaQuery.RateLimit.ResetAt
}
