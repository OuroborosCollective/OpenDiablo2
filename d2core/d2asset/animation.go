package d2asset

import (
	"errors"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math"
)

type playMode int

const (
	playModePause playMode = iota
	playModeForward
	playModeBackward
)

const (
	defaultPlayLength = 1.0

	// ARE/Ouroboros deterministic constants
	kappaInvariant     = 1000
	fieldLogicVersion  = "ARE_FIELDLOGIC_V1"
	collectiveMarkgraf = "COLLECTIVE_MARKGRAF"
	ouroborosSeal      = "OUROBOROS_LOOP"
)

type AnimationAxiom string

const (
	AxiomKappaStable      AnimationAxiom = "KAPPA_POS_1000_STABLE"
	AxiomFrameBounded     AnimationAxiom = "FRAME_INDEX_BOUNDED"
	AxiomDirectionBounded AnimationAxiom = "DIRECTION_INDEX_BOUNDED"
	AxiomNoTimeMutation   AnimationAxiom = "NO_EXTERNAL_TIME_MUTATION"
	AxiomReplayable       AnimationAxiom = "REPLAYABLE_FIELDLOGIC"
)

type FieldLogicSignature struct {
	Version        string
	Collective     string
	Ouroboros      string
	KappaPos       int
	DirectionIndex int
	FrameIndex     int
	FrameCount     int
	PlayedCount    int
	PlayMode       playMode
	PlayLoop       bool
	HasSubLoop     bool
	SubStart       int
	SubEnd         int
	LogicString    string
	Hash64         uint64
	Axioms         []AnimationAxiom
}

type animationFrame struct {
	decoded bool

	width   int
	height  int
	offsetX int
	offsetY int

	image d2interface.Surface
}

type animationDirection struct {
	decoded bool
	frames  []animationFrame
}

// static check that we implement the animation interface
var _ d2interface.Animation = &Animation{}

// Animation has directionality, play modes, deterministic ARE field logic,
// frame counting, optional subloops, and deterministic replay signatures.
type Animation struct {
	renderer       d2interface.Renderer
	onBindRenderer func(renderer d2interface.Renderer) error

	directions []animationDirection

	effect   d2enum.DrawEffect
	colorMod color.Color

	frameIndex     int
	directionIndex int

	lastFrameTime float64
	playedCount   int

	playMode   playMode
	playLength float64

	subStartingFrame int
	subEndingFrame   int

	originAtBottom bool
	playLoop       bool
	hasSubLoop     bool
	hasShadow      bool

	// ARE / Ouroboros metadata
	kappaPos        int
	fieldLogicName  string
	collectiveStamp string
	lastLogicHash   uint64
}

// ensureFieldDefaults keeps older constructors safe.
// This avoids requiring every caller to know about the new ARE fields.
func (a *Animation) ensureFieldDefaults() {
	if a.playLength <= 0 {
		a.playLength = defaultPlayLength
	}

	if a.kappaPos == 0 {
		a.kappaPos = kappaInvariant
	}

	if a.fieldLogicName == "" {
		a.fieldLogicName = fieldLogicVersion
	}

	if a.collectiveStamp == "" {
		a.collectiveStamp = collectiveMarkgraf
	}
}

// validateKappa enforces the deterministic ARE invariant.
func (a *Animation) validateKappa() error {
	a.ensureFieldDefaults()

	if a.kappaPos != kappaInvariant {
		return fmt.Errorf("kappa invariant violated: got=%d want=%d", a.kappaPos, kappaInvariant)
	}

	return nil
}

// validateDirection checks whether the active direction is safe.
func (a *Animation) validateDirection() error {
	if len(a.directions) == 0 {
		return errors.New("animation has no directions")
	}

	if a.directionIndex < 0 || a.directionIndex >= len(a.directions) {
		return fmt.Errorf("invalid direction index: %d", a.directionIndex)
	}

	return nil
}

// validateFrame checks whether the active frame is safe.
func (a *Animation) validateFrame() error {
	if err := a.validateDirection(); err != nil {
		return err
	}

	frameCount := len(a.directions[a.directionIndex].frames)
	if frameCount == 0 {
		return errors.New("animation direction has no frames")
	}

	if a.frameIndex < 0 || a.frameIndex >= frameCount {
		return fmt.Errorf("invalid frame index: %d", a.frameIndex)
	}

	return nil
}

