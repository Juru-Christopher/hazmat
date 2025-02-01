package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("Activating EvilEye.....")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("@$$$33N03V!L$$@")

	err := os.MkdirAll("captured_images", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	recMode := "time" //shots

	// Capture screenshots based on the user's choice
	switch recMode {
	case "time":
		// Capture screenshots for 5 seconds
		fmt.Println("Capturing screenshots for 5 seconds...")
		captureFor5Seconds()
	case "shots":
		//fmt.Println("Taking 5 screenshots...")
		take5Screenshots()
	default:
		fmt.Println("$$#$Y$73WC0NF1G3RR#$$> CAPMODE")
		recMode = "shots"
	}
}

func captureFor5Seconds() {
	// Set the duration for capturing screenshots
	duration := 5 * time.Second
	endTime := time.Now().Add(duration)

	// Capture screenshots every second until the time is up
	for time.Now().Before(endTime) {
		// Capture the screen
		img := robotgo.CaptureScreen()
		if img == nil {
			fmt.Println("Error capturing screen!")
			return
		}

		// Save the image to the 'captured_images' folder
		saveImage(robotgo.ToBitmap(img))
		time.Sleep(1 * time.Second) // Wait 1 second before taking the next screenshot
	}

	fmt.Println("Finished capturing screenshots for 5 seconds.")
}

// Function to take 5 screenshots
func take5Screenshots() {
	for i := 1; i <= 5; i++ {
		// Capture the screen
		img := robotgo.CaptureScreen()
		if img == nil {
			fmt.Println("Error capturing screen!")
			return
		}

		// Save the image to the 'captured_images' folder
		saveImage(robotgo.ToBitmap(img))

		// Wait 1 second before taking the next screenshot
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Finished taking 5 screenshots.")
}

// Function to save the captured image to the 'captured_images' folder
func saveImage(img robotgo.Bitmap) {
	// Create a file name for the image
	fileName := fmt.Sprintf("captured_images/screenshot_%d.png", time.Now().UnixNano())

	// Save the image as a PNG file
	if err := robotgo.SaveCapture(fileName); err != nil {
		fmt.Printf("Error saving screenshot: %v\n", err)
	} else {
		fmt.Printf("Screenshot sa
		ved as: %s\n", fileName)
	}
}
