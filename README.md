Proof of concept: [https://numerology.mysticalriver.com](https://numerology.mysticalriver.com)

## About the Project

Using either Pythagorean or Chaldean number systems, this calculator 
does various common numerological name and date calculations:

- Destiny / Expression
- Soul's Urge / Heart's Desire
- Personality
- Hidden Passions
- Karmic Lessons
- Life Path

The truly unique aspect of this package, however, is that it also 
provides methods to search for either names or dates that satisfy 
various numerological criteria. This is useful if one is trying 
to find a name for a baby, or a wedding date that has the desirable 
numerological properties. 

The way that it does this is using a precomputed table of names 
with corresponding numerological properties. This allows the table 
to be searched for specific names that satisfy given constraints. 
This method is surprisingly efficient. A database of 100,000 names 
with all the necessary pre-calculations is only 13.5 MB and searches 
take only hundredths of a second.

This calculator does not provide any numerological interpretations
as part of the package. Just the raw numbers.

See [godoc.org](https://godoc.org/github.com/LanderTome/numerologyCalculator/numerology)
for additional documentation.

## Getting Started

`go get "github.com/LanderTome/numerologyCalculator"`

### Important numerological concepts

When doing a name calculation there are three properties that need
to be specified that affect the calculation output:

1. Number System
2. Master Numbers
3. Reduce Words

The `Number System` is either one of the two most common systems
that are used to assign numbers to the letters of alphabet; Pythagorean
and Chaldean. 

The `Master Numbers` are numbers that are considered to have special 
numerological importance, and are treated in differently in 
calculations. The most commonly cited master numbers 11, 22, and 33. 
However, some consider ALL repeated numbers to be special, ex. 44, 55, 
66, etc. The master numbers that one wants to use can be customized. 

When `Reduce Words` is true, the calculation
process sums and reduces each name individually before summing and
reducing those results. When it is false, all the parts of a name
are summed and reduced together. The most common way is to sum and
reduce each name individually.

### Name calculations

Start a name calculation with the `Name` function.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
    numberSystem := numerology.Pythagorean
    masterNumbers := []int{11, 22, 33}
    reduceWords := true
    name := numerology.Name("John Doe", numberSystem, masterNumbers, reduceWords)
}
```

The returned struct is now the basis for calculating the numerological
properties of the given name.

Now the various numerological numbers can be calculated by calling 
the various methods.

```go
result := name.Full()
destiny := name.Destiny()           // alias for name.Full()
expression := name.Expression()     // alias for name.Full()

result := name.Vowels()
soulsUrge := name.SoulsUrge()       // alias for name.Vowels()
heartsDesire := name.HeartsDesire() // alias for name.Vowels()

result := name.Consonants()
personality := name.Personality()  // alias for name.Consonants()

hiddenPassions := name.HiddenPassions()
karmicLessons := name.KarmicLessons()
```

`Full, Vowels, and Consonants` are used instead of the common
numerological terms because there are sometimes multiple acceptable
terms, and it can make is easier to understand which calculation is
being done because you know what letters of the name are being used.

Aliases for the above methods are provided for convenience.

The returned `NumerologicalResult` struct contains information about 
the calculation. `Value` is the final reduced numerological number. 
`ReduceSteps` is a slice of numbers that are the numerological values 
at each stage of the final reducing. `Breakdown` contains the calculated values for each letter and each 
name that are used as the basis of the calculations.

`NumerologicalResult` has a method `Debug()` that returns a multiline
string that contains a simplistic breakdown of the calculation steps.
It can be used to visualize the conversion.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
	numberSystem := numerology.Pythagorean
	masterNumbers := []int{11, 22, 33}
	reduceWords := true
	name := numerology.Name("Janet Audrey Doe", numberSystem, masterNumbers, reduceWords)
	personality := name.Consonants()
	println(personality.Debug())
}

Output:
J a n e t
1 · 5 · 2 = 8
A u d r e y
· · 4 9 · 7 = 20 = 2
D o e
4 · · = 4
Reduce: 14 = 5
```

### Searching names

Where this calculator stands apart from other ones is in its ability
to search for names that match given numerological criteria. The idea
is to scan a database of names, and find ones that, when combined with
given names, have the numerological properties that are specified.
This task can be done with brute force methods, but several tricks 
as utilized here to filter out names that we know won't work before
actually querying the database.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
	numberSystem := numerology.Pythagorean
	masterNumbers := []int{11, 22, 33}
	reduceWords := true
	name := numerology.Name("Jane ? Doe", numberSystem, masterNumbers, reduceWords)
	searchOpts := numerology.NameSearchOpts{
		Count:          10,
		Offset:         0,
		Seed:           0,
		Database:       "sqlite://test_names.db",
		Dictionary:     "usa_census",
		Gender:         numerology.Female,
		Sort:           numerology.CommonSort,
		Full:           []int{1, 8, -13},
		Vowels:         []int{4, 5, 6},
		Consonants:     []int{5},
		HiddenPassions: nil,
		KarmicLessons:  nil,
	}
	results, offset, err := name.Search(searchOpts)
}
```

`NameSearchOpts` contains the necessary information to do the
search. `Count` is the number of results to return in each batch.
`Offset` is used to return the next batch of results. `Seed` is
used for generating the "random" sort results. The seed needs to
be known if batching through random results in order to keep the
results consistent. `Dictionary` is the name of the table in the 
database that will be searched. *The included test files are from 
the US Census.* `Gender` allows the names to be filtered by gender, 
so the results will be more male or female sounding. `Database` 
is the database connection string that connects to the database
that has the table to be searched.

Numerological properties: `Full`, `Vowels`, `Consonants`, 
`HiddenPassions`, and `KarmicLessons`. These are all slices of
integers that specify what criteria we want. Positive numbers
are numbers that are acceptable. Negative numbers are numbers 
that are to be excluded. A common use for negative numbers
would be to exclude, what are referred to as, Karmic Debt numbers 
(13, 14, 16, 17). By specifying `[]int{-13, -14, -16, -19}`, 
the search results will not include names with those numbers.

The return consists of a slice of `NameNumerology` results, 
an offset number that indicates where the next batch of results
should begin, and any errors messages that occur. To get the
next batch use the given offset number in the `NameSearchOpts`.

### Date calculations

Date calculation are done with the `Date` function.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
    dt := numerology.NewDate(2021, 1, 1)
    masterNumbers := []int{11, 22, 33}
    date := numerology.Date(dt, masterNumbers)
}
```

