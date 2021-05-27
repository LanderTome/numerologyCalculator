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
	"bufio"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/spkg/bom"
	"github.com/xo/dburl"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Requires Cgo when compiling and therefore doesn't work well for some target architectures.
	// "github.com/cloudquery/sqlite" // Cgo-free version of sqlite. Doesn't yet support all features of native sqlite but works fine here.
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// DB holds the database connection used of name searches. Gorm is used which means that only SQLite, MySQL,
// and PostgreSQL are supported out-of-the-box. Variable is exposed in case someone wanted to hack on other
// database solutions for a project.
var DB *gorm.DB
var dbLogger = logger.Default.LogMode(logger.Silent)

// uint8 (1-byte [0 to 255]) is used because it should be large enough for almost any situation and can save
// space in most databases (not sqlite, however).
//
// The longest English word 'pneumonoultramicroscopicsilicovolcanoconiosis'
// maxes out at only 218.
type precalculatedNumerology struct {
	Id                    int64 `gorm:"primaryKey"`
	Name                  string
	Gender                string `gorm:"index;type:varchar(1)"`
	PythagoreanFull       uint8  `gorm:"index"`
	PythagoreanVowels     uint8  `gorm:"index"`
	PythagoreanConsonants uint8  `gorm:"index"`
	ChaldeanFull          uint8  `gorm:"index"`
	ChaldeanVowels        uint8  `gorm:"index"`
	ChaldeanConsonants    uint8  `gorm:"index"`
	P1                    uint8  // Pythagorean count for number 1
	P2                    uint8  // Pythagorean count for number 2
	P3                    uint8  // Pythagorean count for number 3
	P4                    uint8  // Pythagorean count for number 4
	P5                    uint8  // Pythagorean count for number 5
	P6                    uint8  // Pythagorean count for number 6
	P7                    uint8  // Pythagorean count for number 7
	P8                    uint8  // Pythagorean count for number 8
	P9                    uint8  // Pythagorean count for number 9
	C1                    uint8  // Chaldean count for number 1
	C2                    uint8  // Chaldean count for number 2
	C3                    uint8  // Chaldean count for number 3
	C4                    uint8  // Chaldean count for number 4
	C5                    uint8  // Chaldean count for number 5
	C6                    uint8  // Chaldean count for number 6
	C7                    uint8  // Chaldean count for number 7
	C8                    uint8  // Chaldean count for number 8
}

// connectToDatabase parses a given DSN and establishes a connection to the database using Gorm. Only SQLite,
// PostgreSQL, and MySQL are currently supported.
func connectToDatabase(dsn string) error {
	if DB == nil {
		u, err := dburl.Parse(dsn)
		if err != nil {
			return errors.New("unable to connect to parse database connection string. dns=" + dsn)
		}
		switch u.OriginalScheme {
		case "sqlite":
			DB, err = gorm.Open(sqlite.Open(u.DSN), &gorm.Config{Logger: dbLogger})
		case "postgres":
			DB, err = gorm.Open(postgres.Open(u.DSN), &gorm.Config{Logger: dbLogger})
		case "mysql":
			DB, err = gorm.Open(mysql.Open(u.DSN), &gorm.Config{Logger: dbLogger})
		default:
			return errors.New("unsupported database. " + u.OriginalScheme)
		}
		return err
	}
	return nil
}

// setStructField allows us to add values to a struct by using a constructed string name for the field.
// This is used for the P1, P2, C1, C2, etc. columns of the database.
func setStructField(pcn *precalculatedNumerology, field string, value uint8) {
	v := reflect.ValueOf(pcn).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetUint(uint64(value))
	}
}

// namePopularity is a sortable slice used for ordering names before putting in the database.
type namePopularity []nameEntry

