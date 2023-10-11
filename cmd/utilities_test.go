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
	"testing"
)

func Test_validatePRspec(t *testing.T) {
	type args struct {
		prSpec string
	}
	tests := []struct {
		name        string
		args        args
		wantOrg     string
		wantProject string
		wantPrNbr   int
		wantErr     bool
	}{
		{
			"happy case",
			args{prSpec: "on4kjm/FLEcli/1"},
			"on4kjm", "FLEcli", 1, false,
		},
		{
			"non numeric PR",
			args{prSpec: "on4kjm/FLEcli/aa"},
			"", "", -1, true,
		},
		{
			"empty first field",
			args{prSpec: "/FLEcli/1"},
			"", "", -1, true,
		},
		{
			"empty second field",
			args{prSpec: "on4kjm//1"},
			"", "", -1, true,
		},
		{
			"empty third field",
			args{prSpec: "on4kjm/FLEcli/"},
			"", "", -1, true,
		},
		{
			"too short",
			args{prSpec: "on4kjm/FLEcli"},
			"", "", -1, true,
		},
		{
			"too long",
			args{prSpec: "on4kjm/FLEcli/1/zzz"},
			"", "", -1, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrg, gotProject, gotPrNbr, err := validatePRspec(tt.args.prSpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePRspec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOrg != tt.wantOrg {
				t.Errorf("validatePRspec() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
			if gotProject != tt.wantProject {
				t.Errorf("validatePRspec() gotProject = %v, want %v", gotProject, tt.wantProject)
			}
			if gotPrNbr != tt.wantPrNbr {
				t.Errorf("validatePRspec() gotPrNbr = %v, want %v", gotPrNbr, tt.wantPrNbr)
			}
		})
	}
}

func Test_fileExist(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Happy case",
			args{"../test-data/big-submission-list.csv"},
			true,
		},
		{
			"File does not exist",
			args{"unexistantFile.txt"},
			false,
		},
		{
			"File is a directory in fact",
			args{"../test-data"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExist(tt.args.fileName); got != tt.want {
				t.Errorf("fileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanBody(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{
			"Happy case",
			args{input: "aa aaa \nbbb bbb\n"},
			"aa aaa  bbb bbb ",
		},
		{
			"Empty string",
			args{input: ""},
			"",
		},
		{
			"No return",
			args{input: "aaaa bbbb ccc"},
			"aaaa bbbb ccc",
		},
		{
			"Truncate string",
			args{input: "aaaa bbbb cccc dddd eeee ffff gggg hhhh iiii jjjj"},
			"aaaa bbbb cccc dddd eeee ffff gggg hhhh...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := cleanBody(tt.args.input); gotOutput != tt.wantOutput {
				t.Errorf("cleanBody() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

// func Test_isValidMonthFormat(t *testing.T) {
// 	type args struct {
// 		input string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			"Happy case",
// 			args{input: "2023-09"},
// 			true,
// 		},
// 		{
// 			"just junk",
// 			args{input: "junk"},
// 			false,
// 		},
// 		{
// 			"space",
// 			args{input: " "},
// 			false,
// 		},
// 		{
// 			"empty",
// 			args{input: ""},
// 			false,
// 		},
// 		{
// 			"invalid month",
// 			args{input: "2023-16"},
// 			false,
// 		},
// 		{
// 			"invalid year",
// 			args{input: "1515-06"},
// 			false,
// 		},
// 		{
// 			"too long",
// 			args{input: "2023-09-13"},
// 			false,
// 		},

// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Test_isValidMonthFormat(tt.args.input); got != tt.want {
// 				t.Errorf("Test_isValidMonthFormat() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_isValidMonthFormat(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Happy case",
			args{input: "2023-09"},
			true,
		},
		{
			"just junk",
			args{input: "junk"},
			false,
		},
		{
			"space",
			args{input: " "},
			false,
		},
		{
			"empty",
			args{input: ""},
			false,
		},
		{
			"invalid month",
			args{input: "2023-16"},
			false,
		},
		{
			"invalid year",
			args{input: "1515-06"},
			false,
		},
		{
			"too long",
			args{input: "2023-09-13"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidMonthFormat(tt.args.input); got != tt.want {
				t.Errorf("isValidMonthFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidOrgFormat(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Happy case",
			args{input: "jenkinsci"},
			true,
		},
		{
			"space",
			args{input: " "},
			false,
		},
		{
			"empty",
			args{input: ""},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidOrgFormat(tt.args.input); got != tt.want {
				t.Errorf("isValidOrgFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStartAndEndOfMonth(t *testing.T) {
	type args struct {
		shortMonth string
	}
	tests := []struct {
		name          string
		args          args
		wantStartDate string
		wantEndDate   string
	}{
		{
			"happy case",
			args{shortMonth: "2023-09"},
			"2023-09-01", "2023-09-30",
		},
		{
			"happy case2",
			args{shortMonth: "2023-02"},
			"2023-02-01", "2023-02-28",
		},
		{
			"Rubbish input",
			args{shortMonth: "blaahhh"},
			"0001-01-01", "0001-01-31",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartDate, gotEndDate := getStartAndEndOfMonth(tt.args.shortMonth)
			if gotStartDate != tt.wantStartDate {
				t.Errorf("getStartAndEndOfMonth() gotStartDate = %v, want %v", gotStartDate, tt.wantStartDate)
			}
			if gotEndDate != tt.wantEndDate {
				t.Errorf("getStartAndEndOfMonth() gotEndDate = %v, want %v", gotEndDate, tt.wantEndDate)
			}
		})
	}
}

func Test_splitPeriodForMaxQueryItem(t *testing.T) {
	type args struct {
		totalNbrIssue      int
		shortMonth         string
		requestedIteration int
	}
	tests := []struct {
		name              string
		args              args
		wantStartDate     string
		wantEndDate       string
		wantMoreIteration bool
	}{
		{
			"Below limit",
			args{
				totalNbrIssue:      800,
				shortMonth:         "2023-09",
				requestedIteration: 0,
			},
			"2023-09-01", "2023-09-30", false,
		},
		{
			"invalid input - negative iteration",
			args{
				totalNbrIssue:      800,
				shortMonth:         "2023-09",
				requestedIteration: -1,
			},
			"", "", false,
		},
		{
			"invalid input - negative total issues",
			args{
				totalNbrIssue:      -1,
				shortMonth:         "2023-09",
				requestedIteration: 0,
			},
			"", "", false,
		},
		{
			"invalid input - unprocessable number of issues",
			args{
				totalNbrIssue:      28001,
				shortMonth:         "2023-09",
				requestedIteration: 0,
			},
			"", "", false,
		},
		{
			"Below limit - invalid iteration number",
			args{
				totalNbrIssue:      800,
				shortMonth:         "2023-09",
				requestedIteration: 10,
			},
			"", "", false,
		},
		{
			"Below limit - junk short date",
			args{
				totalNbrIssue:      800,
				shortMonth:         "blaah",
				requestedIteration: 0,
			},
			"", "", false,
		},
		{
			"Above limit - junk short date",
			args{
				totalNbrIssue:      1400,
				shortMonth:         "blaah",
				requestedIteration: 0,
			},
			"", "", false,
		},
		{
			"Above limit (1400) - first iteration",
			args{
				totalNbrIssue:      1400,
				shortMonth:         "2023-09",
				requestedIteration: 0,
			},
			"2023-09-01", "2023-09-15", true,
		},
		{
			"Above limit (1400) - second and last iteration",
			args{
				totalNbrIssue:      1400,
				shortMonth:         "2023-09",
				requestedIteration: 1,
			},
			"2023-09-16", "2023-09-30", false,
		},
		{
			"Above limit (1900) - second and last iteration",
			args{
				totalNbrIssue:      1900,
				shortMonth:         "2023-09",
				requestedIteration: 1,
			},
			"2023-09-16", "2023-09-30", false,
		},
		{
			"Above limit (1900) - out of bound iteration",
			args{
				totalNbrIssue:      1900,
				shortMonth:         "2023-09",
				requestedIteration: 4,
			},
			"", "", false,
		},
		{
			"Above limit (2500) - need 3 iterations - 1",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-09",
				requestedIteration: 0,
			},
			"2023-09-01", "2023-09-10", true,
		},
		{
			"Above limit (2500) - need 3 iterations - 2",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-09",
				requestedIteration: 1,
			},
			"2023-09-11", "2023-09-20", true,
		},
		{
			"Above limit (2500) - need 3 iterations - 3",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-09",
				requestedIteration: 2,
			},
			"2023-09-21", "2023-09-30", false,
		},
		{
			"Above limit (1400) 31d - second and last iteration",
			args{
				totalNbrIssue:      1400,
				shortMonth:         "2023-08",
				requestedIteration: 1,
			},
			"2023-08-16", "2023-08-30", false,
		},
		{
			"Above limit (2500) 31d - need 3 iterations - 1",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-08",
				requestedIteration: 0,
			},
			"2023-08-01", "2023-08-10", true,
		},
		{
			"Above limit (2500) 31d - need 3 iterations - 2",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-08",
				requestedIteration: 1,
			},
			"2023-08-11", "2023-08-20", true,
		},
		{
			"Above limit (2500) 31d- need 3 iterations - 3",
			args{
				totalNbrIssue:      2500,
				shortMonth:         "2023-08",
				requestedIteration: 2,
			},
			"2023-08-21", "2023-08-31", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartDate, gotEndDate, gotMoreIteration := splitPeriodForMaxQueryItem(tt.args.totalNbrIssue, tt.args.shortMonth, tt.args.requestedIteration)
			if gotStartDate != tt.wantStartDate {
				t.Errorf("splitPeriodForMaxQueryItem() gotStartDate = %v, want %v", gotStartDate, tt.wantStartDate)
			}
			if gotEndDate != tt.wantEndDate {
				t.Errorf("splitPeriodForMaxQueryItem() gotEndDate = %v, want %v", gotEndDate, tt.wantEndDate)
			}
			if gotMoreIteration != tt.wantMoreIteration {
				t.Errorf("splitPeriodForMaxQueryItem() gotMoreIteration = %v, want %v", gotMoreIteration, tt.wantMoreIteration)
			}
		})
	}
}
