package d2emergent

import (
	"math"
	"sync"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// CombatModifier represents a deterministic combat modification
type CombatModifier struct {
	DamageMultiplier float64
	DefenseBonus    float64
	AccuracyBonus   float64
	CriticalChance  float64
	ResonanceFactor  float64
}

// CombatEmergentSystem handles deterministic combat with resonance feedback
type CombatEmergentSystem struct {
	mu              sync.RWMutex
	logger          *d2util.Logger
	ouroboros       *OuroborosLogikSystem
	baalAal         *d2script.BaalAalEngine
	combatHistory   []CombatEvent
	activeBuffs     map[string]*CombatBuff
}

// CombatEvent records a combat interaction
type CombatEvent struct {
	SourceID      string
	TargetID      string
	Damage        float64
	Modifier      CombatModifier
	Resonance     float64
	Timestamp     int64
	IsCritical    bool
}

// CombatBuff represents a time-limited combat effect
type CombatBuff struct {
	EntityID      string
	Type          string
	Modifier      CombatModifier
	RemainingTick uint64
	StartTick     uint64
}

// NewCombatEmergentSystem creates a new combat emergent system
func NewCombatEmergentSystem(ouroboros *OuroborosLogikSystem, baalAal *d2script.BaalAalEngine, logLevel d2util.LogLevel) *CombatEmergentSystem {
	system := &CombatEmergentSystem{
		logger:        d2util.NewLogger(),
		ouroboros:     ouroboros,
		baalAal:       baalAal,
		combatHistory: []CombatEvent{},
		activeBuffs:   make(map[string]*CombatBuff),
	}

	system.logger.SetPrefix("Combat-Emergent")
	system.logger.SetLevel(logLevel)

	// Register combat handlers
	system.registerCombatHandlers()

	system.logger.Info("Combat Emergent System initialized")

	return system
}

// registerCombatHandlers binds combat event handlers
func (s *CombatEmergentSystem) registerCombatHandlers() {
	// Subscribe to skill cast events
	s.baalAal.EventBus.Subscribe("Combat-Cast", func(event *d2script.IAxiomaticEvent) {
		s.handleSkillCast(event)
	})

	// Register combat rules
	s.baalAal.RegisterRule("9", s.handleSkillCast)
}

// handleSkillCast processes skill casting with deterministic modifiers
func (s *CombatEmergentSystem) handleSkillCast(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	// Compute deterministic combat modifier
	modifier := s.computeCombatModifier(event)

	// Record combat event
	s.recordCombatEvent(event, modifier)

	// Apply resonance feedback to event
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	event.Metadata["combat_modifier"] = modifier
	event.Metadata["resonance_feedback"] = modifier.ResonanceFactor
}

// computeCombatModifier calculates deterministic combat modifications
func (s *CombatEmergentSystem) computeCombatModifier(event *d2script.IAxiomaticEvent) CombatModifier {
	tick := s.ouroboros.areLogik.Tick
	resonance := s.ouroboros.areLogik.GlobalResonance
	entropy := s.ouroboros.areLogik.Entropy
	expansion := s.ouroboros.areLogik.Expansion

	// Markgraf harmonic for deterministic variation
	harmonic := math.Sin(float64(tick)*MarkgrafHarmonicFrequency) +
		0.5*math.Cos(float64(tick)*MarkgrafHarmonicFrequency*0.618)

	// Base modifier from resonance
	baseResonance := 1.0 + (resonance-0.5)*0.2

	// Entropy-based critical chance
	criticalChance := entropy * 0.1

	// Expansion affects accuracy
	accuracyBonus := math.Log(expansion) * 0.01

	// Ouroboros feedback
	feedback := s.ouroboros.resonanceOffset * MarkgrafFeedbackWeight

	return CombatModifier{
		DamageMultiplier: baseResonance + (harmonic * 0.1),
		DefenseBonus:     entropy * 0.05,
		AccuracyBonus:    accuracyBonus,
		CriticalChance:   criticalChance,
		ResonanceFactor:  feedback,
	}
}

// recordCombatEvent records a combat event for history
func (s *CombatEmergentSystem) recordCombatEvent(event *d2script.IAxiomaticEvent, modifier CombatModifier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Extract entity IDs from metadata
	sourceID := ""
	targetID := ""

	if event.Metadata != nil {
		if sid, ok := event.Metadata["client_id"].(string); ok {
			sourceID = sid
		}
		if tid, ok := event.Metadata["target_id"].(string); ok {
			targetID = tid
		}
	}

	// Check for critical hit
	isCritical := math.Random() < modifier.CriticalChance

	combatEvent := CombatEvent{
		SourceID:   sourceID,
		TargetID:   targetID,
		Damage:     modifier.DamageMultiplier * 100, // Base damage
		Modifier:   modifier,
		Resonance:  s.ouroboros.areLogik.GlobalResonance,
		Timestamp:  event.Timestamp,
		IsCritical: isCritical,
	}

	s.combatHistory = append(s.combatHistory, combatEvent)

	// Limit history size
	if len(s.combatHistory) > 1000 {
		s.combatHistory = s.combatHistory[len(s.combatHistory)-1000:]
	}
}

// ApplyBuff applies a combat buff to an entity
func (s *CombatEmergentSystem) ApplyBuff(entityID, buffType string, durationTicks uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	modifier := s.computeBuffModifier(buffType)

	s.activeBuffs[entityID+buffType] = &CombatBuff{
		EntityID:      entityID,
		Type:          buffType,
		Modifier:      modifier,
		RemainingTick: durationTicks,
		StartTick:     s.ouroboros.areLogik.Tick,
	}
}

// computeBuffModifier calculates buff effects
func (s *CombatEmergentSystem) computeBuffModifier(buffType string) CombatModifier {
	resonance := s.ouroboros.areLogik.GlobalResonance

	switch buffType {
	case "damage_boost":
		return CombatModifier{
			DamageMultiplier: 1.0 + resonance*0.5,
			ResonanceFactor:  resonance,
		}
	case "defense_boost":
		return CombatModifier{
			DefenseBonus:    resonance * 0.3,
			ResonanceFactor: resonance,
		}
	case "speed_boost":
		return CombatModifier{
			AccuracyBonus:   resonance * 0.2,
			ResonanceFactor: resonance,
		}
	default:
		return CombatModifier{
			ResonanceFactor: resonance,
		}
	}
}

// Advance processes combat system updates
func (s *CombatEmergentSystem) Advance() {
	s.mu.Lock()
	defer s.mu.Unlock()

	tick := s.ouroboros.areLogik.Tick

	// Expire buffs
	for key, buff := range s.activeBuffs {
		if tick-buff.StartTick >= buff.RemainingTick {
			delete(s.activeBuffs, key)
		}
	}
}

// GetActiveBuffs returns active buffs for an entity
func (s *CombatEmergentSystem) GetActiveBuffs(entityID string) []*CombatBuff {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var buffs []*CombatBuff
	for _, buff := range s.activeBuffs {
		if buff.EntityID == entityID {
			buffs = append(buffs, buff)
		}
	}
	return buffs
}

// GetCombatHistory returns recent combat events
func (s *CombatEmergentSystem) GetCombatHistory() []CombatEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]CombatEvent, len(s.combatHistory))
	copy(result, s.combatHistory)
	return result
}

