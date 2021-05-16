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

// KarmicLessonResults packages the key information from the Numbers function.
type KarmicLessonResults struct {
	// Numbers is a slice of int of all the numbers that are not found in the Lookup.
	Numbers []int `json:"karmic_lessons"`
}

// karmicLessons calculates the numerology number(s) that do not show up in a given name.
func karmicLessons(counts map[int32]int) (results KarmicLessonResults) {
	// Iterate over counts to find the numbers that are missing.
	r := []int{}
	for k, v := range counts {
		// Look for numbers with a count of 0.
		if v == 0 {
			r = append(r, int(k))
		}
	}
	sort.Ints(r)
	return KarmicLessonResults{r}
}
