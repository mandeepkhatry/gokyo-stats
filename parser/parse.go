package parser

import (
	"fmt"
	"gokyo-stats/models"
	"gokyo-stats/utils"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Params struct {
	MinWeeks int
}

type Parser struct {
	file      *excelize.File
	sheetName string
	sheets    []string
}

func BuildParser(fileName string) (Parser, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return Parser{}, err
	}
	var p Parser
	p.file = f
	p.sheetName = f.GetSheetName(0)

	p.sheets = f.GetSheetList()

	return p, nil
}

func (p *Parser) GetSheetValueFromAxis(axis string) string {

	val, err := p.file.GetCellValue(p.sheetName, axis)
	if err != nil {
		log.Println("ERR : ", err.Error())
		return ""
	}
	return val
}

func (p *Parser) GetSheetValueFromAxisForSheet(sheetName string, axis string) string {
	val, err := p.file.GetCellValue(sheetName, axis)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return val
}

func (p *Parser) GetAllRowsForSheet(sheetName string) [][]string {
	data, err := p.file.GetRows(sheetName)
	if err != nil {
		return [][]string{}
	}
	return data
}

func (p *Parser) GetAllRows() [][]string {
	data, err := p.file.GetRows(p.sheetName)
	if err != nil {
		return [][]string{}
	}
	return data
}

var blockSeperator = "-"
var endSeperator = ""

type BlockSeperators struct {
	FirstBoundary  int
	SecondBoundary int
}

func (p *Parser) GetBlockSeperators() []BlockSeperators {
	wholeBlock := p.GetAllRows()
	blockSeperators := make([]BlockSeperators, 0)

	for i, eachRow := range wholeBlock {
		val := eachRow[0]
		if val == blockSeperator {

			if i+1 < len(wholeBlock) {
				if wholeBlock[i+1][0] == endSeperator {
					lengthOfBlockSeperators := len(blockSeperators)
					if lengthOfBlockSeperators != 0 {
						lastBlockSeperator := lengthOfBlockSeperators - 1

						//Add SecondBoundary to just before block seperator
						blockSeperators[lastBlockSeperator] = BlockSeperators{
							FirstBoundary:  blockSeperators[lastBlockSeperator].FirstBoundary,
							SecondBoundary: i - 1,
						}

					}
					break
				}
			}
			lengthOfBlockSeperators := len(blockSeperators)
			if lengthOfBlockSeperators != 0 {
				lastBlockSeperator := lengthOfBlockSeperators - 1

				//Add SecondBoundary to just before block seperator
				blockSeperators[lastBlockSeperator] = BlockSeperators{
					FirstBoundary:  blockSeperators[lastBlockSeperator].FirstBoundary,
					SecondBoundary: i - 1,
				}

			}
			blockSeperators = append(blockSeperators, BlockSeperators{
				FirstBoundary: i + 1,
			})

		}
	}

	return blockSeperators
}

func GetFloatValue(data string) float64 {
	d, _ := strconv.ParseFloat(data, 64)
	return d
}

var intToMonth = map[int]string{
	1:  time.January.String(),
	2:  time.February.String(),
	3:  time.March.String(),
	4:  time.April.String(),
	5:  time.May.String(),
	6:  time.June.String(),
	7:  time.July.String(),
	8:  time.August.String(),
	9:  time.September.String(),
	10: time.October.String(),
	11: time.November.String(),
	12: time.December.String(),
}

func GetMonth(date string) string {

	monthInt := 1

	if i, err := strconv.Atoi(strings.Split(date, "-")[0]); err == nil {
		monthInt = i
	}

	return intToMonth[monthInt]

}

// type FourthData struct {
// 	OrganizationalUnit          string
// 	CalendarDate                string
// 	MaNo                        string
// 	DateOfAccession             string
// 	JobCategoryForGrade         string
// 	AbsenseHours                string
// 	AttendanceHoursWithHolidays string
// }

