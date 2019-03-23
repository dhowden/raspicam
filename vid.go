// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raspicam

import (
	"time"
)

// DefaultRaspiVidCommmand is the default command for capturing video.
var DefaultRaspiVidCommmand = "raspivid"

// Vid represents the the configuration used to call raspivid.
type Vid struct {
	Timeout       time.Duration // Delay before image is taken
	Width, Height int           // Image dimensions
	Bitrate       int           // Requested bitrate
	Framerate     int           // Requested framerate (fps)
	IntraPeriod   int           // Intra-refresh period (key frame rate)

	// Flag to specify whether encoder works in place or creates a new buffer.
	// Result is preview can display either the camera output or the encoder
	// output (with compression artifacts).
	ImmutableInput bool
	Camera         Camera
	Preview        Preview

	// The command to run when making a video capture.  If left blank, the default
	// command is used.
	Command string

	// Additional arguments.  Default is empty.
	Args []string
}

// The default Vid setup.
// TODO: Framerate is set via a macro, should really call raspivid to get default.
var defaultVid = Vid{
	Timeout:        5 * time.Second,
	Width:          1920,
	Height:         1080,
	Bitrate:        17000000,
	Framerate:      30,
	ImmutableInput: true,
	Camera:         defaultCamera,
	Preview:        defaultPreview,
}

// String returns the parameter string for the given Vid struct.
func (v *Vid) String() string {
	return paramString(v)
}

func (v *Vid) params() []string {
	var out params
	out.add("--output", "-")
	if v.Timeout != defaultVid.Timeout {
		out.addInt64("--timeout", int64(v.Timeout/time.Millisecond))
	}
	if v.Width != defaultVid.Width {
		out.addInt("--width", v.Width)
	}
	if v.Height != defaultVid.Height {
		out.addInt("--height", v.Height)
	}
	if v.Bitrate != defaultVid.Bitrate {
		out.addInt("--bitrate", v.Bitrate)
	}
	if v.Framerate != defaultVid.Framerate {
		out.addInt("--framerate", v.Framerate)
	}
	if v.IntraPeriod != defaultVid.IntraPeriod {
		out.addInt("--intra", v.IntraPeriod)
	}
	out.add(v.Camera.params()...)
	out.add(v.Preview.params()...)
	return out
}

// Cmd returns the raspicam command for a Vid.
func (v *Vid) Cmd() string {
	if v.Command != "" {
		return v.Command
	}
	return DefaultRaspiVidCommmand
}

// Params returns the parameters to be used in the command execution.
func (v *Vid) Params() []string {
	return append(v.params(), v.Args...)
}

// NewVid returns a new *Vid struct setup with the default configuration.
func NewVid() *Vid {
	newVid := defaultVid
	return &newVid
}