The available methods for date numerology are `Event()` and
`LifePath()`.

```go
event := name.Event()
lifePath := name.LifePath()
```

The only real difference between an event calculation and a 
life path calculation is that life path takes master numbers
into account. The returned `NumerologicalResult` similar to
that which is returned by the `Name` methods.

### Searching dates

The calculator is able to search for dates with specific
numerological criteria. Unlike the name search, there are 
no clever optimizations to this process. It simply loops 
through each date, and checks that it meets the criteria.
Usually, the search space for dates is only a few hundred
days, so this does not significantly impact performance.

Life Path numbers have to do with one's date of birth so
searching does not make a lot of sense because it is quite
difficult to control when someone is born. Searching dates
is better suited for finding ideal wedding or event dates.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
	dt := numerology.NewDate(2021, 1, 1)
	masterNumbers := []int{11, 22, 33}
	date := numerology.Date(dt, masterNumbers)

	searchOpts := numerology.DateSearchOpts{
		Count:         10,
		Offset:        0,
		Match:         []int{3, 6},
		MonthsForward: 12,
		Dow:           []int{numerology.Friday, numerology.Saturday},
		LifePath:      false,
	}
	results, offset := date.Search(searchOpts)
}
```

`DateSearchOpts` contains the necessary information to do the
search. `Count` is the number of results to return in each batch.
`Offset` is used to return the next batch of results. `Match`
contains the numerological numbers to find. Positive numbers indicate
numbers that are acceptable, and Negative numbers indicate
numbers that will be avoided. `MonthsForward` is the number of
months to search. It does not necessarily make sense to search
5 years out for a wedding that you want to have in 1 or 2 years.
`Dow` are days of the week that you want results on. This helps
if you are only interested in events that occur on particular
days; like weekends. `LifePath` indicates whether the dates
should take Master Numbers into account.

## Creating the Database

Before using the name search functionality, a database needs to
be created and populated.

```go
package main

import "github.com/LanderTome/numerologyCalculator/numerology"

func main() {
	dsn := "sqlite://file::memory:" // Example for testing
	namesDir := "test_names"
	if err := numerology.CreateDatabase(dsn, namesDir); err != nil {
		println(err.Error())
	}
}
```

### Connecting to database

DSN (data source name) is the connection string for the database
to use. This library uses [Gorm](https://gorm.io) for the database
connection. It currently supports MySQL, PostgreSQL, and SQLite, however
it is anticipated that SQLite would be the most common backend used.
Go package [xo/dburl](https://github.com/xo/dburl) is used for dsn
validation, so it is useful to look at the documentation there to see
acceptable connection string formats.

*The database connection is exposed in a package var `numerology.DB`.
It could theoretically be used to manually attach a custom Gorm backend,
although this feature is untested.*

### Source files

#### CSV structure

Each csv file needs to have three columns, however **NO HEADERS**.

The columns are *name*, *gender*, *popularity*.

```csv
mary,F,100
michael,M,98
john,M,95
susan,F,94
```

The popularity field is used so that names are sorted by how common
they are. Often census compilations include information like this.
It is useful to be able to sort by how common a name is because some
names that people give their children are simply bizarre, and it
becomes hard to sift through quality names without some control of
that.

*In order to save space in the database, the raw popularity numbers
are not stored. They are simply used to order the names before
putting them in the database. Once there, the inserted order of the
database is used to establish popularity.*

#### Directory layout

The database is populated by names coming from one or more CSV files
based on a particular directory layout.

```
baseDir\
  ├── table1\
  │   ├── file1.csv
  │   └── file2.csv
  └── table2\
      ├── file1.csv
      └── file2.csv
