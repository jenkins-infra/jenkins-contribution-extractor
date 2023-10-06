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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// CSV header record
var csvHeader = []string{"PR_ref", "commenter", "month"}

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
func writeCSVtoFile(out *os.File, isAppend bool, isNoHeader bool, csv_output_slice [][]string) {

	localIsNoHeader := isNoHeader

	//create a csv writer
	csv_out := csv.NewWriter(out)

	// Add the CSV header record, unless explicitly asked not to add it
	if !localIsNoHeader {
		headerWriteError := csv_out.Write(csvHeader)
		if headerWriteError != nil {
			log.Fatal(headerWriteError)
		}
		csv_out.Flush()
	}

	// write all the records in memory in one swoop
	write_err := csv_out.WriteAll(csv_output_slice)
	if write_err != nil {
		log.Fatal(write_err)
	}
	csv_out.Flush()
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
	token := os.Getenv(envVariableName)
	if token == "" {
		fmt.Println("Unauthorized: No token present")
		//This is a major error: we crash out of the program
		os.Exit(0)
	}
	return token
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

// Checks whether the input is formatted as a GitHub user or Organization
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

// Validates whether the input is correctly formatted as a GitHub user or oganisation
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