// validateSubLoop validates exclusive end-frame subloop ranges.
func (a *Animation) validateSubLoop() error {
	if !a.hasSubLoop {
		return nil
	}

	frameCount := a.GetFrameCount()
	if frameCount <= 0 {
		return errors.New("cannot validate subloop without frames")
	}

	if a.subStartingFrame < 0 {
		return errors.New("subloop start frame cannot be negative")
	}

	if a.subEndingFrame <= a.subStartingFrame {
		return errors.New("subloop end frame must be greater than start frame")
	}

	if a.subEndingFrame > frameCount {
		return fmt.Errorf("subloop end frame %d exceeds frame count %d", a.subEndingFrame, frameCount)
	}

	return nil
}

// SetSubLoop sets a sub loop for the animation.
// endFrame is exclusive, matching the original OpenDiablo2 behavior.
func (a *Animation) SetSubLoop(startFrame, endFrame int) {
	a.ensureFieldDefaults()

	frameCount := a.GetFrameCount()
	if frameCount <= 0 {
		a.hasSubLoop = false
		return
	}

	if startFrame < 0 {
		startFrame = 0
	}

	if endFrame > frameCount {
		endFrame = frameCount
	}

	if endFrame <= startFrame {
		a.hasSubLoop = false
		return
	}

	a.subStartingFrame = startFrame
	a.subEndingFrame = endFrame
	a.hasSubLoop = true
}

// Advance advances the animation state deterministically.
// elapsed must be supplied by the caller's deterministic tick system.
func (a *Animation) Advance(elapsed float64) error {
	a.ensureFieldDefaults()

	if err := a.validateKappa(); err != nil {
		return err
	}

	if a.playMode == playModePause {
		return nil
	}

	if elapsed < 0 {
		return errors.New("elapsed time cannot be negative")
	}

	frameCount := a.GetFrameCount()
	if frameCount <= 0 {
		return errors.New("cannot advance animation without frames")
	}

	if err := a.validateFrame(); err != nil {
		return err
	}

	if err := a.validateSubLoop(); err != nil {
		return err
	}

	if a.playLength <= 0 {
		return errors.New("play length must be greater than zero")
	}

	frameLength := a.playLength / float64(frameCount)
	if frameLength <= 0 || math.IsNaN(frameLength) || math.IsInf(frameLength, 0) {
		return errors.New("invalid frame length")
	}

	a.lastFrameTime += elapsed

	framesAdvanced := int(a.lastFrameTime / frameLength)
	a.lastFrameTime -= float64(framesAdvanced) * frameLength

	for i := 0; i < framesAdvanced; i++ {
		startIndex := 0
		endIndex := frameCount

		if a.hasSubLoop && a.playedCount > 0 {
			startIndex = a.subStartingFrame
			endIndex = a.subEndingFrame
		}

		switch a.playMode {
		case playModeForward:
			a.frameIndex++

			if a.frameIndex >= endIndex {
				a.playedCount++

				if a.playLoop {
					a.frameIndex = startIndex
				} else {
					a.frameIndex = endIndex - 1
					a.updateLogicHash()
					return nil
				}
			}

		case playModeBackward:
			a.frameIndex--

			if a.frameIndex < startIndex {
				a.playedCount++

				if a.playLoop {
					a.frameIndex = endIndex - 1
				} else {
					a.frameIndex = startIndex
					a.updateLogicHash()
					return nil
				}
			}

		default:
			return errors.New("invalid play mode")
		}
	}

	a.updateLogicHash()

	return nil
}

const (
	full = 1.0
	half = 0.5
	zero = 0.0
)

func (a *Animation) currentFrame() (*animationFrame, error) {
	if err := a.validateFrame(); err != nil {
		return nil, err
	}

	return &a.directions[a.directionIndex].frames[a.frameIndex], nil
}

func (a *Animation) renderShadow(target d2interface.Surface) {
	frame, err := a.currentFrame()
	if err != nil {
		log.Print(err)
		return
	}

	if frame.image == nil {
		return
	}

	target.PushFilter(d2enum.FilterLinear)
	defer target.Pop()

	target.PushTranslation(frame.offsetX, int(float64(frame.offsetY)*half))
	defer target.Pop()

	target.PushScale(full, half)
	defer target.Pop()

	target.PushSkew(half, zero)
	defer target.Pop()

	target.PushEffect(d2enum.DrawEffectPctTransparency25)
	defer target.Pop()

	target.PushBrightness(zero)
	defer target.Pop()

	target.Render(frame.image)
}

// GetCurrentFrameSurface returns the surface for the current frame of the animation.
func (a *Animation) GetCurrentFrameSurface() d2interface.Surface {
	frame, err := a.currentFrame()
	if err != nil {
		log.Print(err)
		return nil
	}

	return frame.image
}

