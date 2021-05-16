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

func Test_inIntSlice(t *testing.T) {
	type args struct {
		v     int
		slice []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Found", args{16, []int{3, 14, 7, 16, 0}}, true},
		{"Not Found", args{16, []int{3, 14, 7, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inIntSlice(tt.args.v, tt.args.slice); got != tt.want {
				t.Errorf("inIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isMasterNumber(t *testing.T) {
	type args struct {
		v             int
		masterNumbers []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Master Number", args{33, []int{11, 22, 33}}, true},
		{"Not Master Number", args{4, []int{11, 22, 33}}, false},
		{"Empty Master Number List", args{1, []int{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMasterNumber(tt.args.v, tt.args.masterNumbers); got != tt.want {
				t.Errorf("isMasterNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_numberOfQuestionMarks(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"3", args{"Bla? Bl?h B?ah"}, 3},
		{"0", args{"No question marks."}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numberOfQuestionMarks(tt.args.s); got != tt.want {
				t.Errorf("numberOfQuestionMarks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDupicateSpaces(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Remove spaces", args{"Too   much   space"}, "Too much space"},
		{"No removal", args{"Too much space"}, "Too much space"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicateSpaces(tt.args.s); got != tt.want {
				t.Errorf("removeDuplicateSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitStringOfNumbers(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"473809", args{473809}, []int{4, 7, 3, 8, 0, 9}},
		{"1", args{1}, []int{1}},
		{"201", args{201}, []int{2, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitNumber(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNumberSystem(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want NumberSystem
	}{
		{"Chaldean", args{"chaldeAn"}, Chaldean},
		{"Pythagorean", args{"PythaGorean"}, Pythagorean},
		{"Default", args{"BJeilsd"}, NumberSystem{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetNumberSystem(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNumberSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountNumbers(t *testing.T) {
	type args struct {
		name         string
		numberSystem NumberSystem
	}
	tests := []struct {
		name  string
		args  args
		want  map[int32]int
		want1 int
	}{
		{"Digits", args{"11233344567778990", Pythagorean},
			map[int32]int{1: 2, 2: 1, 3: 3, 4: 2, 5: 1, 6: 1, 7: 3, 8: 1, 9: 2}, 3},
		{"DigitsWithSomeMissing", args{"1123335677780", Pythagorean},
			map[int32]int{1: 2, 2: 1, 3: 3, 4: 0, 5: 1, 6: 1, 7: 3, 8: 1, 9: 0}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := countNumerologicalNumbers(tt.args.name, tt.args.numberSystem)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countNumerologicalNumbers() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("countNumerologicalNumbers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestToAscii(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"AsciiTest1", args{" 北亰 "}, "Bei Jing"},
		{"AsciiTest2", args{"\tétude  "}, "etude"},
		{"AsciiTest3", args{"Hello    world"}, "Hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAscii(tt.args.s); got != tt.want {
				t.Errorf("ToAscii() = %v, want %v", got, tt.want)
			}
		})
	}
}
