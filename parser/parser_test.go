package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/c-Brooks/find-open-restaurants/hours"
)

type testCase struct {
	row      string
	expected hours.Hours
}

func TestParse(t *testing.T) {
	testCases := prepareData()
	for i, tc := range testCases {
		actual, err := FromString(tc.row)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("%2d expected %s to be %s", i, actual, tc.expected)
			t.Fail()
		}
	}
}

func prepareData() []testCase {
	elevenAM, _ := time.Parse("03:04pm", "11:00am")
	elevenPM, _ := time.Parse("03:04pm", "11:00pm")
	tenPM, _ := time.Parse("03:04pm", "10:00pm")
	elevenToEleven := hours.TimeRange{StartTime: hours.TimeOfDay{elevenAM}, EndTime: hours.TimeOfDay{elevenPM}}
	elevenToTen := hours.TimeRange{StartTime: hours.TimeOfDay{elevenAM}, EndTime: hours.TimeOfDay{tenPM}}
	testData := []testCase{
		testCase{
			row: "Mon-Sat 11 am - 11 pm  / Sun 11 am - 10 pm",
			expected: hours.Hours{
				time.Monday:    elevenToEleven,
				time.Tuesday:   elevenToEleven,
				time.Wednesday: elevenToEleven,
				time.Thursday:  elevenToEleven,
				time.Friday:    elevenToEleven,
				time.Saturday:  elevenToEleven,
				time.Sunday:    elevenToTen,
			},
		},
	}

	return testData
}
