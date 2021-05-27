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

// Package numerology provides methods to calculate various numerological numbers from both
// names and dates. It can calculate Expression/Destiny, Soul's Urge/Heart's Desire, Personality, Hidden Passion,
// Karmic Lesson, and Life Path numbers using both Pythagorean and Chaldean number systems. It allows of for
// various methods of calculation with custom Master Numbers and the option to reduce each word or only whole
// names. It also provides summaries of the steps of the calculation process that can be used to display how the
// calculations were derived.
//
// The truly unique aspect of this package is also provides methods to search for either names or dates that
// satisfy various numerological criteria. This is useful if one is trying to find a name for a baby or a wedding
// date that has the desirable numerological properties. The way that it does this is using a precomputed table
// of names with corresponding numerological properties. This allows the table to be searched for specific names
// that satisfy given constraints. This method is surprisingly efficient. A database of 100,000 names with all
// the necessary pre-calculations is only 13.5 MB and searches take only hundredths of a second.
package numerology

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// AcceptableUnknownCharacters is a regexp that matches characters that should be considered acceptable in names.
// This includes dashes (Hernandez-Johnson), periods (Jr.), and white space characters. This is used before
// Marshaling to JSON. Replace this variable if a different set of AcceptableUnknownCharacters is needed.
var AcceptableUnknownCharacters = regexp.MustCompile("[-\\s.']")

// unknownCharacters contains the set of characters that were not defined in the given numbers system. Used to
// Function characters that are being used, but do not seem to have a numerological value. ex. !#$*
//
// The purpose of this is to give an indication when the numerological conversion may not be accurate. There is
// no generally agreed upon conversion for non-alphanumeric symbols. This is not an error, per-se, but there
// needs to be some method of discovering what characters were ignored. Some common characters in names that
// should not be necessarily be considered errors are dashes (Hernandez-Johnson), periods (Jr.), and white
// space characters.
type unknownCharacters struct {
	// Set is used because we want the unknownCharacters to act like a set by removing duplicates.
	mapset.Set
}

// Add adds another character to the list of unknownCharacters if it is not already in the list.
func (ec unknownCharacters) Add(r rune) unknownCharacters {
	set := ec.Set
	set.Add(unicode.ToLower(r))
	return unknownCharacters{set}
}

// Union joins two sets of unknownCharacters and removes duplicates.
func (ec unknownCharacters) Union(uc unknownCharacters) unknownCharacters {
	oldSet := ec.Set
	newSet := oldSet.Union(uc.Set)
	return unknownCharacters{newSet}
}

// ToRuneSlice returns the set as a slice of runes.
func (ec unknownCharacters) ToRuneSlice() []rune {
	if ec.Set == nil {
		return []rune{}
	}
	var runes []rune
	for _, i := range ec.ToSlice() {
		runes = append(runes, i.(rune))
	}
	return runes
}

// ToStringSlice returns the set as a slice of strings.
func (ec unknownCharacters) ToStringSlice() []string {
	if ec.Set == nil {
		return []string{}
	}
	s := []string{}
	for _, i := range ec.ToSlice() {
		s = append(s, string(i.(rune)))
	}
	return s
}

// Unacceptable takes a regexp of characters that should be considered acceptable and returns the remaining
// characters as unacceptable. The purpose of this method is to filter out unconverted characters that should
// be considered acceptable. It is implemented in a way that the acceptable characters can be customized.
func (ec unknownCharacters) Unacceptable() (unknowns unknownCharacters) {
	errorChars := mapset.NewSet()
	for _, r := range ec.ToRuneSlice() {
		if !AcceptableUnknownCharacters.MatchString(string(r)) {
			errorChars.Add(r)
		}
	}
	return unknownCharacters{errorChars}
}

func (ec unknownCharacters) MarshalJSON() ([]byte, error) {
	unacceptable := ec.Unacceptable()
	return json.Marshal(unacceptable.ToStringSlice())
}

