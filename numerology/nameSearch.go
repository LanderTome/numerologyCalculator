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
	"gorm.io/gorm"
	"math/rand"
	"sort"
	"strings"
	"time"
	"unicode"
)

var largestNameValueInTable = map[string]int{}

// Constants that represent the name sorting methods.
const (
	CommonSort   = "common"
	UncommonSort = "uncommon"
	RandomSort   = "random"
)

type queryLookup struct {
	MinimumSearchNumber int
	MaximumSearchNumber int
	TargetNumbers       []int
	ColumnName          string
	MasterNumbers       []int
	ReduceWords         bool
}

func init() {
	// Initialize the random module.
	rand.Seed(time.Now().UnixNano())
}

// generateLookupNums finds the values to look for in the database that will result in names with the correct
// numerological properties. The key is to find unreduced numbers that will reduce down to the value that we
// want. Since all the values in the database are stored unreduced we can then find names that will work.
func generateLookupNums(minSearchNumber int, maxSearchNumber int, numerologyNums []int, masterNumbers []int, reduceWords bool) []int {
	var nums []int
	// Split up positive and negative numbers into separate lists and make negative numbers positive for comparison.
	positiveNums := []int{}
	negativeNums := []int{}
	for _, n := range numerologyNums {
		if n >= 0 {
			positiveNums = append(positiveNums, n)
		} else {
			negativeNums = append(negativeNums, -n)
		}
	}
	// Iterate through possible numbers looking for matches.
	for i := 1; minSearchNumber+i <= maxSearchNumber; i++ {
		// i needs to be in reduced form for proper calculation.
		reducedI := reduceNumbers(i, masterNumbers, []int{})
		var idx int
		if reduceWords {
			idx = len(reducedI) - 1
		} else {
			idx = 0
		}
		// Create the reduced value of this hypothetical name.
		validNum := reduceNumbers(minSearchNumber+reducedI[idx], masterNumbers, []int{})
		// Check if the validNum satisfies our numerological criteria.
		var noMatch bool
		// If there are positive numbers then noMatch is only false when we specifically find a number.
		// The alternative is that any number is acceptable as long as it isn't a negative number.
		if len(positiveNums) > 0 {
			noMatch = true
		}
		for _, v := range validNum {
			if inIntSlice(v, negativeNums) {
				noMatch = true
				break
			}
			if inIntSlice(v, positiveNums) {
				noMatch = false
				break
			}
		}
		if !noMatch {
			nums = append(nums, i)
		}
	}
	return nums
}

func addQueryLookup(query *gorm.DB, q queryLookup) {
	if len(q.TargetNumbers) == 0 {
		return
	}
	lookupNums := generateLookupNums(q.MinimumSearchNumber, q.MaximumSearchNumber, q.TargetNumbers, q.MasterNumbers, q.ReduceWords)
	if len(lookupNums) > 0 {
		query = query.Where(q.ColumnName+" IN ?", lookupNums)
	} else {
		query = query.Where("FALSE")
	}
	return
}

func addKarmicDebtLookup(query *gorm.DB, q queryLookup) {
	if len(q.TargetNumbers) == 0 {
		return
	}
	lookupNums := generateLookupNums(q.MinimumSearchNumber, q.MaximumSearchNumber, q.TargetNumbers, q.MasterNumbers, q.ReduceWords)
	if len(lookupNums) > 0 {
		query = query.Where(q.ColumnName+" NOT IN ?", lookupNums)
	}
	return
}

func addQueryGender(query *gorm.DB, gender rune) {
	// Add gender parameter if specified. Otherwise include all genders.
	if unicode.ToUpper(gender) == 'M' || unicode.ToUpper(gender) == 'F' {
		query = query.Where("gender = ?", string(unicode.ToUpper(gender)))
	}
	return
}

