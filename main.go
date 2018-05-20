package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/c-Brooks/find-open-restaurants/hours"
	"github.com/c-Brooks/find-open-restaurants/parser"
	"github.com/gogap/logrus"
)

// Given the attached CSV data file, write a function findOpenRestaurants(csvFilename, searchDatetime)
// which take as parameters a filename and a DateTime object and returns a list of restaurant names which are open
// on that date and time.
//
// Assumptions:
// * If a day of the week is not listed, the restaurant is closed on that day
// * All times are local — don’t worry about timezone-awareness
// * The CSV file will be well-formed

type restaurant struct {
	name  string
	hours hours.Hours
}

func main() {
	t := time.Now()

	// t, err := time.Parse(time.ANSIC, "Sat May 21 01:00:00 2018")
	// if err != nil {
	// 	panic(err)
	// }

	openRestaurants := findOpenRestaurants("restaurants.csv", t)

	logrus.Infof("found %d restaurants open for %s:", len(openRestaurants), t.Format(time.RFC822))
	for i, r := range openRestaurants {
		logrus.Infof("%3d | %s", (i + 1), r.name)
	}

}

func findOpenRestaurants(csvFilename string, searchDatetime time.Time) []restaurant {
	fileLoc := "./data/" + csvFilename
	logrus.Infof("reading from file %s", fileLoc)
	rows, err := readCSV(fileLoc)
	if err != nil {
		logrus.Fatalf("could not parse csv: %v", err)
	}
	restaurants, err := parseCSV(rows)
	if err != nil {
		logrus.Fatal(err)
	}

	openRestaurants := make([]restaurant, 0)
	day := time.Weekday(searchDatetime.Weekday())

	for _, r := range restaurants {

		// assumption: no hours entry => the restaurant is closed that day
		tr, ok := r.hours[day]
		if !ok {
			continue
		}

		if tr.Contains(hours.TimeOfDay{searchDatetime}) {
			openRestaurants = append(openRestaurants, r)
		}
	}
	return openRestaurants
}

// parseCSV parses the csv in the data/ directory
// and returns a 2-D slice
func readCSV(filename string) ([][]string, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read from file %s: %v", filename, err)
	}

	r := csv.NewReader(strings.NewReader(string(dat)))
	data, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not parse CSV: %v", err)
	}
	return data, nil
}

// parseCSV parses the csv in the data/ directory
// and returns a 2-D slice
func parseCSV(csvRows [][]string) ([]restaurant, error) {
	restaurants := make([]restaurant, len(csvRows))
	c := make(chan restaurant)
	defer close(c)
	var wg sync.WaitGroup
	wg.Add(len(csvRows))

	// transform csv rows into slice of restaurants
	for i, row := range csvRows {

		// note: seems kinda overkill for this small dataset but I wanna show off
		// ¯\_(ツ)_/¯
		go func(i int, name, rawHours string) {
			defer wg.Done()
			rest, err := toRestaurant(name, rawHours)
			logrus.Debugf("%2d | %s", i, rest.name)
			if err != nil {
				// log the error & skip
				logrus.Errorf("could not get restaurant: %v", err)
				return
			}

			restaurants[i] = *rest
		}(i, row[0], row[1])
	}

	wg.Wait()

	return restaurants, nil
}

func toRestaurant(name string, rawHours string) (*restaurant, error) {
	hrs, err := parser.HoursFromString(rawHours)
	if err != nil {
		return nil, err
	}

	r := &restaurant{
		name:  name,
		hours: hrs,
	}

	return r, nil
}
