// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raspicam

import (
	"fmt"
	"strings"
)

const raspiStillCommand = "raspistill"
const raspiStillYUVCommand = "raspiyuv"

// Encoding represents an enumeration of the supported encoding types
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

// BaseStill represents the common elements between a Still and StillYUV
// as described in their equivalents found in RaspiStill.c and RaspiStillYUV.c
// respectively.
type BaseStill struct {
	Timeout       int // Delay (milliseconds) before image is taken
	Width, Height int // Image dimensions
	Timelapse     int // The number of pictures to take in the given timeout interval
	Camera        Camera
	Preview       Preview
}

// The default BaseStill setup
var defaultBaseStill = BaseStill{Timeout: 5000, Width: 2592, Height: 1944,
	Camera: defaultCamera, Preview: defaultPreview}

// String returns the parameter string for the given BaseStill
func (s *BaseStill) String() string {
	output := "--output -"
	if s.Timeout != defaultStill.Timeout {
		output += fmt.Sprintf(" --timeout %v", s.Timeout)
	}
	if s.Width != defaultStill.Width {
		output += fmt.Sprintf(" --width %v", s.Width)
	}
	if s.Height != defaultStill.Height {
		output += fmt.Sprintf(" --height %v", s.Height)
	}
	output += " " + s.Camera.String()
	output += " " + s.Preview.String()
	return strings.TrimSpace(output)
}

// The default Still setup
var defaultStill = Still{BaseStill: defaultBaseStill, Quality: 85,
	Encoding: EncodingJPEG}

// Still represents the configuration necessary to call raspistill
type Still struct {
	BaseStill
	Quality  int  // Quality (for lossy encoding)
	Raw      bool // Want a raw image?
	Encoding Encoding
}

// String returns the parameter string for the given Still struct
func (s *Still) String() string {
	output := s.BaseStill.String()
	if s.Quality != defaultStill.Quality {
		output += fmt.Sprintf(" --quality %v", s.Quality)
	}
	if s.Raw {
		output += " --raw"
	}
	if s.Encoding != defaultStill.Encoding {
		output += fmt.Sprintf(" --encoding %v", s.Encoding)
	}
	return strings.TrimSpace(output)
}

// cmd returns the raspicam command associated with this struct
func (s *Still) cmd() string {
	return raspiStillCommand
}

// params returns the parameters to be used in the command execution
func (s *Still) params() []string {
	return strings.Fields(s.String())
}

// NewStill returns a *Still with the default values set by the raspistill command
// (see userland/linux/apps/raspicam/RaspiStill.c)
func NewStill() *Still {
	newStill := defaultStill
	return &newStill
}

// StillYUV represents the configuration necessary to call raspistillYUV
type StillYUV struct {
	BaseStill
	UseRGB bool // Output RGB data rather than YUV
}

// The default StillYUV setup
var defaultStillYUV = StillYUV{BaseStill: defaultBaseStill}

// String returns the parameter string for the given StillYUV struct
func (s *StillYUV) String() string {
	output := s.BaseStill.String()
	if s.UseRGB {
		output += " --rgb"
	}
	return strings.TrimSpace(output)
}

// cmd returns the raspicam command associated with this struct
func (s *StillYUV) cmd() string {
	return raspiStillYUVCommand
}

// params returns the parameters to be used in the command execution
func (s *StillYUV) params() []string {
	return strings.Fields(s.String())
}

// NewStill returns a *StillYUV with the default values set by the raspistillYUV command
// (see userland/linux/apps/raspicam/RaspiStillYUV.c)
func NewStillYUV() *StillYUV {
	newStillYUV := defaultStillYUV
	return &newStillYUV
}
