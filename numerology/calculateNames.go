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
	"strings"
)

// NameNumerology is used as a struct to store the name and configuration information that is required
// to calculate the numerological values of names.
type NameNumerology struct {
	Name string
	*NameOpts
	*NameSearchOpts
	mask     *maskStruct
	counts   *map[int32]int
	unknowns *unknownCharacters
}

// initMask builds the maskStruct as it is needed.
func (n *NameNumerology) initMask() {
	if n.mask == nil {
		m := maskConstructor(n.Name)
		n.mask = &m
	}
}

// Full contains the numerological calculations done using all the letters of the given name. Sometimes referred to as
// the Destiny or Expression number.
func (n NameNumerology) Full() (result NumerologicalResult) {
	n.initMask()
	return calculateCoreNumber(n.Name, n.MasterNumbers, n.ReduceWords, n.NumberSystem, n.mask.Full())
}

// Destiny is an alias for Full().
func (n NameNumerology) Destiny() (destiny NumerologicalResult) {
	return n.Full()
}

// Expression is an alias for Full().
func (n NameNumerology) Expression() (expression NumerologicalResult) {
	return n.Full()
}

// Vowels contains the numerological calculations done using just the vowel letters of the given name. Sometimes
// referred to as Heart's Desire or Soul's Urge number.
func (n NameNumerology) Vowels() (result NumerologicalResult) {
	n.initMask()
	return calculateCoreNumber(n.Name, n.MasterNumbers, n.ReduceWords, n.NumberSystem, n.mask.Vowels())
}

// SoulsUrge is an alias for Vowels().
func (n NameNumerology) SoulsUrge() (soulsUrge NumerologicalResult) {
	return n.Vowels()
}

// HeartsDesire is an alias for Vowels().
func (n NameNumerology) HeartsDesire() (heartsDesire NumerologicalResult) {
	return n.Vowels()
}

// Consonants contains the numerological calculations done using just the consonant letters of the given name.
// Sometimes referred to as Personality number.
func (n NameNumerology) Consonants() (result NumerologicalResult) {
	n.initMask()
	return calculateCoreNumber(n.Name, n.MasterNumbers, n.ReduceWords, n.NumberSystem, n.mask.Consonants())
}

// Personality is an alias for Consonants().
func (n NameNumerology) Personality() (personality NumerologicalResult) {
	return n.Consonants()
}

// HiddenPassions contains the calculation of the numerological number(s) that repeat the most in the given name.
func (n NameNumerology) HiddenPassions() (result HiddenPassionResults) {
	return hiddenPassions(n.Counts())
}

// KarmicLessons contains the calculation of the numerological number(s) that do not appear in the given name.
func (n NameNumerology) KarmicLessons() (result KarmicLessonResults) {
	return karmicLessons(n.Counts())
}

// Search searches the name database for names that satisfy given numerological criteria.
func (n NameNumerology) Search(opts NameSearchOpts) (results []NameNumerology, offset int64, err error) {
	switch strings.Count(n.Name, "?") {
	case 0:
		return []NameNumerology{}, 0, errors.New("missing '?' in name")
	case 1:
		return nameSearch(n.Name, n.NumberSystem, n.MasterNumbers, n.ReduceWords, opts)
	default:
		return []NameNumerology{}, 0, errors.New("too many '?' (only able to search for one '?' at a time)")
	}
}

// Counts returns a map of each numerological value and how many times it appears in the name.
func (n *NameNumerology) Counts() (counts map[int32]int) {
	if n.counts == nil {
		counts, _, unknowns := countNumerologicalNumbers(n.Name, n.NumberSystem)
		n.counts = &counts
		n.unknowns = &unknowns
	}
	return *n.counts
}

// UnknownCharacters is a set of characters that cannot be converted to numerological values. They are
// ignored in calculations.
func (n *NameNumerology) UnknownCharacters() (unknowns []string) {
	if n.unknowns == nil {
		counts, _, unknowns := countNumerologicalNumbers(n.Name, n.NumberSystem)
		n.counts = &counts
		n.unknowns = &unknowns
	}
	return n.unknowns.Unacceptable().ToStringSlice()
}

// Name calculates various numerological numbers from a given name. Argument numberSystem indicates what
// calculation method is used (Pythagorean or Chaldean). Argument reduceWords determines whether each part of
// a name is reduced independently before being merged in a final number.
func Name(name string, numberSystem NumberSystem, masterNumbers []int, reduceWords bool) (result NameNumerology) {
	asciiName := ToAscii(name)
	result = NameNumerology{
		Name: asciiName,
		NameOpts: &NameOpts{
			NumberSystem:  numberSystem,
			MasterNumbers: masterNumbers,
			ReduceWords:   reduceWords,
		},
	}
	return
}

// Names calculates various numerological numbers from a list of given names. Argument numberSystem indicates what
// calculation method is used (Pythagorean or Chaldean). Argument reduceWords determines whether each part of
// a name is reduced independently before being merged in a final number.
func Names(names []string, numberSystem NumberSystem, masterNumbers []int, reduceWords bool) (results []NameNumerology) {
	opts := NameOpts{
		NumberSystem:  numberSystem,
		MasterNumbers: masterNumbers,
		ReduceWords:   reduceWords,
	}

	for _, n := range names {
		name := ToAscii(n)
		results = append(results, NameNumerology{name, &opts, nil, nil, nil, nil})
	}
	return
}
