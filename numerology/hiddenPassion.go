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
	"sort"
)

// HiddenPassionResults packages the key information from the Numbers function.
type HiddenPassionResults struct {
	// Numbers is a slice of int of all the numbers that have the highest or equally highest count.
	Numbers []int `json:"numbers"`

	// MaxCount is the count of the highest number(s).
	MaxCount int `json:"max_count"`
}

// hiddenPassions calculates the numerology number(s) that are repeated the most in the given name.
func hiddenPassions(counts map[int32]int) (results HiddenPassionResults) {
	// Iterate over counts to find the numbers that match the largest count.
	passions := []int{}
	max := 0
	for k, v := range counts {
		if v > max { // If we find a new max then update the var and replace passions with new number.
			max = v
			passions = []int{int(k)}
		} else if v == max { // If we find an equal to max then add the number to slice of passions.
			passions = append(passions, int(k))
		}
	}
	sort.Ints(passions)
	return HiddenPassionResults{passions, max}
}