// ItemEmergentSystem handles deterministic item spawning and properties
type ItemEmergentSystem struct {
	mu              sync.RWMutex
	logger          *d2util.Logger
	ouroboros       *OuroborosLogikSystem
	baalAal         *d2script.BaalAalEngine
	itemRegistry    map[string]*ItemEmergentState
	lootTables      map[string][]LootEntry
}

// ItemEmergentState tracks item emergence properties
type ItemEmergentState struct {
	ItemID       string
	KappaX, Y    int32
	ChunkID      string
	Generation   uint64
	Quality      int // 0-7 for normal/magic/rare/set/unique
	Properties   []ItemProperty
	Resonance    float64
	Timestamp    int64
}

// ItemProperty represents an item modifier
type ItemProperty struct {
	Name   string
	Value  float64
	Source string
}

// LootEntry defines a potential loot drop
type LootEntry struct {
	ItemCode     string
	BaseProb     float64
	ResonanceMod float64
	EntropyMod   float64
}

// NewItemEmergentSystem creates a new item emergent system
func NewItemEmergentSystem(ouroboros *OuroborosLogikSystem, baalAal *d2script.BaalAalEngine, logLevel d2util.LogLevel) *ItemEmergentSystem {
	system := &ItemEmergentSystem{
		logger:       d2util.NewLogger(),
		ouroboros:    ouroboros,
		baalAal:      baalAal,
		itemRegistry: make(map[string]*ItemEmergentState),
		lootTables:   make(map[string][]LootEntry),
	}

	system.logger.SetPrefix("Item-Emergent")
	system.logger.SetLevel(logLevel)

	// Register item handlers
	system.registerItemHandlers()

	// Initialize default loot tables
	system.initializeLootTables()

	system.logger.Info("Item Emergent System initialized")

	return system
}

