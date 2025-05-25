package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Create a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to signal when to print system info
	printChan := make(chan struct{})

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start a goroutine to handle the printing
	wg.Add(1)
	go func() {
		defer wg.Done()
		printSystemInfo(printChan)
	}()

	// Create a ticker to signal every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Main loop to handle signals and ticker
	// as well as keep the main goroutine running
	for {
		select {
		case <-ticker.C:
			printChan <- struct{}{}
		case <-stopChan:
			fmt.Println("Received stop signal, stopping...")
			close(printChan)
			wg.Wait()
			return
		}
	}
}

func printSystemInfo(printChan chan struct{}) {
	for range printChan {
		cpuUsage, _ := GetCPUStats()

		memoryUsage, _ := GetMEMStats()

		runningProcesses, _ := GetProcesses(10)

		fmt.Println("CPU Percentage    :", cpuUsage)
		fmt.Println("Memory Percentage :", memoryUsage)
		fmt.Println("Running Processes :", runningProcesses)
	}
}
