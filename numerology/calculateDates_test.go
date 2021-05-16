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
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCalculateDate(t *testing.T) {
	type args struct {
		date          time.Time
		masterNumbers []int
	}
	tests := []struct {
		name       string
		args       args
		wantResult DateNumerology
	}{
		{"DateNumerology", args{NewDate(2021, 1, 1), []int{11, 22, 33}},
			DateNumerology{NewDate(2021, 1, 1), &DateOpts{[]int{11, 22, 33}}, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Date(tt.args.date, tt.args.masterNumbers); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Date() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestCalculateDates(t *testing.T) {
	type args struct {
		dates         []time.Time
		masterNumbers []int
	}
	tests := []struct {
		name        string
		args        args
		wantResults []DateNumerology
	}{
		{"DatesNumerology", args{[]time.Time{NewDate(2021, 2, 1)}, []int{11, 22}},
			[]DateNumerology{{NewDate(2021, 2, 1), &DateOpts{[]int{11, 22}}, nil}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResults := Dates(tt.args.dates, tt.args.masterNumbers); !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("Dates() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestDateNumerology_Event(t *testing.T) {
	type fields struct {
		Date           time.Time
		DateOpts       *DateOpts
		DateSearchOpts *DateSearchOpts
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult int
	}{
		{"1970-01-01", fields{NewDate(1970, 1, 1), &DateOpts{[]int{11, 22, 33}}, nil},
			1},
		{"1970-01-02", fields{NewDate(1970, 1, 2), &DateOpts{[]int{11, 22, 33}}, nil},
			2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DateNumerology{
				Date:           tt.fields.Date,
				DateOpts:       tt.fields.DateOpts,
				DateSearchOpts: tt.fields.DateSearchOpts,
			}
			if gotResult := d.Event().Value; !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Event() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDateNumerology_LifePath(t *testing.T) {
	type fields struct {
		Date           time.Time
		DateOpts       *DateOpts
		DateSearchOpts *DateSearchOpts
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult int
	}{
		{"2010-05-04", fields{NewDate(2010, 5, 4), &DateOpts{[]int{11, 22, 33}}, nil},
			3},
		{"1970-01-02", fields{NewDate(1970, 1, 2), &DateOpts{[]int{11, 22, 33}}, nil},
			11},
		{"1993-11-11", fields{NewDate(1993, 11, 11), &DateOpts{[]int{11, 22, 33, 44}}, nil},
			44},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DateNumerology{
				Date:           tt.fields.Date,
				DateOpts:       tt.fields.DateOpts,
				DateSearchOpts: tt.fields.DateSearchOpts,
			}
			if gotResult := d.LifePath().Value; !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("LifePath() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDateNumerology_Search(t *testing.T) {
	type fields struct {
		Date           time.Time
		DateOpts       *DateOpts
		DateSearchOpts *DateSearchOpts
	}
	type args struct {
		opts DateSearchOpts
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult int
		wantOffset int64
	}{
		{"1970-01-02",
			fields{NewDate(1970, 1, 1), &DateOpts{[]int{11, 22, 33}}, nil},
			args{DateSearchOpts{
				Count:         100,
				Offset:        0,
				Match:         []int{11},
				MonthsForward: 48,
				Dow:           []int{0, 5, 6},
				LifePath:      true,
			}},
			39, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DateNumerology{
				Date:           tt.fields.Date,
				DateOpts:       tt.fields.DateOpts,
				DateSearchOpts: tt.fields.DateSearchOpts,
			}
			gotResult, gotOffset := d.Search(tt.args.opts)
			if !reflect.DeepEqual(len(gotResult), tt.wantResult) {
				t.Errorf("Search() gotResult = %v, want %v", len(gotResult), tt.wantResult)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("Search() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func TestNewDate(t *testing.T) {
	type args struct {
		year  int
		month int
		day   int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"2010-09-18", args{2010, 9, 18}, time.Date(2010, 9, 18, 12, 0, 0, 0, time.UTC)},
		{"1934-04-29", args{1934, 4, 29}, time.Date(1934, 4, 29, 12, 0, 0, 0, time.UTC)},
		{"2145-12-08", args{2145, 12, 8}, time.Date(2145, 12, 8, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDate(tt.args.year, tt.args.month, tt.args.day); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_letterValuesFromNumber(t *testing.T) {
	type args struct {
		n             int
		masterNumbers []int
	}
	tests := []struct {
		name string
		args args
		want []letterValue
	}{
		{"MasterNumberTest", args{11, []int{11, 22, 33}},
			[]letterValue{{"11", 11}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letterValuesFromNumber(tt.args.n, tt.args.masterNumbers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("letterValuesFromNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleDateNumerology_Event() {
	date := NewDate(2021, 1, 5)
	result := Date(date, []int{11, 22, 33}).Event()
	fmt.Printf("%v - %v", date.Format("Monday, January 2, 2006"), result.Value)
	// Output: Tuesday, January 5, 2021 - 2
}

func ExampleDateNumerology_LifePath() {
	date := NewDate(2021, 1, 5)
	result := Date(date, []int{11, 22, 33}).LifePath()
	fmt.Printf("%v - %v", date.Format("Monday, January 2, 2006"), result.Value)
	//  Output: Tuesday, January 5, 2021 - 11
}

func ExampleDateNumerology_Search() {
	opts := DateSearchOpts{
		Count:         5,
		Match:         []int{1},
		MonthsForward: 12,
		Dow:           []int{int(time.Monday), int(time.Friday), int(time.Saturday)},
		LifePath:      false,
	}
	date := NewDate(2021, 1, 1)
	results, _ := Date(date, []int{11, 22, 33}).Search(opts)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "Date", "#"})
	for i, result := range results {
		table.Append([]string{fmt.Sprintf("%v", i+1), result.Date.Format("Monday, January 2, 2006"), fmt.Sprintf("%v", result.Event().Value)})
	}
	table.Render()
	// Output:
	// +---+---------------------------+---+
	// |   |           DATE            | # |
	// +---+---------------------------+---+
	// | 1 | Monday, January 4, 2021   | 1 |
	// | 2 | Friday, January 22, 2021  | 1 |
	// | 3 | Friday, February 12, 2021 | 1 |
	// | 4 | Saturday, March 20, 2021  | 1 |
	// | 5 | Monday, March 29, 2021    | 1 |
	// +---+---------------------------+---+
}
