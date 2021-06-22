package cmd

import (
	"encoding/csv"
	"fmt"
	"gokyo-stats/csv_writer"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var yearToDates = map[int]map[string][]time.Time{
	2019: map[string][]time.Time{
		"Monday":    []time.Time{},
		"Tuesday":   []time.Time{},
		"Wednesday": []time.Time{},
		"Thursday":  []time.Time{},
		"Friday":    []time.Time{},
		"Saturday":  []time.Time{},
		"Sunday":    []time.Time{},
	},
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
	GetAllDatesOfYear(2019)
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

var outputColumns = []string{"file_name", "row_number", "ward", "date", "year", "month", "week_num", "day", "time_range", "employment_type", "schedule_type", "presense", "budget_actual", "hours"}

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

func TimeInString(dateTime time.Time) string {
	d := strings.Split(dateTime.Format("15:04:05"), ":")
	return fmt.Sprintf("%s:%s", d[0], d[1])
}

var timeSlot = 30

type TimeFrame struct {
	Frame string
	Diff  int
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

// func SplitTimeRange(from string, to string) []TimeFrame {

// 	fromHour := GetHour(from)
// 	toHour := GetHour(to)

// 	fromMinute := GetMinutes(from)
// 	toMinute := GetMinutes(to)

// 	initialFromDiff := 30
// 	initialToDiff := 30

// 	if float64(fromMinute)/float64(timeSlot) > 0 && float64(fromMinute)/float64(timeSlot) < 1 {
// 		initialFromDiff = timeSlot - fromMinute
// 		fromMinute = 0
// 	} else if float64(fromMinute)/float64(timeSlot) > 1 {
// 		initialFromDiff = fromMinute - timeSlot
// 		fromMinute = timeSlot
// 	}

// 	if float64(toMinute)/float64(timeSlot) > 0 && float64(toMinute)/float64(timeSlot) < 1 {
// 		initialToDiff = timeSlot - toMinute
// 		toMinute = 0
// 	} else if float64(toMinute)/float64(timeSlot) > 1 {
// 		initialToDiff = toMinute - timeSlot
// 		toMinute = timeSlot
// 	}

// 	firstTime := time.Date(2020, 1, 1, fromHour, fromMinute, 0, 0, time.Now().Location())
// 	secondTime := time.Date(2020, 1, 1, toHour, toMinute, 0, 0, time.Now().Location())

// 	timeRanges := make([]TimeFrame, 0)

// 	for firstTime.Before(secondTime) {
// 		firstRange := TimeInString(firstTime)

// 		firstTime = firstTime.Add(30 * time.Minute)

// 		secondRange := TimeInString(firstTime)

// 		var tf = TimeFrame{
// 			Frame: fmt.Sprintf("%s-%s", firstRange, secondRange),
// 			Diff:  30,
// 		}

// 		timeRanges = append(timeRanges, tf)

// 	}

// 	if len(timeRanges) == 0 {
// 		return timeRanges
// 	}

// 	timeRanges[0] = TimeFrame{
// 		Frame: timeRanges[0].Frame,
// 		Diff:  initialFromDiff,
// 	}

// 	timeRanges[len(timeRanges)-1] = TimeFrame{
// 		Frame: timeRanges[len(timeRanges)-1].Frame,
// 		Diff:  initialToDiff,
// 	}

// 	return timeRanges

// }

var intToDay = map[int]string{
	0: "Monday",
	1: "Tuesday",
	2: "Wednesday",
	3: "Thursday",
	4: "Friday",
	5: "Saturday",
	6: "Sunday",
}

func GenerateOutput(cmd *cobra.Command, args []string) {

	csvFileName, err := cmd.Flags().GetString("file")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	outputFileName, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	csvFile, err := os.Open(csvFileName)
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	// "file_name", "ward", "date", "year", "month", "week_num", "day", "time_range", "employment_type", "presense", "budget_actual", "hours"
	var finalOutput = [][]string{outputColumns}

	for _, row := range csvLines[1:] {

		fileName := row[0]
		ward := row[1]
		// date := row[2]
		year := row[3]
		// month := row[4]
		// weekNum := row[5]

		scheduleType := row[7]
		employmentType := row[8]

		monday := row[9]
		tuesday := row[10]
		wednesday := row[11]
		thursday := row[12]
		friday := row[13]
		saturday := row[14]
		sunday := row[15]

		rowNumber := row[19]

		var days = []string{monday, tuesday, wednesday, thursday, friday, saturday, sunday}

		from := row[16]
		to := row[17]

		timeRanges := SplitTimeRange(from, to)

		dayToDates := yearToDates[GetIntFromString(year)]

		if len(timeRanges) != 0 {
			for i, presense := range days {
				weekDay := intToDay[i]
				for _, eachDate := range dayToDates[weekDay] {
					_, w := eachDate.ISOWeek()
					weekNum := GetStringFromInt(w)
					for _, timeRange := range timeRanges {
						hours := (float64(GetFloatFromString(presense)) * float64(timeRange.Diff)) / 60.0
						var eachRow = []string{fileName, rowNumber, ward, eachDate.Format("2006-01-02"), GetStringFromInt(eachDate.Year()), eachDate.Month().String(), weekNum, weekDay, timeRange.Frame, employmentType, scheduleType, presense, "Budget", fmt.Sprintf("%f", hours)}
						finalOutput = append(finalOutput, eachRow)
					}
				}

			}
		}

	}
	err = csv_writer.WriteToCSV(outputFileName, finalOutput)
	if err != nil {
		fmt.Println("CSV generation failed")
		os.Exit(1)
	}

}

func GetStringFromInt(data int) string {
	return strconv.Itoa(data)
}

func GetFloatFromString(data string) float64 {
	if s, err := strconv.ParseFloat(data, 64); err == nil {
		return s
	}
	return 0.0
}

func GetIntFromString(data string) int {
	if i, err := strconv.Atoi(data); err == nil {
		return i
	}

	return 0
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Output File",
	Long:  `Generate Output file from give CSV data file`,
	Run: func(cmd *cobra.Command, args []string) {
		GenerateOutput(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("file", "f", "result.csv", "csv data file")
	generateCmd.Flags().StringP("output", "o", "output.csv", "output file name")

}