func (p *Parser) GetNaeData(params Params) []models.FourthData {

	fromRow := 1

	rows := p.GetAllRows()[fromRow:]

	actual := fromRow + 1

	totalData := make([]models.FourthData, 0)

	for i, row := range rows {

		organizationalUnit := row[0]
		calendarDate := row[1]
		maNo := row[2]
		dateOfAccession := row[3]
		jobCategory := row[4]
		absenseHours := row[5]
		attendanceHoursWithHolidays := row[6]

		rowNo := actual + i

		if calendarDate == "Resultat" || maNo == "Resultat" {

		} else {
			data := models.FourthData{
				RowNumber:                   strconv.Itoa(rowNo),
				OrganizationalUnit:          organizationalUnit,
				CalendarDate:                calendarDate,
				MaNo:                        maNo,
				DateOfAccession:             dateOfAccession,
				JobCategoryForGrade:         jobCategory,
				AbsenseHours:                absenseHours,
				AttendanceHoursWithHolidays: attendanceHoursWithHolidays,
			}
			totalData = append(totalData, data)
		}
	}
	return totalData
}

func (p *Parser) GetEngData(params Params) models.ThirdData {

	period := p.GetSheetValueFromAxis("C9")
	sector := p.GetSheetValueFromAxis("F10")

	allRows := p.GetAllRows()

	fromRow := 13

	rows := allRows[fromRow:]

	prevOneTimeBenefit := ""
	prevJobCategory := ""
	prevEmployee := ""

	thirdData := models.ThirdData{
		Period: period,
		Sector: sector,
	}

	thirdDataRows := make([]models.ThirdDataRow, 0)

	actualRow := fromRow + 1

	for i, row := range rows {

		rowNo := actualRow + i

		oneTimeBenefit := row[1]
		if oneTimeBenefit != "" {
			prevOneTimeBenefit = oneTimeBenefit
		}

		if oneTimeBenefit == "Total for perioden" {
			break
		}

		jobCategory := row[3]
		if jobCategory != "" {
			prevJobCategory = jobCategory
		}

		employee := row[7]
		if employee != "" {
			prevEmployee = employee
		}

		fmt.Println("ROW : ", row)

		week := row[9]
		day := row[12]
		date := row[13]
		quantity := row[14]

		if jobCategory == "" && employee == "" && week == "" && day == "" && date == "" && quantity == "" {
			//Do Nothing
		} else {

			if strings.Split(week, " ")[0] != "Total" || strings.Split(employee, " ")[0] != "Total" || strings.Split(jobCategory, " ")[0] != "Total" {
				modelRow := models.ThirdDataRow{
					RowNumber:      strconv.Itoa(rowNo),
					OneTimeBenefit: prevOneTimeBenefit,
					JobCategory:    prevJobCategory,
					Employee:       prevEmployee,
					Week:           week,
					Day:            day,
					Date:           date,
					Quantity:       quantity,
				}

				if date != "" && quantity != "" {
					thirdDataRows = append(thirdDataRows, modelRow)
				}

			}

		}

	}

	thirdData.Rows = thirdDataRows

	return thirdData

}

func (p *Parser) GetArbData(params Params) models.AllSecondData {

	allSecondData := models.AllSecondData{}

	allData := make([]models.SecondData, 0)

	for _, sheetName := range p.sheets {
		sector := p.GetSheetValueFromAxisForSheet(sheetName, "P11")
		period := p.GetSheetValueFromAxisForSheet(sheetName, "B11")
		name := p.GetSheetValueFromAxisForSheet(sheetName, "B17")
		initials := p.GetSheetValueFromAxisForSheet(sheetName, "P17")
		employeeNumber := p.GetSheetValueFromAxisForSheet(sheetName, "W17")

		allRows := p.GetAllRowsForSheet(sheetName)

		fromRow := 21

		rows := allRows[fromRow:]

		if rows[0][1] == "" {
			fromRow = 22
			rows = allRows[fromRow:]

		}

		prevWeekRow := ""
		prevDate := ""
		prevDay := ""
		prevFrom := ""
		prevTo := ""
		prevServiceTag := ""
		prevDataType := ""
		prevDescription := ""

		secondData := models.SecondData{
			Sector:         sector,
			Period:         period,
			Name:           name,
			EmployeeNumber: employeeNumber,
			Initials:       initials,
		}

		finalRows := make([]models.SecondDataRow, 0)

		actualRow := fromRow + 1

		for i, row := range rows {

			weekRow := row[1]
			date := row[2]
			day := row[4]

			from := ""
			to := ""

			rowNo := actualRow + i

			if row[5] != "" {
				from = strings.Split(row[5], "-")[0]
				to = strings.Split(row[5], "-")[1]
			}

			serviceTag := row[10]
			dataType := row[13]
			description := row[18]

			if weekRow != "" {
				prevWeekRow = weekRow
			}

			if date != "" {
				prevDate = date
			}

			if day != "" {
				prevDay = day
			}

			if serviceTag != "" {
				prevServiceTag = serviceTag
			}

			if dataType != "" {
				prevDataType = dataType
			}

			if description != "" {
				prevDescription = description
			}

			if strings.Split(description, ":")[0] == "TOTAL" {
				break
			}

			if from != "" {
				prevFrom = from
			}

			if to != "" {
				prevTo = to
			}

			dataRow := models.SecondDataRow{
				RowNumber:   strconv.Itoa(rowNo),
				Week:        prevWeekRow,
				Date:        prevDate,
				Day:         prevDay,
				From:        prevFrom,
				To:          prevTo,
				ServiceTag:  prevServiceTag,
				Type:        prevDataType,
				Description: prevDescription,
			}

			finalRows = append(finalRows, dataRow)

		}

		secondData.Rows = finalRows

		allData = append(allData, secondData)
	}

	allSecondData.Data = allData

	return allSecondData
}

