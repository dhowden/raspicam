// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package raspicam

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const raspiStillCommand = "raspistill"

// Encoding represents an enumeration of the suported encoding types
type Encoding uint

const (
	EncodingJPEG Encoding = iota
	EncodingBMP
	EncodingGIF
	EncodingPNG
)

var encodings = [...]string{
	"jpg",
	"bmp",
	"gif",
	"png",
}

// String returns the parameter string for the given encoding
func (s Encoding) String() string {
	return encodings[s]
}

// The default still settings
var defaultStill = Still{Timeout: 5000, Width: 2592, Height: 1944, Quality: 85,
	Encoding: EncodingJPEG, Camera: NewCamera()}

// Still represents the configuration necessary to capture a still image
type Still struct {
	Timeout       int  // Delay (milliseconds) before image is taken
	Width, Height int  // Image dimensions
	Quality       int  // Quality (for lossy encoding)
	Raw           bool // Want a raw image?
	Timelapse     int  // The number of pictures to take in the given timeout interval
	Encoding      Encoding
	Camera        Camera
}

// String returns the parameter string for the given Still struct
func (s *Still) String() string {
	output := " --output -"
	if s.Timeout != defaultStill.Timeout {
		output += fmt.Sprintf(" --timeout %v", s.Timeout)
	}
	if s.Width != defaultStill.Width {
		output += fmt.Sprintf(" --width %v", s.Width)
	}
	if s.Height != defaultStill.Height {
		output += fmt.Sprintf(" --height %v", s.Height)
	}
	if s.Quality != defaultStill.Quality {
		output += fmt.Sprintf(" --quality %v", s.Quality)
	}
	if s.Raw != defaultStill.Raw {
		output += " --raw"
	}
	if s.Encoding != defaultStill.Encoding {
		output += fmt.Sprintf(" --encoding %v", s.Encoding)
	}
	output += s.Camera.String()
	return strings.TrimSpace(output)
}

// NewStill returns a *Still with the default values set by the raspistill command
// (see  userland/linux/apps/raspicam/RaspiStill.c)
func NewStill() *Still {
	newStill := defaultStill
	return &newStill
}

// Capture takes a still and writes the result to the given writer. Any
// errors are sent back on the given error channel, which is closed before
// the function returns
func (s *Still) Capture(w io.Writer, errCh chan<- error) {
	done := make(chan struct{})
	defer func() {
		<-done
		close(errCh)
	}()

	cmd := exec.Command(raspiStillCommand, strings.Fields(s.String())...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errCh <- err
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		errCh <- err
		return
	}

	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			errCh <- fmt.Errorf("%v: %v", raspiStillCommand, errScanner.Text())
		}
		if err := errScanner.Err(); err != nil {
			errCh <- err
		}
		close(done)
	}()

	if err := cmd.Start(); err != nil {
		errCh <- fmt.Errorf("starting: %v", err)
		return
	}
	defer func() {
		if err := cmd.Wait(); err != nil {
			errCh <- fmt.Errorf("waiting: %v", err)
		}
	}()

	_, err = io.Copy(w, stdout)
	if err != nil {
		errCh <- err
	}
}
