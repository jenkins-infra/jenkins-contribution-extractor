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
	"bufio"
	"fmt"
	"os"
)

// Loads the list of gitHub users to exclude from the count
func load_exclusions(exclusions_filename string) (error, []string) {

	if len(exclusions_filename) == 0 {
		return fmt.Errorf("No filename provided."), nil
	}

	f, err := os.Open(exclusions_filename)
	if err != nil {
		return fmt.Errorf("Unable to read input file %s: %v\n", exclusions_filename, err), nil
	}
	defer f.Close()

	var loadedFile []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		loadedFile = append(loadedFile, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error loading \"%s\": %v", exclusions_filename, err), nil
	}

	if validationError := validate_loadedFile(loadedFile); validationError != nil {
		return validationError, nil
	} else {
		return nil, loadedFile
	}
}

// Validates whether the supplied string slice is composed of properly formatted GitHub users
func validate_loadedFile(loadedFile []string) error {
	if len(loadedFile) == 0 {
		return fmt.Errorf("Error: empty file")
	}

	for _, githubUserToCheck := range loadedFile {
		if !isValidOrgFormat(githubUserToCheck) {
			return fmt.Errorf("Invalid excluded user \"%s\" (does not match GitHub user syntax)", githubUserToCheck)
		}
	}

	return nil
}
