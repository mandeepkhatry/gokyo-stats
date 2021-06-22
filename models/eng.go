package models

type ThirdData struct {
	Period string
	Sector string
	Rows   []ThirdDataRow
}

type ThirdDataRow struct {
	RowNumber      string
	OneTimeBenefit string
	JobCategory    string
	Employee       string
	Week           string
	Day            string
	Date           string
	Quantity       string
}
