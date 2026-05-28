package ebiten

import (
	"errors"
	"image"
	"image/color"
	"sync"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
)

const (
	screenWidth       = 800
	screenHeight      = 600
	defaultSaturation = 1.0
	defaultBrightness = 1.0
	defaultSkewX      = 0.0
	defaultSkewY      = 0.0
	defaultScaleX     = 1.0
	defaultScaleY     = 1.0
)

type renderCallback = func(surface d2interface.Surface) error

type updateCallback = func() error

// static check that we implement our renderer interface
var _ d2interface.Renderer = &Renderer{}

// Renderer is an implementation of a renderer
type Renderer struct {
	gamma           float64
	contrast        float64
	offscreen       *ebiten.Image
	shader          *ebiten.Shader
	shaderErr       error
	shaderOnce      sync.Once
	updateCallback
	renderCallback
	*d2util.GlyphPrinter
	lastRenderError error
}

// Update calls the game's logical update function (the `Advance` method)
func (r *Renderer) Update() error {
	if r.updateCallback == nil {
		return errors.New("no update callback defined for ebiten renderer")
	}

	return r.updateCallback()
}

const drawError = "no render callback defined for ebiten renderer"

// Draw updates the screen with the given *ebiten.Image
func (r *Renderer) Draw(screen *ebiten.Image) {
	r.lastRenderError = nil

	if r.renderCallback == nil {
		r.lastRenderError = errors.New(drawError)
		return
	}

	// If no gamma or contrast adjustment is needed, draw directly to the screen
	if r.gamma == 1.0 && r.contrast == 1.0 {
		r.lastRenderError = r.renderCallback(createEbitenSurface(r, screen))
		return
	}

	// Use offscreen buffer for post-processing
	w, h := screen.Size()
	if r.offscreen == nil {
		r.offscreen = ebiten.NewImage(w, h)
	} else {
		ow, oh := r.offscreen.Size()
		if ow != w || oh != h {
			r.offscreen = ebiten.NewImage(w, h)
		}
	}

	r.offscreen.Fill(color.Transparent)
	r.lastRenderError = r.renderCallback(createEbitenSurface(r, r.offscreen))
	if r.lastRenderError != nil {
		return
	}

	// Apply effects using shader
	shader, err := r.getShader()
	if err == nil {
		opts := &ebiten.DrawRectShaderOptions{}
		opts.Images[0] = r.offscreen
		opts.Uniforms = map[string]interface{}{
			"Gamma":    float32(r.gamma),
			"Contrast": float32(r.contrast),
		}
		screen.DrawRectShader(w, h, shader, opts)
	} else {
		// Fallback to ColorM for contrast if shader fails
		// Gamma cannot be easily implemented with ColorM
		opts := &ebiten.DrawImageOptions{}
		if r.contrast != 1.0 {
			// Contrast formula in ColorM:
			// R''' = (R - 0.5) * contrast + 0.5 = R * contrast + (0.5 - 0.5 * contrast)
			c := r.contrast
			t := 0.5 - 0.5*c
			opts.ColorM.Scale(c, c, c, 1.0)
			opts.ColorM.Translate(t, t, t, 0)
		}
		screen.DrawImage(r.offscreen, opts)
	}
}

// Layout returns the renderer screen width and height
func (r *Renderer) Layout(_, _ int) (width, height int) {
	return screenWidth, screenHeight
}

// CreateRenderer creates an ebiten renderer instance
func CreateRenderer(cfg *d2config.Configuration) (*Renderer, error) {
	result := &Renderer{
		GlyphPrinter: d2util.NewDebugPrinter(),
		gamma:        1.0,
		contrast:     1.0,
	}

	if cfg != nil {
		config := cfg

		ebiten.SetCursorMode(ebiten.CursorModeHidden)
		ebiten.SetFullscreen(config.FullScreen)
		ebiten.SetRunnableOnUnfocused(config.RunInBackground)
		ebiten.SetVsyncEnabled(config.VsyncEnabled)
		ebiten.SetMaxTPS(config.TicksPerSecond)

		result.gamma = config.Gamma
		result.contrast = config.Contrast
	}

	return result, nil
}

// GetRendererName returns the name of the renderer
func (*Renderer) GetRendererName() string {
	return "Ebiten"
}

