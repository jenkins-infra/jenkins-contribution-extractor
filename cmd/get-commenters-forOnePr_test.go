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

// FIXME: this should be an integration test: it requires a defined token and a set of global (default) values
// Worked accidentally on GitHub Action
func Test_getCommenters(t *testing.T) {
	type args struct {
		prSpec         string
		isAppend       bool
		isNoHeader     bool
		outputFileName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"happy case",
			args{
				prSpec:         "on4kjm/FLEcli/1",
				isAppend:       false,
				isNoHeader:     false,
				outputFileName: "jenkins_commenters_data.csv",
			},
		},
		{
			"happy case - append",
			args{
				prSpec:         "on4kjm/FLEcli/1",
				isAppend:       true,
				isNoHeader:     false,
				outputFileName: "jenkins_commenters_data.csv",
			},
		},
		{
			"ghost user",
			args{
				prSpec:         "jenkinsci/aqua-security-scanner-plugin/51",
				isAppend:       true,
				isNoHeader:     false,
				outputFileName: "jenkins_commenters_data.csv",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getCommenters(tt.args.prSpec, tt.args.isAppend, tt.args.isNoHeader, tt.args.outputFileName)
		})
	}
}

// https://github.com/on4kjm/flecli/pull/1
var testResult1 = []string{
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jmMeessen\",\"2020-07\"",
	"\"on4kjm/flecli/1\",\"jlevesy\",\"2020-07\"",
}

// https://github.com/jenkinsci/aqua-security-scanner-plugin/pull/51
var testResult2 = []string{
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"rajinikanthj\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"rajinikanthj\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
	"\"jenkinsci/aqua-security-scanner-plugin/51\",\"deleted_user\",\"2023-06\"",
}

// https://github.com/jenkins-infra/helm-charts/pull/586
var testResult3 = []string{
	"\"jenkins-infra/helm-charts/586\",\"lemeurherve\",\"2023-08\"",
	"\"jenkins-infra/helm-charts/586\",\"lemeurherve\",\"2023-08\"",
	"\"jenkins-infra/helm-charts/586\",\"dduportal\",\"2023-08\"",
}

// https://github.com/jenkinsci/build-blocker-plugin/pull/19
var testResult4 = []string{
	"\"jenkinsci/build-blocker-plugin/19\",\"olamy\",\"2023-08\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"jglick\",\"2023-08\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"olamy\",\"2023-08\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"jglick\",\"2023-08\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"jonesbusy\",\"2023-08\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"olamy\",\"2023-09\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"Denis1990\",\"2023-09\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"Denis1990\",\"2023-09\"",
	"\"jenkinsci/build-blocker-plugin/19\",\"jglick\",\"2023-08\"",
}

// https://github.com/jenkinsci/credentials-plugin/pull/475
var testResult5 = []string{
	"\"jenkinsci/credentials-plugin/475\",\"jtnord\",\"2023-09\"",
}

//bot test
var testResult6 = []string{
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"dwnusbaum\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"NicuPascu\",\"2020-02\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-03\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"bitwiseman\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"bitwiseman\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"olamy\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-05\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"bitwiseman\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
	"\"jenkinsci/blueocean-plugin/2050\",\"stuartrowe\",\"2020-04\"",
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
		wantOutput     []string
	}{
		{
			"first test",
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
		{
			"PR with bot user",
			args{
				org: "jenkinsci",
				prj: "blueocean-plugin",
				pr:  2050,
			},
			26, testResult6,
		},
		// jenkins-infra/helm-charts/pull/586
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
		{
			"no PR to see here",
			args{
				org: "on4kjm",
				prj: "flecli",
				pr:  4,
			},
			0, nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNbrComment, gotOutput := fetchComments_v4(tt.args.org, tt.args.prj, tt.args.pr)
			if gotNbrComment != tt.wantNbrComment {
				t.Errorf("fetchComments_alt() gotNbrComment = %v, want %v", gotNbrComment, tt.wantNbrComment)
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("fetchComments_alt() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
