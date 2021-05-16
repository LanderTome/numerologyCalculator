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
	"reflect"
	"testing"
)

func Test_maskConstructor(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want letterMask
	}{
		{"One letter Y", args{"Y"}, letterMask{true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskConstructor(tt.args.s); !reflect.DeepEqual(got.Vowels(), tt.want) {
				t.Errorf("maskConstructor() = %v, want %v", got.Vowels(), tt.want)
			}
		})
	}
}

func TestUnknownCharacters_MarshalJSON(t *testing.T) {
	type fields struct {
		Set mapset.Set
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"Marshal", fields{mapset.NewSetFromSlice([]interface{}{'-', '!'})},
			[]byte(`["!"]`), false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := unknownCharacters{
				Set: tt.fields.Set,
			}
			got, err := ec.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_unknownCharacters_Unacceptable(t *testing.T) {
	type fields struct {
		Set mapset.Set
	}
	tests := []struct {
		name         string
		fields       fields
		wantUnknowns unknownCharacters
	}{
		{"Unknown Characters", fields{mapset.NewSetWith(
			interface{}(' '), interface{}('.'), interface{}('-'), interface{}('$'), interface{}('😜')),
		}, unknownCharacters{mapset.NewSetWith(interface{}('$'), interface{}('😜'))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := unknownCharacters{
				Set: tt.fields.Set,
			}
			if gotUnknowns := ec.Unacceptable(); !reflect.DeepEqual(gotUnknowns, tt.wantUnknowns) {
				t.Errorf("Unacceptable() = %v, want %v", gotUnknowns, tt.wantUnknowns)
			}
		})
	}
}

func Test_reduceNumbers(t *testing.T) {
	type args struct {
		n             int
		masterNumbers []int
		steps         []int
	}
	tests := []struct {
		name            string
		args            args
		wantReduceSteps []int
	}{
		{"59338273", args{59338273, []int{11, 22, 33}, []int{}},
			[]int{59338273, 40, 4},
		},
		{"337432", args{337432, []int{11, 22, 33}, []int{}},
			[]int{337432, 22},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReduceSteps := reduceNumbers(tt.args.n, tt.args.masterNumbers, tt.args.steps); !reflect.DeepEqual(gotReduceSteps, tt.wantReduceSteps) {
				t.Errorf("reduceNumbers() = %v, want %v", gotReduceSteps, tt.wantReduceSteps)
			}
		})
	}
}