// SetWindowIcon sets the icon for the window, visible in the chrome of the window
func (*Renderer) SetWindowIcon(fileName string) {
	_, iconImage, err := ebitenutil.NewImageFromFile(fileName)
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{iconImage})
	}
}

// IsDrawingSkipped returns a bool for whether or not the drawing has been skipped
func (r *Renderer) IsDrawingSkipped() bool {
	return r.lastRenderError != nil
}

// Run initializes the renderer
func (r *Renderer) Run(f renderCallback, u updateCallback, width, height int, title string) error {
	r.renderCallback = f
	r.updateCallback = u

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(width, height)

	return ebiten.RunGame(r)
}

// CreateSurface creates a renderer surface from an existing surface
func (r *Renderer) CreateSurface(surface d2interface.Surface) (d2interface.Surface, error) {
	img := surface.(*ebitenSurface).image
	sfcState := surfaceState{
		filter:     ebiten.FilterNearest,
		effect:     d2enum.DrawEffectNone,
		saturation: defaultSaturation,
		brightness: defaultBrightness,
		skewX:      defaultSkewX,
		skewY:      defaultSkewY,
		scaleX:     defaultScaleX,
		scaleY:     defaultScaleY,
	}
	result := createEbitenSurface(r, img, sfcState)

	return result, nil
}

// NewSurface creates a new surface
func (r *Renderer) NewSurface(width, height int) d2interface.Surface {
	img := ebiten.NewImage(width, height)

	return createEbitenSurface(r, img)
}

// IsFullScreen returns a boolean for whether or not the renderer is currently set to fullscreen
func (r *Renderer) IsFullScreen() bool {
	return ebiten.IsFullscreen()
}

// SetFullScreen sets the renderer to fullscreen, given a boolean
func (r *Renderer) SetFullScreen(fullScreen bool) {
	ebiten.SetFullscreen(fullScreen)
}

// SetVSyncEnabled enables vsync, given a boolean
func (r *Renderer) SetVSyncEnabled(vsync bool) {
	ebiten.SetVsyncEnabled(vsync)
}

// GetVSyncEnabled returns a boolean for whether or not vsync is enabled
func (r *Renderer) GetVSyncEnabled() bool {
	return ebiten.IsVsyncEnabled()
}

// GetCursorPos returns the current cursor position x,y coordinates
func (r *Renderer) GetCursorPos() (x, y int) {
	return ebiten.CursorPosition()
}

// CurrentFPS returns the current frames per second of the renderer
func (r *Renderer) CurrentFPS() float64 {
	return ebiten.CurrentFPS()
}

// ShowPanicScreen shows a panic message in a forever loop
func (r *Renderer) ShowPanicScreen(message string) {
	errorScreen := CreatePanicScreen(message)

	err := ebiten.RunGame(errorScreen)
	if err != nil {
		panic(err)
	}
}

// SetGamma sets the gamma for the renderer
func (r *Renderer) SetGamma(gamma float64) {
	if gamma < 0.1 {
		gamma = 0.1
	}
	r.gamma = gamma
}

// SetContrast sets the contrast for the renderer
func (r *Renderer) SetContrast(contrast float64) {
	if contrast < 0.0 {
		contrast = 0.0
	}
	r.contrast = contrast
}

func (r *Renderer) getShader() (*ebiten.Shader, error) {
	r.shaderOnce.Do(func() {
		s, err := ebiten.NewShader([]byte(postProcessShaderCode))
		if err != nil {
			r.shaderErr = err
			return
		}
		r.shader = s
	})

	return r.shader, r.shaderErr
}

const postProcessShaderCode = `
package main

var Gamma float
var Contrast float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	origin, size := imageSrcRegionOnTexture()
	clr := imageSrc0At(texCoord)

	if clr.a == 0.0 {
		return clr
	}

	// Unpremultiply
	rgb := clr.rgb / clr.a

	// Apply Contrast
	// Contrast formula: (color - 0.5) * contrast + 0.5
	rgb = (rgb - 0.5) * Contrast + 0.5
	rgb = clamp(rgb, 0.0, 1.0)

	// Apply Gamma
	// Gamma formula: color ^ (1.0 / gamma)
	rgb = pow(rgb, vec3(1.0 / Gamma))

	return vec4(rgb * clr.a, clr.a)
}
`
