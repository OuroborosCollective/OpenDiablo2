package d2emergent

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// Ouroboros integration constants (KAPPA and timing only - Markgraf constants in emergent_logic.go)
const (
	// KAPPA chunk generation constants
	KappaScale       = 1000 // Fixed-point spatial positioning multiplier
	KappaMoveEpsilon = 1.0  // Minimum movement threshold in KAPPA units

	// Emergent behavior update interval
	EmergenceTickMs = 100
)

// OuroborosLogikSystem binds ARE-Logik to the game system via BaalAal event bus
type OuroborosLogikSystem struct {
	mu              sync.RWMutex
	logger          *d2util.Logger
	areLogik        *ARELogikEngine
	baalAal         *d2script.BaalAalEngine
	lastEmergence   int64
	resonanceOffset float64
	chunkRegistry   map[string]*KappaChunk
	entityStates    map[string]*EntityState
	rules           []DeterministicRule
}

// KappaChunk represents a spatial region in the game world
type KappaChunk struct {
	mu         sync.RWMutex
	ID         string
	X, Y       int32 // KAPPA-space coordinates
	Resonance  float32
	Occupants  []string
	Generation uint64
}

// EntityState tracks an entity's deterministic state
type EntityState struct {
	mu              sync.RWMutex
	EntityID        string
	KappaPos        []int32 // [x, y] in KAPPA units
	LastResonance   float32
	BehaviorProfile int // 0-255 deterministic behavior selector
	Liveliness      float32
	LastUpdateTick  uint64
}

// DeterministicRule defines an axiomatic rule for game behavior
type DeterministicRule struct {
	Name        string
	Priority    int
	Condition   func(*d2script.IAxiomaticEvent, *OuroborosLogikSystem) bool
	Action      func(*d2script.IAxiomaticEvent, *OuroborosLogikSystem)
	Description string
}

// NewOuroborosLogikSystem creates a new Ouroboros-ARE-Logik system bound to BaalAal
func NewOuroborosLogikSystem(baalAal *d2script.BaalAalEngine, logLevel d2util.LogLevel) *OuroborosLogikSystem {
	system := &OuroborosLogikSystem{
		logger:        d2util.NewLogger(),
		areLogik:      CreateARELogikEngine(logLevel),
		baalAal:       baalAal,
		chunkRegistry: make(map[string]*KappaChunk),
		entityStates:  make(map[string]*EntityState),
		rules:         []DeterministicRule{},
	}

	system.logger.SetPrefix("Ouroboros-ARE")
	system.logger.SetLevel(logLevel)

	// Register Ouroboros event handlers with BaalAal
	system.registerEventHandlers()

	// Register deterministic rules
	system.registerDeterministicRules()

	system.logger.Info("Ouroboros-ARE-Logik System initialized")

	return system
}

// registerEventHandlers binds Ouroboros event handlers to the BaalAal event bus
func (s *OuroborosLogikSystem) registerEventHandlers() {
	// Subscribe to world emergence events
	s.baalAal.EventBus.Subscribe("Ouroboros-WorldEmergence", func(event *d2script.IAxiomaticEvent) {
		s.handleWorldEmergence(event)
	})

	// Subscribe to player movement
	s.baalAal.EventBus.Subscribe("Ouroboros-PlayerMove", func(event *d2script.IAxiomaticEvent) {
		s.handlePlayerMove(event)
	})

	// Subscribe to combat events
	s.baalAal.EventBus.Subscribe("Ouroboros-Combat", func(event *d2script.IAxiomaticEvent) {
		s.handleCombatEvent(event)
	})

	// Subscribe to item spawns
	s.baalAal.EventBus.Subscribe("Ouroboros-ItemSpawn", func(event *d2script.IAxiomaticEvent) {
		s.handleItemSpawn(event)
	})

	// Register rules for event processing
	s.baalAal.RegisterRule("WorldEmergence", s.handleWorldEmergence)
	s.baalAal.RegisterRule("PlayerMove", s.handlePlayerMove)
	s.baalAal.RegisterRule("PLAYER_MOVE", s.handlePlayerMove)
}

// registerDeterministicRules defines the axiomatic rules for emergent behavior
func (s *OuroborosLogikSystem) registerDeterministicRules() {
	// Rule: KAPPA chunk generation for player movement
	s.rules = append(s.rules, DeterministicRule{
		Name:        "KappaChunkGen",
		Priority:    100,
		Condition:   s.condKappaChunkNeeded,
		Action:      s.actionGenerateKappaChunk,
		Description: "Generate KAPPA spatial chunks when entities enter new regions",
	})

	// Rule: NPC emergent behavior based on resonance
	s.rules = append(s.rules, DeterministicRule{
		Name:        "NPCEmergentBehavior",
		Priority:    50,
		Condition:   s.condNPCShouldAct,
		Action:      s.actionNPCEmergentBehavior,
		Description: "Trigger NPC emergent behavior based on global resonance",
	})

	// Rule: Combat resonance feedback
	s.rules = append(s.rules, DeterministicRule{
		Name:        "CombatResonanceFeedback",
		Priority:    75,
		Condition:   s.condCombatEvent,
		Action:      s.actionCombatResonanceFeedback,
		Description: "Apply resonance feedback to combat outcomes",
	})
}

