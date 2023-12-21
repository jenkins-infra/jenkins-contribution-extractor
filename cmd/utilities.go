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
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// validates that the supplied string is a valid PR specification
// in the form of "org/project/pr_nbr"
func validatePRspec(prSpec string) (org string, project string, prNbr int, err error) {
	splittedString := strings.Split(strings.TrimSpace(prSpec), "/")

	if len(splittedString) != 3 {
		return "", "", -1, fmt.Errorf("Invalid number of elements in PR Specification (\"org/project/pr\"). (expecting 3, found %v)\n", len(splittedString))
	}

	work_Org := splittedString[0]
	work_Project := splittedString[1]
	prString := splittedString[2]

	if strings.TrimSpace(work_Org) == "" {
		return "", "", -1, fmt.Errorf("Organization element in PR Specification is empty\n")
	}
	if strings.TrimSpace(work_Project) == "" {
		return "", "", -1, fmt.Errorf("Project element in PR Specification is empty\n")
	}
	if strings.TrimSpace(prString) == "" {
		return "", "", -1, fmt.Errorf("PR element in PR Specification is empty\n")
	}

	work_prNbr, err := strconv.Atoi(strings.TrimSpace(prString))
	if err != nil {
		return "", "", -1, fmt.Errorf("PR part of PR Specification is not numerical (%v)\n", err)
	}
	return work_Org, work_Project, work_prNbr, nil
}

