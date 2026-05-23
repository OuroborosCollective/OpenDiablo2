package d2script

import (
	"encoding/binary"
	"math"
	"sync"
)

// KAPPA_SCALE is the fixed-point spatial positioning multiplier (val * 1000)
const KAPPA_SCALE = 1000

// IAxiomaticEvent represents a system-wide axiomatic event.
type IAxiomaticEvent struct {
	ID         string
	SequenceID uint64
	Type       string
	Timestamp  int64
	Payload    interface{}
	Metadata   map[string]interface{}
}

// AREStateData represents the fundamental atomic state of an entity.
type AREStateData struct {
	L uint32    // Logic Index
	K []int32   // Kappa-space coordinates
	R []float32 // Resonance values
}

// AREStateCompiler handles serialization and deterministic logic.
type AREStateCompiler struct{}

func (c *AREStateCompiler) ToKappa(val float64) int32 {
	return int32(math.Floor(val * KAPPA_SCALE))
}

func (c *AREStateCompiler) GetDeterministicResonance(tick uint64, entityID uint32) float32 {
	seed := (tick * 15485863) ^ (uint64(entityID) * 2038074743)
	return float32(seed%10000) / 10000.0
}

// AxiomaticEventBus is the central event hub.
type AxiomaticEventBus struct {
	sync.RWMutex
	ledger           []*IAxiomaticEvent
	maxLedgerSize    int
	writePointer     int
	isFull           bool
	globalSequenceID uint64
	resonanceState   float64
	subscribers      map[string]func(*IAxiomaticEvent)
}

func NewAxiomaticEventBus(size int) *AxiomaticEventBus {
	return &AxiomaticEventBus{
		ledger:        make([]*IAxiomaticEvent, size),
		maxLedgerSize: size,
		subscribers:   make(map[string]func(*IAxiomaticEvent)),
	}
}

// Subscribe registers a new subscriber for events.
func (b *AxiomaticEventBus) Subscribe(id string, fn func(*IAxiomaticEvent)) {
	b.Lock()
	defer b.Unlock()
	b.subscribers[id] = fn
}

// Unsubscribe removes a subscriber.
func (b *AxiomaticEventBus) Unsubscribe(id string) {
	b.Lock()
	defer b.Unlock()
	delete(b.subscribers, id)
}

func (b *AxiomaticEventBus) Publish(event *IAxiomaticEvent) {
	b.Lock()
	defer b.Unlock()

	event.SequenceID = b.globalSequenceID
	b.globalSequenceID++

	// Recursive BaalAal logic: The state "eats" its own tail (resonance feedback)
	// We'll simulate this by incorporating the previous resonance state into the new event
	resonance := b.calculateResonance(event)
	b.resonanceState = math.Mod(b.resonanceState+resonance, 2147483647)

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}
	event.Metadata["resonance"] = resonance
	event.Metadata["baalaal_cycle"] = b.resonanceState // The self-eating cycle

	b.ledger[b.writePointer] = event
	b.writePointer = (b.writePointer + 1) % b.maxLedgerSize
	if b.writePointer == 0 {
		b.isFull = true
	}

	// Notify subscribers
	for _, fn := range b.subscribers {
		fn(event)
	}
}

func (b *AxiomaticEventBus) calculateResonance(event *IAxiomaticEvent) float64 {
	// Simplified hash for resonance calculation
	return float64(event.SequenceID % 10000) / 10000.0
}

func (b *AxiomaticEventBus) GetHistory() []*IAxiomaticEvent {
	b.RLock()
	defer b.RUnlock()

	history := make([]*IAxiomaticEvent, 0, b.maxLedgerSize)
	if !b.isFull {
		history = append(history, b.ledger[:b.writePointer]...)
	} else {
		history = append(history, b.ledger[b.writePointer:]...)
		history = append(history, b.ledger[:b.writePointer]...)
	}
	return history
}

