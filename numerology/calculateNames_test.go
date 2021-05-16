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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func init() {
	// Change working directory because testing messes it up.
	_, b, _, _ := runtime.Caller(0)
	if err := os.Chdir(filepath.Join(filepath.Dir(b), "..")); err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	type args struct {
		name          string
		numberSystem  NumberSystem
		masterNumbers []int
		reduceWords   bool
	}
	tests := []struct {
		name       string
		args       args
		wantResult NameNumerology
	}{
		{"Arthur Von Black", args{"Arthur Von Black", Pythagorean, []int{11, 22, 33}, true},
			NameNumerology{
				Name: "Arthur Von Black",
				NameOpts: &NameOpts{
					NumberSystem:  Pythagorean,
					MasterNumbers: []int{11, 22, 33},
					ReduceWords:   true,
				},
				NameSearchOpts: nil,
				mask:           nil,
				counts:         nil,
				unknowns:       nil,
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Name(tt.args.name, tt.args.numberSystem, tt.args.masterNumbers, tt.args.reduceWords); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Name() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestNames(t *testing.T) {
	type args struct {
		names         []string
		numberSystem  NumberSystem
		masterNumbers []int
		reduceWords   bool
	}
	tests := []struct {
		name        string
		args        args
		wantResults []NameNumerology
	}{
		{"Name 1, Name 2, Name 3", args{[]string{"Name 1", "Name 2", "Name 3"}, Pythagorean, []int{11, 22, 33}, true},
			[]NameNumerology{
				{
					Name: "Name 1",
					NameOpts: &NameOpts{
						NumberSystem:  Pythagorean,
						MasterNumbers: []int{11, 22, 33},
						ReduceWords:   true,
					},
					NameSearchOpts: nil,
					mask:           nil,
					counts:         nil,
					unknowns:       nil,
				},
				{
					Name: "Name 2",
					NameOpts: &NameOpts{
						NumberSystem:  Pythagorean,
						MasterNumbers: []int{11, 22, 33},
						ReduceWords:   true,
					},
					NameSearchOpts: nil,
					mask:           nil,
					counts:         nil,
					unknowns:       nil,
				},
				{
					Name: "Name 3",
					NameOpts: &NameOpts{
						NumberSystem:  Pythagorean,
						MasterNumbers: []int{11, 22, 33},
						ReduceWords:   true,
					},
					NameSearchOpts: nil,
					mask:           nil,
					counts:         nil,
					unknowns:       nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := Names(tt.args.names, tt.args.numberSystem, tt.args.masterNumbers, tt.args.reduceWords); !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("Names() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func ExampleNameNumerology_Full() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Full().Debug())
	// Output:
	// K e v i n
	// 2 5 4 9 5 = 25 = 7
	// N o r w o o d
	// 5 6 9 5 6 6 4 = 41 = 5
	// B a c o n
	// 2 1 3 6 5 = 17 = 8
	// Reduce: 20 = 2
}

func ExampleNameNumerology_Destiny() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Destiny().Debug())
	// Output:
	// K e v i n
	// 2 5 4 9 5 = 25 = 7
	// N o r w o o d
	// 5 6 9 5 6 6 4 = 41 = 5
	// B a c o n
	// 2 1 3 6 5 = 17 = 8
	// Reduce: 20 = 2
}

func ExampleNameNumerology_Expression() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Expression().Debug())
	// Output:
	// K e v i n
	// 2 5 4 9 5 = 25 = 7
	// N o r w o o d
	// 5 6 9 5 6 6 4 = 41 = 5
	// B a c o n
	// 2 1 3 6 5 = 17 = 8
	// Reduce: 20 = 2
}

func ExampleNameNumerology_Vowels() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Vowels().Debug())
	// Output:
	// K e v i n
	// · 5 · 9 · = 14 = 5
	// N o r w o o d
	// · 6 · · 6 6 · = 18 = 9
	// B a c o n
	// · 1 · 6 · = 7
	// Reduce: 21 = 3
}

func ExampleNameNumerology_SoulsUrge() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.SoulsUrge().Debug())
	// Output:
	// K e v i n
	// · 5 · 9 · = 14 = 5
	// N o r w o o d
	// · 6 · · 6 6 · = 18 = 9
	// B a c o n
	// · 1 · 6 · = 7
	// Reduce: 21 = 3
}

func ExampleNameNumerology_HeartsDesire() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.HeartsDesire().Debug())
	// Output:
	// K e v i n
	// · 5 · 9 · = 14 = 5
	// N o r w o o d
	// · 6 · · 6 6 · = 18 = 9
	// B a c o n
	// · 1 · 6 · = 7
	// Reduce: 21 = 3
}

