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
		name    string
		args    args
		wantErr bool
	}{
		{
			"happy case",
			args{prSpec: "on4kjm/FLEcli/1"},
			false,
		},
		{
			"non numeric PR",
			args{prSpec: "on4kjm/FLEcli/aa"},
			true,
		},
		{
			"empty first field",
			args{prSpec: "/FLEcli/1"},
			true,
		},
		{
			"empty second field",
			args{prSpec: "on4kjm//1"},
			true,
		},
		{
			"empty third field",
			args{prSpec: "on4kjm/FLEcli/"},
			true,
		},
		{
			"too short",
			args{prSpec: "on4kjm/FLEcli"},
			true,
		},
		{
			"too long",
			args{prSpec: "on4kjm/FLEcli/1/zzz"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePRspec(tt.args.prSpec); (err != nil) != tt.wantErr {
				t.Errorf("validatePRspec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//FIXME: this should be an integration test: it requires a defined token and a set of global (default) values
func Test_getCommenters(t *testing.T) {
	type args struct {
		prSpec string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"happy case",
			args{prSpec: "on4kjm/FLEcli/1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getCommenters(tt.args.prSpec)
		})
	}
}
