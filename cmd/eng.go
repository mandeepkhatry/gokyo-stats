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

var thirdDataColumns = []string{"file_name", "row_number", "period", "sector", "one_time_benefit", "job_category", "employee", "week", "day", "date", "quantity"}

func NewTestOutput(cmd *cobra.Command, args []string) {
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
	data = append(data, thirdDataColumns)

	//var thirdDataColumns = []string{"period", "sector", "one_time_benefit", "job_category", "employee", "week", "day", "date", "quantity"}

	for _, file := range allFiles {
		p, _ := parser.BuildParser(file)
		thirdData := p.GetEngData(parser.Params{})

		for _, row := range thirdData.Rows {
			data = append(data, []string{file, row.RowNumber, thirdData.Period, thirdData.Sector, row.OneTimeBenefit, row.JobCategory, row.Employee, row.Week, row.Day, row.Date, row.Quantity})
		}

	}

	err = csv_writer.WriteToCSV(outputFile, data)
	if err != nil {
		NewTestOutput(cmd, args)
	}

}

// naeCmd represents the nae command
var endCmd = &cobra.Command{
	Use:   "eng",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		NewTestOutput(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(endCmd)
	endCmd.Flags().StringP("folder", "f", "eng", "select folder containing XLSX files")
	endCmd.Flags().StringP("output", "o", "eng_result.csv", "output file name csv")
}
