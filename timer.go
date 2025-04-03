package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"
)

type Timer struct {
	Pose     int
	Duration time.Duration
}

type SessionData struct {
	DayStreak        int
	TotalOverallTime time.Duration
	LastRunDate      string
}

func main() {
	// Define the sequence of poses with their respective durations
	timers := []Timer{
		{Pose: 10, Duration: 30 * time.Second},
		{Pose: 5, Duration: 1 * time.Minute},
		{Pose: 2, Duration: 5 * time.Minute},
		{Pose: 1, Duration: 10 * time.Minute},
	}

	// Load session data from file
	sessionData := loadSessionData()

	// Check if it's a new day
	currentDate := time.Now().Format("2006-01-02")
	if sessionData.LastRunDate != currentDate {
		sessionData.DayStreak++ // New day, increment the streak
		sessionData.LastRunDate = currentDate
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
					fmt.Printf("\rTime remaining: %v (Press 'p' to pause, 'r' to restart, 'n' to skip to next, 'q' to quit)", timer.Duration-elapsed.Round(time.Second))
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
						// Round the total time elapsed to the nearest second before quitting
						fmt.Printf("\nTotal time elapsed: %v\n", totalElapsed.Round(time.Second))
						fmt.Printf("Total time spent across all sessions: %v\n", sessionData.TotalOverallTime.Round(time.Second))
						fmt.Printf("Day streak: %d days\n", sessionData.DayStreak)
						fmt.Println("Exiting timer.")
						// Save session data to file
						sessionData.TotalOverallTime += totalElapsed
						saveSessionData(sessionData)
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
							// Round the total time elapsed to the nearest second before quitting
							fmt.Printf("\nTotal time elapsed: %v\n", totalElapsed.Round(time.Second))
							fmt.Printf("Total time spent across all sessions: %v\n", sessionData.TotalOverallTime.Round(time.Second))
							fmt.Printf("Day streak: %d days\n", sessionData.DayStreak)
							fmt.Println("Exiting timer.")
							// Save session data to file
							sessionData.TotalOverallTime += totalElapsed
							saveSessionData(sessionData)
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

			// Countdown before moving to the next pose
			if i == timer.Pose { // Countdown only after completing the last pose of the current set
				fmt.Println("Starting countdown to next pose:")
				for j := 5; j > 0; j-- {
					fmt.Printf("%d\n", j)
					time.Sleep(1 * time.Second)
				}
			}
		}
		currentPoseIndex++
	}

	// If we finish the entire sequence, print the total time
	fmt.Printf("All poses completed. Total time elapsed: %v\n", totalElapsed.Round(time.Second))
	// Save session data to file
	sessionData.TotalOverallTime += totalElapsed
	saveSessionData(sessionData)
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

// Function to load the session data from file
func loadSessionData() SessionData {
	file, err := os.Open("session_data.gob")
	if err != nil {
		// If file does not exist, return default values
		return SessionData{
			DayStreak:        0,
			TotalOverallTime: 0 * time.Second,
			LastRunDate:      "",
		}
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var sessionData SessionData
	err = decoder.Decode(&sessionData)
	if err != nil {
		fmt.Println("Error decoding session data:", err)
		return SessionData{}
	}

	return sessionData
}

// Function to save the session data to file
func saveSessionData(sessionData SessionData) {
	file, err := os.Create("session_data.gob")
	if err != nil {
		fmt.Println("Error saving session data:", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(sessionData)
	if err != nil {
		fmt.Println("Error encoding session data:", err)
	}
}
