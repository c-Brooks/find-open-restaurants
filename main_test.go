package main

import (
	"testing"
	"time"
)

type testCase struct {
	fileName        string
	searchDatetime  time.Time
	expectedResults []restaurant
}

func TestFindOpenRestaurants(t *testing.T) {
	testCases := prepareData()
	for _, tc := range testCases {
		csvFilename := tc.fileName
		searchDatetime := tc.searchDatetime
		restaurants := findOpenRestaurants(csvFilename, searchDatetime)
		if !testRestaurantNames(tc.expectedResults, restaurants) {
			t.Errorf("expected %v to equal %v", tc.expectedResults, restaurants)
			t.Fail()
		}
	}
}

func prepareData() []testCase {
	saturdayOneAM, err := time.Parse(time.ANSIC, "Sat May 21 01:00:00 2018")
	fridayNoon, err := time.Parse(time.ANSIC, "Fri May 20 12:00:00 2018")
	if err != nil {
		panic(err)
	}
	return []testCase{
		testCase{
			fileName:        "restaurants_test.csv",
			searchDatetime:  saturdayOneAM,
			expectedResults: []restaurant{restaurant{name: "Late Night Restaurant"}},
		},
		testCase{
			fileName:        "restaurants_test.csv",
			searchDatetime:  fridayNoon,
			expectedResults: []restaurant{restaurant{name: "Lunch Restaurant 1"}, restaurant{name: "Lunch Restaurant 2"}},
		},
	}
}

func testRestaurantNames(expected, actual []restaurant) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i, v := range expected {
		if v.name != actual[i].name {
			return false
		}
	}
	return true
}
