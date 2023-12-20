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
