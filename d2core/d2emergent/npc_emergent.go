package d2emergent

import (
	"math"
	"sync"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// NPCBehaviorType represents deterministic NPC behavior patterns
type NPCBehaviorType int

const (
	NPCBehaviorIdle       NPCBehaviorType = 0
	NPCBehaviorPatrol     NPCBehaviorType = 1
	NPCBehaviorInteract   NPCBehaviorType = 2
	NPCBehaviorFlee       NPCBehaviorType = 3
	NPCBehaviorAggro     NPCBehaviorType = 4
	NPCBehaviorEmergent1 NPCBehaviorType = 5
	NPCBehaviorEmergent2 NPCBehaviorType = 6
	NPCBehaviorEmergent3 NPCBehaviorType = 7
)

// NPCEmergentState tracks deterministic NPC behavior state
type NPCEmergentState struct {
	mu               sync.RWMutex
	EntityID         string
	CurrentBehavior  NPCBehaviorType
	BehaviorTick     uint64
	ResonanceInfluence float32
	ActionHistory    []NPCAction
	LastAction      NPCAction
}

// NPCAction represents a deterministic NPC action
type NPCAction struct {
	Type       string
	Timestamp  int64
	Resonance  float64
	TargetX, Y int32
}

// NPCEmergentSystem handles deterministic emergent NPC behavior
type NPCEmergentSystem struct {
	mu              sync.RWMutex
	logger          *d2util.Logger
	ouroboros       *OuroborosLogikSystem
	baalAal         *d2script.BaalAalEngine
	npcStates       map[string]*NPCEmergentState
	behaviorRules   []NPCBehaviorRule
}

// NPCBehaviorRule defines a deterministic rule for NPC behavior
type NPCBehaviorRule struct {
	Name        string
	Check       func(*NPCEmergentState, *NPCEmergentSystem) bool
	Apply       func(*NPCEmergentState, *NPCEmergentSystem) NPCBehaviorType
	Priority    int
	Description string
}

// NewNPCEmergentSystem creates a new NPC emergent behavior system
func NewNPCEmergentSystem(ouroboros *OuroborosLogikSystem, baalAal *d2script.BaalAalEngine, logLevel d2util.LogLevel) *NPCEmergentSystem {
	system := &NPCEmergentSystem{
		logger:        d2util.NewLogger(),
		ouroboros:     ouroboros,
		baalAal:       baalAal,
		npcStates:     make(map[string]*NPCEmergentState),
		behaviorRules: []NPCBehaviorRule{},
	}

	system.logger.SetPrefix("NPC-Emergent")
	system.logger.SetLevel(logLevel)

	// Register NPC-specific event handlers
	system.registerNPCHandlers()

	// Register deterministic behavior rules
	system.registerBehaviorRules()

	system.logger.Info("NPC Emergent System initialized")

	return system
}

// registerNPCHandlers binds NPC event handlers to the BaalAal event bus
func (s *NPCEmergentSystem) registerNPCHandlers() {
	// Subscribe to player move for NPC proximity checks
	s.baalAal.EventBus.Subscribe("NPC-ProximityCheck", func(event *d2script.IAxiomaticEvent) {
		s.handleProximityCheck(event)
	})

	// Subscribe to NPC-specific events
	s.baalAal.EventBus.Subscribe("NPC-Interact", func(event *d2script.IAxiomaticEvent) {
		s.handleNPCInteract(event)
	})
}

// registerBehaviorRules defines deterministic NPC behavior rules
func (s *NPCEmergentSystem) registerBehaviorRules() {
	// Rule 1: Resonance-based idle behavior
	s.behaviorRules = append(s.behaviorRules, NPCBehaviorRule{
		Name: "ResonanceIdle",
		Check: func(state *NPCEmergentState, sys *NPCEmergentSystem) bool {
			resonance := sys.ouroboros.areLogik.GlobalResonance
			return resonance > 0.4 && resonance < 0.6
		},
		Apply: func(state *NPCEmergentState, sys *NPCEmergentSystem) NPCBehaviorType {
			return NPCBehaviorIdle
		},
		Priority:    10,
		Description: "NPC idles when resonance is balanced",
	})

	// Rule 2: High resonance patrol behavior
	s.behaviorRules = append(s.behaviorRules, NPCBehaviorRule{
		Name: "HighResonancePatrol",
		Check: func(state *NPCEmergentState, sys *NPCEmergentSystem) bool {
			resonance := sys.ouroboros.areLogik.GlobalResonance
			return resonance >= 0.6
		},
		Apply: func(state *NPCEmergentState, sys *NPCEmergentSystem) NPCBehaviorType {
			return NPCBehaviorPatrol
		},
		Priority:    20,
		Description: "NPC patrols when resonance is high",
	})

	// Rule 3: Low resonance flee behavior
	s.behaviorRules = append(s.behaviorRules, NPCBehaviorRule{
		Name: "LowResonanceFlee",
		Check: func(state *NPCEmergentState, sys *NPCEmergentSystem) bool {
			resonance := sys.ouroboros.areLogik.GlobalResonance
			return resonance < 0.3
		},
		Apply: func(state *NPCEmergentState, sys *NPCEmergentSystem) NPCBehaviorType {
			return NPCBehaviorFlee
		},
		Priority:    30,
		Description: "NPC flees when resonance is very low",
	})

	// Rule 4: Emergent behavior based on entropy
	s.behaviorRules = append(s.behaviorRules, NPCBehaviorRule{
		Name: "EntropyEmergent",
		Check: func(state *NPCEmergentState, sys *NPCEmergentSystem) bool {
			entropy := sys.ouroboros.areLogik.Entropy
			return entropy > 0.7 || entropy < 0.3
		},
		Apply: func(state *NPCEmergentState, sys *NPCEmergentSystem) NPCBehaviorType {
			// Deterministic emergent behavior selection using entropy
			entropy := sys.ouroboros.areLogik.Entropy
			tick := sys.ouroboros.areLogik.Tick

			// Use entropy and tick to deterministically select emergent behavior
			selection := int((entropy + float64(tick%10)*0.1) * 3)
			selection = selection % 3

			switch selection {
			case 0:
				return NPCBehaviorEmergent1
			case 1:
				return NPCBehaviorEmergent2
			default:
				return NPCBehaviorEmergent3
			}
		},
		Priority:    5,
		Description: "NPC exhibits emergent behaviors when entropy is extreme",
	})
}

// handleProximityCheck processes player proximity to NPCs
func (s *NPCEmergentSystem) handleProximityCheck(event *d2script.IAxiomaticEvent) {
	if event == nil || event.Metadata == nil {
		return
	}

	clientID, ok := event.Metadata["client_id"].(string)
	if !ok {
		return
	}

	// Check if player is near any NPC
	nearbyNPCs := s.findNearbyNPCs(event)

	for _, npcID := range nearbyNPCs {
		s.triggerNPCResponse(npcID, clientID, event)
	}
}

// handleNPCInteract processes NPC interaction events
func (s *NPCEmergentSystem) handleNPCInteract(event *d2script.IAxiomaticEvent) {
	if event == nil || event.Metadata == nil {
		return
	}

	npcID, ok := event.Metadata["npc_id"].(string)
	if !ok {
		return
	}

	s.mu.RLock()
	state, exists := s.npcStates[npcID]
	s.mu.RUnlock()

	if !exists {
		state = s.createNPCState(npcID)
	}

	// Record interaction
	s.recordNPCHistory(state, NPCAction{
		Type:      "Interact",
		Timestamp: event.Timestamp,
		Resonance: s.ouroboros.areLogik.GlobalResonance,
	})

	// Publish NPC interaction event
	s.baalAal.EventBus.Publish(&d2script.IAxiomaticEvent{
		ID:   "NPC-Interaction-" + npcID,
		Type: "NPCInteraction",
		Metadata: map[string]interface{}{
			"npc_id":    npcID,
			"behavior":  state.CurrentBehavior,
			"resonance": s.ouroboros.areLogik.GlobalResonance,
		},
	})
}

// createNPCState creates a new NPC state entry
func (s *NPCEmergentSystem) createNPCState(entityID string) *NPCEmergentState {
	state := &NPCEmergentState{
		EntityID:          entityID,
		CurrentBehavior:   NPCBehaviorIdle,
		BehaviorTick:      s.ouroboros.areLogik.Tick,
		ResonanceInfluence: 0.5,
		ActionHistory:     []NPCAction{},
	}

	s.mu.Lock()
	s.npcStates[entityID] = state
	s.mu.Unlock()

	return state
}

// findNearbyNPCs finds NPCs near the event location
func (s *NPCEmergentSystem) findNearbyNPCs(event *d2script.IAxiomaticEvent) []string {
	// Extract player position
	var playerKappaX, playerKappaY int32

	if kx, ok := event.Metadata["kappa_x"].(int32); ok {
		playerKappaX = kx
	} else if x, ok := event.Metadata["x"].(float64); ok {
		playerKappaX = s.ouroboros.ToKappa(x)
	}

	if ky, ok := event.Metadata["kappa_y"].(int32); ok {
		playerKappaY = ky
	} else if y, ok := event.Metadata["y"].(float64); ok {
		playerKappaY = s.ouroboros.ToKappa(y)
	}

	// Find chunks within proximity
	proximityRange := int32(50) // KAPPA units
	var nearbyNPCs []string

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Iterate through all entity states and check proximity
	for _, entityState := range s.ouroboros.entityStates {
		if len(entityState.KappaPos) != 2 {
			continue
		}

		dx := entityState.KappaPos[0] - playerKappaX
		dy := entityState.KappaPos[1] - playerKappaY

		// Simple distance check (Manhattan distance for efficiency)
		if math.Abs(float64(dx)) < float64(proximityRange) && math.Abs(float64(dy)) < float64(proximityRange) {
			nearbyNPCs = append(nearbyNPCs, entityState.EntityID)
		}
	}

	return nearbyNPCs
}

// triggerNPCResponse determines and applies NPC response to player presence
func (s *NPCEmergentSystem) triggerNPCResponse(npcID, playerID string, event *d2script.IAxiomaticEvent) {
	state := s.getOrCreateNPCState(npcID)

	// Determine new behavior based on deterministic rules
	newBehavior := s.determineBehavior(state)

	if newBehavior != state.CurrentBehavior {
		s.transitionBehavior(state, newBehavior)
	}

	// Update resonance influence
	state.mu.Lock()
	state.ResonanceInfluence = float32(s.ouroboros.areLogik.GlobalResonance)
	state.mu.Unlock()

	// Record action
	s.recordNPCHistory(state, NPCAction{
		Type:      "ProximityResponse",
		Timestamp: event.Timestamp,
		Resonance: s.ouroboros.areLogik.GlobalResonance,
		TargetX:   0,
		TargetY:   0,
	})
}

// determineBehavior applies deterministic rules to select NPC behavior
func (s *NPCEmergentSystem) determineBehavior(state *NPCEmergentState) NPCBehaviorType {
	// Sort rules by priority (higher priority first)
	// For simplicity, iterate through and find first matching rule

	for i := len(s.behaviorRules) - 1; i >= 0; i-- {
		rule := s.behaviorRules[i]
		if rule.Check(state, s) {
			return rule.Apply(state, s)
		}
	}

	return NPCBehaviorIdle
}

// transitionBehavior handles NPC behavior state transitions
func (s *NPCEmergentSystem) transitionBehavior(state *NPCEmergentState, newBehavior NPCBehaviorType) {
	state.mu.Lock()
	defer state.mu.Unlock()

	oldBehavior := state.CurrentBehavior
	state.CurrentBehavior = newBehavior
	state.BehaviorTick = s.ouroboros.areLogik.Tick

	s.logger.Debugf("NPC %s: behavior %v -> %v at tick %d",
		state.EntityID, oldBehavior, newBehavior, state.BehaviorTick)

	// Publish behavior change event
	s.baalAal.EventBus.Publish(&d2script.IAxiomaticEvent{
		ID:   "NPC-BehaviorChange-" + state.EntityID,
		Type: "NPCBehaviorChange",
		Metadata: map[string]interface{}{
			"npc_id":         state.EntityID,
			"old_behavior":   oldBehavior,
			"new_behavior":   newBehavior,
			"tick":           state.BehaviorTick,
			"resonance":      s.ouroboros.areLogik.GlobalResonance,
			"entropy":        s.ouroboros.areLogik.Entropy,
		},
	})
}

// recordNPCHistory records an action in NPC history
func (s *NPCEmergentSystem) recordNPCHistory(state *NPCEmergentState, action NPCAction) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.ActionHistory = append(state.ActionHistory, action)
	state.LastAction = action

	// Limit history size
	if len(state.ActionHistory) > 100 {
		state.ActionHistory = state.ActionHistory[len(state.ActionHistory)-100:]
	}
}