// registerItemHandlers binds item event handlers
func (s *ItemEmergentSystem) registerItemHandlers() {
	// Subscribe to item spawn events
	s.baalAal.EventBus.Subscribe("Item-Spawn", func(event *d2script.IAxiomaticEvent) {
		s.handleItemSpawn(event)
	})

	// Register item rules
	s.baalAal.RegisterRule("10", s.handleItemSpawn)
}

// initializeLootTables sets up deterministic loot generation
func (s *ItemEmergentSystem) initializeLootTables() {
	// Basic loot table with deterministic probabilities
	s.lootTables["common"] = []LootEntry{
		{ItemCode: "misc_gld", BaseProb: 0.3, ResonanceMod: 0.1, EntropyMod: 0.0},
		{ItemCode: "misc_pot", BaseProb: 0.25, ResonanceMod: 0.05, EntropyMod: 0.0},
		{ItemCode: "misc_key", BaseProb: 0.05, ResonanceMod: 0.02, EntropyMod: 0.0},
	}

	s.lootTables["magic"] = []LootEntry{
		{ItemCode: "weap_norm_001", BaseProb: 0.15, ResonanceMod: 0.1, EntropyMod: 0.05},
		{ItemCode: "armor_norm_001", BaseProb: 0.15, ResonanceMod: 0.1, EntropyMod: 0.05},
	}

	s.lootTables["rare"] = []LootEntry{
		{ItemCode: "weap_magic_001", BaseProb: 0.08, ResonanceMod: 0.15, EntropyMod: 0.1},
		{ItemCode: "armor_magic_001", BaseProb: 0.08, ResonanceMod: 0.15, EntropyMod: 0.1},
	}
}

// handleItemSpawn processes item spawning with emergent properties
func (s *ItemEmergentSystem) handleItemSpawn(event *d2script.IAxiomaticEvent) {
	if event == nil {
		return
	}

	// Extract position data
	var kappaX, kappaY int32

	if event.Metadata != nil {
		if kx, ok := event.Metadata["kappa_x"].(int32); ok {
			kappaX = kx
		} else if x, ok := event.Metadata["x"].(float64); ok {
			kappaX = s.ouroboros.ToKappa(x)
		}

		if ky, ok := event.Metadata["kappa_y"].(int32); ok {
			kappaY = ky
		} else if y, ok := event.Metadata["y"].(float64); ok {
			kappaY = s.ouroboros.ToKappa(y)
		}
	}

	// Determine chunk and generation
	chunkID := s.ouroboros.getChunkID(kappaX, kappaY)
	generation := s.computeItemGeneration(event)

	// Determine quality based on resonance and entropy
	quality := s.determineItemQuality(event)

	// Create item state
	itemState := &ItemEmergentState{
		ItemID:     event.ID,
		KappaX:     kappaX,
		KappaY:     kappaY,
		ChunkID:    chunkID,
		Generation: generation,
		Quality:    quality,
		Properties: s.generateItemProperties(quality),
		Resonance:  s.ouroboros.areLogik.GlobalResonance,
		Timestamp:  event.Timestamp,
	}

	s.mu.Lock()
	s.itemRegistry[event.ID] = itemState
	s.mu.Unlock()

	// Update chunk occupancy
	if chunk, exists := s.ouroboros.chunkRegistry[chunkID]; exists {
		chunk.mu.Lock()
		chunk.Occupants = append(chunk.Occupants, event.ID)
		chunk.Generation = generation
		chunk.mu.Unlock()
	}

	// Update event with emergent properties
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	event.Metadata["item_quality"] = quality
	event.Metadata["item_generation"] = generation
	event.Metadata["chunk_id"] = chunkID
	event.Metadata["item_properties"] = itemState.Properties
}

