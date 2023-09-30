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
	"reflect"
	"testing"
)

// https://github.com/on4kjm/flecli/pull/1
var testResult1 = [][]string{
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jmMeessen", "2020-07"},
	{"on4kjm/flecli/1", "jlevesy", "2020-07"},
}

// https://github.com/jenkinsci/aqua-security-scanner-plugin/pull/51
var testResult2 = [][]string{
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "rajinikanthj", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "rajinikanthj", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
	{"jenkinsci/aqua-security-scanner-plugin/51", "deleted_user", "2023-06"},
}

// https://github.com/jenkins-infra/helm-charts/pull/586
var testResult3 = [][]string{
	{"jenkins-infra/helm-charts/586", "lemeurherve", "2023-08"},
	{"jenkins-infra/helm-charts/586", "lemeurherve", "2023-08"},
	{"jenkins-infra/helm-charts/586", "dduportal", "2023-08"},
}

// https://github.com/jenkinsci/build-blocker-plugin/pull/19
var testResult4 = [][]string{
	{"jenkinsci/build-blocker-plugin/19", "olamy", "2023-08"},
	{"jenkinsci/build-blocker-plugin/19", "jglick", "2023-08"},
	{"jenkinsci/build-blocker-plugin/19", "olamy", "2023-08"},
	{"jenkinsci/build-blocker-plugin/19", "jglick", "2023-08"},
	{"jenkinsci/build-blocker-plugin/19", "jonesbusy", "2023-08"},
	{"jenkinsci/build-blocker-plugin/19", "olamy", "2023-09"},
	{"jenkinsci/build-blocker-plugin/19", "Denis1990", "2023-09"},
	{"jenkinsci/build-blocker-plugin/19", "Denis1990", "2023-09"},
	{"jenkinsci/build-blocker-plugin/19", "jglick", "2023-08"},
}

// https://github.com/jenkinsci/credentials-plugin/pull/475
var testResult5 = [][]string{
	{"jenkinsci/credentials-plugin/475","jtnord","2023-09"},
}

func Test_fetchComments_alt(t *testing.T) {
	type args struct {
		org string
		prj string
		pr  int
	}
	tests := []struct {
		name           string
		args           args
		wantNbrComment int
		wantOutput     [][]string
	}{
		{
			"Blank test",
			args{
				org: "on4kjm",
				prj: "flecli",
				pr:  1,
			},
			57, testResult1,
		},
		{
			"PR with deleted user",
			args{
				org: "jenkinsci",
				prj: "aqua-security-scanner-plugin",
				pr:  51,
			},
			9, testResult2,
		},
		// enkins-infra/helm-charts/pull/586
		{
			"random PR",
			args{
				org: "jenkins-infra",
				prj: "helm-charts",
				pr:  586,
			},
			3, testResult3,
		},
		//https://github.com/jenkinsci/embeddable-build-status-plugin/pull/229
		{
			"PR with no comments",
			args{
				org: "jenkinsci",
				prj: "embeddable-build-status-plugin",
				pr:  229,
			},
			0, nil,
		},
		// https://github.com/jenkinsci/build-blocker-plugin/pull/19
		{
			"Random PR #2",
			args{
				org: "jenkinsci",
				prj: "build-blocker-plugin",
				pr:  19,
			},
			9, testResult4,
		},
		//https://github.com/jenkinsci/credentials-plugin/pull/475
		{
			"Random PR #3",
			args{
				org: "jenkinsci",
				prj: "credentials-plugin",
				pr:  475,
			},
			1, testResult5,
		},
		// unexisting PR
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNbrComment, gotOutput := fetchComments_alt(tt.args.org, tt.args.prj, tt.args.pr)
			if gotNbrComment != tt.wantNbrComment {
				t.Errorf("fetchComments_alt() gotNbrComment = %v, want %v", gotNbrComment, tt.wantNbrComment)
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("fetchComments_alt() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
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
			args{input: "aaaa bbbb cccc dddd eeee ffff"},
			"aaaa bbbb cccc dddd...",
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
