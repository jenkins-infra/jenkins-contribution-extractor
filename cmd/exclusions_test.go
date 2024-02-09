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
	"reflect"
	"testing"
)

func Test_load_exclusions(t *testing.T) {
	type args struct {
		exclusions_filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Empty filename provided",
			args{
				exclusions_filename: "",
			},
			nil,
			true,
		},
		{
			"filename does not exist",
			args{
				exclusions_filename: "inexistentFile.txt",
			},
			nil,
			true,
		},
		{
			"happy case",
			args{
				exclusions_filename: "../test-data/exclusions.txt",
			},
			[]string{"user1", "user2"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, got := load_exclusions(tt.args.exclusions_filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("load_exclusions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("load_exclusions() = %v, want %v", got, tt.want)
			}
		})
	}
}

var emptyStringList []string

func Test_validate_loadedFile(t *testing.T) {
	type args struct {
		loadedFile []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty list",
			args{
				loadedFile: emptyStringList,
			},
			true,
		},
		{
			"Happy Case",
			args{
				loadedFile: []string{"user1", "user2"},
			},
			false,
		},
		{
			"Happy Case - single user",
			args{
				loadedFile: []string{"user1"},
			},
			false,
		},
		{
			"space separated users on one line",
			args{
				loadedFile: []string{"user1", "user2 user3 user4"},
			},
			true,
		},
		{
			"Bad github user",
			args{
				loadedFile: []string{"user1", "user%2"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate_loadedFile(tt.args.loadedFile); (err != nil) != tt.wantErr {
				t.Errorf("validate_loadedFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_removeComments(t *testing.T) {
	type args struct {
		rawList []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no comment",
			args{rawList: []string{"user1", "user2"}},
			[]string{"user1", "user2"},
		},
		{
			"line comment",
			args{rawList: []string{"# comment", "user1", "user2"}},
			[]string{"user1", "user2"},
		},
		{
			"empty line",
			args{rawList: []string{" ", "user1", "user2"}},
			[]string{"user1", "user2"},
		},
		{
			"empty line 2",
			args{rawList: []string{"", "user1", "user2"}},
			[]string{"user1", "user2"},
		},
		{
			"inline comment 1",
			args{rawList: []string{"", "user1 #comment", "user2"}},
			[]string{"user1", "user2"},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeComments(tt.args.rawList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isCommentedLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"no comment",
			args{line: "user"},
			false,
		},
		{
			"Commented 1",
			args{line: "# this is a comment"},
			true,
		},
		{
			"Commented 2",
			args{line: " # this is a comment"},
			true,
		},
		{
			"Commented 3",
			args{line: "#this is a comment"},
			true,
		},
		{
			"Commented 4",
			args{line: "#this is #a comment"},
			true,
		},
		{
			"Empty line 1",
			args{line: " "},
			true,
		},
		{
			"Empty line 2",
			args{line: ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCommentedLine(tt.args.line); got != tt.want {
				t.Errorf("isCommentedLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
