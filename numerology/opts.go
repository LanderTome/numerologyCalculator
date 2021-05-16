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
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// NameOpts contains the options that are required in order to get a numerological value from a name.
type NameOpts struct {
	// NumberSystem is the numerological number system to use to convert the name into numbers.
	NumberSystem NumberSystem `json:"number_system"`

	// MasterNumbers are two (or more) digit numbers that should not be reduced because they are
	// considered to have special numerological value. The most commonly considered master numbers
	// are repeating digits. ex 11, 22, 33, 44, 55
	MasterNumbers []int `json:"master_numbers"`

	// ReduceWords determines whether a name is first separated into individual names that are
	// reduced independently before summing together and reduce a final time or if the whole
	// name is reduced all at once.
	ReduceWords bool `json:"reduce_words"`
}

// NameSearchOpts contains the search options specific to searching for numerological names.
type NameSearchOpts struct {
	// Count is the number of results to return.
	Count int `json:"count,omitempty"`

	// Offset is used to page through the search results. Its main use is with an web service.
	Offset int `json:"offset,omitempty"`

	// Seed is used to generate random search results for names. Seed is necessary in order to page through
	// random results by maintaining the random order.
	Seed int64 `json:"seed,omitempty"`

	// Database is the DSN string that connects to the database that has the precalculated numerological names.
	Database string `json:"db,omitempty"`

	// Dictionary is the name of the database table to search.
	Dictionary string `json:"dictionary,omitempty"`

	// Gender is a filter to limit search results to names that are generally (m)ale, (f)emale, or (b)oth.
	Gender Gender `json:"gender,omitempty"`

	// Sort is the method of sorting used when return names. The options are "common", "uncommon", and "random".
	Sort string `json:"sort,omitempty"`

	// Full, Vowels, and Consonants are  the numerological numbers to look for while searching. They are calculated
	// using all the letters of the name, just the vowels, and just the consonants, respectively. There are various
	// common numerological names for these values; destiny, express, heart's desire, soul's urge, personality, etc.
	// These names were chosen in order to avoid favoring one particular designation, and to make it clear what
	// calculation each one is performing.

	// Full is the numerological numbers you want to search for that are calculated with all the letters of the given
	// name. Often referred to as "Destiny" or "Expression" number, Full was chosen in order to avoid favoring
	// one particular designation, and to make it clear what calculation is being performing.
	Full []int `json:"full,omitempty"`

	// Vowels is the numerological numbers you want to search for that are calculated with just the vowel letters of
	// the given name. Often referred to as "Soul's Urge" or "Heart's Desire" number, Vowels was chosen in order to
	// avoid favoring one particular designation, and to make it clear what calculation is being performing.
	//
	// This calculation is complicated by the fact that the letter "Y" sometimes acts like a vowel. Because the
	// database is precomputed, a consistent set rules need to be used in order to get consistent results. The
	// method used will generally be correct, but names with "Y" in them are probably the least reliable results.
	Vowels []int `json:"vowels,omitempty"`

	// Consonants is the numerological numbers you want to search for that are calculated with just the consonant
	// letters of the given name. Sometimes referred to as "Personality", Consonants was chosen in order to avoid
	// favoring one particular designation, and to make it clear what calculation is being performing.
	//
	// This calculation is complicated by the fact that the letter "Y" sometimes acts like a vowel. Because the
	// database is precomputed, a consistent set rules need to be used in order to get consistent results. The
	// method used will generally be correct, but names with "Y" in them are probably the least reliable results.
	Consonants []int `json:"consonants,omitempty"`

	// HiddenPassions contains the numerological number(s) that occur the most frequently in the name.
	// Positive and  negative numbers can be used. Positive numbers are inclusive; they ensure that
	// number will exist in the final grouping. Negative numbers are exclusive, they exclude that number from
	// existing in the final grouping.
	HiddenPassions []int `json:"hidden_passions,omitempty"`

	// KarmicLessons contains the numerological number(s) that do not appear at all in the name.
	// Positive and  negative numbers can be used. Positive numbers are inclusive; they ensure that
	// number will exist in the final grouping. Negative numbers are exclusive, they exclude that number from
	// existing in the final grouping.
	KarmicLessons []int `json:"karmic_lessons,omitempty"`
}

// DateOpts contains the options that are required in order to get a numerological value from a date.
type DateOpts struct {
	// MasterNumbers are two (or more) digit numbers that should not be reduced because they are
	// considered to have special numerological value. The most commonly considered master numbers
	// are repeating digits. ex 11, 22, 33, 44, 55
	MasterNumbers []int `json:"master_numbers"`
}

// DateSearchOpts contains the search options specific to searching for numerological dates.
type DateSearchOpts struct {
	// Count is the number of results to return.
	Count int `json:"count,omitempty"`

	// Offset is used to page through the search results. Its main use is with an web service.
	Offset int `json:"offset,omitempty"`

	// Match is a slice of ints that are the numerological values that we want to find.
	Match []int `json:"match,omitempty"`

	// MonthsForward is the number of months in which to search. Usually one is trying to find a date within
	// some reasonable time frame; like for a wedding or event.
	MonthsForward int `json:"months_forward,omitempty"`

	// Dow, or Days of the Week, limits the results to ones that fall on certain days.
	// 0=Sun 1=Mon 2=Tue 3=Wed 4=Thu 5=Fri 6=Sat
	Dow []int `json:"dow,omitempty"`

	// LifePath is generally used with ones date of birth. If false, then Master Numbers are ignored
	// for the calculation.
	LifePath bool `json:"life_path"`
}

// Gender constants for use when searching for names of a particular gender.
const (
	Male   = Gender('M')
	Female = Gender('F')
	Both   = Gender('B')
)

// Gender is custom type that is needed so we can customize the Marshaling and Unmarshaling.
type Gender rune

func (g Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToUpper(string(g)))
}

// JSON will come in as a string, but we want to convert it to a rune. Only accepted options are (M)ale,
// (F)emale, and (B)oth. If the option is not one of these then it defaults to (B)oth.
func (g *Gender) UnmarshalJSON(value []byte) error {
	v := strings.Trim(string(value), `"`)
	r := unicode.ToUpper(rune(v[0]))
	genderTest := Gender(r)
	if genderTest == Male || genderTest == Female || genderTest == Both {
		*g = genderTest
	} else {
		return errors.New(fmt.Sprintf("unknown gender code: '%v'", v))
	}
	return nil
}
