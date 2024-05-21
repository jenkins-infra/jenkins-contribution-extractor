/*
Copyright © 2024 Jean-Marc Meessen jean-marc@meessen-web.org

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
	"encoding/csv"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var honorDataDir string
var honorOutput string

// honorCmd represents the honor command
var honorCmd = &cobra.Command{
	Use:   "honor <month>",
	Short: "Gets a contributor to honor",
	Long: `A command to get a random submitter from a given month and
format his data in such a way that it can be used to format an honoring
message at the bottom of the https://contributors.jenkins.io/ page.

\"month\" is a required parameter. It is in YYYY-MM format.`,
	Args: func(cmd *cobra.Command, args []string) error {
		//call requires two parameters (org and month)
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			if err.Error() == "requires at least 1 arg(s), only received 0" {
				return fmt.Errorf("\"month\" argument is missing.")
			} else {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return performHonorContributorSelection(honorDataDir, honorOutput, args[0])
	},
}

// Initialize command parameters and defaults
func init() {
	rootCmd.AddCommand(honorCmd)
	honorCmd.Flags().StringVarP(&honorDataDir, "data_dir", "", "data", "Directory containing the data to be read")
	honorCmd.Flags().StringVarP(&honorOutput, "output", "", "", "File to output the data to (default: \"[data_dir]/honored_contributor.csv\")")
}

// Command processing entry point
func performHonorContributorSelection(dataDir string, suppliedOutputFileName string, monthToSelectFrom string) error {
	// validate the month
	if !isValidMonthFormat(monthToSelectFrom) {
		return fmt.Errorf("\"%s\" is not a valid month.", monthToSelectFrom)
	}

	// does the dataDir exist ?
	if !isValidDir(dataDir) {
		return fmt.Errorf("Supplied DataDir \"%s\" does not exist.", dataDir)
	}

	// if output is not defined, build it
	honorOutputFileName := ""
	if suppliedOutputFileName == "" {
		honorOutputFileName = filepath.Join(dataDir, "honored_contributor.csv")
	} else {
		honorOutputFileName = suppliedOutputFileName
	}
	if isVerbose {
		fmt.Println("Output file: " + honorOutputFileName + "\n")
	}

	//compute the correct input filename (pr_per_submitter-YYYY-MM.csv)
	inputFileName := filepath.Join(dataDir, "pr_per_submitter-"+monthToSelectFrom+".csv")

	// fail if the file does not exist else open the file
	f, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("Unable to read input file "+inputFileName+"\n", err)
	}
	defer f.Close()

	// validate that it has the correct format (CSV and column names)
	r := csv.NewReader(f)

	headerLine, err1 := r.Read()
	if err1 != nil {
		return fmt.Errorf("Unexpected error loading"+inputFileName+"\n", err)
	}

	if isVerbose {
		fmt.Println("Checking input file " + inputFileName)
	}

	referencePrPerSubmitterHeader := []string{"user", "PR"}
	if !validateHeader(headerLine, referencePrPerSubmitterHeader, isVerbose) {
		return fmt.Errorf(" Error: header is incorrect.")
	} else {
		if isVerbose {
			fmt.Printf("  - Header is correct\n")
		}
	}

	// load the file in memory
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("Unexpected error loading \""+inputFileName+"\"\n", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("Error: No data available after the header\n")
	}
	if isVerbose {
		fmt.Println("  - At least one Submitter data available")
	}

	// pick a data line randomly
	nbrOfRecordsLoaded := len(records) - 1

	randomRecordNumber := rand.IntN(nbrOfRecordsLoaded)
	submittersName := records[randomRecordNumber][0]
	submittersPRs := records[randomRecordNumber][1]
	// fmt.Printf("[%d] - %s - %s PRs\n", randomRecordNumber, records[randomRecordNumber][0], records[randomRecordNumber][1])
	if isVerbose {
		fmt.Printf("  - Picked record %d : %s - %s PRs\n", randomRecordNumber, submittersName, submittersPRs)
	}

	// make a GitHub query to retrieve the contributors information (URL, avatar) and PRs
	if err := getSubmittersPRfromGH(submittersName, submittersPRs, monthToSelectFrom); err != nil {
		return err
	}

	// TODO: format the output with the gathered data
	// TODO: output the file

	return nil
}

var uniqueRepoSet = make(map[string]bool)
var uniqueRepoSlice = []string{}

// func main() {
//     add("aaa")
//     add("bbb")
//     add("bbb")
//     add("ccc")
// }

// Adds an item to the slice only if it is not there yet. See https://stackoverflow.com/questions/33207197/how-can-i-create-an-array-that-contains-unique-strings
func addUniqueItem(s string) {
    if uniqueRepoSet[s] {
        return // Already in the map
    }
    uniqueRepoSlice = append(uniqueRepoSlice, s)
    uniqueRepoSet[s] = true
}

/*****
 ***** Github query definition
 *****/

