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
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		githubUser     string
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
				githubUser:     "olamy",
			},
			cleanedSubmittersList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanCsvList(tt.args.csvToCleanList, tt.args.githubUser); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanCsvList() = %v, want %v", got, tt.want)
			}
		})
	}
}
