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
	ResonanceState   float64
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

	event.SequenceID = b.globalSequenceID
	b.globalSequenceID++

	// Recursive BaalAal logic: The state "eats" its own tail (resonance feedback)
	// We'll simulate this by incorporating the previous resonance state into the new event
	resonance := b.calculateResonance(event)
	b.ResonanceState = math.Mod(b.ResonanceState+resonance, 2147483647)

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}
	event.Metadata["resonance"] = resonance
	event.Metadata["baalaal_cycle"] = b.ResonanceState // The self-eating cycle

	b.ledger[b.writePointer] = event
	b.writePointer = (b.writePointer + 1) % b.maxLedgerSize
	if b.writePointer == 0 {
		b.isFull = true
	}

	// Snapshot subscribers to notify outside of main lock to prevent deadlocks
	subs := make([]func(*IAxiomaticEvent), 0, len(b.subscribers))
	for _, fn := range b.subscribers {
		subs = append(subs, fn)
	}
	b.Unlock()

	for _, fn := range subs {
		fn(event)
	}
}

func (b *AxiomaticEventBus) calculateResonance(event *IAxiomaticEvent) float64 {
	// Simplified hash for resonance calculation
	return float64(event.SequenceID%10000) / 10000.0
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
	sync.RWMutex
	engine    *BaalAalEngine
	Positions map[string][]int32
	Compiler  *AREStateCompiler
}

func NewKappaSystem(engine *BaalAalEngine) *KappaSystem {
	return &KappaSystem{
		Positions: make(map[string][]int32),
		Compiler:  &AREStateCompiler{},
		engine:    engine,
	}
}

