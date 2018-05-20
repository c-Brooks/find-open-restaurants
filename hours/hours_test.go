package hours

import (
	"testing"
	"time"
)

type testCase struct {
	timeRange     TimeRange
	time          TimeOfDay
	shouldContain bool
}

func TestContains(t *testing.T) {
	testCases := prepareData()

	for i, tc := range testCases {
		if tc.shouldContain && !tc.timeRange.Contains(tc.time) {
			t.Errorf("%2d | expected %s to contain %s", i, tc.timeRange, tc.time.Format(time.RFC822))
			t.Failed()
		} else if !tc.shouldContain && tc.timeRange.Contains(tc.time) {
			t.Errorf("%2d | expected %s to not contain %s", i, tc.timeRange, tc.time.Format(time.RFC822))
			t.Failed()
		}
	}

}

func prepareData() []testCase {
	eightAM, _ := time.Parse(time.Kitchen, "8:00AM")
	sixPM, _ := time.Parse(time.Kitchen, "6:00PM")
	sevenPM, _ := time.Parse(time.Kitchen, "7:00PM")
	oneAM, _ := time.Parse(time.Kitchen, "01:00AM")
	twoAM, _ := time.Parse(time.Kitchen, "02:00AM")
	noon, _ := time.Parse(time.Kitchen, "12:00PM")
	nineToFive := TimeRange{StartTime: TimeOfDay{eightAM}, EndTime: TimeOfDay{sixPM}}
	sevenToTwo := TimeRange{StartTime: TimeOfDay{sevenPM}, EndTime: TimeOfDay{twoAM}}

	testData := []testCase{
		testCase{timeRange: nineToFive, time: TimeOfDay{eightAM}, shouldContain: false},
		testCase{timeRange: nineToFive, time: TimeOfDay{noon}, shouldContain: true},
		testCase{timeRange: nineToFive, time: TimeOfDay{eightAM}, shouldContain: false},
		testCase{timeRange: nineToFive, time: TimeOfDay{twoAM}, shouldContain: false},
		testCase{timeRange: sevenToTwo, time: TimeOfDay{oneAM}, shouldContain: true},
		testCase{timeRange: sevenToTwo, time: TimeOfDay{noon}, shouldContain: false},
	}

	return testData
}
