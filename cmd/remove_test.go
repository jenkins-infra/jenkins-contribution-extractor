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
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ComputeRemoveBackupFile(t *testing.T) {
	result_regexp, _ := regexp.Compile(`^(.|test_data)/removeBackup_20[0-3][0-9][0-1][0-9][0-3][0-9]_[0-2][0-9][0-5][0-9][0-5][0-9]__testFile\.csv$`)

	backupFileName := compute_removeBackupFileName("testFile.csv")
	assert.True(t, result_regexp.MatchString(backupFileName), "Backup file name (%s) doesn't have the expected format", backupFileName)

	backupFileName = compute_removeBackupFileName("test_data/testFile.csv")
	assert.True(t, result_regexp.MatchString(backupFileName), "Backup file name (%s) doesn't have the expected format", backupFileName)
}

func Test_ExecuteMustHaveTwoArguments(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"remove", "A"})
	error := rootCmd.Execute()

	assert.Error(t, error, "Function call should have failed")

	//Error is expected
	expectedMsg := "Error: requires at least 2 arg(s), only received 1"
	lines := strings.Split(actual.String(), "\n")
	assert.Equal(t, expectedMsg, lines[0], "Function did not fail for the expected cause")
}

// This is an end to end test
func Test_ExecuteIntegrationTest(t *testing.T) {

	// Setup environment
	tempDir := t.TempDir()
	tempFileName, err := duplicateFile("../test-data/submissions-2023-08.csv", tempDir, true)

	assert.NoError(t, err, "Unexpected File duplication error")
	assert.NotEmpty(t, tempFileName, "Unexpected empty temporary filename")

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	var commandArguments []string
	commandArguments = append(commandArguments, "remove", "File:../test-data/test-exclusion.txt", tempFileName, "-b")
	rootCmd.SetArgs(commandArguments)
	error := rootCmd.Execute()

	assert.NoError(t, error, "Function should not have failed")

	// Compare output with reference (golden) file
	goldenFileName := "../test-data/submissions-2023-08_cleaned.csv"
	assert.True(t, isFileEquivalent(tempFileName, goldenFileName))

	//Does the backup file exist?
	backupFileName := compute_removeBackupFileName(tempFileName)
	assert.FileExistsf(t, backupFileName, "Backup file (%s) has not been created.", backupFileName)
}

