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
	"errors"
	mapset "github.com/deckarep/golang-set"
	"github.com/mozillazg/go-unidecode"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var duplicateSpaces = regexp.MustCompile("\\s+")
var removeNewLines = regexp.MustCompile("[\r\n\t]")

func removeDuplicateSpaces(s string) string {
	return duplicateSpaces.ReplaceAllString(s, " ")
}

// ToAscii transliterates a unicode string to ascii string and removes extraneous
// newline chars and whitespace. This conversion is not 100% due to complexities
// of language. See this article by Sean M. Burke for more information.
// https://interglacial.com/~sburke/tpj/as_html/tpj22.html
func ToAscii(s string) string {
	ascii := unidecode.Unidecode(s)
	noNewLines := removeNewLines.ReplaceAllString(ascii, "")
	noDupWhiteSpace := duplicateSpaces.ReplaceAllString(noNewLines, " ")
	return strings.TrimSpace(noDupWhiteSpace)
}

func numberOfQuestionMarks(s string) int {
	return strings.Count(s, "?")
}

func isMasterNumber(v int, masterNumbers []int) bool {
	for _, n := range masterNumbers {
		if v == n {
			return true
		}
	}
	return false
}

func inIntSlice(v int, slice []int) bool {
	for _, n := range slice {
		if v == n {
			return true
		}
	}
	return false
}

// Convert a number into a slice of numbers. 1234 -> [1,2,3,4]
func splitNumber(i uint64) []int {
	listNumbers := []int{}
	convertedToString := strconv.FormatUint(i, 10)
	// Borrows from standard library for strings.Atoi. Figured we could make it faster by skipping some of
	// their checks because we know that the number is guaranteed to be between 0 and 9.
	for _, ch := range []byte(convertedToString) {
		ch -= '0'
		listNumbers = append(listNumbers, int(ch))
	}
	return listNumbers
}

// countNumerologicalNumbers converts the letters of a name into the corresponding numerological values based on the given
// numberSystem argument. Those numbers are then compiled into a map that counts the occurrences of each number.
func countNumerologicalNumbers(name string, numberSystem NumberSystem) (counts map[int32]int, maxCount int, unknownChars unknownCharacters) {
	counts = map[int32]int{}
	unknownChars = unknownCharacters{mapset.NewSet()}
	// Initialize map with all valid numbers.
	for _, i := range numberSystem.ValidNumbers {
		counts[int32(i)] = 0
	}

	for _, n := range name {
		num := int32(numberSystem.NumberMapping[unicode.ToLower(n)])
		// Increment the counts of the numerological numbers.
		if c, ok := counts[num]; ok {
			newCount := c + 1
			counts[num] = newCount
			if newCount > maxCount {
				maxCount = newCount
			}
		} else {
			// If a number isn't valid then skip it. Generally this would be a '0'.
			unknownChars = unknownChars.Add(n)
			continue
		}
	}
	return counts, maxCount, unknownChars
}

// GetNumberSystem returns the appropriate NumberSystem type from the number systems name. (Pythagorean, Chaldean).
func GetNumberSystem(s string) (NumberSystem, error) {
	switch strings.ToLower(s) {
	case "pythagorean":
		return Pythagorean, nil
	case "chaldean":
		return Chaldean, nil
	default:
		return NumberSystem{}, errors.New("unknown number system: " + s)
	}
}
