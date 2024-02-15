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
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_get_quota(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"Happy case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get_quota()
		})
	}
}

func Test_get_quota_data_v4(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Happy case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get_quota_data_v4()
		})
	}
}

func Test_waitForReset(t *testing.T) {
	time1 := time.Now()

	seconds_toWait := 10
	waitForReset(seconds_toWait)

	time2 := time.Now()
	difference := time2.Sub(time1)

	assert.EqualValues(t, seconds_toWait, int(difference.Seconds()))

}

func Test_checkIfSufficientQuota(t *testing.T) {
	isRootDebug = true

	checkIfSufficientQuota(15)

	//TODO: How do we know that the result was expected ? =>very louzy test
}