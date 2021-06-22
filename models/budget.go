package models

type Data struct {
	Section string             `json:"section"` //afsnit
	Blocks  []BlockInformation `json:"blocks"`
	Date    string
	Month   string
	Year    string
}

//Information per block
type BlockInformation struct {
	NoOfWeeks   int          `json:"no_of_weeks"` //antal uger
	ActiveName  string       `json:"active_name"` //Norm.aktiv. or Lav aktiv.
	InnerBlocks []InnerBlock `json:"inner_blocks"`
}

type InnerBlock struct {
	Name          string         `json:"name"`
	InnerSections []InnerSection `json:"inner_blocks"`
}

type InnerSection struct {
	RowNumber       int        `json:"row_number"`
	Name            string     `json:"name"`
	AttendanceCount Attendance `json:"attendance"`
	From            string     `json:"from"`
	To              string     `json:"to"`
	Timer           float64    `json:"timer"`
}

type Attendance struct {
	Monday    float64 `json:"monday"`    //ma
	Tuesday   float64 `json:"tuesday"`   //ti
	Wednesday float64 `json:"wednesday"` //on
	Thursday  float64 `json:"thursday"`  //to
	Friday    float64 `json:"friday"`    //fr
	Saturday  float64 `json:"saturday"`  //lø
	Sunday    float64 `json:"sunday"`    //sø
}
