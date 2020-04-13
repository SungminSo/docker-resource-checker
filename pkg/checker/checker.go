package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/SungminSo/docker-resource-checker/models"
	"io"
	"strings"
	"time"
)

// Find docker container ID.
// By using container name.
// Varaible "target" means the docker container name
func FindContainerID(target string) string {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, container := range containers {
		if container.Names[0] == "/" + target {
			return container.ID
		}
	}

	return ""
}

// Find ID with a name if a docker container exist.
// if so, run the "docker stats" with the corresponding container ID.
// Parse the execution results for each item.
func CheckContainerStat(target string) (models.Info, error){
	containerID := FindContainerID(target)
	if containerID == "" {
		fmt.Println("No docker container with that name could be found")
		return models.Info{}, nil
	}

	returnInfo := models.Info{}

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err.Error())
		return returnInfo, err
	}

	stats, err := cli.ContainerStats(context.Background(), containerID, false)
	if err != nil {
		fmt.Println(err.Error())
		return returnInfo, err
	}
	defer stats.Body.Close()


	decoder := json.NewDecoder(stats.Body)
	for {
	var v *types.StatsJSON
	if err := decoder.Decode(&v); err == io.EOF {
		break
	} else if err != nil {
		decoder = json.NewDecoder(io.MultiReader(decoder.Buffered(), stats.Body))
		time.Sleep(100 * time.Millisecond)
		continue
	}

	osType := stats.OSType

	if osType != "windows" {
		returnInfo.CPU = calculateCPUPercentUnix(v.PreCPUStats, v.CPUStats)
		returnInfo.MEM = float64(v.MemoryStats.Usage - v.MemoryStats.Stats["cache"])
		returnInfo.DISK.Read, returnInfo.DISK.Write = calculateDiskRW(v.BlkioStats)
	} else {
		returnInfo.CPU = calculateCPUPercentWindows(v)
		returnInfo.MEM = float64(v.MemoryStats.PrivateWorkingSet)
		returnInfo.DISK.Read = float64(v.StorageStats.ReadSizeBytes)
		returnInfo.DISK.Write = float64(v.StorageStats.WriteSizeBytes)
	}
	returnInfo.NET.In, returnInfo.NET.Out = calculateNetworkIO(v.Networks)
	}

	return returnInfo, nil
}

func calculateCPUPercentUnix(previousCPU, presentCPU types.CPUStats) float64 {
	var cpuPercent float64
	cpuDelta := float64(presentCPU.CPUUsage.TotalUsage) - float64(previousCPU.CPUUsage.TotalUsage)
	systemDelta := float64(presentCPU.SystemUsage) - float64(previousCPU.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(presentCPU.CPUUsage.PercpuUsage)) * 100
	}

	return cpuPercent
}

func calculateCPUPercentWindows(v *types.StatsJSON) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                         // Convert to number of 100ns intervals
	possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors

	// Intervals used
	intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals) * 100.0
	}
	return 0.00
}

func calculateDiskRW(blkIO types.BlkioStats) (float64, float64) {
	diskRead := float64(0)
	diskWrite := float64(0)

	for _, bioEntry := range blkIO.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			diskRead += float64(bioEntry.Value)
		case "write":
			diskWrite += float64(bioEntry.Value)
		}
	}
	return diskRead, diskWrite
}

func calculateNetworkIO(network map[string]types.NetworkStats) (float64, float64) {
	rx := float64(0)  // in
	tx := float64(0)  // out

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}