func ExampleNameNumerology_Consonants() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Consonants().Debug())
	// Output:
	// K e v i n
	// 2 · 4 · 5 = 11
	// N o r w o o d
	// 5 · 9 5 · · 4 = 23 = 5
	// B a c o n
	// 2 · 3 · 5 = 10 = 1
	// Reduce: 17 = 8
}

func ExampleNameNumerology_Personality() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.Personality().Debug())
	// Output:
	// K e v i n
	// 2 · 4 · 5 = 11
	// N o r w o o d
	// 5 · 9 5 · · 4 = 23 = 5
	// B a c o n
	// 2 · 3 · 5 = 10 = 1
	// Reduce: 17 = 8
}

func ExampleNameNumerology_HiddenPassions() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, false)
	fmt.Printf("%v (%v)", result.Name, result.HiddenPassions().Numbers)
	// Output: Kevin Norwood Bacon ([5])
}

func ExampleNameNumerology_KarmicLessons() {
	result := Name("Kevin Norwood Bacon", Pythagorean, []int{11, 22, 33}, false)
	fmt.Printf("%v (%v)", result.Name, result.KarmicLessons().Numbers)
	// Output: Kevin Norwood Bacon ([7 8])
}

func ExampleNameNumerology_Search() {
	opts := NameSearchOpts{
		Count:      10,
		Dictionary: "usa_census",
		Gender:     'M',
		Sort:       "common",
		Full:       []int{3, 5, 8},
		Database:   "sqlite://file::memory:?cache=shared",
	}
	results, _, err := Name("Abraham ? Lincoln", Chaldean, []int{11, 22, 33}, true).Search(opts)
	if err != nil {
		log.Fatal(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Name", "#"})
	for i, result := range results {
		table.Append([]string{fmt.Sprintf("%v", i+1), result.Name, fmt.Sprintf("%v", result.Full().Value)})
	}
	table.Render()
	// Output:
	// +----+-------------------------+---+
	// |    |          NAME           | # |
	// +----+-------------------------+---+
	// |  1 | Abraham James Lincoln   | 3 |
	// |  2 | Abraham David Lincoln   | 5 |
	// |  3 | Abraham William Lincoln | 8 |
	// |  4 | Abraham Daniel Lincoln  | 8 |
	// |  5 | Abraham Joshua Lincoln  | 3 |
	// |  6 | Abraham Anthony Lincoln | 8 |
	// |  7 | Abraham Andrew Lincoln  | 3 |
	// |  8 | Abraham Kevin Lincoln   | 8 |
	// |  9 | Abraham Steven Lincoln  | 8 |
	// | 10 | Abraham Jacob Lincoln   | 3 |
	// +----+-------------------------+---+
}

func TestNameNumerology_Search(t *testing.T) {
	type fields struct {
		Name           string
		NameOpts       *NameOpts
		NameSearchOpts *NameSearchOpts
		mask           *maskStruct
	}
	type args struct {
		opts NameSearchOpts
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults int
		wantOffset  int64
		wantErr     bool
	}{
		{"John ? Doe", fields{
			Name: "John ? Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   false,
			},
		}, args{NameSearchOpts{
			Count:          0,
			Offset:         5001,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "uncommon",
			Full:           []int{4},
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{7, -3},
			KarmicLessons:  []int{-6},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 3, 0, false},
		{"HiddenPassionsSearch", fields{
			Name: "John ? Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   false,
			},
		}, args{NameSearchOpts{
			Count:          25,
			Offset:         0,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "random",
			Full:           nil,
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{-3},
			KarmicLessons:  []int{},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 25, 25, false},
		{"MultipleHiddenPassionsSearch", fields{
			Name: "John ? Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   false,
			},
		}, args{NameSearchOpts{
			Count:          25,
			Offset:         0,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "random",
			Full:           nil,
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{1, 3},
			KarmicLessons:  []int{},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 25, 25, false},
		{"KarmicLessonsSearch", fields{
			Name: "John M? Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   true,
			},
		}, args{NameSearchOpts{
			Count:          25,
			Offset:         1,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "common",
			Full:           nil,
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{4},
			KarmicLessons:  []int{7, -3},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 25, 10108, false},
		{"Error John Doe", fields{
			Name: "John Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   false,
			},
		}, args{NameSearchOpts{
			Count:          0,
			Offset:         5001,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "uncommon",
			Full:           []int{4},
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{7, -3},
			KarmicLessons:  []int{-6},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 0, 0, true},
		{"Error John ?? Doe", fields{
			Name: "John ?? Doe",
			NameOpts: &NameOpts{
				NumberSystem:  Pythagorean,
				MasterNumbers: []int{11},
				ReduceWords:   false,
			},
		}, args{NameSearchOpts{
			Count:          0,
			Offset:         5001,
			Seed:           0,
			Dictionary:     "usa_census",
			Gender:         'M',
			Sort:           "uncommon",
			Full:           []int{4},
			Vowels:         nil,
			Consonants:     nil,
			HiddenPassions: []int{7, -3},
			KarmicLessons:  []int{-6},
			Database:       "sqlite://file::memory:?cache=shared",
		}}, 0, 0, true},
	}
	for _, tt := range tests {
		DB = nil
		t.Run(tt.name, func(t *testing.T) {
			n := NameNumerology{
				Name:           tt.fields.Name,
				NameOpts:       tt.fields.NameOpts,
				NameSearchOpts: tt.fields.NameSearchOpts,
				mask:           tt.fields.mask,
			}
			gotResults, gotOffset, err := n.Search(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(gotResults), tt.wantResults) {
				t.Errorf("Search() gotResults = %v, want %v", gotResults, tt.wantResults)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("Search() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
	DB = nil
}

func TestNameNumerology_Full(t *testing.T) {
	type fields struct {
		Name           string
		NameOpts       *NameOpts
		NameSearchOpts *NameSearchOpts
		mask           *maskStruct
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult int
	}{
		{"Jane Doe", fields{"Jane Doe", &NameOpts{Chaldean, []int{11}, true}, nil, nil},
			1},
		{"Lightning McQueen", fields{"Lightning McQueen", &NameOpts{Chaldean, []int{11}, true}, nil, nil},
			5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NameNumerology{
				Name:           tt.fields.Name,
				NameOpts:       tt.fields.NameOpts,
				NameSearchOpts: tt.fields.NameSearchOpts,
				mask:           tt.fields.mask,
			}
			if gotResult := n.Full(); !reflect.DeepEqual(gotResult.Value, tt.wantResult) {
				t.Errorf("Full() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestCalculateNames(t *testing.T) {
	type args struct {
		names         []string
		numberSystem  NumberSystem
		masterNumbers []int
		reduceWords   bool
	}
	tests := []struct {
		name        string
		args        args
		wantResults int
	}{
		{"CheckNames", args{
			names:         []string{"Jane Doe", "Lightning McQueen"},
			numberSystem:  Chaldean,
			masterNumbers: []int{11},
			reduceWords:   true,
		}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := Names(tt.args.names, tt.args.numberSystem, tt.args.masterNumbers, tt.args.reduceWords); !reflect.DeepEqual(len(gotResults), tt.wantResults) {
				t.Errorf("Names() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestNameNumerology_UnknownCharacters(t *testing.T) {
	type fields struct {
		Name           string
		NameOpts       *NameOpts
		NameSearchOpts *NameSearchOpts
		mask           *maskStruct
		counts         *map[int32]int
		unknowns       *unknownCharacters
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"UnknownCharacters1", fields{Name: "Pr*blem Name", NameOpts: &NameOpts{
			NumberSystem:  Pythagorean,
			MasterNumbers: []int{11, 22, 33},
			ReduceWords:   true,
		}}, []string{"*"}},
		{"UnknownCharacters2", fields{Name: "Problem N@me", NameOpts: &NameOpts{
			NumberSystem:  Pythagorean,
			MasterNumbers: []int{11, 22, 33},
			ReduceWords:   true,
		}}, []string{"@"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NameNumerology{
				Name:           tt.fields.Name,
				NameOpts:       tt.fields.NameOpts,
				NameSearchOpts: tt.fields.NameSearchOpts,
				mask:           tt.fields.mask,
				counts:         tt.fields.counts,
				unknowns:       tt.fields.unknowns,
			}
			if got := n.UnknownCharacters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnknownCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNumerologicalResult_Debug() {
	result := Name("George Washington", Pythagorean, []int{11, 22, 33}, true)
	fmt.Print(result.HeartsDesire().Debug())
	// Output:
	// G e o r g e
	// · 5 6 · · 5 = 16 = 7
	// W a s h i n g t o n
	// · 1 · · 9 · · · 6 · = 16 = 7
	// Reduce: 14 = 5
}
