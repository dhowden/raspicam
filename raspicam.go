// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package raspicam provides basic Go APIs for interating with the Raspberry Pi
// camera untilities.
//
// Currently this is done by calling the existing raspicam command and
// capturing the output from stdout/stderr.  Eventually we would like
// to call the C APIs directly.
package raspicam

import (
	"fmt"
	"strings"
)

// ExposureMode is an enumeration of supported exposure modes
type ExposureMode uint

const (
	ExposureOff ExposureMode = iota
	ExposureAuto
	ExposureNight
	ExposureNightPreview
	ExposureBacklight
	ExposureSpotlight
	ExposureSports
	ExposureSnow
	ExposureBeach
	ExposureVerylong
	ExposureFixedFPS
	ExposureAntishake
	ExposureFireworks
)

var exposureModes = [...]string{
	"off",
	"auto",
	"night",
	"nightpreview",
	"backlight",
	"spotlight",
	"sports",
	"snow",
	"beach",
	"verylong",
	"fixedfps",
	"antishake",
	"fireworks",
}

// String returns the command line parameter for the given ExposureMode
func (e ExposureMode) String() string { return exposureModes[e] }

// An MeteringMode specificies an exposure metering mode
type MeteringMode uint

const (
	MeteringAverage MeteringMode = iota
	MeteringSpot
	MeteringBacklit
	MeteringMatrix
)

var exposureMeteringMode = [...]string{
	"average",
	"spot",
	"backlit",
	"matrix",
}

// String returns the command line parameter for the given MeteringMode
func (m MeteringMode) String() string { return exposureMeteringMode[m] }

// An AWBMode is an enumeration of the auto white balance modes
type AWBMode uint

const (
	AWBOff AWBMode = iota
	AWBAuto
	AWBSunlight
	AWBCloudy
	AWBShade
	AWBTungsten
	AWBFluorescent
	AWBIncandescent
	AWBFlash
	AWBHorizon
)

var awbModes = [...]string{
	"off",
	"auto",
	"sun",
	"cloud",
	"shade",
	"tungsten",
	"fluorescent",
	"incandescent",
	"flash",
	"horizon",
}

// String returns the command line parameter for the given AWBMode
func (a AWBMode) String() string { return awbModes[a] }

// An ImageFX specifies an image effect for the camera
type ImageFX uint

const (
	FXNone ImageFX = iota
	FXNegative
	FXSolarize
	FXPosterize
	FXWhiteboard
	FXBlackboard
	FXSketch
	FXDenoise
	FXEmboss
	FXOilpaint
	FXHatch
	FXGpen
	FXPastel
	FXWatercolour
	FXFilm
	FXBlur
	FXSaturation
	FXColourswap
	FXWashedout
	FXPosterise
	FXColourpoint
	FXColourbalance
	FXCartoon
)

var imageFXModes = [...]string{
	"none",
	"negative",
	"solarise",
	"sketch",
	"denoise",
	"emboss",
	"oilpaint",
	"hatch",
	"gpen",
	"pastel",
	"watercolour",
	"film",
	"blur",
	"saturation",
	"colourswap",
	"washedout",
	"posterise",
	"colourpoint",
	"colourbalance",
	"cartoon",
}

// String returns the command-line parameter for the given imageFX
func (i ImageFX) String() string { return imageFXModes[i] }

// ColourFX represents colour effects parameters
type ColourFX struct {
	Enabled bool
	U, V    int
}

func (c ColourFX) String() string {
	return fmt.Sprintf("%v:%v", c.U, c.V)
}

// FloatRect contains the information necessary to construct a rectangle
// with dimensions in floating point.
type FloatRect struct {
	X, Y, W, H float64
}

// String returns the command parameter for the given FloatRect
func (r *FloatRect) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", r.X, r.Y, r.W, r.H)
}

var defaultRegionOfInterest = FloatRect{W: 1.0, H: 1.0}

// Camera represents a camera configuration
type Camera struct {
	Sharpness            int // -100 to 100
	Contrast             int // -100 to 100
	Brightness           int // 0 to 100
	Saturation           int // -100 to 100
	ISO                  int // TODO: what range? (see RaspiCamControl.h)
	VideoStabilisation   bool
	ExposureCompensation int // -10 to 10? (see RaspiCamControl.h)
	ExposureMode         ExposureMode
	MeteringMode         MeteringMode
	AWBMode              AWBMode
	ImageEffect          ImageFX
	ColourEffects        ColourFX
	Rotation             int // 0 to 359
	HFlip, VFlip         bool
	RegionOfInterest     FloatRect // Assumes Normalised to [0.0,1.0]
}