// Order the results based on the selected option
func addQuerySort(query *gorm.DB, sort string, offset int, seed int64) {
	switch strings.ToLower(sort) {
	case UncommonSort:
		// Todo: Do not hardcode uncommon skip value
		// Skip down a ways so the names are less common.
		skip := 5000
		if offset > skip {
			skip = offset
		}
		query = query.Where("id >= ?", skip).Order("id asc")
	case RandomSort:
		// Randomizing order with seed. https://stackoverflow.com/a/24511461
		randSource := rand.NewSource(seed)
		newRand := rand.New(randSource)
		multiplier := newRand.Float64()
		query = query.Order(fmt.Sprintf("(substr(id * %v, length(id) + 2))", multiplier))
		// Offset for random uses regular offset function of db because there is no easier way to skip results.
		query = query.Offset(offset)
	default: // Default is a catchall for "common"
		query = query.Where("id >= ?", offset).Order("id asc")
	}
	return
}

// A surprisingly difficult function to find the criteria to match the Hidden Passions criteria. Some of the
// difficulty comes from the fact that negative numbers acting as exclusionary complicates the analysis.
func queryHiddenPassions(query *gorm.DB, name string, hiddenPassions []int, numberSystem NumberSystem) {
	if len(hiddenPassions) == 0 {
		return
	}
	type BuildQuery struct {
		Left       string
		Comparator string
		Right      interface{}
		Target     int
	}
	// Prefix needed to pick column when querying the database.
	prefix := string(strings.ToLower(numberSystem.Name)[0])
	// RunCount all the numerological numbers.
	currentCount, maxCount, _ := countNumerologicalNumbers(name, numberSystem)
	// Sort the numbers from largest to smallest. This is important because we want negative numbers last.
	sort.Sort(sort.Reverse(sort.IntSlice(hiddenPassions)))
	// Select the largest number as the prime number that calculations are based off of.
	prime := hiddenPassions[0]

	// If the largest number is negative then all the numbers are negative. Process differently.
	if prime < 0 {
		for _, hp := range hiddenPassions {
			var where *gorm.DB
			hpCount, _ := currentCount[int32(-hp)]
			// Loop through all valid numbers and find any one that is bigger than the negative prime we do not want.
			for _, i := range numberSystem.ValidNumbers {
				if inIntSlice(-i, hiddenPassions) {
					continue
				}
				iCount, _ := currentCount[int32(i)]
				leftCol := fmt.Sprintf("%v%d", prefix, i)
				rightCol := fmt.Sprintf("%v%d", prefix, -hp)

				if iCount-hpCount == 0 { // If modifier is '0' then no reason for unnecessary addition in query.
					if where == nil { // If our Where clause hasn't been initialized yet then do it now.
						where = DB.Where(fmt.Sprintf("%v > %v", leftCol, rightCol))
					} else { // if already initialized then add an Or clause
						where.Or(fmt.Sprintf("%v > %v", leftCol, rightCol))
					}
				} else {
					if where == nil { // If our Where clause hasn't been initialized yet then do it now.
						where = DB.Where(fmt.Sprintf("%v > %v - ?", leftCol, rightCol), iCount-hpCount)
					} else { // if already initialized then add an Or clause
						where.Or(fmt.Sprintf("%v > %v - ?", leftCol, rightCol), iCount-hpCount)
					}
				}
			}
			// Add grouped query to main query
			query = query.Where(where)
		}
		return
	}

	// Get the count of the prime number.
	primeCount, _ := currentCount[int32(prime)]

	// Loop through all valid number system numbers and build the initial query values.
	buildQuery := map[int]BuildQuery{}
	for _, i := range numberSystem.ValidNumbers {
		if i == prime {
			// Column name is prefix + number. ex p1, p2
			col := fmt.Sprintf("%v%d", prefix, i)
			// The prime number count needs to be as large or larger than the current largest count.
			buildQuery[i] = BuildQuery{col, ">=", maxCount - primeCount, 0}
		} else {
			// Column name is prefix + number. ex p1, p2
			leftCol := fmt.Sprintf("%v%d", prefix, i)
			rightCol := fmt.Sprintf("%v%d", prefix, prime)
			// Get count of number
			count, _ := currentCount[int32(i)]
			// In order to satisfy search condition this number needs to be less than or equal to a target
			// that is equal to the difference between the count of this number and the prime number
			targetCount := primeCount - count
			buildQuery[i] = BuildQuery{leftCol, "<=", rightCol, targetCount}
		}
	}
	// Loop through numbers to adjust query. Skip the first number because it was handled as the prime number.
	for _, p := range hiddenPassions[1:] {
		if p >= 0 {
			bq := buildQuery[p]
			bq.Comparator = "="
			buildQuery[p] = bq
		} else {
			bq := buildQuery[-p]
			bq.Comparator = "<"
			buildQuery[-p] = bq
		}
	}
	// Now that we have sufficiently modified the queries. Add them to the main query.
	for _, v := range buildQuery {
		if v.Target != 0 {
			query = query.Where(fmt.Sprintf("%v %v %v + ?", v.Left, v.Comparator, v.Right), v.Target)
		} else {
			query = query.Where(fmt.Sprintf("%v %v %v", v.Left, v.Comparator, v.Right))
		}
	}
}

