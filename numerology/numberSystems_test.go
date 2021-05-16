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

func TestNumberSystem_MarshalJSON(t *testing.T) {
	type fields struct {
		Name          string
		NumberMapping map[int32]int
		ValidNumbers  []int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"Marshal Pythagorean Number System", fields(Pythagorean), []byte("\"pythagorean\""), false},
		{"Marshal Chaldean Number System", fields(Chaldean), []byte("\"chaldean\""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NumberSystem{
				Name:          tt.fields.Name,
				NumberMapping: tt.fields.NumberMapping,
				ValidNumbers:  tt.fields.ValidNumbers,
			}
			got, err := ns.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberSystem_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Name          string
		NumberMapping map[int32]int
		ValidNumbers  []int
	}
	type args struct {
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Unmarshal Pythagorean", fields(Pythagorean), args{[]byte("\"pythagorean\"")}, false},
		{"Unmarshal Chaldean", fields(Chaldean), args{[]byte("\"chaldean\"")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := &NumberSystem{
				Name:          tt.fields.Name,
				NumberMapping: tt.fields.NumberMapping,
				ValidNumbers:  tt.fields.ValidNumbers,
			}
			if err := ns.UnmarshalJSON(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
