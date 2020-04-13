package models

type Resource struct {
	Infos 	[]Info 		`json:"infos"`
}

type Info struct {
	CPU 	float64 	`json:"cpu"`
	MEM 	float64 	`json:"mem"`
	NET 	IO			`json:"net"`
	DISK 	RW			`json:"disk"`
}

type Summary struct {
	CPU           Stats  `json:"cpu"`
	MEM           Stats  `json:"mem"`
	MEMAvgUnit    string `json:"memAvgUnit"`
	MEMMinUnit    string `json:"memMinUnit"`
	MEMMaxUnit    string `json:"memMaxUnit"`
	NET           IO     `json:"net"`
	NETInUnit     string `json:"netInUnit"`
	NETOutUnit    string `json:"netOutUnit"`
	DISK          RW     `json:"disk"`
	DISKReadUnit  string `json:"diskReadUnit"`
	DISKWriteUNIT string `json:"diskWriteUnit"`
}

type Stats struct {
	Avg 	float64 	`json:"avg"`
	Min 	float64 	`json:"min"`
	Max 	float64 	`json:"max"`
}

type IO struct {
	In		float64		`json:"in"`
	Out 	float64 	`json:"out"`
}

type RW struct {
	Read 	float64 	`json:"read"`
	Write 	float64 	`json:"write"`
}

const (
	B = 1
	KB = B * 1000
	MB = KB * 1000
	GB = MB * 1000

	KiB = B * 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)