// computeItemGeneration calculates deterministic item generation seed
func (s *ItemEmergentSystem) computeItemGeneration(event *d2script.IAxiomaticEvent) uint64 {
	tick := s.ouroboros.areLogik.Tick
	resonance := s.ouroboros.areLogik.GlobalResonance
	entropy := s.ouroboros.areLogik.Entropy

	// Deterministic generation seed using Markgraf harmonics
	harmonic := math.Sin(float64(tick)*MarkgrafHarmonicFrequency) +
		0.5*math.Cos(float64(tick)*MarkgrafHarmonicFrequency*0.618)

	seed := uint64((float64(tick) + harmonic*100) * (resonance + entropy))
	return seed
}

// determineItemQuality determines item quality based on game state
func (s *ItemEmergentSystem) determineItemQuality(event *d2script.IAxiomaticEvent) int {
	resonance := s.ouroboros.areLogik.GlobalResonance
	entropy := s.ouroboros.areLogik.Entropy
	tick := s.ouroboros.areLogik.Tick

	// Use deterministic factors for quality selection
	qualitySeed := (resonance + entropy) * float64(tick%100)

	if qualitySeed > 0.95 {
		return 7 // Unique
	} else if qualitySeed > 0.85 {
		return 5 // Set
	} else if qualitySeed > 0.70 {
		return 4 // Rare
	} else if qualitySeed > 0.50 {
		return 3 // Magic
	} else if qualitySeed > 0.25 {
		return 1 // Superior
	}
	return 0 // Normal
}

// generateItemProperties generates deterministic item properties
func (s *ItemEmergentSystem) generateItemProperties(quality int) []ItemProperty {
	var properties []ItemProperty
	resonance := s.ouroboros.areLogik.GlobalResonance

	// Base property from resonance
	properties = append(properties, ItemProperty{
		Name:   "resonance_bonus",
		Value:  resonance * 0.1,
		Source: "ouroboros",
	})

	// Quality-based properties
	switch {
	case quality >= 7: // Unique
		properties = append(properties,
			ItemProperty{Name: "all_stats", Value: 10 + int(resonance*10), Source: "unique"},
			ItemProperty{Name: "all_resist", Value: 20 + int(resonance*10), Source: "unique"},
		)
	case quality >= 5: // Set
		properties = append(properties,
			ItemProperty{Name: "strength", Value: 5 + int(resonance*5), Source: "set"},
			ItemProperty{Name: "vitality", Value: 5 + int(resonance*5), Source: "set"},
		)
	case quality >= 4: // Rare
		properties = append(properties,
			ItemProperty{Name: "strength", Value: 3 + int(resonance*3), Source: "rare"},
		)
	case quality >= 3: // Magic
		properties = append(properties,
			ItemProperty{Name: "damage", Value: 1 + int(resonance*2), Source: "magic"},
		)
	}

	return properties
}

// GetItemState returns the emergent state of an item
func (s *ItemEmergentSystem) GetItemState(itemID string) *ItemEmergentState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.itemRegistry[itemID]
}

// GetItemsInChunk returns all items in a KAPPA chunk
func (s *ItemEmergentSystem) GetItemsInChunk(chunkID string) []*ItemEmergentState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*ItemEmergentState
	for _, item := range s.itemRegistry {
		if item.ChunkID == chunkID {
			items = append(items, item)
		}
	}
	return items
}

// Advance processes item system updates
func (s *ItemEmergentSystem) Advance() {
	// Item system is mostly passive; no per-frame processing needed
}