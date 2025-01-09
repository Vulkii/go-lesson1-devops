package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		checkServerStats()
	}
}

func checkServerStats() error {
	url := "http://srv.msk01.gigacorp.local/_stats"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	values := strings.Split(strings.TrimSpace(string(body)), ",")
	if len(values) != 7 {
		return fmt.Errorf("There are %d digits, not 7.", len(values))
	}

	loadAverageStr := values[0]
	memTotalStr := values[1]
	memUsedStr := values[2]
	diskTotalStr := values[3]
	diskUsedStr := values[4]
	netTotalStr := values[5]
	netUsedStr := values[6]

	loadAverage, err := strconv.ParseFloat(loadAverageStr, 64)
	if err != nil {
		return fmt.Errorf("Error with load average: %w", err)
	}

	memTotal, err := strconv.ParseFloat(memTotalStr, 64)
	if err != nil {
		return fmt.Errorf("Error with memTotal: %w", err)
	}

	memUsed, err := strconv.ParseFloat(memUsedStr, 64)
	if err != nil {
		return fmt.Errorf("Error with memUsed: %w", err)
	}

	diskTotal, err := strconv.ParseFloat(diskTotalStr, 64)
	if err != nil {
		return fmt.Errorf("Error with парсинга diskTotal: %w", err)
	}

	diskUsed, err := strconv.ParseFloat(diskUsedStr, 64)
	if err != nil {
		return fmt.Errorf("Error with парсинга diskUsed: %w", err)
	}

	netTotal, err := strconv.ParseFloat(netTotalStr, 64)
	if err != nil {
		return fmt.Errorf("Error with парсинга netTotal: %w", err)
	}

	netUsed, err := strconv.ParseFloat(netUsedStr, 64)
	if err != nil {
		return fmt.Errorf("Error with парсинга netUsed: %w", err)
	}

	if loadAverage > 30 {
		fmt.Printf("Load Average is too high: %.2f\n", loadAverage)
	}

	if memTotal > 0 {
		memUsage := memUsed / memTotal * 100
		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %.2f%%\n", memUsage)
		}
	}

	if diskTotal > 0 {
		diskUsedPct := diskUsed / diskTotal * 100
		if diskUsedPct > 90 {
			freeBytes := diskTotal - diskUsed
			freeMB := freeBytes / 1024.0 / 1024.0
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeMB)
		}
	}

	if netTotal > 0 {
		netUsedPct := netUsed / netTotal * 100
		if netUsedPct > 90 {
			freeNet := netTotal - netUsed
			freeMbit := (freeNet * 8) / (1024.0 * 1024.0)
			fmt.Printf("Network bandwidth usage high: %.2f Mbit/s available\n", freeMbit)
		}
	}

	return nil
}