// Compile binary payload
func (c *AREStateCompiler) Compile(states []AREStateData) []byte {
	// Implementation of binary compilation
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(states)))

	for _, state := range states {
		lBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(lBuf, state.L)
		buf = append(buf, lBuf...)

		kLenBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(kLenBuf, uint32(len(state.K)))
		buf = append(buf, kLenBuf...)

		rLenBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(rLenBuf, uint32(len(state.R)))
		buf = append(buf, rLenBuf...)

		for _, k := range state.K {
			kBuf := make([]byte, 4)
			binary.LittleEndian.PutUint32(kBuf, uint32(k))
			buf = append(buf, kBuf...)
		}

		for _, r := range state.R {
			rBuf := make([]byte, 4)
			binary.LittleEndian.PutUint32(rBuf, math.Float32bits(r))
			buf = append(buf, rBuf...)
		}
	}

	return buf
}

// KappaSystem implements deterministic coordinate tracking.
type KappaSystem struct {
	engine *BaalAalEngine
}

func NewKappaSystem(engine *BaalAalEngine) *KappaSystem {
	k := &KappaSystem{
		engine: engine,
	}
	engine.EventBus.Subscribe("KappaSystem", k.onEvent)
	return k
}

func (k *KappaSystem) onEvent(event *IAxiomaticEvent) {
	if event.Type == "PLAYER_MOVE" || event.Type == "PlayerMove" {
		k.processMove(event)
	}
}

func (k *KappaSystem) processMove(event *IAxiomaticEvent) {
	// Implement deterministic coordinate tracking using KAPPA_SCALE
	// This would typically update the entity state in the BaalAal engine
	if event.Payload == nil {
		return
	}

	// Simplified: just log the move in metadata to prove it processed
	if moveData, ok := event.Payload.(map[string]interface{}); ok {
		if x, ok := moveData["x"].(float64); ok {
			event.Metadata["kappa_x"] = k.engine.Compiler.ToKappa(x)
		}
		if y, ok := moveData["y"].(float64); ok {
			event.Metadata["kappa_y"] = k.engine.Compiler.ToKappa(y)
		}
	}
}

// BaalAalEngine wraps the Axiomatic components into a cohesive engine.
type BaalAalEngine struct {
	Compiler    *AREStateCompiler
	EventBus    *AxiomaticEventBus
	KappaSystem *KappaSystem
}

func NewBaalAalEngine() *BaalAalEngine {
	e := &BaalAalEngine{
		Compiler: &AREStateCompiler{},
		EventBus: NewAxiomaticEventBus(50000), // Matching Wasd repo size
	}
	e.KappaSystem = NewKappaSystem(e)
	return e
}

// ProcessCycle represents a single recursive BaalAal cycle.
func (e *BaalAalEngine) ProcessCycle(tick uint64) {
	e.EventBus.Lock()
	defer e.EventBus.Unlock()

	// The boss Baal recursive snake self eating recursive cycle system
	// Incorporate the resonance state back into itself
	if e.EventBus.resonanceState == 0 {
		e.EventBus.resonanceState = 1.0 // Seed if zero
	}

	// Complex recursive logic: The resonance state grows and wraps,
	// but is also influenced by the tick and a harmonic sine wave.
	harmonic := math.Sin(float64(tick)*0.01) * 0.05
	e.EventBus.resonanceState = math.Mod((e.EventBus.resonanceState+harmonic)*1.0001, 2147483647)

	if e.EventBus.resonanceState < 0 {
		e.EventBus.resonanceState = 1.0
	}
}

// GetStatus returns the current resonance and cycle.
func (e *BaalAalEngine) GetStatus() (float64, float64) {
	e.EventBus.RLock()
	defer e.EventBus.RUnlock()

	resonance := math.Mod(e.EventBus.resonanceState, 1.0)
	return resonance, e.EventBus.resonanceState
}
