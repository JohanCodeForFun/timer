package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

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

	var totalElapsed time.Duration
	var currentPoseIndex int
	fmt.Printf("Starting timer for a total duration of %v...\n", totalTime)

	// Loop through the different timer configurations
	for currentPoseIndex < len(timers) {
		timer := timers[currentPoseIndex]
		fmt.Printf("Starting %d poses for %v...\n", timer.Pose, timer.Duration)

		for i := 1; i <= timer.Pose; i++ {
			// Start the timer for the current pose
			fmt.Printf("Pose %d/%d started. Duration: %v\n", i, timer.Pose, timer.Duration)

			// Track the time for the current pose
			poseStart := time.Now()
			paused := false
			movingToNext := false

			// Timer for the current pose
			for elapsed := 0 * time.Second; elapsed < timer.Duration; elapsed = time.Since(poseStart) {
				// Print the remaining time only once per loop
				if !movingToNext {
					fmt.Printf("\rTime remaining: %v (Press 'p' to pause, 'r' to restart, 'n' to skip to next, 'q' to quit)", timer.Duration-elapsed)
				}

				// Sleep for 1 second to update the time remaining
				time.Sleep(1 * time.Second)

				// Handle user input after each second
				keyPress := checkKeyPress() // Properly get user input here
				if keyPress != "" {
					switch keyPress {
					case "p": // Pause
						paused = true
						poseStart = time.Now() // reset the start time
						fmt.Println("\nTimer paused. Press 'r' to resume, 'n' to move to next pose, or 'q' to quit.")
					case "r": // Restart the current pose
						poseStart = time.Now() // reset the start time
						fmt.Println("\nTimer restarted.")
					case "n": // Move to next pose
						if !movingToNext {
							movingToNext = true
							fmt.Println("\nMoving to next pose.")
							totalElapsed += time.Since(poseStart)
						}
					case "q": // Quit timer
						fmt.Printf("\nTotal time elapsed: %v\n", totalElapsed)
						fmt.Println("Exiting timer.")
						return
					}
				}

				// Pause handling
				if paused {
					// Wait for the user to resume
					for {
						keyPress := checkKeyPress()
						if keyPress == "r" {
							paused = false
							break
						}
						if keyPress == "q" {
							fmt.Printf("\nTotal time elapsed: %v\n", totalElapsed)
							fmt.Println("Exiting timer.")
							return
						}
					}
				}

				// Exit if we are moving to the next pose
				if movingToNext {
					break
				}
			}

			// Timer for the pose completed
			if !movingToNext {
				fmt.Printf("\nPose %d/%d completed.\n", i, timer.Pose)
				totalElapsed += time.Since(poseStart)
			}
		}
		currentPoseIndex++
	}

	// If we finish the entire sequence, print the total time
	fmt.Printf("All poses completed. Total time elapsed: %v\n", totalElapsed)
}

// Function to check for key presses from the user
func checkKeyPress() string {
	// Create a reader to capture user input
	reader := bufio.NewReader(os.Stdin)
	// Prompt the user to press a key
	fmt.Print("\nPress 'p' to pause, 'r' to restart, 'n' to move to next, or 'q' to quit: ")
	// Capture the input
	input, _ := reader.ReadString('\n')
	// Return the trimmed input to avoid newlines
	return strings.TrimSpace(input)
}
