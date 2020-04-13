package utils

import (
	"encoding/json"
	"fmt"
	"github.com/SungminSo/docker-resource-checker/models"
	"io/ioutil"
	"time"
)

// If the arguments are incorrect during execution,
// Print the following for instructions on use.
func Guide() {
	fmt.Println("Usage:")
	fmt.Println("    ./monitor -t|--time <time> -c|--container <docker container name>")
	fmt.Println("\nFlags:")
	fmt.Println("    -h, --help          help for docker resource checker")
	fmt.Println("    -t, --time          set checking time")
	fmt.Println("    -c, --container     name of docker container for checking")
	fmt.Println("\nReturns:")
	fmt.Println("    - Print resource usage summary to console")
	fmt.Println("    - Save the JSON file named resource_<timestamp>.json")
}

// Print out the execution results for each item.
// And save the contents as a JSON file.
func PrintAndSave(result *models.Summary) {
	fmt.Println("CPU: ")
	fmt.Println("    AVG: ", result.CPU.Avg, "%")
	fmt.Println("    MIN: ", result.CPU.Min, "%")
	fmt.Println("    MAX: ", result.CPU.Max, "%")
	fmt.Println("MEMORY: ")
	fmt.Println("    AVG: ", result.MEM.Avg, result.MEMAvgUnit)
	fmt.Println("    MIN: ", result.MEM.Min, result.MEMMinUnit)
	fmt.Println("    MAX: ", result.MEM.Max, result.MEMMaxUnit)
	fmt.Println("NETWORK: ")
	fmt.Println("    In: ", result.NET.In, result.NETInUnit)
	fmt.Println("    OUT: ", result.NET.Out, result.NETOutUnit)
	fmt.Println("DISK: ")
	fmt.Println("    READ: ", result.DISK.Read, result.DISKReadUnit)
	fmt.Println("    WRITE: ", result.DISK.Write, result.DISKWriteUNIT)

	fileData, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err.Error())
	}

	fileName := "resource_" + time.Now().String() + ".json"
	ioutil.WriteFile(fileName, fileData, 0644)
}

func SetByteUnit(value float64) (float64, string) {
	unit := ""

	switch {
	case value >= models.GB:
		unit = "GB"
		value = value / models.GB
	case value >= models.MB:
		unit = "MB"
		value = value / models.MB
	case value >= models.KB:
		unit = "KB"
		value = value / models.KB
	case value >= models.B:
		unit = "B"
	}

	return value, unit
}

func SetIByteUnit(value float64) (float64, string) {
	unit := ""

	switch {
	case value >= models.GiB:
		unit = "GiB"
		value = value / models.GiB
	case value >= models.MiB:
		unit = "MiB"
		value = value / models.MiB
	case value >= models.KiB:
		unit = "KiB"
		value = value / models.KiB
	case value >= models.B:
		unit = "B"
	}

	return value, unit
}