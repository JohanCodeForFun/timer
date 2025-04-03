package main

import (
	"fmt"
	"time"
)

// Timer structure to hold the pose and duration
type Timer struct {
	Pose     int
	Duration time.Duration
}

func main() {
	// Define the sequence of poses with their respective durations
	timers := []Timer{
		{Pose: 10, Duration: 30 * time.Second},
		{Pose: 5, Duration: 1 * time.Minute},
		{Pose: 2, Duration: 5 * time.Minute},
		{Pose: 1, Duration: 10 * time.Minute},
	}

	totalTime := 0 * time.Minute
	// Add up the total time for the entire sequence
	for _, t := range timers {
		totalTime += t.Duration
	}

	fmt.Printf("Starting timer for a total duration of %v...\n", totalTime)

	// Loop through the different timer configurations
	for _, timer := range timers {
		fmt.Printf("Starting %d poses for %v...\n", timer.Pose, timer.Duration)
		for i := 1; i <= timer.Pose; i++ {
			fmt.Printf("Pose %d/%d started. Duration: %v\n", i, timer.Pose, timer.Duration)
			time.Sleep(timer.Duration)
			fmt.Printf("Pose %d/%d completed.\n", i, timer.Pose)
		}
	}

	// Finish the timer
	fmt.Printf("All poses completed. Total time elapsed: %v\n", totalTime)
}
