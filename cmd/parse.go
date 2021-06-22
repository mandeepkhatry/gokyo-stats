package cmd

import (
	"gokyo-stats/csv_writer"
	"gokyo-stats/parser"
	"gokyo-stats/utils"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"fmt"

	"github.com/spf13/cobra"
)

var columns = []string{"file_name", "section", "date", "year", "month", "no_of_weeks", "active_name", "inner_block_name", "inner_section_name", "attendance_monday", "attendance_tuesday", "attendance_wednesday", "attendance_thursday", "attendance_friday", "attendance_saturday", "attendance_sunday", "from", "to", "timer", "row_number"}

var re = regexp.MustCompile(`^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`)

func ParseFiles(cmd *cobra.Command, args []string) {

	folderName, err := cmd.Flags().GetString("folder")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	outputFileName, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	minWeeks, err := cmd.Flags().GetInt("weeks")
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	outputFile := utils.GetCSVFileName(outputFileName)

	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	allFiles := make([]string, 0)
	for _, f := range files {
		allFiles = append(allFiles, fmt.Sprintf("%s/%s", folderName, f.Name()))
	}

	data := make([][]string, 0)
	data = append(data, columns)

	for _, file := range allFiles {

		fmt.Println("File name : ", file)

		p, _ := parser.BuildParser(file)

		d := p.GetData(parser.Params{
			MinWeeks: minWeeks,
		})

		section := d.Section
		for _, block := range d.Blocks {
			activeName := block.ActiveName
			noOfWeeks := utils.GetStringFromInt(block.NoOfWeeks)

			for _, innerBlock := range block.InnerBlocks {
				innerBlockName := innerBlock.Name
				for _, innerSection := range innerBlock.InnerSections {

					innerSectionName := innerSection.Name
					attendanceMonday := utils.GetStringFromFloat(innerSection.AttendanceCount.Monday)
					attendanceTuesday := utils.GetStringFromFloat(innerSection.AttendanceCount.Tuesday)
					attendanceWednesday := utils.GetStringFromFloat(innerSection.AttendanceCount.Wednesday)
					attendanceThursday := utils.GetStringFromFloat(innerSection.AttendanceCount.Thursday)
					attendanceFriday := utils.GetStringFromFloat(innerSection.AttendanceCount.Friday)
					attendanceSaturday := utils.GetStringFromFloat(innerSection.AttendanceCount.Saturday)
					attendanceSunday := utils.GetStringFromFloat(innerSection.AttendanceCount.Sunday)

					from := innerSection.From
					to := innerSection.To

					timer := utils.GetStringFromFloat(innerSection.Timer)

					var actualData = []string{file, section, d.Date, d.Year, d.Month, noOfWeeks, activeName, innerBlockName, innerSectionName, attendanceMonday, attendanceTuesday, attendanceWednesday, attendanceThursday, attendanceFriday, attendanceSaturday, attendanceSunday, from, to, timer, GetStringFromInt(innerSection.RowNumber)}

					if re.MatchString(from) && re.MatchString(to) {
						data = append(data, actualData)
					}

				}
			}

		}

		// for _, k := range p.GetData().Blocks {
		// 	fmt.Println("**************************************************")
		// 	fmt.Println("No of weeks : ", k.NoOfWeeks)
		// 	fmt.Println("Active Name : ", k.ActiveName)
		// 	for _, v := range k.InnerBlocks {
		// 		fmt.Println("------------------------------------------------")
		// 		fmt.Println("------------------------------------------------")
		// 		fmt.Println("Inner section name : ", v.Name)
		// 		fmt.Println("------------------------------------------------")
		// 		for _, d := range v.InnerSections {
		// 			fmt.Println("Name : ", d.Name)
		// 			fmt.Println("From : ", d.From)
		// 			fmt.Println("To : ", d.To)
		// 			fmt.Println("Timer : ", d.Timer)
		// 			fmt.Println("Attendance : ", d.AttendanceCount)
		// 		}
		// 		fmt.Println("------------------------------------------------")
		// 		fmt.Println("------------------------------------------------")
		// 		fmt.Println("------------------------------------------------")
		// 	}

		// }
	}

	err = csv_writer.WriteToCSV(outputFile, data)
	if err != nil {
		fmt.Println("CSV generation failed")
		os.Exit(1)
	}

}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse XLSX files and convert into standard CSV format",
	Long: `Parse XLSX files and convert into standard CSV format. 
	For example:

	gokyo-stats parse folder files`,
	Run: func(cmd *cobra.Command, args []string) {
		ParseFiles(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringP("folder", "f", "files", "select folder containing XLSX files")
	parseCmd.Flags().StringP("output", "o", "result.csv", "output file name csv")
	parseCmd.Flags().IntP("weeks", "w", 2, "minimum number of weeks")
}
