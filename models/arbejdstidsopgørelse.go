package models

type AllSecondData struct {
	Data []SecondData
}

type SecondData struct {
	Sector         string
	Period         string
	Name           string
	EmployeeNumber string
	Initials       string
	Rows           []SecondDataRow
}

type SecondDataRow struct {
	RowNumber   string
	Week        string
	Date        string
	Day         string
	From        string
	To          string
	ServiceTag  string
	Type        string
	Description string
}