var defaultCamera = Camera{Brightness: 50, ISO: 400, ExposureMode: ExposureAuto,
	MeteringMode: MeteringAverage, AWBMode: AWBAuto, ImageEffect: FXNone,
	ColourEffects: ColourFX{U: 128, V: 128}, RegionOfInterest: defaultRegionOfInterest}

// NewCamera constructs a default instance of the Camera struct based on
// `raspicamcontrol_set_defaults` in RaspiCamControl.c
func NewCamera() Camera {
	return defaultCamera
}

// String returns the parameters necessary to construct the
// equivalent command line arguments for the raspicam tools
func (c *Camera) String() string {
	output := ""
	if c.Sharpness != defaultCamera.Sharpness {
		output += fmt.Sprintf(" --sharpness %v", c.Sharpness)
	}
	if c.Contrast != defaultCamera.Contrast {
		output += fmt.Sprintf(" --contrast %v", c.Contrast)
	}
	if c.Brightness != defaultCamera.Brightness {
		output += fmt.Sprintf(" --brightness %v", c.Brightness)
	}
	if c.Saturation != defaultCamera.Saturation {
		output += fmt.Sprintf(" --saturation %v", c.Saturation)
	}
	if c.ISO != defaultCamera.ISO {
		output += fmt.Sprintf(" --ISO %v", c.ISO)
	}
	if c.VideoStabilisation {
		output += " --vstab"
	}
	if c.ExposureCompensation != defaultCamera.ExposureCompensation {
		output += fmt.Sprintf(" --ev %v", c.ExposureCompensation)
	}
	if c.ExposureMode != defaultCamera.ExposureMode {
		output += fmt.Sprintf(" --exposure %v", c.ExposureMode)
	}
	if c.MeteringMode != defaultCamera.MeteringMode {
		output += fmt.Sprintf(" --metering %v", c.MeteringMode)
	}
	if c.AWBMode != defaultCamera.AWBMode {
		output += fmt.Sprintf(" --awb %v", c.AWBMode)
	}
	if c.ImageEffect != defaultCamera.ImageEffect {
		output += fmt.Sprintf(" --imxfx %v", c.ImageEffect)
	}
	if c.ColourEffects.Enabled {
		output += fmt.Sprintf(" --colfx %v", c.ColourEffects)
	}
	if c.MeteringMode != defaultCamera.MeteringMode {
		output += fmt.Sprintf(" --metering %v", c.MeteringMode)
	}
	if c.Rotation != defaultCamera.Rotation {
		output += fmt.Sprintf(" --rotation %v", c.Rotation)
	}
	if c.HFlip {
		output += " --hflip"
	}
	if c.VFlip {
		output += " --vflip"
	}
	if c.RegionOfInterest != defaultCamera.RegionOfInterest {
		output += fmt.Sprintf(" --roi %v", c.RegionOfInterest)
	}
	return strings.TrimSpace(output)
}

type Rect struct {
	X, Y, Width, Height uint32
}

// PreviewMode represents an enumeration of preview modes
type PreviewMode int

const (
	PreviewDisabled PreviewMode = iota
	PreviewEnabled              // Preview enabled
	PreviewFullscreen
)

var previewModes = [...]string{
	"--nopreview",
	"--preview",
	"--fullscreen",
}

// String returns the parameter string for the given PreviewMode
func (p PreviewMode) String() string { return previewModes[p] }

// Preview contains the settings for the camera previews
type Preview struct {
	Mode    PreviewMode
	Opacity int  // Opacity of window - 0 = transparent, 255 = opaque
	Rect    Rect // Used when Mode is PreviewEnabled
}

var defaultPreview = Preview{Mode: PreviewFullscreen, Opacity: 255,
	Rect: Rect{X: 0, Y: 0, Width: 1024, Height: 768}}

// String returns the parameter string for the given Preview
func (p *Preview) String() string {
	output := ""
	if p.Mode != defaultPreview.Mode {
		output += " " + p.String()
	}
	if p.Opacity != defaultPreview.Opacity {
		output += fmt.Sprintf(" --opacity %v", p.Opacity)
	}
	return strings.TrimSpace(output)
}
