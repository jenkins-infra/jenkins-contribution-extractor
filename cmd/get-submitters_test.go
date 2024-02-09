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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTotalNumberOfItems(t *testing.T) {
	type args struct {
		searchedOrg   string
		searchedMonth string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"Below 1K",
			args{
				searchedOrg:   "jenkinsci",
				searchedMonth: "2023-09",
			},
			692, false,
		},
		{
			"Above 1K",
			args{
				searchedOrg:   "jenkinsci",
				searchedMonth: "2020-01",
			},
			1233, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTotalNumberOfItems(tt.args.searchedOrg, tt.args.searchedMonth)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTotalNumberOfItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTotalNumberOfItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_performSearch(t *testing.T) {
	type args struct {
		searchedOrg   string
		searchedMonth string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test run for debug",
			args{
				searchedOrg:   "on4kjm",
				searchedMonth: "2020-01",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := performSearch(tt.args.searchedOrg, tt.args.searchedMonth); (err != nil) != tt.wantErr {
				t.Errorf("performSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ExecuteGetSubmitterProcessExcludeIfPresent(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"get", "submitters", "jenkins", "2024-01", "-x", "nonExistingFile.txt"})
	error := rootCmd.Execute()

	assert.Error(t, error, "Function call should have failed")

	//Error is expected
	expectedMsg := "Error: invalid excluded user list => Unable to read input file nonExistingFile.txt: open nonExistingFile.txt: no such file or directory"
	lines := strings.Split(actual.String(), "\n")
	assert.Equal(t, expectedMsg, lines[0], "Function did not fail for the expected cause")
}