func (a namePopularity) Len() int           { return len(a) }
func (a namePopularity) Less(i, j int) bool { return a[i].Popularity < a[j].Popularity }
func (a namePopularity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type nameEntry struct {
	Name       string
	Gender     uint8
	Popularity int
}

// Extract all the names from the CSV files in the directory, merge the names, and create a sorted slice of results.
func extractNamesFromFiles(directory string) (namePopularity, error) {
	log.Printf("Extracting names from %v", directory)
	namePopularityMap := map[uint8]map[string]int{
		'M': {},
		'F': {},
	}
	files, err := fs.Glob(os.DirFS(directory), "*.csv")
	if err != nil {
		return namePopularity{}, errors.New("unable to scan directory")
	}
	totalFiles := len(files)
	bar := pb.Full.Start(totalFiles)
	for i, fn := range files {
		bar.Increment()
		// Weighting comes into play when there are multiple files imported together. Older names are weighted less than modern names.
		// The assumption is that each file is a year and earlier years are sorted in ascending order.
		weight := float64(i+1) / float64(totalFiles)

		file, err := os.Open(filepath.Join(directory, fn))
		if err != nil {
			log.Printf("Error opening file %v. --Skipping--", filepath.Join(directory, fn))
			continue
		}
		// Look through files and extract the data
		// bom.NewReader gets rid of UTF-8 byte order marks that can cause problems.
		fscanner := bufio.NewScanner(bom.NewReader(file))
		for fscanner.Scan() {
			cols := strings.Split(fscanner.Text(), ",")
			lname, gender, c := strings.TrimSpace(cols[0]), strings.TrimSpace(strings.ToUpper(cols[1]))[0], strings.TrimSpace(cols[2])
			/*
				----- This section is taken out because it only applies to some datasets. -----
				// USA Census file names are truncated at 15 characters. There are only a few dozen and they are all
				// combination names like ChristopherJohn and MariaDelRosario. Just ignore them.
				if len(lname) >= 15 {
					continue
				}
			*/
			origCount, err := strconv.Atoi(c)
			if err != nil {
				// If name count cannot be converted to a number then just skip the name.
				log.Printf("Unable to convert popularity from string to number. %v,%v,%v", lname, gender, c)
				continue
			}
			// Weight the popularity of the name.
			count := int(math.Ceil(float64(origCount) * weight))
			if c, ok := namePopularityMap[gender][lname]; ok {
				namePopularityMap[gender][lname] = c + count
			} else {
				namePopularityMap[gender][lname] = count
			}
		}
	}
	bar.Finish()

	// Put names in a slice of structs so we can sort it using the standard library.
	var Names namePopularity
	for gender, v := range namePopularityMap {
		for keyName, count := range v {
			Names = append(Names, nameEntry{
				Name:       keyName,
				Gender:     gender,
				Popularity: count,
			})
			// After name is put in struct, delete it from the map to conserve resources.
			delete(namePopularityMap[gender], keyName)
		}
	}
	// Sort in descending order of popularity.
	sort.Sort(sort.Reverse(Names))
	return Names, nil
}

// getAllDirectories gets all the directories in the baseDir folder.
func getAllDirectories(baseDir string) []string {
	directories := []string{}
	items, _ := os.ReadDir(baseDir)
	for _, item := range items {
		if item.IsDir() {
			directories = append(directories, item.Name())
		}
	}
	return directories
}

// Create the table if it is not already created.
func setupDatabaseTable(table string) error {
	if err := DB.Table(table).AutoMigrate(&precalculatedNumerology{}); err != nil {
		return err
	}

	// Because we are using the same struct for all our tables, Gorm uses the same name for
	// all the indexes. This causes and error since you can only have one index of each name.
	// Manually rename the indexes after GORM creates them.
	idxPrefix := "idx_precalculated_numerologies_"
	indexes := []string{
		"chaldean_consonants",
		"chaldean_full",
		"chaldean_vowels",
		"pythagorean_consonants",
		"pythagorean_full",
		"pythagorean_vowels",
		"gender",
	}
	for _, idx := range indexes {
		if err := DB.Table(table).Migrator().RenameIndex(
			&precalculatedNumerology{},
			idxPrefix+idx, table+"_idx_"+idx,
		); err != nil {
			return err
		}
		if err := DB.Table(table).Migrator().DropIndex(&precalculatedNumerology{}, idxPrefix+idx); err != nil {
			return err
		}
	}
	return nil
}

// CreateDatabase function creates and populates the database table with the pre-populated numerological
// calculations. The argument dsn is the connection string for the database that will utilized. The argument
// baseDir is the directory where the CSV files are stored that contain the names that will populate the
// database. Each folder in the baseDir becomes a table in the database. This allows for multiple name
// sources.
//
// The format for the CSV files is name, gender, popularity with no header. Gender is just a letter 'M' for
// male or 'F' female. Popularity is used to determine the sort order of the names. Each reoccurrence of the
// same name aggregates the popularity.
//
//  john,M,10000
//  sara,F,9000
//  jack,M,8000
func CreateDatabase(dsn string, baseDir string) error {
	directories := getAllDirectories(baseDir)

	log.Println("Connecting to database...")
	if err := connectToDatabase(dsn); err != nil {
		return errors.New("unable to connect to database. " + err.Error())
	}
	// Iterate over each of the folders and make a separate db table for each.
	for _, dir := range directories {
		// Create the table if it is not already created.
		if err := setupDatabaseTable(dir); err != nil {
			return err
		}

		// Make sure the table is empty. If it is not, then adding entries could mess it up.
		var count int64
		DB.Table(dir).Count(&count)
		if count > 0 {
			log.Printf("Table %v is not empty. Skipping.", dir)
			continue
		}

		names, err := extractNamesFromFiles(filepath.Join(baseDir, dir))
		if err != nil {
			return err
		}

		log.Printf("Populating database table %v", dir)
		bar := pb.Full.Start(len(names))
		for _, entry := range names {
			bar.Increment()
			// Precalculate the numerological numbers we want to put in the database.
			pythagorean := Name(entry.Name, Pythagorean, []int{}, false)
			chaldean := Name(entry.Name, Chaldean, []int{}, false)

			// Check for unacceptable characters
			if len(pythagorean.UnknownCharacters()) > 0 || len(chaldean.UnknownCharacters()) > 0 {
				log.Println(fmt.Sprintf("Skipping name with unacceptable characters: %v", entry.Name))
				continue
			}

			dbEntry := precalculatedNumerology{
				Name:                  entry.Name,
				Gender:                string(entry.Gender),
				PythagoreanFull:       uint8(pythagorean.Full().Breakdown[0].ReduceSteps[0]),
				PythagoreanVowels:     uint8(pythagorean.Vowels().Breakdown[0].ReduceSteps[0]),
				PythagoreanConsonants: uint8(pythagorean.Consonants().Breakdown[0].ReduceSteps[0]),
				ChaldeanFull:          uint8(chaldean.Full().Breakdown[0].ReduceSteps[0]),
				ChaldeanVowels:        uint8(chaldean.Vowels().Breakdown[0].ReduceSteps[0]),
				ChaldeanConsonants:    uint8(chaldean.Consonants().Breakdown[0].ReduceSteps[0]),
			}

			// Use reflection to populate these fields in the struct. There may be a better way to do this.
			pCounts := pythagorean.Counts()
			for k, v := range pCounts {
				field := fmt.Sprintf("P%v", k)
				setStructField(&dbEntry, field, uint8(v))
			}
			cCounts := chaldean.Counts()
			for k, v := range cCounts {
				field := fmt.Sprintf("C%v", k)
				setStructField(&dbEntry, field, uint8(v))
			}

			// Insert the record into the database.
			if err := DB.Table(dir).Create(&dbEntry).Error; err != nil {
				log.Println(fmt.Sprintf("unable to insert record into database: %v", dbEntry))
				continue
			}
		}
		bar.Finish()
	}
	log.Println("Vacuuming database to complete process.")
	// Vacuum the database to make sure any extra space is reclaimed.
	DB.Raw("VACUUM;")
	return nil
}
