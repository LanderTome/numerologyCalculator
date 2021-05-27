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
	mapset "github.com/deckarep/golang-set"
	"testing"
)

func Test_nameSearchRandomCheck(t *testing.T) {
	first, offset, _ := nameSearch("Random ? Test", Pythagorean, []int{11, 22, 33}, true, NameSearchOpts{
		Count:      5,
		Offset:     0,
		Seed:       3384983,
		Dictionary: "usa_census",
		Gender:     Gender('M'),
		Sort:       "random",
		Database:   "sqlite://file::memory:?cache=shared",
	})
	second, _, _ := nameSearch("Random ? Test", Pythagorean, []int{11, 22, 33}, true, NameSearchOpts{
		Count:      5,
		Offset:     int(offset - 1),
		Seed:       3384983,
		Dictionary: "usa_census",
		Gender:     Gender('M'),
		Sort:       "random",
		Database:   "sqlite://file::memory:?cache=shared",
	})
	if first[4].Name != second[0].Name {
		t.Error("Random offset in name search is producing bad results.")
	}
}

func Test_nameSearchOffsetCheck(t *testing.T) {
	names, offset, _ := nameSearch("Ath? Doe", Pythagorean, []int{11, 22, 33}, true, NameSearchOpts{
		Count:      50,
		Offset:     0,
		Seed:       3384983,
		Dictionary: "usa_census",
		Gender:     Gender('F'),
		Sort:       "random",
		Database:   "sqlite://file::memory:?cache=shared",
	})
	if len(names) > 0 && offset > 0 {
		t.Error("Offset error.")
	}
}

func Test_nameSearchDupCheck(t *testing.T) {
	names, _, _ := nameSearch("John ? Doe", Pythagorean, []int{11, 22, 33, 44, 55, 66, 77, 88, 99}, true, NameSearchOpts{
		Count:      50000,
		Offset:     0,
		Seed:       3384983,
		Full:       []int{6},
		Dictionary: "usa_census",
		Gender:     Gender('M'),
		Sort:       "common",
		Database:   "sqlite://file::memory:?cache=shared",
	})
	for i := range names {
		_ = names[i].Full()
		for j := range names {
			if i != j && names[i] == names[j] {
				t.Error("Duplicate name found when there shouldn't be any.")
			}
		}
	}
	DB = nil
}

func Test_KarmicDebtSearch(t *testing.T) {
	type args struct {
		n             string
		numberSystem  NumberSystem
		masterNumbers []int
		reduceWords   bool
		opts          NameSearchOpts
	}
	tests := []struct {
		name        string
		args        args
		wantResults []NameNumerology
		wantOffset  int64
		wantErr     bool
	}{
		{"KarmicDebtSearch", args{"John ? Doe", Pythagorean, []int{11, 22, 33}, true, NameSearchOpts{
			Count:          200,
			Offset:         0,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         Gender('M'),
			Sort:           "common",
			Full:           []int{-13, -14, -16, -19},
			Vowels:         []int{-13, -14, -16, -19},
			Consonants:     []int{-13, -14, -16, -19},
			HiddenPassions: nil,
			KarmicLessons:  nil,
			Database:       "sqlite://file::memory:?cache=shared",
		}}, []NameNumerology{}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, _, _ := nameSearch(tt.args.n, tt.args.numberSystem, tt.args.masterNumbers, tt.args.reduceWords, tt.args.opts)
			set := mapset.NewSet()
			for _, result := range gotResults {
				for _, i := range append(append(result.Full().ReduceSteps, result.Vowels().ReduceSteps...), result.Consonants().ReduceSteps...) {
					set.Add(i)
				}
			}
			for _, kd := range []int{13, 14, 16, 19} {
				for _, value := range set.ToSlice() {
					if value.(int) == kd {
						t.Errorf("nameSearch() gotResults = %v, did not want %v", value, kd)
					}
				}
			}
		})
	}
	DB = nil
}

func Test_nameSearch(t *testing.T) {
	baseSearch := NameSearchOpts{
		Count:          500,
		Offset:         0,
		Seed:           0,
		Dictionary:     "usa_census",
		Gender:         Gender('B'),
		Sort:           "random",
		Full:           nil,
		Vowels:         nil,
		Consonants:     nil,
		HiddenPassions: nil,
		KarmicLessons:  nil,
		Database:       "sqlite://file::memory:?cache=shared",
	}
	tests := [][]int{
		{1, 4},
		{2, 7, 9},
		{3, -1, 2},
		{11, 8, -9},
		{3, -13, -14, 6},
		{22},
		{3, 8, 9},
	}
	// Full test
	for _, test := range tests {
		search := NameSearchOpts{
			Count:      baseSearch.Count,
			Offset:     baseSearch.Offset,
			Seed:       0,
			Dictionary: baseSearch.Dictionary,
			Gender:     baseSearch.Gender,
			Sort:       baseSearch.Sort,
			Full:       test,
			Database:   baseSearch.Database,
		}
		results, _, _ := Name("?", Pythagorean, []int{11, 22, 33}, true).Search(search)
		if len(results) == 0 {
			t.Errorf("No results when results expected. %v", search)
		}
		for _, r := range results {
			var ok bool
			for _, n := range r.Full().ReduceSteps {
				if inIntSlice(n, test) {
					ok = true
					break
				}
				if inIntSlice(-n, test) {
					ok = false
					break
				}
			}
			if !ok {
				t.Errorf("Unexpected result. reduced: %v | test: %v", r.Full().ReduceSteps, test)
			}
		}
	}
	// Vowels test
	for _, test := range tests {
		search := NameSearchOpts{
			Count:      baseSearch.Count,
			Offset:     baseSearch.Offset,
			Seed:       1,
			Dictionary: baseSearch.Dictionary,
			Gender:     baseSearch.Gender,
			Sort:       baseSearch.Sort,
			Vowels:     test,
			Database:   baseSearch.Database,
		}
		results, _, _ := Name("?", Pythagorean, []int{11, 22, 33}, true).Search(search)
		if len(results) == 0 {
			t.Errorf("No results when results expected. %v", search)
		}
		for _, r := range results {
			var ok bool
			for _, n := range r.Vowels().ReduceSteps {
				if inIntSlice(n, test) {
					ok = true
					break
				}
				if inIntSlice(-n, test) {
					ok = false
					break
				}
			}
			if !ok {
				t.Errorf("Unexpected result. reduced: %v | test: %v", r.Vowels().ReduceSteps, test)
			}
		}
	}
	// Consonants test
	for _, test := range tests {
		search := NameSearchOpts{
			Count:      baseSearch.Count,
			Offset:     baseSearch.Offset,
			Seed:       2,
			Dictionary: baseSearch.Dictionary,
			Gender:     baseSearch.Gender,
			Sort:       baseSearch.Sort,
			Consonants: test,
			Database:   baseSearch.Database,
		}
		results, _, _ := Name("?", Pythagorean, []int{11, 22, 33}, true).Search(search)
		if len(results) == 0 {
			t.Errorf("No results when results expected. %v", search)
		}
		for _, r := range results {
			var ok bool
			for _, n := range r.Consonants().ReduceSteps {
				if inIntSlice(n, test) {
					ok = true
					break
				}
				if inIntSlice(-n, test) {
					ok = false
					break
				}
			}
			if !ok {
				t.Errorf("Unexpected result. reduced: %v | test: %v", r.Consonants().ReduceSteps, test)
			}
		}
	}
	DB = nil
}