// ToKappa converts float64 world coordinates to KAPPA int32
func (s *OuroborosLogikSystem) ToKappa(val float64) int32 {
	return int32(math.Floor(val * KappaScale))
}

// ToWorld converts KAPPA int32 coordinates back to world float64
func (s *OuroborosLogikSystem) ToWorld(val int32) float64 {
	return float64(val) / KappaScale
}

// handleWorldEmergence processes ARE-Logik emergence events
func (s *OuroborosLogikSystem) handleWorldEmergence(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	// Update ARE-Logik engine state
	s.areLogik.ProcessEmergence()

	// Calculate new resonance with Ouroboros feedback
	s.mu.Lock()
	s.resonanceOffset = s.areLogik.GlobalResonance
	s.mu.Unlock()

	// Process deterministic rules
	s.processRules(event)
}

// handlePlayerMove processes player movement with KAPPA chunk generation
func (s *OuroborosLogikSystem) handlePlayerMove(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Extract or compute KAPPA coordinates
	clientID, _ := event.Metadata["client_id"].(string)
	var kappaX, kappaY int32

	if kx, ok := event.Metadata["kappa_x"].(int32); ok {
		kappaX = kx
	} else if x, ok := event.Metadata["x"].(float64); ok {
		kappaX = s.ToKappa(x)
	}

	if ky, ok := event.Metadata["kappa_y"].(int32); ok {
		kappaY = ky
	} else if y, ok := event.Metadata["y"].(float64); ok {
		kappaY = s.ToKappa(y)
	}

	// Update or create entity state
	if clientID != "" {
		state, exists := s.entityStates[clientID]
		if !exists {
			state = &EntityState{
				EntityID:        clientID,
				KappaPos:        []int32{kappaX, kappaY},
				BehaviorProfile: s.computeBehaviorProfile(event),
				Liveliness:      0.5,
			}
			s.entityStates[clientID] = state
		} else {
			state.mu.Lock()
			state.KappaPos = []int32{kappaX, kappaY}
			state.LastUpdateTick = s.areLogik.Tick
			state.mu.Unlock()
		}
	}

	// Generate KAPPA chunk if needed
	chunkID := s.getChunkID(kappaX, kappaY)
	if _, exists := s.chunkRegistry[chunkID]; !exists {
		s.createKappaChunk(chunkID, kappaX, kappaY)
	}
}

// handleCombatEvent processes combat with resonance feedback
func (s *OuroborosLogikSystem) handleCombatEvent(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Apply Ouroboros feedback to combat resonance
	resonance := s.areLogik.GlobalResonance
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	event.Metadata["ouroboros_resonance"] = resonance
	event.Metadata["baal_resonance"] = s.baalAal.EventBus.ResonanceState

	// Compute deterministic combat modifier
	combatMod := s.computeDeterministicCombat(event)
	event.Metadata["deterministic_combat_mod"] = combatMod
}

// handleItemSpawn processes item spawning with KAPPA chunk association
func (s *OuroborosLogikSystem) handleItemSpawn(event *d2script.IAxiomaticEvent) {
	if event == nil || event.Payload == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Extract position data
	var kappaX, kappaY int32

	if payload, ok := event.Payload.(map[string]interface{}); ok {
		if x, ok := payload["x"].(float64); ok {
			kappaX = s.ToKappa(x)
		}
		if y, ok := payload["y"].(float64); ok {
			kappaY = s.ToKappa(y)
		}
	}

	// Associate with KAPPA chunk
	chunkID := s.getChunkID(kappaX, kappaY)
	chunk, exists := s.chunkRegistry[chunkID]
	if !exists {
		chunk = s.createKappaChunk(chunkID, kappaX, kappaY)
	}

	// Add emergent generation seed
	genSeed := s.computeGenerationSeed(event)
	event.Metadata["kappa_generation"] = genSeed
	event.Metadata["chunk_id"] = chunkID
}

// processRules evaluates and executes deterministic rules
func (s *OuroborosLogikSystem) processRules(event *d2script.IAxiomaticEvent) {
	for _, rule := range s.rules {
		if rule.Condition(event, s) {
			rule.Action(event, s)
		}
	}
}

// Deterministic rule conditions and actions

func (s *OuroborosLogikSystem) condKappaChunkNeeded(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) bool {
	if event == nil {
		return false
	}

	eventType := event.Type
	return eventType == "PlayerMove" || eventType == "PLAYER_MOVE"
}

