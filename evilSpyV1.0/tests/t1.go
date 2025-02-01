package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/go-vgo/robotgo"
)

func main() {
	// Define the output file path
	outputFile := "captured_videos/output.mp4"

	// Get screen width and height dynamically
	w, h := robotgo.GetScreenSize()
	screen := fmt.Sprintf("%dx%d", w, h)

	// Create the output directory if it doesn't exist
	err := os.MkdirAll("captured_videos", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	fmt.Println("Activating EvilEye.....")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("@$$$33N03V!L$$@")

	// Set up the ffmpeg command for screen capture
	cmd := exec.Command("ffmpeg",
		"-y",            // Overwrite output file if exists
		"-f", "x11grab", // Use X11 grab for Linux screen capture
		"-s", screen, // Screen resolution (width x height)
		"-i", ":0.0", // Screen display identifier
		"-vcodec", "libx264", // Video codec
		"-framerate", "30", // Frame rate
		"-preset", "fast", // Encoding speed (ultrafast, fast, medium, etc.)
		"-crf", "23", // Constant Rate Factor for quality (lower is better)
		"-t", "00:00:10", // Duration of capture (10 seconds)
		outputFile, // Output file path
	)

	// Start the ffmpeg command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting ffmpeg command:", err)
		return
	}

	// Initialize keyboard input handling
	err = keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	// Wait for the F10 key to be pressed to exit
	fmt.Println("Press F10 to exit...")

	// Listen for keypresses
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				fmt.Println("Error reading keyboard input:", err)
				return
			}
			if key == keyboard.KeyF10 { // Check if F10 key is pressed
				// Exit the program when F10 is pressed
				fmt.Println("F10 key pressed, exiting...")
				cmd.Process.Kill() // Kill the ffmpeg process
				os.Exit(0)         // Exit the program
			}
			_ = char // We can ignore the char since we only care about the key
		}
	}()

	// Wait for the ffmpeg command to finish execution or exit signal from goroutine
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error during ffmpeg execution:", err)
		return
	}

	// Print completion message
	fmt.Println("Screen capture completed. Video saved to", outputFile)
}