func queryKarmicLessons(query *gorm.DB, name string, n []int, numberSystem NumberSystem) {
	if len(n) == 0 {
		return
	}
	// Prefix needed to pick column when querying the database.
	prefix := string(strings.ToLower(numberSystem.Name)[0])
	// RunCount all the numerological numbers.
	currentCount, _, _ := countNumerologicalNumbers(name, numberSystem)
	for _, i := range n {
		if i < 0 {
			col := fmt.Sprintf("%v%d", prefix, -i)
			count, _ := currentCount[int32(-i)]
			query = query.Where(col+" + ? > 0", count)
		} else {
			col := fmt.Sprintf("%v%d", prefix, i)
			count, _ := currentCount[int32(i)]
			// This will either be 0 or less than 0.
			// Less than 0 is impossible match and therefore acts to excludes these rows.
			query = query.Where(col+" <= ?", 0-count)
		}
	}
}

// This function does all the heavy lifting for searching names.
func nameSearch(n string, numberSystem NumberSystem, masterNumbers []int, reduceWords bool, opts NameSearchOpts) (results []NameNumerology, offset int64, err error) {
	requiredOpts := NameOpts{
		NumberSystem:  numberSystem,
		MasterNumbers: masterNumbers,
		ReduceWords:   reduceWords,
	}
	searchOpts := NameSearchOpts{
		Count:          opts.Count,
		Offset:         opts.Offset,
		Seed:           opts.Seed,
		Dictionary:     opts.Dictionary,
		Gender:         opts.Gender,
		Sort:           opts.Sort,
		Full:           opts.Full,
		Vowels:         opts.Vowels,
		Consonants:     opts.Consonants,
		HiddenPassions: opts.HiddenPassions,
		KarmicLessons:  opts.KarmicLessons,
		Database:       opts.Database,
	}

	if DB == nil {
		err = connectToDatabase(opts.Database)
		if err != nil {
			return []NameNumerology{}, 0, err
		}
	}

	table := strings.ToLower(opts.Dictionary)
	// Make sure table is valid and has names.
	var count int64
	DB.Table(table).Count(&count)
	if count == 0 {
		return []NameNumerology{}, 0, fmt.Errorf("database table %v is empty", opts.Dictionary)
	}

	// Lowercase the number system name for use later.
	nsName := strings.ToLower(numberSystem.Name)

	splitNames := strings.Split(n, " ")
	var nameToSearch string
	for i, n := range splitNames {
		if numberOfQuestionMarks(n) > 0 {
			// Save the name for use later and replace it with just a ?. This will help the new name construction later.
			nameToSearch = n
			splitNames[i] = "?"
		}
	}
	reconstructedName := strings.Join(splitNames, " ")
	nonSearchNameResults := NameNumerology{reconstructedName, &requiredOpts, nil, nil, nil, nil}

	// Begin constructing the query
	if opts.Count == 0 {
		opts.Count = 25
	}
	query := DB.Table(table).Select("min(id) as id, name").Group("name")
	// Increase query count by 1 because we will use the last result as an indication that there are more results
	// that can be paged through. The extra result will be dropped in the return.
	query = query.Limit(opts.Count + 1)

	// If there are letters around the ? then we need to do a LIKE search.
	if len(nameToSearch) > 1 {
		query = query.Where("LOWER(name) LIKE ?", strings.Replace(strings.ToLower(nameToSearch), "?", "%", -1))
	}

	addQuerySort(query, opts.Sort, opts.Offset, opts.Seed)
	addQueryGender(query, rune(opts.Gender))

	queryHiddenPassions(query, reconstructedName, opts.HiddenPassions, numberSystem)
	queryKarmicLessons(query, reconstructedName, opts.KarmicLessons, numberSystem)

	// Get the largest name value from the database and cache it so we don't have to look it up again.
	var largestNameValueInDb int
	var ok bool
	if largestNameValueInDb, ok = largestNameValueInTable[table]; !ok {
		if err := DB.Table(table).Select("max(max(pythagorean_full), max(chaldean_full)) as largest").Find(&largestNameValueInDb).Error; err != nil {
			largestNameValueInDb = 100
		}
		largestNameValueInTable[table] = largestNameValueInDb
	}

	// Lookup for Full numbers
	addQueryLookup(query, queryLookup{
		MinimumSearchNumber: nonSearchNameResults.Full().ReduceSteps[0],
		MaximumSearchNumber: nonSearchNameResults.Full().ReduceSteps[0] + largestNameValueInDb,
		TargetNumbers:       opts.Full,
		ColumnName:          nsName + "_full",
		MasterNumbers:       masterNumbers,
		ReduceWords:         reduceWords,
	})
	// Lookup for Vowel numbers
	addQueryLookup(query, queryLookup{
		MinimumSearchNumber: nonSearchNameResults.Vowels().ReduceSteps[0],
		MaximumSearchNumber: nonSearchNameResults.Vowels().ReduceSteps[0] + largestNameValueInDb,
		TargetNumbers:       opts.Vowels,
		ColumnName:          nsName + "_vowels",
		MasterNumbers:       masterNumbers,
		ReduceWords:         reduceWords,
	})
	// Lookup for Consonant numbers
	addQueryLookup(query, queryLookup{
		MinimumSearchNumber: nonSearchNameResults.Consonants().ReduceSteps[0],
		MaximumSearchNumber: nonSearchNameResults.Consonants().ReduceSteps[0] + largestNameValueInDb,
		TargetNumbers:       opts.Consonants,
		ColumnName:          nsName + "_consonants",
		MasterNumbers:       masterNumbers,
		ReduceWords:         reduceWords,
	})

	var selectedNames []precalculatedNumerology
	func() {
		query.Find(&selectedNames)
	}()

	results = []NameNumerology{}
	for i, r := range selectedNames {
		// In order to find out if there are more results, we search for 1 extra and make a note of its id
		// to use as an offset later. Then exclude the final result from what is returned.
		if len(selectedNames) <= opts.Count || i < len(selectedNames)-1 {
			// Use reconstructedName because we want to replace the whole ? name, and not accidentally include additional letters.
			// John Da? Doe would come out as John DaDavid Doe. reconstructedName avoids this.
			newName := strings.Replace(reconstructedName, "?", r.Name, 1)
			// Calculate the numerology results using the new full name.
			results = append(results, NameNumerology{newName, &requiredOpts, &searchOpts, nil, nil, nil})
		} else {
			offset = r.Id
		}
	}
	// If sort is random then we need to derive the offset a different way.
	if offset > 0 && strings.ToLower(opts.Sort) == RandomSort {
		offset = int64(opts.Offset + opts.Count)
	}
	return results, offset, nil
}
