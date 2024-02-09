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

func Test_loadPrListFile(t *testing.T) {
	type args struct {
		fileName  string
		isVerbose bool
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			"test empty file",
			args{
				fileName:  "../test-data/empty-submission-list.csv",
				isVerbose: true,
			},
			nil, true,
		},
		// {
		// 	"test file with one data line",
		// 	args{
		// 		fileName: "../test-data/oneLine-submission-list.csv",
		// 		isVerbose: true,
		// 	},
		// 	nil,true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := loadPrListFile(tt.args.fileName, tt.args.isVerbose)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadPrListFile() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("loadPrListFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_ExecuteGetCommenterProcessExcludeIfPresent(t *testing.T) {
	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"get", "commenters", "../test-data/empty-submission-list.csv", "-x", "nonExistingFile.txt"})
	error := rootCmd.Execute()

	assert.Error(t, error, "Function call should have failed")

	//Error is expected
	expectedMsg := "Error: invalid excluded user list => Unable to read input file nonExistingFile.txt: open nonExistingFile.txt: no such file or directory"
	lines := strings.Split(actual.String(), "\n")
	assert.Equal(t, expectedMsg, lines[0], "Function did not fail for the expected cause")
}
