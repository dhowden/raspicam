// Copyright 2013, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package raspicam provides basic Go APIs for interacting with the Raspberry Pi
// camera.
//
// All captures are prepared by first creating a CaptureCommand (Still, StillYUV or
// Vid structs via calls to the NewStill, NewStillYUV or NewVid functions respectively).
// The Capture function can then be used to perform the capture.
package raspicam

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ExposureMode is an enumeration of supported exposure modes.
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

// String returns the command line parameter for the given ExposureMode.
func (e ExposureMode) String() string { return exposureModes[e] }

// A MeteringMode specificies an exposure metering mode.
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

// String returns the command line parameter for the given MeteringMode.
func (m MeteringMode) String() string { return exposureMeteringMode[m] }

// An AWBMode is an enumeration of the auto white balance modes.
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

// String returns the command line parameter for the given AWBMode.
func (a AWBMode) String() string { return awbModes[a] }

// An ImageFX specifies an image effect for the camera.
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
	FXColourSwap
	FXWashedOut
	FXPosterise
	FXColourPoint
	FXColourBalance
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

// String returns the command-line parameter for the given imageFX.
func (i ImageFX) String() string { return imageFXModes[i] }

// ColourFX represents colour effects parameters.
type ColourFX struct {
	Enabled bool
	U, V    int
}

// String returns the command parameter for the given ColourFX.
func (c ColourFX) String() string {
	return fmt.Sprintf("%v:%v", c.U, c.V)
}

// FloatRect contains the information necessary to construct a rectangle
// with dimensions in floating point.
type FloatRect struct {
	X, Y, W, H float64
}

// String returns the command parameter for the given FloatRect.
func (r *FloatRect) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", r.X, r.Y, r.W, r.H)
}

// The default RegionOfInterest setup.
var defaultRegionOfInterest = FloatRect{W: 1.0, H: 1.0}

// Camera represents a camera configuration.
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
	ShutterSpeed         time.Duration
}

// The default Camera setup.
var defaultCamera = Camera{
	Brightness:       50,
	ISO:              400,
	ExposureMode:     ExposureAuto,
	MeteringMode:     MeteringAverage,
	AWBMode:          AWBAuto,
	ImageEffect:      FXNone,
	ColourEffects:    ColourFX{U: 128, V: 128},
	RegionOfInterest: defaultRegionOfInterest,
}

// String returns the parameters necessary to construct the
// equivalent command line arguments for the raspicam tools.
func (c *Camera) String() string {
	return paramString(c)
}

// params is a wrapper around a string slice which adds convenience
// methods for adding different types of parameters
type params []string

func (ps *params) add(xs ...string)           { *ps = append(*ps, xs...) }
func (ps *params) addInt(x string, n int)     { *ps = append(*ps, x, strconv.Itoa(n)) }
func (ps *params) addInt64(x string, n int64) { *ps = append(*ps, x, strconv.FormatInt(n, 10)) }

func paramString(x interface{ params() []string }) string {
	return strings.Join(x.params(), " ")
}

func (c *Camera) params() []string {
	var out params
	if c.Sharpness != defaultCamera.Sharpness {
		out.addInt("--sharpness", c.Sharpness)
	}
	if c.Contrast != defaultCamera.Contrast {
		out.addInt("--contrast", c.Contrast)
	}
	if c.Brightness != defaultCamera.Brightness {
		out.addInt("--brightness", c.Brightness)
	}
	if c.Saturation != defaultCamera.Saturation {
		out.addInt("--saturation", c.Saturation)
	}
	if c.ISO != defaultCamera.ISO {
		out.addInt("--ISO", c.ISO)
	}
	if c.VideoStabilisation {
		out.add("--vstab")
	}
	if c.ExposureCompensation != defaultCamera.ExposureCompensation {
		out.addInt("--ev", c.ExposureCompensation)
	}
	if c.ExposureMode != defaultCamera.ExposureMode {
		out.add("--exposure", c.ExposureMode.String())
	}
	if c.MeteringMode != defaultCamera.MeteringMode {
		out.add("--metering", c.MeteringMode.String())
	}
	if c.AWBMode != defaultCamera.AWBMode {
		out.add("--awb", c.AWBMode.String())
	}
	if c.ImageEffect != defaultCamera.ImageEffect {
		out.add("--imxfx", c.ImageEffect.String())
	}
	if c.ColourEffects.Enabled {
		out.add("--colfx", c.ColourEffects.String())
	}
	if c.MeteringMode != defaultCamera.MeteringMode {
		out.add("--metering", c.MeteringMode.String())
	}
	if c.Rotation != defaultCamera.Rotation {
		out.addInt("--rotation", c.Rotation)
	}
	if c.HFlip {
		out.add("--hflip")
	}
	if c.VFlip {
		out.add("--vflip")
	}
	if c.RegionOfInterest != defaultCamera.RegionOfInterest {
		out.add("--roi", c.RegionOfInterest.String())
	}
	if c.ShutterSpeed != defaultCamera.ShutterSpeed {
		out.addInt64("--shutter", int64(c.ShutterSpeed/time.Microsecond))
	}
	return out
}

// Rect represents a rectangle defined by integer parameters.
type Rect struct {
	X, Y, Width, Height uint32
}

// String returns the parameter string for the given Rect.
func (r *Rect) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", r.X, r.Y, r.Width, r.Height)
}

// PreviewMode represents an enumeration of preview modes.
type PreviewMode uint

const (
	PreviewFullscreen PreviewMode = iota // Enabled by default
	PreviewWindow
	PreviewDisabled
)

var previewModes = [...]string{
	"fullscreen",
	"preview",
	"nopreview",
}

// String returns the parameter string for the given PreviewMode.
func (p PreviewMode) String() string { return previewModes[p] }

// Preview contains the settings for the camera previews.
type Preview struct {
	Mode    PreviewMode
	Opacity int  // Opacity of window (0 = transparent, 255 = opaque)
	Rect    Rect // Used when Mode is PreviewWindow
}

// The default Preview setup.
var defaultPreview = Preview{
	Mode:    PreviewFullscreen,
	Opacity: 255,
	Rect:    Rect{X: 0, Y: 0, Width: 1024, Height: 768},
}

// String returns the parameter string for the given Preview.
func (p *Preview) String() string {
	return paramString(p)
}

func (p *Preview) params() []string {
	var out params
	if p.Mode == PreviewWindow {
		out.add("--"+p.Mode.String(), p.Rect.String())
	} else {
		if p.Mode != defaultPreview.Mode {
			out.add("--" + p.Mode.String())
		}
	}
	if p.Opacity != defaultPreview.Opacity {
		out.addInt("--opacity", p.Opacity)
	}
	return out
}

// CaptureCommand represents a prepared capture command.
type CaptureCommand interface {
	Cmd() string
	Params() []string
}

// Capture runs the given CaptureCommand and writes the result to the given
// writer. Any errors are sent back on the given error channel, which is closed
// before the function returns.
func Capture(c CaptureCommand, w io.Writer, errCh chan<- error) {
	done := make(chan struct{})
	defer func() {
		<-done
		close(errCh)
	}()

	cmd := exec.Command(c.Cmd(), c.Params()...)

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
			errCh <- fmt.Errorf("%v: %v", c.Cmd(), errScanner.Text())
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