type letterValue struct {
	Letter string `json:"letter"`
	Value  int    `json:"value"`
}

// NumerologicalResult contains the results of the conversion of a name.
type NumerologicalResult struct {
	// Value contains the reduced numerological value of the name.
	Value int `json:"value"`

	// ReduceSteps contains the conversion at each step until the final number.
	ReduceSteps []int `json:"reduce_steps"`

	// Breakdown contains each individual Breakdown that contributed to the final calculation.
	Breakdown []Breakdown `json:"breakdown"`
}

// Debug returns a printable string that offers a simplistic summary of the numerological conversion.
func (n NumerologicalResult) Debug() (debug string) {
	var words []string
	for _, breakdown := range n.Breakdown {
		var letters []string
		var calcs []string
		for _, l := range breakdown.LetterValues {
			letters = append(letters, l.Letter)
			if l.Value > 0 {
				calcs = append(calcs, strconv.Itoa(l.Value))
			} else {
				calcs = append(calcs, "·")
			}
		}
		for _, reduce := range breakdown.ReduceSteps {
			calcs = append(calcs, "=", strconv.Itoa(reduce))
		}
		words = append(words, fmt.Sprintf("%v\n%v", strings.Join(letters, " "), strings.Join(calcs, " ")))
	}
	var r []string
	for _, reduce := range n.ReduceSteps {
		r = append(r, strconv.Itoa(reduce))
	}
	words = append(words, "Reduce: "+strings.Join(r, " = "))
	return strings.Join(words, "\n")
}

// Breakdown contains information from the conversion of a name to its numerological equivalent.
type Breakdown struct {
	// Value contains the reduced numerological value of the name.
	Value int `json:"value"`

	// ReduceSteps contains the conversion at each step until the final number.
	ReduceSteps []int `json:"reduce_steps"`

	// LetterValues shows the numerological value of each letter in the name.
	LetterValues []letterValue `json:"letter_values"`
}

type letterMask []bool

// maskStruct holds the masks used to single out parts of the name that are used in particular
// calculations. The struct does the calculations as necessary, but caches the results so repeated
// calculations are quicker.
type maskStruct struct {
	letterMask
}

func (m *maskStruct) Full() (mask letterMask) {
	mask = make(letterMask, len(m.letterMask))
	for i := 0; i < len(mask); i++ {
		mask[i] = true
	}
	return mask
}

func (m *maskStruct) Vowels() (mask letterMask) {
	return m.letterMask
}

func (m *maskStruct) Consonants() (mask letterMask) {
	mask = make(letterMask, len(m.letterMask))
	vowels := m.Vowels()
	for i := 0; i < len(mask); i++ {
		mask[i] = !vowels[i]
	}
	return mask
}

func maskConstructor(s string) maskStruct {
	var mask letterMask
	isVowel := func(r rune) bool {
		vowels := map[rune]bool{
			'a': true,
			'e': true,
			'i': true,
			'o': true,
			'u': true,
		}
		_, ok := vowels[r]
		return ok
	}
	runes := []rune(strings.ToLower(s))
	for i, letter := range runes {
		// Special rules for "Y" from https://www.worldnumerology.com/numerology-Y-vowel-consonant.htm
		var m bool
		if letter == 'y' {
			switch {
			// If the name is only one character long then "Y" must be a vowel. Weird name, though.
			case len(runes) == 1:
				m = true

			// If the Y is the first letter of a name and it is followed by a consonant, it is vowel. (Yvonne, Ylsa, Yvette)
			// If the Y is the first letter of a name and it is followed by another vowel, it is a consonant. (Yolanda, Yammy)
			case i == 0:
				m = !isVowel(runes[i+1])

			// If the Y is the last letter of a name and it comes after a consonant, it is vowel. (Barry, Tommy, Jimmy, etc.)
			// If the Y is the last letter of a name and it comes after a vowel, it is a consonant. (Mulrooney, Mickey)
			case i == len(runes)-1:
				m = !isVowel(runes[i-1])

			// If the Y is found between two consonants, it’s a vowel (Kyle, Tyson)
			// If the Y is found between two vowels, it’s a consonant (Eyarta)
			case i != 0 && i != len(runes)-1 && isVowel(runes[i-1]) == isVowel(runes[i+1]):
				m = !(isVowel(runes[i-1]) && isVowel(runes[i+1]))

			// If the Y is found between a consonant and a vowel, the default should be a vowel, but there are exceptions
			// If the Y is found between a vowel and a consonant, the default should be a consonant, but there are exceptions
			case i != 0 && i != len(runes)-1 && isVowel(runes[i-1]) != isVowel(runes[i+1]):
				m = isVowel(runes[i+1])
			}
		} else { // If letter is not a "Y" then just look up if it is a vowel.
			m = isVowel(letter)
		}
		mask = append(mask, m)
	}
	return maskStruct{mask}
}