func (k *KappaSystem) HandleMove(event *IAxiomaticEvent) {
	if event.Payload == nil && event.Metadata == nil {
		return
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	// Try to get data from Metadata (used in extra_test)
	clientID, ok := event.Metadata["client_id"].(string)
	if ok {
		x, xOk := event.Metadata["x"].(float64)
		y, yOk := event.Metadata["y"].(float64)
		if xOk && yOk {
			kx := k.Compiler.ToKappa(x)
			ky := k.Compiler.ToKappa(y)
			event.Metadata["kappa_x"] = kx
			event.Metadata["kappa_y"] = ky
			k.Lock()
			k.Positions[clientID] = []int32{kx, ky}
			k.Unlock()
			return
		}
	}

	// Try to get data from Payload (used in main test)
	if moveData, ok := event.Payload.(map[string]interface{}); ok {
		x, xOk := moveData["x"].(float64)
		y, yOk := moveData["y"].(float64)

		if xOk && yOk {
			kx := k.Compiler.ToKappa(x)
			ky := k.Compiler.ToKappa(y)

			if event.Metadata == nil {
				event.Metadata = make(map[string]interface{})
			}
			event.Metadata["kappa_x"] = kx
			event.Metadata["kappa_y"] = ky

			if clientID, ok := event.Metadata["client_id"].(string); ok {
				k.Lock()
				k.Positions[clientID] = []int32{kx, ky}
				k.Unlock()
			}
		}
	}
}

// BaalAalEngine wraps the Axiomatic components into a cohesive engine.
type BaalAalEngine struct {
	Compiler           *AREStateCompiler
	EventBus           *AxiomaticEventBus
	KappaSystem        *KappaSystem
	CombatSystem       *CombatSystem
	ItemSystem         *ItemSystem
	WorldSystem        *WorldSystem
	rules              map[string][]func(*IAxiomaticEvent)
	lastProcessedIndex int
}

func NewBaalAalEngine() *BaalAalEngine {
	e := &BaalAalEngine{
		Compiler: &AREStateCompiler{},
		EventBus: NewAxiomaticEventBus(50000), // Matching Wasd repo size
		rules:    make(map[string][]func(*IAxiomaticEvent)),
	}

	e.KappaSystem = NewKappaSystem(e)
	e.CombatSystem = &CombatSystem{Compiler: e.Compiler}
	e.ItemSystem = &ItemSystem{Compiler: e.Compiler}
	e.WorldSystem = NewWorldSystem()

	e.EventBus.Subscribe("KappaSystem", func(event *IAxiomaticEvent) {
		if event.Type == "PLAYER_MOVE" || event.Type == "PlayerMove" {
			e.KappaSystem.HandleMove(event)
		}
	})

	e.RegisterRule("PlayerMove", e.KappaSystem.HandleMove)
	e.RegisterRule("PLAYER_MOVE", e.KappaSystem.HandleMove)
	e.RegisterRule("WorldEmergence", e.WorldSystem.HandleEmergence)

	// Network packet types as strings
	e.RegisterRule("9", e.CombatSystem.HandleCast)
	e.RegisterRule("10", e.ItemSystem.HandleSpawn)

	return e
}

func (e *BaalAalEngine) RegisterRule(eventType string, handler func(*IAxiomaticEvent)) {
	e.rules[eventType] = append(e.rules[eventType], handler)
}

// ProcessCycle represents a single recursive BaalAal cycle.
func (e *BaalAalEngine) ProcessCycle(tick uint64) {
	e.EventBus.Lock()

	// The boss Baal recursive snake self eating recursive cycle system
	// Incorporate the resonance state back into itself
	if e.EventBus.ResonanceState == 0 {
		e.EventBus.ResonanceState = 1.0 // Seed if zero
	}

	// Complex recursive logic: The resonance state grows and wraps,
	// but is also influenced by the tick and a harmonic sine wave.
	harmonic := math.Sin(float64(tick)*0.01) * 0.05
	e.EventBus.ResonanceState = math.Mod((e.EventBus.ResonanceState+harmonic)*1.0001, 2147483647)

	if e.EventBus.ResonanceState < 0 {
		e.EventBus.ResonanceState = 1.0
	}

	// Collect new events from ledger under lock
	var newEvents []*IAxiomaticEvent
	for e.lastProcessedIndex != e.EventBus.writePointer {
		event := e.EventBus.ledger[e.lastProcessedIndex]
		if event != nil {
			newEvents = append(newEvents, event)
		}
		e.lastProcessedIndex = (e.lastProcessedIndex + 1) % e.EventBus.maxLedgerSize
	}
	e.EventBus.Unlock()

	// Process new events outside of main lock to prevent deadlocks
	for _, event := range newEvents {
		if handlers, ok := e.rules[event.Type]; ok {
			for _, handler := range handlers {
				handler(event)
			}
		}
	}
}

// CombatSystem handles axiomatic combat events.
type CombatSystem struct {
	Compiler *AREStateCompiler
}

func (c *CombatSystem) HandleCast(event *IAxiomaticEvent) {
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	// Use SequenceID as pseudo-tick for deterministic resonance
	var entityID uint32
	if clientID, ok := event.Metadata["client_id"].(string); ok {
		for _, char := range clientID {
			entityID += uint32(char)
		}
	}

	resonance := c.Compiler.GetDeterministicResonance(event.SequenceID, entityID)
	event.Metadata["deterministic_resonance"] = resonance
}

// ItemSystem handles axiomatic item events.
type ItemSystem struct {
	Compiler *AREStateCompiler
}

func (i *ItemSystem) HandleSpawn(event *IAxiomaticEvent) {
	if event.Payload == nil {
		return
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	// Try to extract coordinates from payload and convert to Kappa
	if data, ok := event.Payload.(map[string]interface{}); ok {
		if x, ok := data["x"].(float64); ok {
			event.Metadata["kappa_x"] = i.Compiler.ToKappa(x)
		}
		if y, ok := data["y"].(float64); ok {
			event.Metadata["kappa_y"] = i.Compiler.ToKappa(y)
		}
	}
}

// WorldSystem tracks global deterministic world state.
type WorldSystem struct {
	sync.RWMutex
	GlobalResonance float64
	Expansion       float64
	Entropy         float64
}

func NewWorldSystem() *WorldSystem {
	return &WorldSystem{
		GlobalResonance: 1.0,
		Expansion:       1.0,
		Entropy:         0.5,
	}
}

func (w *WorldSystem) HandleEmergence(event *IAxiomaticEvent) {
	w.Lock()
	defer w.Unlock()

	if resonance, ok := event.Payload.(float64); ok {
		w.GlobalResonance = resonance
	}

	if expansion, ok := event.Metadata["expansion"].(float64); ok {
		w.Expansion = expansion
	}

	if entropy, ok := event.Metadata["entropy"].(float64); ok {
		w.Entropy = entropy
	}
}

// GetStatus returns the current resonance and cycle.
func (e *BaalAalEngine) GetStatus() (float64, float64) {
	e.EventBus.RLock()
	defer e.EventBus.RUnlock()

	resonance := math.Mod(e.EventBus.ResonanceState, 1.0)
	return resonance, e.EventBus.ResonanceState
}
