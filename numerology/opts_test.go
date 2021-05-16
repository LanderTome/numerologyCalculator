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

func TestGender_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		g       Gender
		want    []byte
		wantErr bool
	}{
		{"Marshal Gender M", 'M', []byte("\"M\""), false},
		{"Marshal Gender F", 'F', []byte("\"F\""), false},
		{"Marshal Gender B", 'B', []byte("\"B\""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.MarshalJSON()
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

func TestGender_UnmarshalJSON(t *testing.T) {
	type args struct {
		value []byte
	}
	tests := []struct {
		name    string
		g       Gender
		args    args
		wantErr bool
	}{
		{"Unmarshal Gender", 'F', args{[]byte("\"F\"")}, false},
		{"Unmarshal Male", 'M', args{[]byte("\"Male\"")}, false},
		{"Unmarshal Female", 'F', args{[]byte("\"Female\"")}, false},
		{"Unmarshal Both", 'B', args{[]byte("\"B\"")}, false},
		{"Unmarshal Unknown", 'B', args{[]byte("\"Something\"")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.UnmarshalJSON(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
