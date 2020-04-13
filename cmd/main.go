package main

import (
	"flag"
	"fmt"
	"github.com/SungminSo/docker-resource-checker/models"
	"github.com/SungminSo/docker-resource-checker/pkg/checker"
	"github.com/SungminSo/docker-resource-checker/pkg/utils"
	"os"
	"time"
)

// Check arguments length is valid.
func validateArgs() {
	if len(os.Args) < 2 {
		utils.Guide()
		os.Exit(1)
	}
}

// Check the stats every 2 seconds from the start to the end time.
// Calculate the recorded values at the end time.
func Checker(nextTime time.Time, endTime time.Time, containerName string, resource *models.Resource) {
	if  nextTime.Unix() < endTime.Unix() {
		time.Sleep(time.Until(nextTime))

		info, err := checker.CheckContainerStat(containerName)
		if err != nil {
			fmt.Println(err.Error())
		}
		resource.Infos = append(resource.Infos, info)
		nextTime = nextTime.Add(time.Second * 2)
		Checker(nextTime, endTime, containerName, resource)
	} else {
		result := &models.Summary{
			CPU:  models.Stats{
				Avg: 0,
				Min: 0,
				Max: 0,
			},
			MEM:  models.Stats{
				Avg: 0,
				Min: 0,
				Max: 0,
			},
			NET:  models.IO{
				In: 0,
				Out: 0,
			},
			DISK: models.RW{
				Read: 0,
				Write: 0,
			},
		}

		for _, info := range resource.Infos {
			result.CPU.Avg += info.CPU
			if result.CPU.Min != 0 && result.CPU.Min > info.CPU {
				result.CPU.Min = info.CPU
			} else if result.CPU.Min == 0 {
				result.CPU.Min = info.CPU
			}
			if result.CPU.Max != 0 && result.CPU.Max < info.CPU {
				result.CPU.Max = info.CPU
			} else if result.CPU.Max == 0 {
				result.CPU.Max = info.CPU
			}

			result.MEM.Avg += info.MEM
			if result.MEM.Min != 0 && result.MEM.Min > info.MEM {
				result.MEM.Min = info.MEM
			} else if result.MEM.Min == 0 {
				result.MEM.Min = info.MEM
			}
			if result.MEM.Max != 0 && result.MEM.Max < info.MEM {
				result.MEM.Max = info.MEM
			} else if result.MEM.Max == 0 {
				result.MEM.Max = info.MEM
			}

			result.NET.In += info.NET.In
			result.NET.Out += info.NET.Out

			result.DISK.Write += info.DISK.Write
			result.DISK.Read += info.DISK.Read
		}

		result.CPU.Avg = result.CPU.Avg / float64(len(resource.Infos))
		result.MEM.Avg, result.MEMAvgUnit = utils.SetIByteUnit(result.MEM.Avg / float64(len(resource.Infos)))
		result.MEM.Min, result.MEMMinUnit = utils.SetIByteUnit(result.MEM.Min / float64(len(resource.Infos)))
		result.MEM.Max, result.MEMMaxUnit = utils.SetIByteUnit(result.MEM.Max / float64(len(resource.Infos)))
		result.NET.In, result.NETInUnit = utils.SetByteUnit(result.NET.In / float64(len(resource.Infos)))
		result.NET.Out, result.NETOutUnit = utils.SetByteUnit(result.NET.Out / float64(len(resource.Infos)))
		result.DISK.Read, result.DISKReadUnit = utils.SetByteUnit(result.DISK.Read / float64(len(resource.Infos)))
		result.DISK.Write, result.DISKWriteUNIT = utils.SetByteUnit(result.DISK.Write / float64(len(resource.Infos)))

		utils.PrintAndSave(result)
	}
}

func main() {
	validateArgs()

	timeShort := flag.Int("t", 0, "How long to check resources")
	timeLong := flag.Int("time", 0, "How long to check resources")
	containerShort := flag.String("c", "", "Container name")
	containerLong := flag.String("container", "", "Container name")

	flag.Parse()

	var timeInterval int
	var containerName string

	if *timeShort != 0 {
		timeInterval = *timeShort
	} else if *timeLong != 0 {
		timeInterval = *timeLong
	} else {
		utils.Guide()
		os.Exit(1)
	}

	if *containerShort != "" {
		containerName = *containerShort
	} else if *containerLong != "" {
		containerName = *containerLong
	} else {
		utils.Guide()
		os.Exit(1)
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Second * time.Duration(timeInterval))

	fmt.Println("Start at:", startTime)
	fmt.Println("End at:", endTime)

	resource := &models.Resource{}

	Checker(startTime, endTime, containerName, resource)
}