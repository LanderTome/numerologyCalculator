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
	"strconv"
	"time"
)

// Days of the week values for use in date searches. Values are similar to time.Weekday.
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// DateNumerology stores required information to calculate the numerological values
// of dates.
type DateNumerology struct {
	// Date is either the date to use to calculate numerological values or the date
	// used to start searches.
	Date time.Time

	// Pointers are used so that these options can be shared amongst a number of
	// DateNumerology objects.
	*DateOpts
	*DateSearchOpts
}

// Event calculates the numerological number for a given date. Unlike LifePath calculations,
// Master Numbers are always reduced. This calculation is generally used for events like
// weddings or other special occasions.
func (d DateNumerology) Event() (event NumerologicalResult) {
	return calculateDate(d.Date, d.MasterNumbers, false)
}

// LifePath calculates the numerological number for a given date, but stops reducing at
// Master Numbers. This calculation is generally used with one's date of birth.
func (d DateNumerology) LifePath() (lifePath NumerologicalResult) {
	return calculateDate(d.Date, d.MasterNumbers, true)
}

// Search executes a forward looking search of dates to find ones that satisfy given
// numerological criteria. The argument opts contains the searching criteria. Offset in the
// output is the offset to be used to get the next batch of results using the same query. If
// offset is 0 then there are no more results.
func (d DateNumerology) Search(opts DateSearchOpts) (searchResults []DateNumerology, offset int64) {
	return dateSearch(d.Date, d.MasterNumbers, &opts)
}

// NewDate is a wrapper to easily create a time.Time variable without entering hour, min, sec, nsec, loc
// since they are not important for any date calculation. Note: Golang time.Time always has a timezone.
// UTC is used, and the hour used is midday (12:00). This should help avoid situations where time.Time
// gets converted to a new time zone accidentally so the date isn't inadvertently changed.
func NewDate(year, month, day int) (newDate time.Time) {
	return time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
}

// letterValuesFromNumbers is a simplified function to map the numbers from a date into a breakdown.
func letterValuesFromNumber(n int, masterNumbers []int) (letterValues []letterValue) {
	// Do not treat Master Numbers as a string of numbers
	if isMasterNumber(n, masterNumbers) {
		return append(letterValues, letterValue{Letter: strconv.Itoa(n), Value: n})
	}
	// Split number into array of individual numbers. ex 2020 -> 2,0,2,0
	ints := splitNumber(uint64(n))
	for _, i := range ints {
		letterValues = append(letterValues, letterValue{Letter: strconv.Itoa(i), Value: i})
	}
	return
}

// calculateDate outputs the numerological result of a date calculation.
func calculateDate(date time.Time, masterNumbers []int, lifePath bool) (result NumerologicalResult) {
	var totalValue int
	var breakdown []Breakdown

	for _, val := range []int{date.Year(), int(date.Month()), date.Day()} {
		reduceSteps := reduceNumbers(val, masterNumbers, []int{})
		calc := Breakdown{
			Value:        reduceSteps[len(reduceSteps)-1],
			ReduceSteps:  reduceSteps,
			LetterValues: letterValuesFromNumber(val, masterNumbers),
		}
		breakdown = append(breakdown, calc)
		totalValue += calc.Value
	}
	var reduceSteps []int
	if lifePath {
		reduceSteps = reduceNumbers(totalValue, masterNumbers, []int{})
	} else {
		reduceSteps = reduceNumbers(totalValue, []int{}, []int{})
	}
	return NumerologicalResult{
		Value:       reduceSteps[len(reduceSteps)-1],
		ReduceSteps: reduceSteps,
		Breakdown:   breakdown,
	}
}

// Date returns a DateNumerology object that can be used to calculate
// numerological values or search for dates with numerological significance.
func Date(date time.Time, masterNumbers []int) (result DateNumerology) {
	return DateNumerology{Date: date, DateOpts: &DateOpts{MasterNumbers: masterNumbers}}
}

// Dates returns a slice of DateNumerology objects that can be used to
// calculate numerological values or search for dates with numerological significance.
func Dates(dates []time.Time, masterNumbers []int) (results []DateNumerology) {
	opts := DateOpts{MasterNumbers: masterNumbers}
	results = []DateNumerology{}
	for _, date := range dates {
		results = append(results, DateNumerology{Date: date, DateOpts: &opts})
	}
	return results
}

// dateSearch is the function that does the actual date searching.
func dateSearch(startDate time.Time, masterNumbers []int, opts *DateSearchOpts) (searchResults []DateNumerology, offset int64) {
	endDate := startDate.AddDate(0, opts.MonthsForward, 0)
	for d := startDate.Add(time.Duration(24*opts.Offset) * time.Hour); d.Before(endDate); d = d.AddDate(0, 0, 1) {
		if len(opts.Dow) > 0 && !inIntSlice(int(d.Weekday()), opts.Dow) {
			continue
		}
		calc := calculateDate(d, masterNumbers, opts.LifePath)
		if inIntSlice(calc.Value, opts.Match) || len(opts.Match) == 0 {
			if len(searchResults) == opts.Count {
				offset = int64(d.Sub(startDate).Hours() / 24)
				break
			}
			searchResults = append(searchResults, DateNumerology{
				Date:           d,
				DateOpts:       &DateOpts{masterNumbers},
				DateSearchOpts: opts,
			})
		}
	}
	return searchResults, offset
}