// getOrCreateNPCState gets or creates NPC state
func (s *NPCEmergentSystem) getOrCreateNPCState(entityID string) *NPCEmergentState {
	s.mu.RLock()
	state, exists := s.npcStates[entityID]
	s.mu.RUnlock()

	if !exists {
		state = s.createNPCState(entityID)
	}

	return state
}

// Advance processes NPC emergent behavior
func (s *NPCEmergentSystem) Advance() {
	s.mu.RLock()
	tick := s.ouroboros.areLogik.Tick
	s.mu.RUnlock()

	s.mu.RLock()
	for _, state := range s.npcStates {
		state.mu.Lock()

		// Check if behavior should be re-evaluated
		if tick-state.BehaviorTick > 60 { // ~1 second at 60fps
			newBehavior := s.determineBehavior(state)
			if newBehavior != state.CurrentBehavior {
				state.mu.Unlock()
				s.transitionBehavior(state, newBehavior)
				state.mu.Lock()
			}
		}

		state.mu.Unlock()
	}
	s.mu.RUnlock()
}

// GetNPCState returns the state of an NPC
func (s *NPCEmergentSystem) GetNPCState(entityID string) *NPCEmergentState {
	return s.getOrCreateNPCState(entityID)
}

// GetAllNPCStates returns all NPC states
func (s *NPCEmergentSystem) GetAllNPCStates() map[string]*NPCEmergentState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*NPCEmergentState)
	for k, v := range s.npcStates {
		result[k] = v
	}
	return result
}