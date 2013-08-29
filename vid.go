// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raspicam

import (
	"fmt"
	"strings"
)

const raspiVidCommmand = "raspivid"

// Vid represents the the configuration used to call raspivid
type Vid struct {
	Timeout       int // Delay (milliseconds) before image is taken
	Width, Height int // Image dimensions
	Bitrate       int // Requested bitrate
	Framerate     int // Requested framerate (fps)
	IntraPeriod   int // Intra-refresh period (key frame rate)

	// Flag to specify whether encoder works in place or creates a new buffer.
	// Result is preview can display either the camera output or the encoder
	// output (with compression artifacts).
	ImmutableInput bool
	Camera         Camera
	Preview        Preview
}

// The default Vid setup
// TODO: Framerate is set via a macro, should really call raspivid to get default
var defaultVid = Vid{Timeout: 5000, Width: 1920, Height: 1080, Bitrate: 17000000,
	Framerate: 30, ImmutableInput: true, Camera: defaultCamera, Preview: defaultPreview}

// String returns the parameter string for the given Vid struct
func (v *Vid) String() string {
	output := "--output -"
	if v.Timeout != defaultVid.Timeout {
		output += fmt.Sprintf(" --timeout %v", v.Timeout)
	}
	if v.Width != defaultVid.Width {
		output += fmt.Sprintf(" --width %v", v.Width)
	}
	if v.Height != defaultVid.Height {
		output += fmt.Sprintf(" --height %v", v.Height)
	}
	if v.Bitrate != defaultVid.Bitrate {
		output += fmt.Sprintf(" --bitrate %v", v.Bitrate)
	}
	if v.Framerate != defaultVid.Framerate {
		output += fmt.Sprintf(" --framerate %v", v.Framerate)
	}
	if v.IntraPeriod != defaultVid.IntraPeriod {
		output += fmt.Sprintf(" --intra %v", v.IntraPeriod)
	}
	output += " " + v.Camera.String()
	output += " " + v.Preview.String()
	return strings.TrimSpace(output)
}

// cmd returns the raspicam command associated with this struct
func (v *Vid) cmd() string {
	return raspiVidCommmand
}

// params returns the parameters to be used in the command execution
func (v *Vid) params() []string {
	return strings.Fields(v.String())
}

// NewVid returns a new *Vid struct setup with the default configuration
func NewVid() *Vid {
	newVid := defaultVid
	return &newVid
}
