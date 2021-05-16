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
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	if err := CreateDatabase("sqlite://file::memory:?cache=shared", "test_names"); err != nil {
		panic("unable to create in-memory sqlite database for testing.")
	}
}

func TestConnectToDatabase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Connect Postgres", args{"postgres://user:pass@localhost/dbname"}, true},
		{"Connect MySql", args{"mysql://user:pass@localhost/dbname"}, true},
		{"Connect Sqlite", args{"sqlite://file::memory:?cache=shared"}, false},
		{"Unknown", args{"oracle://user:pass@somehost.com/sid"}, true},
		{"Nonsense", args{"nonsense://file"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DB = nil
			if err := connectToDatabase(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("connectToDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	DB = nil
}

func TestCreateDatabase(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	err := os.Chdir(filepath.Join(filepath.Dir(b), ".."))
	if err != nil {
		panic(err)
	}
	type args struct {
		dsn     string
		baseDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"CannotWriteToDB", args{"sqlite://file::memory:?mode=ro", "test_names"}, true},
		// {"CoverageCheck", args{"sqlite://file::memory:?cache=shared", "test_names"}, false},
		{"DatabaseNotEmpty", args{"sqlite://file::memory:?cache=shared", "test_names"}, false},
		{"UnableToConnect", args{"postgres://errorurl:5432", "test_names"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DB = nil
			if err := CreateDatabase(tt.args.dsn, tt.args.baseDir); (err != nil) != tt.wantErr {
				t.Errorf("CreateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	DB = nil
}