func Test_performRemove(t *testing.T) {
	type args struct {
		githubUser       string
		fileToClean_name string
		isBackup         bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Invalid GitHub user",
			args{
				githubUser:       "ax_4!",
				fileToClean_name: "",
				isBackup:         true,
			},
			true,
		},
		{
			"Non existent file",
			args{
				githubUser:       "jenkinsci",
				fileToClean_name: "unexistantFile.txt",
				isBackup:         true,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := performRemove(tt.args.githubUser, tt.args.fileToClean_name, tt.args.isBackup); (err != nil) != tt.wantErr {
				t.Errorf("performRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var expectedSubmittersList = []string{
	"org,repository,number,url,state,created_at,merged_at,user.login,month_year,title",
	"\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
	"\"jenkinsci\",\"ldap-plugin\",248,\"https://github.com/jenkinsci/ldap-plugin/pull/248\",\"closed\",\"2023-08-12T12:09:11Z\",\"2023-09-22T16:21:31Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
	"\"jenkinsci\",\"ecu-test-execution-plugin\",54,\"https://github.com/jenkinsci/ecu-test-execution-plugin/pull/54\",\"closed\",\"2023-08-07T10:06:24Z\",\"2023-09-22T09:03:34Z\",\"MxEh-TT\",\"2023-08\",\"inital package check implementation (#53)\"",
	"\"jenkinsci\",\"build-blocker-plugin\",19,\"https://github.com/jenkinsci/build-blocker-plugin/pull/19\",\"closed\",\"2023-08-07T06:35:02Z\",\"2023-09-18T13:42:06Z\",\"olamy\",\"2023-08\",\"add @Symbol to be able to easily use the plugin in a declarative pipeline\"",
	"\"jenkinsci\",\"credentials-plugin\",475,\"https://github.com/jenkinsci/credentials-plugin/pull/475\",\"closed\",\"2023-08-12T08:16:01Z\",\"2023-09-21T16:16:52Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
	"\"jenkinsci\",\"ssh-credentials-plugin\",179,\"https://github.com/jenkinsci/ssh-credentials-plugin/pull/179\",\"closed\",\"2023-08-12T08:32:14Z\",\"2023-09-21T16:12:07Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
}

// Removed "olamy"
var cleanedSubmittersList = []string{
	"org,repository,number,url,state,created_at,merged_at,user.login,month_year,title",
	"\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
	"\"jenkinsci\",\"ldap-plugin\",248,\"https://github.com/jenkinsci/ldap-plugin/pull/248\",\"closed\",\"2023-08-12T12:09:11Z\",\"2023-09-22T16:21:31Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
	"\"jenkinsci\",\"ecu-test-execution-plugin\",54,\"https://github.com/jenkinsci/ecu-test-execution-plugin/pull/54\",\"closed\",\"2023-08-07T10:06:24Z\",\"2023-09-22T09:03:34Z\",\"MxEh-TT\",\"2023-08\",\"inital package check implementation (#53)\"",
	"\"jenkinsci\",\"credentials-plugin\",475,\"https://github.com/jenkinsci/credentials-plugin/pull/475\",\"closed\",\"2023-08-12T08:16:01Z\",\"2023-09-21T16:16:52Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
	"\"jenkinsci\",\"ssh-credentials-plugin\",179,\"https://github.com/jenkinsci/ssh-credentials-plugin/pull/179\",\"closed\",\"2023-08-12T08:32:14Z\",\"2023-09-21T16:12:07Z\",\"NotMyFault\",\"2023-08\",\"Test on Java 21\"",
}

func Test_loadCSVtoClean(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Happy Case",
			args{fileName: "../test-data/small-submission-list.csv"},
			expectedSubmittersList,
			false,
		},
		{
			"Empty File",
			args{fileName: "../test-data/empty-submission-list.csv"},
			nil,
			true,
		},
		{
			"Inexistent File",
			args{fileName: "blaahhh.csv"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, got := loadCSVtoClean(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadCSVtoClean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadCSVtoClean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanCsvList(t *testing.T) {
	type args struct {
		csvToCleanList []string
		githubUserList []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"happy case",
			args{
				csvToCleanList: expectedSubmittersList,
				githubUserList: []string{"olamy"},
			},
			cleanedSubmittersList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanCsvList(tt.args.csvToCleanList, tt.args.githubUserList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanCsvList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listItemContainedInLine(t *testing.T) {
	type args struct {
		line     string
		userList []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"detected user from single user list",
			args{
				line:     "\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
				userList: []string{"MarkEWaite"},
			},
			true,
		},
		{
			"detected user from multi user list",
			args{
				line:     "\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
				userList: []string{"user1", "MarkEWaite"},
			},
			true,
		},
		{
			"undetected user from multi user list",
			args{
				line:     "\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
				userList: []string{"user1", "oLamy"},
			},
			false,
		},
		{
			"empty user list",
			args{
				line:     "\"jenkinsci\",\"embeddable-build-status-plugin\",229,\"https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229\",\"closed\",\"2023-08-11T21:18:19Z\",\"2023-08-12T03:55:01Z\",\"MarkEWaite\",\"2023-08\",\"Test with Java 21\"",
				userList: []string{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listItemContainedInLine(tt.args.line, tt.args.userList); got != tt.want {
				t.Errorf("listItemContainedInLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFileSpec(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"not a file",
			args{input: "gitHubUser"},
			"",
		},
		{
			"with filespec",
			args{input: "file:thisIsAFile.txt"},
			"thisIsAFile.txt",
		},
		{
			"with mixed case filespec",
			args{input: "File:thisIsAFile.txt"},
			"thisIsAFile.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFileSpec(tt.args.input); got != tt.want {
				t.Errorf("isFileSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ------------------------------
//
// Test Utilities
//
// ------------------------------

// duplicate test file as a temporary file.
// The temporary directory should be created in the calling test so that it gets cleaned at test completion.
func duplicateFile(originalFileName, targetDir string, generateFilename bool) (tempFileName string, err error) {

	//Check the status and size of the original file
	sourceFileStat, err := os.Stat(originalFileName)
	if err != nil {
		return "", err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", originalFileName)
	}
	sourceFileSize := sourceFileStat.Size()

	//Open the original file
	source, err := os.Open(originalFileName)
	if err != nil {
		return "", err
	}
	defer source.Close()

	// generate temporary file name in temp directory if requested
	if generateFilename {
	file, err := os.CreateTemp(targetDir, "testData.*.csv")
	if err != nil {
		return "", err
	}
	tempFileName = file.Name()
} else {
	//we want to keep the original filename
	_, file := filepath.Split(originalFileName)
	tempFileName = filepath.Join(targetDir,file)
}

	// create the new file duplication
	destination, err := os.Create(tempFileName)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	// Do the actual copy
	bytesCopied, err := io.Copy(destination, source)
	if err != nil {
		return tempFileName, err
	}
	if bytesCopied != sourceFileSize {
		return tempFileName, fmt.Errorf("Source and destination file size do not match after copy (%s is %d bytes and %s is %d bytes", originalFileName, sourceFileSize, tempFileName, bytesCopied)
	}

	// All went well
	return tempFileName, nil
}

func isFileEquivalent(tempFileName, goldenFileName string) bool {

	// Is the size the same
	tempFileSize := getFileSize(tempFileName)
	goldenFileSize := getFileSize(goldenFileName)

	if tempFileSize == 0 || goldenFileSize == 0 {
		fmt.Printf("0 byte file length\n")
		return false
	}

	if tempFileSize != goldenFileSize {
		fmt.Printf("Files are of different sizes: found %d bytes while expecting reference %d bytes \n", tempFileSize, goldenFileSize)
		return false
	}

	// load both files
	err, tempFile_List := loadCSVtoClean(tempFileName)
	if err != nil {
		fmt.Printf("Unexpected error loading %s : %v \n", tempFileName, err)
		return false
	}

	err, goldenFile_List := loadCSVtoClean(goldenFileName)
	if err != nil {
		fmt.Printf("Unexpected error loading %s : %v \n", goldenFileName, err)
		return false
	}

	//Compare the two lists
	for index, line := range tempFile_List {
		if line != goldenFile_List[index] {
			fmt.Printf("Compare failure: line %d do not match\n", index)
			return false
		}
	}

	//If we reached this, we are all good
	return true
}

// Gets the size of a file
func getFileSize(fileName string) int64 {
	tempFileStat, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("Unexpected error getting details of %s: %v\n", fileName, err)
		return 0
	}
	if !tempFileStat.Mode().IsRegular() {
		fmt.Printf("%s is not a regular file\n", fileName)
		return 0
	}
	return tempFileStat.Size()
}
