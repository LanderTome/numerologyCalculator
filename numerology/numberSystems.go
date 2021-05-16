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
	"encoding/json"
	"strings"
)

// NumberSystem contains the information necessary to convert letters to their numerological value.
type NumberSystem struct {
	// Name of the NumberSystem that is used to translate to and from JSON.
	Name string

	// NumberMapping is a map of letters and their corresponding value.
	NumberMapping map[int32]int

	// ValidNumbers shows all the numbers that the number system accepts. The Chaldean number system,
	// in particular, does not use the number 9 in conversions.
	ValidNumbers []int
}

// MarshalJSON returns the name of the number system because it doesn't make sense to encode the whole struct
// for output to JSON.
func (ns NumberSystem) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(ns.Name))
}

// UnmarshalJSON translates the name of the number system to the appropriate NumberSystem struct.
func (ns *NumberSystem) UnmarshalJSON(value []byte) error {
	numberSystem, err := GetNumberSystem(strings.Trim(string(value), `"`))
	*ns = numberSystem
	return err
}

// Todo: Add non-English characters like é, ö, å, etc.

// The conversion table for the Chaldean number system.
var Chaldean = NumberSystem{
	Name: "Chaldean",
	NumberMapping: map[int32]int{
		'a': 1, 'i': 1, 'j': 1, 'q': 1, 'y': 1,
		'b': 2, 'k': 2, 'r': 2,
		'c': 3, 'g': 3, 'l': 3, 's': 3,
		'd': 4, 'm': 4, 't': 4,
		'e': 5, 'h': 5, 'n': 5, 'x': 5,
		'u': 6, 'v': 6, 'w': 6,
		'o': 7, 'z': 7,
		'f': 8, 'p': 8,
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
		'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		' ': 0, '.': 0, '-': 0, // Common acceptable characters in names that have no value.
	},
	ValidNumbers: []int{1, 2, 3, 4, 5, 6, 7, 8},
}

// The conversion table for the Pythagorean number system.
var Pythagorean = NumberSystem{
	Name: "Pythagorean",
	NumberMapping: map[int32]int{
		'a': 1, 'j': 1, 's': 1,
		'b': 2, 'k': 2, 't': 2,
		'c': 3, 'l': 3, 'u': 3,
		'd': 4, 'm': 4, 'v': 4,
		'e': 5, 'n': 5, 'w': 5,
		'f': 6, 'o': 6, 'x': 6,
		'g': 7, 'p': 7, 'y': 7,
		'h': 8, 'q': 8, 'z': 8,
		'i': 9, 'r': 9,
		'0': 0, '1': 1, '2': 2,
		'3': 3, '4': 4, '5': 5,
		'6': 6, '7': 7, '8': 8,
		'9': 9,
		' ': 0, '.': 0, '-': 0, // Commonly acceptable characters in names that have no value.
	},
	ValidNumbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
}
