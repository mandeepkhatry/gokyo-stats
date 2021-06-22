package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var yearToDates = map[int]map[string][]time.Time{
	2020: map[string][]time.Time{
		"Monday":    []time.Time{},
		"Tuesday":   []time.Time{},
		"Wednesday": []time.Time{},
		"Thursday":  []time.Time{},
		"Friday":    []time.Time{},
		"Saturday":  []time.Time{},
		"Sunday":    []time.Time{},
	},
	2021: map[string][]time.Time{
		"Monday":    []time.Time{},
		"Tuesday":   []time.Time{},
		"Wednesday": []time.Time{},
		"Thursday":  []time.Time{},
		"Friday":    []time.Time{},
		"Saturday":  []time.Time{},
		"Sunday":    []time.Time{},
	},
}

func init() {
	GetAllDatesOfYear(2020)
	GetAllDatesOfYear(2021)
}

func GetAllDatesOfYear(year int) {

	firstDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	nextYear := year + 1

	for {
		yearToDates[year][firstDate.Weekday().String()] = append(yearToDates[year][firstDate.Weekday().String()], firstDate)

		firstDate = firstDate.Add(24 * time.Hour)

		if firstDate.Year() == nextYear {
			break
		}
	}

}

var timeSlot = 30

type TimeFrame struct {
	Frame string
	Diff  int
}

func GetHour(timeRange string) int {

	if i, err := strconv.Atoi(strings.Split(timeRange, ":")[0]); err == nil {
		return i
	}

	return 0
}

func GetMinutes(timeRange string) int {

	if i, err := strconv.Atoi(strings.Split(timeRange, ":")[1]); err == nil {
		return i
	}

	return 0
}

func SplitTimeRange(from string, to string) []TimeFrame {

	fromHour := GetHour(from)
	toHour := GetHour(to)

	fromMinute := GetMinutes(from)
	toMinute := GetMinutes(to)

	initialFromDiff := 30
	initialToDiff := 30

	if float64(fromMinute)/float64(timeSlot) > 0 && float64(fromMinute)/float64(timeSlot) < 1 {
		initialFromDiff = timeSlot - fromMinute
		fromMinute = 0
	} else if float64(fromMinute)/float64(timeSlot) > 1 {
		initialFromDiff = fromMinute - timeSlot
		fromMinute = timeSlot
	}

	if float64(toMinute)/float64(timeSlot) > 0 && float64(toMinute)/float64(timeSlot) < 1 {
		initialToDiff = timeSlot - toMinute
		toMinute = 0
	} else if float64(toMinute)/float64(timeSlot) > 1 {
		initialToDiff = toMinute - timeSlot
		toMinute = timeSlot
	}

	firstTime := time.Date(2020, 1, 1, fromHour, fromMinute, 0, 0, time.Now().Location())
	secondTime := time.Date(2020, 1, 1, toHour, toMinute, 0, 0, time.Now().Location())

	timeRanges := make([]TimeFrame, 0)

	firstTimeString := firstTime.Format("15:04:05")
	secondTimeString := secondTime.Format("15:04:05")

	for firstTimeString != secondTimeString {
		firstRange := TimeInString(firstTime)

		firstTime = firstTime.Add(30 * time.Minute)

		secondRange := TimeInString(firstTime)

		firstTimeString = firstTime.Format("15:04:05")

		// fmt.Println(firstTime)
		// fmt.Println(secondTime)
		// time.Sleep(1 * time.Second)

		var tf = TimeFrame{
			Frame: fmt.Sprintf("%s-%s", firstRange, secondRange),
			Diff:  30,
		}

		timeRanges = append(timeRanges, tf)

	}

	if len(timeRanges) == 0 {
		return timeRanges
	}

	timeRanges[0] = TimeFrame{
		Frame: timeRanges[0].Frame,
		Diff:  initialFromDiff,
	}

	timeRanges[len(timeRanges)-1] = TimeFrame{
		Frame: timeRanges[len(timeRanges)-1].Frame,
		Diff:  initialToDiff,
	}

	return timeRanges

}

func TimeInString(dateTime time.Time) string {
	d := strings.Split(dateTime.Format("15:04:05"), ":")
	return fmt.Sprintf("%s:%s", d[0], d[1])
}

func main() {
	for _, d := range SplitTimeRange("22:45", "3:07") {
		fmt.Println("Frame : ", d)
	}

}
