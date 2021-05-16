// Copyright 2021 Robert D. Wukmir
// This file is subject to the terms and conditions defined in
// the LICENSE file, which is part of this source code package.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the
// License.

package numerology

import (
	"reflect"
	"testing"
)

func Test_hiddenPassion(t *testing.T) {
	type args struct {
		name         string
		numberSystem NumberSystem
	}
	tests := []struct {
		name string
		args args
		want HiddenPassionResults
	}{
		{"Alphabet", args{"abCdefghiJklmnOpqrstuvwxyz", Pythagorean}, HiddenPassionResults{
			Numbers:  []int{1, 2, 3, 4, 5, 6, 7, 8},
			MaxCount: 3,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counts, _, _ := countNumerologicalNumbers(tt.args.name, tt.args.numberSystem)
			if got := hiddenPassions(counts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hiddenPassions() = %v, want %v", got, tt.want)
			}
		})
	}
}
