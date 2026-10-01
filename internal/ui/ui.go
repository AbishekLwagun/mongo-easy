package ui

import (
	"fmt"
	"sync"
	"time"
)

// These characters are shown one at a time to make the import look active.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Protects spinner startup so two spinners do not start at the same time.
var spinnerMu sync.Mutex

type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
}

func Header() {
	fmt.Println()
	fmt.Println("Mongo Easy")
	fmt.Println()
}

func Success(message string) {
	fmt.Printf("✔ %s\n", message)
}

func Error(message string) {
	fmt.Printf("✖ %s\n", message)
}

func Info(message string) {
	fmt.Printf("• %s\n", message)
}

func Warning(message string) {
	fmt.Printf("⚠ %s\n", message)
}

func Done() {
	fmt.Println()
	fmt.Println("Done.")
}

func StartSpinner(message string) *Spinner {
	// Only allow one spinner to be started at a time.
	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	s := &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}

	// Run the spinner in the background so the import can keep running.
	go s.run()

	return s
}

func (s *Spinner) run() {
	// Update the spinner every 100ms so it looks smooth without printing too much.
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Let StopSuccess or StopError know when the spinner goroutine has stopped.
	defer close(s.done)

	frame := 0

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\r%s %s", spinnerFrames[frame], s.message)
			frame = (frame + 1) % len(spinnerFrames)

		case <-s.stop:
			return
		}
	}
}

func (s *Spinner) StopSuccess(message string) {
	// Tell the spinner goroutine to stop and wait until it has finished.
	close(s.stop)
	<-s.done

	// Clear the old spinner line before showing the final result.
	fmt.Printf("\r\033[K")
	Success(message)
}

func (s *Spinner) StopError(message string) {
	// Stop the spinner before printing the error message.
	close(s.stop)
	<-s.done

	// Clear the old spinner line so the error is shown cleanly.
	fmt.Printf("\r\033[K")
	Error(message)
}