// Render renders the animation to the given surface.
func (a *Animation) Render(target d2interface.Surface) {
	if target == nil {
		return
	}

	if a.renderer == nil {
		a.BindRenderer(target.Renderer())
	}

	frame, err := a.currentFrame()
	if err != nil {
		log.Print(err)
		return
	}

	if frame.image == nil {
		return
	}

	target.PushTranslation(frame.offsetX, frame.offsetY)
	defer target.Pop()

	target.PushEffect(a.effect)
	defer target.Pop()

	target.PushColor(a.colorMod)
	defer target.Pop()

	target.Render(frame.image)
}

// BindRenderer binds the given renderer to the animation so that it can initialize
// the required surfaces.
func (a *Animation) BindRenderer(r d2interface.Renderer) {
	if r == nil {
		return
	}

	if a.onBindRenderer != nil {
		if err := a.onBindRenderer(r); err != nil {
			log.Println(err)
			return
		}
	}

	a.renderer = r
}

// RenderFromOrigin renders the animation from the animation origin.
func (a *Animation) RenderFromOrigin(target d2interface.Surface, shadow bool) {
	if target == nil {
		return
	}

	if a.renderer == nil {
		a.BindRenderer(target.Renderer())
	}

	if err := a.validateFrame(); err != nil {
		log.Print(err)
		return
	}

	if a.originAtBottom {
		frame := a.directions[a.directionIndex].frames[a.frameIndex]

		target.PushTranslation(0, -frame.height)
		defer target.Pop()
	}

	if shadow && !a.effect.Transparent() && a.hasShadow {
		_, height := a.GetFrameBounds()
		height = int(math.Abs(float64(height)))
		halfHeight := height / 2

		target.PushTranslation(-halfHeight, 0)
		a.renderShadow(target)
		target.Pop()
	}

	// Important fix:
	// Original logic returned after drawing the shadow.
	// This renders shadow first, then the actual sprite.
	a.Render(target)
}

// RenderSection renders the section of the animation frame enclosed by bounds.
func (a *Animation) RenderSection(target d2interface.Surface, bound image.Rectangle) {
	if target == nil {
		return
	}

	if a.renderer == nil {
		a.BindRenderer(target.Renderer())
	}

	frame, err := a.currentFrame()
	if err != nil {
		log.Print(err)
		return
	}

	if frame.image == nil {
		return
	}

	target.PushTranslation(frame.offsetX, frame.offsetY)
	defer target.Pop()

	target.PushEffect(a.effect)
	defer target.Pop()

	target.PushColor(a.colorMod)
	defer target.Pop()

	target.RenderSection(frame.image, bound)
}

// GetFrameSize gets the Size(width, height) of an indexed frame.
func (a *Animation) GetFrameSize(frameIndex int) (width, height int, err error) {
	if err := a.validateDirection(); err != nil {
		return 0, 0, err
	}

	direction := a.directions[a.directionIndex]

	if frameIndex < 0 || frameIndex >= len(direction.frames) {
		return 0, 0, errors.New("invalid frame index")
	}

	frame := direction.frames[frameIndex]

	return frame.width, frame.height, nil
}

// GetCurrentFrameSize gets the Size(width, height) of the current frame.
func (a *Animation) GetCurrentFrameSize() (width, height int) {
	width, height, err := a.GetFrameSize(a.frameIndex)
	if err != nil {
		log.Print(err)
	}

	return width, height
}

// GetFrameBounds gets maximum Size(width, height) of all frames.
func (a *Animation) GetFrameBounds() (maxWidth, maxHeight int) {
	if err := a.validateDirection(); err != nil {
		log.Print(err)
		return 0, 0
	}

	direction := a.directions[a.directionIndex]

	for _, frame := range direction.frames {
		maxWidth = d2math.MaxInt(maxWidth, frame.width)
		maxHeight = d2math.MaxInt(maxHeight, frame.height)
	}

	return maxWidth, maxHeight
}

// GetCurrentFrame gets index of current frame in animation.
func (a *Animation) GetCurrentFrame() int {
	return a.frameIndex
}

// GetFrameCount gets number of frames in animation.
func (a *Animation) GetFrameCount() int {
	if len(a.directions) == 0 {
		return 0
	}

	if a.directionIndex < 0 || a.directionIndex >= len(a.directions) {
		return 0
	}

	return len(a.directions[a.directionIndex].frames)
}