func (p *Parser) GetData(params Params) models.Data {
	//Section

	section := p.GetSheetValueFromAxis("B4")
	//Section value of data

	sectionValue := strings.Split(section, ":")[1]
	allRows := p.GetAllRows()

	date := allRows[0][16]
	year := allRows[3][16]
	month := allRows[2][16]

	if i, err := strconv.Atoi(month); err == nil {
		month = intToMonth[i]
	}

	var data = models.Data{
		Section: sectionValue,
		Blocks:  []models.BlockInformation{},
		Date:    date,
		Year:    year,
		Month:   month,
	}

	//TESTING

	for _, eachBlockSeperator := range p.GetBlockSeperators() {
		firstBoundary := eachBlockSeperator.FirstBoundary
		secondBoundary := eachBlockSeperator.SecondBoundary

		var blockInfo models.BlockInformation
		block := allRows[firstBoundary : secondBoundary+1]

		active := block[0][0]
		noOfWeeks, _ := strconv.Atoi(block[2][1])

		actualRow := firstBoundary + 1

		if noOfWeeks >= params.MinWeeks {

			blockInfo.ActiveName = active
			blockInfo.NoOfWeeks = noOfWeeks

			dataBlock := block[3:]

			actualRow += 3

			var innerBlock models.InnerBlock

			start := false
			for j, d := range dataBlock {
				if d[3] == "-" && d[4] == "-" && d[0] != "" {
					if start {
						if len(innerBlock.InnerSections) != 0 {
							blockInfo.InnerBlocks = append(blockInfo.InnerBlocks, innerBlock)
						}

					}
					innerBlock = models.InnerBlock{}
					innerBlock.Name = d[0]
					start = true
				} else if d[3] == "" && d[4] == "" && d[5] == "" && d[6] == "" && d[7] == "" && d[8] == "" && d[9] == "" {
				} else {
					//Ignore empty rows
					monday := GetFloatValue(d[3])
					tuesday := GetFloatValue(d[4])
					wednesday := GetFloatValue(d[5])
					thursday := GetFloatValue(d[6])
					friday := GetFloatValue(d[7])
					saturday := GetFloatValue(d[8])
					sunday := GetFloatValue(d[9])

					var innerSection = models.InnerSection{}
					innerSection.Name = d[1]
					innerSection.From = utils.FromToFilter(d[11])
					innerSection.To = utils.FromToFilter(d[12])
					timer := GetFloatValue(d[13])
					innerSection.Timer = timer

					var attendance models.Attendance
					attendance.Monday = monday
					attendance.Tuesday = tuesday
					attendance.Wednesday = wednesday
					attendance.Thursday = thursday
					attendance.Friday = friday
					attendance.Saturday = saturday
					attendance.Sunday = sunday
					innerSection.AttendanceCount = attendance

					innerSection.RowNumber = actualRow + j

					innerBlock.InnerSections = append(innerBlock.InnerSections, innerSection)
				}

			}

			data.Blocks = append(data.Blocks, blockInfo)
		}

	}

	return data
}
