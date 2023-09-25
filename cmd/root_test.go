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

func Test_validateHeader(t *testing.T) {
	type args struct {
		header          []string
		referenceHeader []string
		isVerbose       bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Not expected number of fields",
			args{
				header:          []string{"field1", "field2"},
				referenceHeader: []string{"field1", "field2", "field3"},
				isVerbose:       true,
			},
			false,
		},
		{
			"Not expected field name",
			args{
				header:          []string{"field1", "FIELD2", "field3"},
				referenceHeader: []string{"field1", "field2", "field3"},
				isVerbose:       true,
			},
			false,
		},
		{
			"Happy case",
			args{
				header:          []string{"field1", "field2", "field3"},
				referenceHeader: []string{"field1", "field2", "field3"},
				isVerbose:       true,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateHeader(tt.args.header, tt.args.referenceHeader, tt.args.isVerbose); got != tt.want {
				t.Errorf("validateHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