func reduceNumbers(n int, masterNumbers []int, steps []int) (reduceSteps []int) {
	steps = append(steps, n)
	// If totalValue is less than 10 then return because no more reducing can be done.
	if n < 10 {
		return steps
	}
	// Check for master numbers
	if isMasterNumber(n, masterNumbers) {
		return steps
	}
	// Sum and reduce result
	var totalValue int
	for _, val := range splitNumber(uint64(n)) {
		totalValue += val
	}
	return reduceNumbers(totalValue, masterNumbers, steps)
}

func numerologyCalculation(s string, masterNumbers []int, numberSystem NumberSystem, reduce bool, mask []bool) (results Breakdown) {
	var nameSteps []letterValue
	var totalValue int

	// Iterate over letters and convert them to numbers based on our chosen number system
	for i, letter := range strings.TrimSpace(s) {
		// Check if letter is in our pre-made number system and that it is not ignored by our letter mask
		if letterVal, ok := numberSystem.NumberMapping[unicode.ToLower(letter)]; ok && mask[i] == true {
			totalValue += letterVal
			nameSteps = append(nameSteps, letterValue{string(letter), letterVal})
		} else {
			// If the letter either not in the number system or it must be masked so it is zero.
			nameSteps = append(nameSteps, letterValue{string(letter), 0})
		}
	}
	var reduceSteps []int
	if reduce {
		reduceSteps = reduceNumbers(totalValue, masterNumbers, []int{})
	} else {
		reduceSteps = []int{totalValue}
	}
	return Breakdown{
		Value:        reduceSteps[len(reduceSteps)-1],
		ReduceSteps:  reduceSteps,
		LetterValues: nameSteps,
	}
}

func calculateCoreNumber(name string, masterNumbers []int, reduceWords bool, numberSystem NumberSystem, mask letterMask) (results NumerologicalResult) {
	var totalValue int
	var breakdown []Breakdown
	var idxStart int
	// Manually look for spaces so we can sync the name with the mask.
	for idxEnd, n := range name + " " {
		// Look for a space to use to break the name.
		if n != ' ' {
			continue
		}
		// If start and end indexes are the same then it means the name has 0 length. No reason to process those.
		if idxStart != idxEnd {
			calc := numerologyCalculation(name[idxStart:idxEnd], masterNumbers, numberSystem, reduceWords, mask[idxStart:idxEnd])
			breakdown = append(breakdown, calc)
			totalValue += calc.Value
		}
		idxStart = idxEnd + 1
	}
	var reduceSteps []int
	// If there is only one breakdown then the reducing steps are the same.
	if len(breakdown) == 1 && reduceWords {
		reduceSteps = breakdown[0].ReduceSteps
	} else {
		reduceSteps = reduceNumbers(totalValue, masterNumbers, []int{})
	}
	return NumerologicalResult{
		Value:       reduceSteps[len(reduceSteps)-1],
		ReduceSteps: reduceSteps,
		Breakdown:   breakdown,
	}
}