```

The folder names within the `baseDir` will be used as table names in
the database. The names in each csv file in the table directory will 
be iterated in ascending order and aggregated before being added to 
the database. 

*Multiple csv files can be used when inserting into the database.
The original source of names was a yearly compilation of USA Census
names. They consisted of the popularity of each name for each year
from past to present. The names are weighted so that current
popularity is valued higher than older popularity. As the CSV files
are iterated, the popularity is summed up with a weighted scale
based on how many files there are.*

## Other notes

### Is 'Y' a consonant or a vowel?

One problem that often comes up when doing name numerology is that 
sometimes the letter 'Y' is treated as a consonant, and sometimes 
it is treated like a vowel. How to handle these situations is, 
unfortunately, not so straightforward for an algorithm. The obvious 
answer is that 'Y' is a consonant when it is used as a consonant in 
a word, and a vowel when used as a vowel. But how does one know
the difference?

In the English language, every word can be broken down in to
syllables. The rule is that each syllable must have a vowel. That
means that when there is a syllable with no obvious vowels, but there
is a 'Y', the 'Y' is treated like a vowel.

Take a name like Sydney. The syllables would be `Syd` and `ney`. In
`Syd`, there is no vowel, so the 'Y' functions as a vowel. However, in
`ney`, the 'E' is the vowel, so the 'Y' defaults to a consonant.

This idea is pretty easy to grasp, but turns out it is hard to program
because there is no good algorithm for separating words (especially 
names) into syllables. You essentially need a dictionary of the English
language (and perhaps others, as well) because there are different 
ways to pronounce similarly written words. This is even more difficult 
for names because great liberties are often taken with their pronunciation.

One route is the use machine learning, and train it to deduce the
syllable structure of words. This is a lot of work and overhead for a 
solution that won't really even be 100% accurate anyway.

Some numerology calculators allow for special characters to signify 
when to treat 'Y' differently. Something like `S*yndey` would use
the * symbol to denote that the following 'Y' should be a vowel
instead of a consonant. This is difficult for this library because
in order to search it needs to precompute a table of names. In order
to do that, the value of 'Y' needs to be picked ahead of time and in
a consistent way. Arbitrarily choosing the value of a 'Y' could lead 
to inconsistent name results.

Luckily, the 'Y' situation does not come up very often, and there
are several rules-of-thumb that can be used to pick the right value
most of the time. This is what this library does, however, it does
mean that names with 'Y' ought to be looked at with extra scrutiny.

### Transliteration of non-ASCII characters

The Pythagorean and Chaldean numbers systems are defined for the
Latin alphabet A-Z. Of course, some languages have characters that
do not fit into this alphabet. The general consensus is to
transliterate those characters to Latin equivalent characters.
This is easy for some characters such as ä, é, or ñ. However, some
others are not so easy or obvious; ex. 张.

Most calculators just ignore these characters, but we can do better.
This library includes a Go package [mozillazg/go-unidecode](https://github.com/mozillazg/go-unidecode)
that is inspired by a Python port of a well-known Perl package `Text::Unidecode`.
The author of the Perl package has some excellent [information](https://interglacial.com/~sburke/tpj/as_html/tpj22.html)
on the difficulties of transliterating languages. It is worth the
read in order to understand the limitations of the automatic
conversion.

When in doubt, manual transliteration to ASCII characters is
recommended.

### Unknown Characters

Even with transliteration of alphabetic characters, it is inevitable
that names pop up that have symbols or emojis or other weird
characters. What should a calculator do in this situation? Symbols
do not have numerological value, so it seems like there are two
solutions: 1) Error on unknown characters or 2) Give them a zero value.

Returning an error does not seem quite appropriate because sometimes
you want to use strange characters in a name, and purposefully ignore
them. Commonly occurring characters of this sort are spaces, dashes,
and periods. It is not unusual to have a hyphenated name or to include
a period in something like 'Jr.'.

If someone had a name like 'Crazy! Smith', we could force them to drop
the exclamation point to do the calculation, or we could allow it and
just give it no value. That way it calculates with no extra work, and
the symbols are reflected as part of the name.

The problem with just zeroing out characters is that if someone does
not know that a character is being interpreted as an unknown
character then they won't be able to know if the numerological
calculation is what they expect. Perhaps someone expects that 
`$amuel $mith` is similar enough to `Samuel Smith` that it is
automatically converted.

The compromise used here is that unknown characters are processed
just fine, albeit with a zero value, and all unknown characters are
stored with the name and can be checked with the `UnknownCharacters()`
method.