func (s *OuroborosLogikSystem) actionGenerateKappaChunk(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) {
	// KAPPA chunk generation is handled in handlePlayerMove
}

func (s *OuroborosLogikSystem) condNPCShouldAct(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) bool {
	if event == nil {
		return false
	}

	// NPCs act based on resonance threshold
	resonance := sys.areLogik.GlobalResonance
	return resonance > 0.8 || resonance < 0.2
}

func (s *OuroborosLogikSystem) actionNPCEmergentBehavior(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) {
	// Generate emergent NPC behavior based on deterministic rules
	resonance := sys.areLogik.GlobalResonance

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	// Use Markgraf harmonic for deterministic behavior selection
	tick := sys.areLogik.Tick
	behaviorSeed := math.Sin(float64(tick)*MarkgrafHarmonicFrequency) * 127.5 // -127 to +127
	event.Metadata["npc_emergent_behavior"] = behaviorSeed
	event.Metadata["npc_resonance_influence"] = resonance
}

func (s *OuroborosLogikSystem) condCombatEvent(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) bool {
	if event == nil {
		return false
	}

	return event.Type == "9" || event.Type == "CombatCast"
}

func (s *OuroborosLogikSystem) actionCombatResonanceFeedback(event *d2script.IAxiomaticEvent, sys *OuroborosLogikSystem) {
	// Combat resonance is handled in handleCombatEvent
}

// Helper methods

func (s *OuroborosLogikSystem) getChunkID(x, y int32) string {
	// Use chunk size of 100 KAPPA units
	chunkX := x / 100
	chunkY := y / 100
	return fmt.Sprintf("chunk_%d_%d", chunkX, chunkY)
}

func (s *OuroborosLogikSystem) createKappaChunk(id string, x, y int32) *KappaChunk {
	chunk := &KappaChunk{
		ID:         id,
		X:          x,
		Y:          y,
		Resonance:  float32(s.areLogik.GlobalResonance),
		Occupants:  []string{},
		Generation: s.areLogik.Tick,
	}
	s.chunkRegistry[id] = chunk
	return chunk
}

func (s *OuroborosLogikSystem) computeBehaviorProfile(event *d2script.IAxiomaticEvent) int {
	if event == nil {
		return 128
	}

	// Deterministic behavior profile based on event hash
	seed := uint64(len(event.ID))
	if event.SequenceID > 0 {
		seed ^= event.SequenceID * 15485863
	}
	return int(seed % 256)
}

func (s *OuroborosLogikSystem) computeDeterministicCombat(event *d2script.IAxiomaticEvent) float64 {
	if event == nil {
		return 1.0
	}

	// Deterministic combat modifier using Markgraf harmonic
	tick := s.areLogik.Tick
	harmonic := math.Sin(float64(tick)*MarkgrafHarmonicFrequency) +
		0.5*math.Cos(float64(tick)*MarkgrafHarmonicFrequency*0.618)

	return 1.0 + (harmonic * MarkgrafFeedbackWeight)
}

func (s *OuroborosLogikSystem) computeGenerationSeed(event *d2script.IAxiomaticEvent) uint64 {
	if event == nil {
		return 0
	}

	// Generate deterministic seed from event data
	seed := (s.areLogik.Tick * 15485863) ^ (event.SequenceID * 2038074743)
	return seed
}

// GetStatus returns current Ouroboros-ARE system status
func (s *OuroborosLogikSystem) GetStatus() (resonance, expansion, entropy float64, tick uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.areLogik.GlobalResonance,
		s.areLogik.Expansion,
		s.areLogik.Entropy,
		s.areLogik.Tick
}

// GetChunkInfo returns information about a KAPPA chunk
func (s *OuroborosLogikSystem) GetChunkInfo(id string) *KappaChunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.chunkRegistry[id]
}

// GetEntityState returns the deterministic state of an entity
func (s *OuroborosLogikSystem) GetEntityState(entityID string) *EntityState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.entityStates[entityID]
}

// Advance processes a single emergence tick
func (s *OuroborosLogikSystem) Advance() {
	now := time.Now().UnixMilli()

	s.mu.RLock()
	needsUpdate := (now - s.lastEmergence) > EmergenceTickMs
	s.mu.RUnlock()

	if needsUpdate {
		s.mu.Lock()
		s.lastEmergence = now
		s.mu.Unlock()

		// Process emergence cycle
		event := s.areLogik.ProcessEmergence()
		s.baalAal.EventBus.Publish(event)
	}
}

// ProcessEvent processes an event through the Ouroboros system
func (s *OuroborosLogikSystem) ProcessEvent(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	// Publish to BaalAal for standard processing
	s.baalAal.EventBus.Publish(event)

	// Process through Ouroboros rules
	s.processRules(event)
}