//GitHub Graphql query. Test at https://docs.github.com/en/graphql/overview/explorer
/*
{
	user(login: "basil"){
    	name
    	company
    	avatarUrl
    	url
  }
	search(query: "org:jenkinsci org:jenkins-infra is:pr author:dduportal created:2024-04-01..2024-04-30", type: ISSUE, first: 100) {
    issueCount
    edges {
      node {
        ... on PullRequest {
          author {
            login
            avatarUrl
            url
          }
          url
          title
          createdAt
          repository {
            name
          }
        }
      }
    }
  }
}
*/

//******************************

// Gets all the PRs in the given month for the submitters
func getSubmittersPRfromGH(submittersName string, submittersPRs string, monthToSelectFrom string) error {

	// Setup the GH query client
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	startDate, endDate := getStartAndEndOfMonth(monthToSelectFrom)

	var prQuery3 struct {
		Search struct {
			IssueCount int
			Edges      []struct {
				Node struct {
					PullRequest struct {
						Url        string
						Title      string
						CreatedAt  time.Time
						Repository struct {
							Name  string
							Owner struct {
								Login string
							}
						}
						Author struct {
							Login        string
							AvatarUrl    string
							Url          string
							ResourcePath string
						}
					} `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(first: $count, query: $searchQuery, type: ISSUE)"`
	}

	variables := map[string]interface{}{
		"searchQuery": githubv4.String(
			fmt.Sprintf(`org:%s org:%s is:pr author:%s created:%s..%s`,
				githubv4.String("jenkinsci"),
				githubv4.String("jenkins-infra"),
				// githubv4.String("dduportal"),
				githubv4.String(submittersName),
				githubv4.String(startDate),
				githubv4.String(endDate),
			),
		),
		"count": githubv4.Int(100),
	}

	if err := client.Query(context.Background(), &prQuery3, variables); err != nil {
		return fmt.Errorf("Error performing query: %v\n", err)
	}

	totalPRs := prQuery3.Search.IssueCount
	//FIXME: check if the count equals the one passed to func

	fmt.Printf("PRs found: %d\n", totalPRs)

	for ii, singlePr := range prQuery3.Search.Edges {
		foundAuthor := singlePr.Node.PullRequest.Author.Login
		authorURL := singlePr.Node.PullRequest.Author.Url
		authorAvatarUrl := singlePr.Node.PullRequest.Author.AvatarUrl
		repositoryName := singlePr.Node.PullRequest.Repository.Owner.Login + "/" + singlePr.Node.PullRequest.Repository.Name


		if ii == 0 {
			fmt.Printf("Author: %s\n", foundAuthor)
			fmt.Printf("URL:    %s\n", authorURL)
			fmt.Printf("Avatar: %s\n", authorAvatarUrl)
		}
		// fmt.Printf("%s ", repositoryName)
		addUniqueItem(repositoryName)
	}
	fmt.Print("\nrepos: ")
	fmt.Println(uniqueRepoSlice)

	return nil
}