// IsOnFirstFrame gets if the animation is on its first frame.
func (a *Animation) IsOnFirstFrame() bool {
	return a.frameIndex == 0
}

// IsOnLastFrame gets if the animation is on its last frame.
func (a *Animation) IsOnLastFrame() bool {
	frameCount := a.GetFrameCount()
	if frameCount <= 0 {
		return false
	}

	return a.frameIndex == frameCount-1
}

// GetDirectionCount gets the number of animation directions.
func (a *Animation) GetDirectionCount() int {
	return len(a.directions)
}

// SetDirection places the animation in the direction of an animation.
func (a *Animation) SetDirection(directionIndex int) error {
	a.ensureFieldDefaults()

	const smallestInvalidDirectionIndex = 64

	if directionIndex < 0 || directionIndex >= smallestInvalidDirectionIndex {
		return errors.New("invalid direction index")
	}

	if len(a.directions) == 0 {
		return errors.New("animation has no directions")
	}

	mappedDirection := d2dcc.Dir64ToDcc(directionIndex, len(a.directions))
	if mappedDirection < 0 || mappedDirection >= len(a.directions) {
		return errors.New("mapped direction index out of range")
	}

	a.directionIndex = mappedDirection
	a.frameIndex = 0
	a.lastFrameTime = 0
	a.updateLogicHash()

	return nil
}

// GetDirection gets the current animation direction.
func (a *Animation) GetDirection() int {
	return a.directionIndex
}

// SetCurrentFrame sets animation at a specific frame.
func (a *Animation) SetCurrentFrame(frameIndex int) error {
	if frameIndex < 0 || frameIndex >= a.GetFrameCount() {
		return errors.New("invalid frame index")
	}

	a.frameIndex = frameIndex
	a.lastFrameTime = 0
	a.updateLogicHash()

	return nil
}

// Rewind animation to beginning.
func (a *Animation) Rewind() {
	if err := a.SetCurrentFrame(0); err != nil {
		log.Print(err)
	}

	a.ResetPlayedCount()
}

// PlayForward plays animation forward.
func (a *Animation) PlayForward() {
	a.playMode = playModeForward
	a.lastFrameTime = 0
	a.updateLogicHash()
}

// PlayBackward plays animation backward.
func (a *Animation) PlayBackward() {
	a.playMode = playModeBackward
	a.lastFrameTime = 0
	a.updateLogicHash()
}

// Pause animation.
func (a *Animation) Pause() {
	a.playMode = playModePause
	a.lastFrameTime = 0
	a.updateLogicHash()
}

// SetPlayLoop sets whether to loop the animation.
func (a *Animation) SetPlayLoop(loop bool) {
	a.playLoop = loop
	a.updateLogicHash()
}

// SetPlaySpeed sets play speed of the animation.
// In this implementation playSpeed means seconds per frame,
// preserving the original behavior.
func (a *Animation) SetPlaySpeed(playSpeed float64) {
	if playSpeed <= 0 || math.IsNaN(playSpeed) || math.IsInf(playSpeed, 0) {
		return
	}

	frameCount := a.GetFrameCount()
	if frameCount <= 0 {
		return
	}

	a.SetPlayLength(playSpeed * float64(frameCount))
}

// SetPlayLength sets the Animation's play length in seconds.
func (a *Animation) SetPlayLength(playLength float64) {
	if playLength <= 0 || math.IsNaN(playLength) || math.IsInf(playLength, 0) {
		playLength = defaultPlayLength
	}

	a.playLength = playLength
	a.lastFrameTime = 0
	a.updateLogicHash()
}

// SetColorMod sets the Animation's color mod.
func (a *Animation) SetColorMod(colorMod color.Color) {
	a.colorMod = colorMod
}

// GetPlayedCount gets the number of times the animation played.
func (a *Animation) GetPlayedCount() int {
	return a.playedCount
}

// ResetPlayedCount resets the play count.
func (a *Animation) ResetPlayedCount() {
	a.playedCount = 0
	a.updateLogicHash()
}

// SetEffect sets the draw effect for the animation.
func (a *Animation) SetEffect(e d2enum.DrawEffect) {
	a.effect = e
}

// SetShadow sets bool for whether or not to draw a shadow.
func (a *Animation) SetShadow(shadow bool) {
	a.hasShadow = shadow
}

// ForceKappaInvariant hard-resets the animation into ARE-compliant state.
func (a *Animation) ForceKappaInvariant() {
	a.kappaPos = kappaInvariant
	a.fieldLogicName = fieldLogicVersion
	a.collectiveStamp = collectiveMarkgraf
	a.updateLogicHash()
}

