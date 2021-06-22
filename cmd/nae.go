package cmd

import (
	"fmt"
	"gokyo-stats/csv_writer"
	"gokyo-stats/parser"
	"gokyo-stats/utils"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var fourthDataColumns = []string{"file_name", "row_number", "organization_unit", "calendar_date", "ma_no", "date_of_accession", "job_category_for_grade", "absense_hours", "attendance_hours_with_holidays"}

func FourthTestOutput(cmd *cobra.Command, args []string) {
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
	data = append(data, fourthDataColumns)

	// var fourthDataColumns = []string{"organization_unit", "calendar_date", "ma_no", "date_of_accession", "job_category_for_grade", "absense_hours", "attendance_hours_with_holidays"}

	for _, file := range allFiles {
		p, _ := parser.BuildParser(file)
		fourthData := p.GetNaeData(parser.Params{})

		for _, row := range fourthData {
			data = append(data, []string{file, row.RowNumber, row.OrganizationalUnit, row.CalendarDate, row.MaNo, row.DateOfAccession, row.JobCategoryForGrade, row.AbsenseHours, row.AttendanceHoursWithHolidays})
		}
	}

	err = csv_writer.WriteToCSV(outputFile, data)
	if err != nil {
		FourthTestOutput(cmd, args)
	}

}

// naeCmd represents the nae command
var naeCmd = &cobra.Command{
	Use:   "nae",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		FourthTestOutput(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(naeCmd)
	naeCmd.Flags().StringP("folder", "f", "nae", "select folder containing XLSX files")
	naeCmd.Flags().StringP("output", "o", "nae_result.csv", "output file name csv")

}