// Write the string slice to a file formatted as a CSV
func writeCSVtoFile(out *os.File, isAppend bool, isNoHeader bool, header string, csv_output_slice []string) {

	localIsNoHeader := isNoHeader

	datawriter := bufio.NewWriter(out)

	// Add the CSV header record, unless explicitly asked not to add it
	if !localIsNoHeader {
		_, headerWriteError := datawriter.WriteString(header + "\n")
		if headerWriteError != nil {
			log.Fatal(headerWriteError)
		}
		datawriter.Flush()
	}

	// write all the records in memory
	for _, data := range csv_output_slice {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
}

// creates or opens for append (if the file exists) the output file
// If no append is requested and the file exists, it is overwritten
func openOutputCSV(outFname string, isAppend bool, isNoHeader bool) (*os.File, bool) {

	isExisting := fileExist(outputFileName)
	localIsNoHeader := isNoHeader

	var isAppendString string
	isNoHeaderString := "without"
	if !localIsNoHeader {
		isNoHeaderString = "with"
	}

	var out *os.File
	var open_error error

	if isExisting {
		if isAppend {
			// Open for append
			out, open_error = os.OpenFile(outFname, os.O_APPEND|os.O_WRONLY, 0644)
			if open_error != nil {
				log.Fatal(open_error)
			}

			isAppendString = "(appending"
			// no Header forced
			isNoHeaderString = "without"
			localIsNoHeader = true
		} else {
			// overwrite output file
			out, open_error = os.Create(outFname)
			if open_error != nil {
				log.Fatal(open_error)
			}
			isAppendString = "(overwriting"
			// honor the noheader setting
		}
	} else {
		//create output file
		out, open_error = os.Create(outFname)
		if open_error != nil {
			log.Fatal(open_error)
		}
		isAppendString = "(creating"
		// honor noHeader setting
	}

	if isVerbose {
		fmt.Printf("Writing data to \"%s\" %s %s header)\n", outputFileName, isAppendString, isNoHeaderString)
	}

	return out, localIsNoHeader
}

// Validates that the input file is a real file (and not a directory)
func fileExist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Load the GitHub token from the specified environment variable
func loadGitHubToken(envVariableName string) string {
	token, found := os.LookupEnv(envVariableName)
	if !found {
		fmt.Println("Unauthorized: No token present")
		//This is a major error: we crash out of the program
		log.Fatal("GitHub token not found!")
	}
	return token
}

// Removes and truncates a Body or BodyText element
func cleanBody(input string) (output string) {
	re := regexp.MustCompile(`\r?\n`)
	temp := re.ReplaceAllString(input, " ")

	re2 := regexp.MustCompile(`\"`)
	temp2 := re2.ReplaceAllString(temp, "'")

	output = truncateString(temp2, 40)
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

// Checks whether the input is a month in the expected format
func isValidMonthFormat(input string) bool {
	if input == "" {
		if isVerbose {
			fmt.Print("Empty month\n")
		}
		return false
	}

	regexpMonth := regexp.MustCompile(`^20[12][0-9]-(0[1-9]|1[0-2])$`)
	if !regexpMonth.MatchString(input) {
		if isVerbose {
			fmt.Printf("Supplied data (%s) is not in a valid month format. Should be \"YYYY-MM\" and later than 2010\n", input)
		}
		return false
	}
	return true
}

// Validates whether the input is correctly formatted as a GitHub user or organisation
func isValidOrgFormat(input string) bool {
	if input == "" {
		if isVerbose {
			fmt.Print("Empty Org\n")
		}
		return false
	}

	//The GitHub user validation regexp (see https://stackoverflow.com/questions/58726546/github-username-convention-using-regex)
	// should be regexp.Compile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`). But the dataset contains "invalid" data: username ending with a "-" or
	// a double "-" in the name.
	name_regexp := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	if !name_regexp.MatchString(input) {
		if isVerbose {
			fmt.Printf("Supplied data (%s) is not in a valid GitHub user/org format.\n", input)
		}
		return false
	}

	return true
}

// checks whether the user is an application based on the URL
func isUserBot(url string) bool {
	if strings.HasPrefix(strings.ToLower(url), "https://github.com/apps/") {
		return true
	} else {
		return false
	}
}

// Computes the start and end date based on the total number of issues returned by query
func splitPeriodForMaxQueryItem(totalNbrIssue int, shortMonth string, requestedIteration int) (startDate string, endDate string, moreIteration bool) {
	queryLimit := 1000 // constant

	// Parameters validation
	if requestedIteration < 0 {
		if isRootDebug {
			loggers.debug.Printf("Error: requested iteration (%d) is negative\n", requestedIteration)
		}
		return "", "", false
	}

	if totalNbrIssue < 0 {
		if isRootDebug {
			loggers.debug.Printf("Error: total number of issues (%d) is negative\n", totalNbrIssue)
		}
		return "", "", false
	}

	if totalNbrIssue > (28 * queryLimit) {
		if isRootDebug {
			loggers.debug.Printf("Error: requested iteration (%d) is greater than what we can handle (28 * %d)\n", requestedIteration, queryLimit)
		}
		return "", "", false
	}

	//load short month in a time structure and implicitly validate it
	inputDate, err := time.Parse("2006-01", shortMonth)
	if err != nil {
		if isRootDebug {
			loggers.debug.Printf("Unexpected error parsing short month (%v)\n", err)
		}
		return "", "", false
	}

	// *** Let's go ****

	//retrieve the year and month in time structure
	inputYear, inputMonth, _ := inputDate.Date()

	// Get the first and last day of the month we are looking at
	currentLocation := inputDate.Location()
	firstOfMonth := time.Date(inputYear, inputMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	if totalNbrIssue < queryLimit {
		// we convert the time structs to strings
		startDate = firstOfMonth.Format("2006-01-02")
		endDate = lastOfMonth.Format("2006-01-02")
		moreIteration = false

		//There must be a single iteration
		if requestedIteration != 0 {
			if isRootDebug {
				loggers.debug.Printf("Unexpected error: iteration number (%d) should be equal to 0\n", requestedIteration)
			}
			return "", "", false
		}
	} else {
		// compute the total number of iteration required
		totalIterations := int(totalNbrIssue/queryLimit) + 1
		if requestedIteration > totalIterations {
			if isRootDebug {
				loggers.debug.Printf("Error: requested iteration (%d) is greater than the total number of iteration (%d)\n", requestedIteration, totalIterations)
			}
			return "", "", false
		}
		numberOfDaysInMonth := lastOfMonth.Day()
		daysPerIterations := int(numberOfDaysInMonth / totalIterations)

		//compute the iteration start date
		iterationStartDay := (daysPerIterations * requestedIteration) + 1
		startOfIterationDate := time.Date(inputYear, inputMonth, iterationStartDay, 0, 0, 0, 0, currentLocation)
		startDate = startOfIterationDate.Format("2006-01-02")

		//compute the iteration end date
		iterationEndDay := daysPerIterations + (daysPerIterations * requestedIteration)
		endOfIterationDate := time.Date(inputYear, inputMonth, iterationEndDay, 0, 0, 0, 0, currentLocation)
		endDate = endOfIterationDate.Format("2006-01-02")

		//did we reach the last iteration?
		if (requestedIteration + 1) == totalIterations {
			moreIteration = false
			// in the last iteration, we catch up any rounding errors by forcing the months's last day
			endOfIterationDate := time.Date(inputYear, inputMonth, numberOfDaysInMonth, 0, 0, 0, 0, currentLocation)
			endDate = endOfIterationDate.Format("2006-01-02")
		} else {
			moreIteration = true
		}

	}

	if isRootDebug {
		loggers.debug.Printf("Iteration start: %s end: %s has more iterations: %v\n", startDate, endDate, moreIteration)
	}
	return startDate, endDate, moreIteration
}

// returns the start and end day for a given month (YYYY-MM)
func getStartAndEndOfMonth(shortMonth string) (startDate string, endDate string) {
	//load short month in a time structure
	inputDate, _ := time.Parse("2006-01", shortMonth)

	//retrieve the year and month in time structure
	inputYear, inputMonth, _ := inputDate.Date()

	//Build the dates we are looking for
	currentLocation := inputDate.Location()
	firstOfMonth := time.Date(inputYear, inputMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	//convert first and last days into a string in the expected format
	firstOfMonthString := firstOfMonth.Format("2006-01-02")
	lastOfMonthString := lastOfMonth.Format("2006-01-02")

	return firstOfMonthString, lastOfMonthString
}

// TODO: test this
// Checks whether the retrieved header is equivalent to the reference header
func validateHeader(header []string, referenceHeader []string, isVerbose bool) bool {
	if len(header) != len(referenceHeader) {
		if isVerbose {
			fmt.Printf(" Error: field number mismatch (found %d, wanted %d)\n", len(header), len(referenceHeader))
		}
		return false
	}
	for i, v := range header {
		if v != referenceHeader[i] {
			if isVerbose {
				fmt.Printf(" Error: not the expected header field at column %d (found \"%v\", wanted \"%v\")\n", i+1, v, referenceHeader[i])
			}
			return false
		}
	}
	return true
}
