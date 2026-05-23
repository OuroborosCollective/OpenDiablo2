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
}

func NewAxiomaticEventBus(size int) *AxiomaticEventBus {
	return &AxiomaticEventBus{
		ledger:        make([]*IAxiomaticEvent, size),
		maxLedgerSize: size,
	}
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

func (b *AxiomaticEventBus) GetResonance() float64 {
	b.RLock()
	defer b.RUnlock()
	return b.resonanceState
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

// BaalAalEngine wraps the Axiomatic components into a cohesive engine.
type BaalAalEngine struct {
	sync.Mutex
	Compiler *AREStateCompiler
	EventBus *AxiomaticEventBus
	rules    map[string]func(event *IAxiomaticEvent)
}

func NewBaalAalEngine() *BaalAalEngine {
	return &BaalAalEngine{
		Compiler: &AREStateCompiler{},
		EventBus: NewAxiomaticEventBus(50000), // Matching Wasd repo size
		rules:    make(map[string]func(event *IAxiomaticEvent)),
	}
}

// RegisterRule binds a deterministic logic handler to an event type.
func (e *BaalAalEngine) RegisterRule(eventType string, handler func(event *IAxiomaticEvent)) {
	e.Lock()
	defer e.Unlock()
	e.rules[eventType] = handler
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
	e.EventBus.resonanceState = math.Mod(e.EventBus.resonanceState*1.0001, 2147483647)

	// Process pending events through rules
	history := e.EventBus.GetHistory()
	if len(history) > 0 {
		lastEvent := history[len(history)-1]
		if handler, ok := e.rules[lastEvent.Type]; ok {
			handler(lastEvent)
		}
	}
}
