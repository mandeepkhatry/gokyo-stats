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

var secondDataColumns = []string{"file_name", "row_number", "sector", "period", "name", "employee_number", "initials", "week", "date", "day", "from", "to", "service_tag", "type", "description"}

func TestOutput(cmd *cobra.Command, args []string) {
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
	data = append(data, secondDataColumns)

	//[]string{"sector", "period", "name", "employee_number", "initials", "week", "date", "day", "from", "to", "service_tag", "type", "description"}

	for _, file := range allFiles {
		p, _ := parser.BuildParser(file)
		secondData := p.GetArbData(parser.Params{})

		for _, d := range secondData.Data {
			for _, row := range d.Rows {
				data = append(data, []string{file, row.RowNumber, d.Sector, file, d.Period, d.Name, d.EmployeeNumber, d.Initials, row.Week, row.Date, row.Day, row.From, row.To, row.ServiceTag, row.Type, row.Description})
			}
		}

	}

	err = csv_writer.WriteToCSV(outputFile, data)
	if err != nil {
		fmt.Println("CSV generation failed")
		os.Exit(1)
	}

}

// arbeCmd represents the arbe command
var arbeCmd = &cobra.Command{
	Use:   "arbe",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		TestOutput(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(arbeCmd)
	arbeCmd.Flags().StringP("folder", "f", "arb", "select folder containing XLSX files")
	arbeCmd.Flags().StringP("output", "o", "arb_result.csv", "output file name csv")

}