// AutoGenerateFieldLogic creates a deterministic state signature.
// This is safe for replay logs, debugging, networking checks, and AI audit trails.
func (a *Animation) AutoGenerateFieldLogic() FieldLogicSignature {
	a.ensureFieldDefaults()

	frameCount := a.GetFrameCount()

	logic := strings.Builder{}
	logic.WriteString(a.fieldLogicName)
	logic.WriteString("|")
	logic.WriteString(a.collectiveStamp)
	logic.WriteString("|")
	logic.WriteString(ouroborosSeal)
	logic.WriteString("|KAPPA=")
	logic.WriteString(fmt.Sprintf("%d", a.kappaPos))
	logic.WriteString("|DIR=")
	logic.WriteString(fmt.Sprintf("%d", a.directionIndex))
	logic.WriteString("|FRAME=")
	logic.WriteString(fmt.Sprintf("%d", a.frameIndex))
	logic.WriteString("|COUNT=")
	logic.WriteString(fmt.Sprintf("%d", frameCount))
	logic.WriteString("|PLAYED=")
	logic.WriteString(fmt.Sprintf("%d", a.playedCount))
	logic.WriteString("|MODE=")
	logic.WriteString(fmt.Sprintf("%d", a.playMode))
	logic.WriteString("|LOOP=")
	logic.WriteString(fmt.Sprintf("%t", a.playLoop))
	logic.WriteString("|SUB=")
	logic.WriteString(fmt.Sprintf("%t:%d:%d", a.hasSubLoop, a.subStartingFrame, a.subEndingFrame))

	logicString := logic.String()
	hash := deterministicHash64(logicString)

	return FieldLogicSignature{
		Version:        a.fieldLogicName,
		Collective:     a.collectiveStamp,
		Ouroboros:      ouroborosSeal,
		KappaPos:       a.kappaPos,
		DirectionIndex: a.directionIndex,
		FrameIndex:     a.frameIndex,
		FrameCount:     frameCount,
		PlayedCount:    a.playedCount,
		PlayMode:       a.playMode,
		PlayLoop:       a.playLoop,
		HasSubLoop:     a.hasSubLoop,
		SubStart:       a.subStartingFrame,
		SubEnd:         a.subEndingFrame,
		LogicString:    logicString,
		Hash64:         hash,
		Axioms: []AnimationAxiom{
			AxiomKappaStable,
			AxiomFrameBounded,
			AxiomDirectionBounded,
			AxiomNoTimeMutation,
			AxiomReplayable,
		},
	}
}

// GetLogicString returns the current deterministic ARE/Ouroboros logic string.
func (a *Animation) GetLogicString() string {
	return a.AutoGenerateFieldLogic().LogicString
}

// GetLogicHash returns the current deterministic hash.
func (a *Animation) GetLogicHash() uint64 {
	sig := a.AutoGenerateFieldLogic()
	return sig.Hash64
}

// GetLastLogicHash returns the last cached hash updated by state-changing calls.
func (a *Animation) GetLastLogicHash() uint64 {
	return a.lastLogicHash
}

func (a *Animation) updateLogicHash() {
	a.lastLogicHash = a.GetLogicHash()
}

func deterministicHash64(input string) uint64 {
	h := fnv.New64a()

	_, _ = h.Write([]byte(input))

	return h.Sum64()
}

// ValidateAxioms validates the animation against deterministic safety axioms.
func (a *Animation) ValidateAxioms() error {
	if err := a.validateKappa(); err != nil {
		return err
	}

	if err := a.validateDirection(); err != nil {
		return err
	}

	if err := a.validateFrame(); err != nil {
		return err
	}

	if err := a.validateSubLoop(); err != nil {
		return err
	}

	return nil
}

// Clone creates a copy of the Animation.
// Surfaces are intentionally shared references because duplicating GPU surfaces
// here would be renderer-specific and unsafe.
func (a *Animation) Clone() d2interface.Animation {
	clone := *a

	if a.directions != nil {
		clone.directions = make([]animationDirection, len(a.directions))

		for directionIndex := range a.directions {
			srcDirection := a.directions[directionIndex]

			clone.directions[directionIndex] = animationDirection{
				decoded: srcDirection.decoded,
				frames:  nil,
			}

			if srcDirection.frames != nil {
				clone.directions[directionIndex].frames = make([]animationFrame, len(srcDirection.frames))
				copy(clone.directions[directionIndex].frames, srcDirection.frames)
			}
		}
	}

	clone.updateLogicHash()

	return &clone
}
