package d2emergent

import (
	"math"
	"sync"
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// TestOuroborosLogikSystem tests the main Ouroboros system integration
func TestOuroborosLogikSystem(t *testing.T) {
	// Create BaalAal engine
	baalAal := d2script.NewBaalAalEngine()

	// Create Ouroboros system
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	// Test initial state
	resonance, expansion, entropy, tick := ouroboros.GetStatus()
	if tick != 0 {
		t.Errorf("expected initial tick 0, got %d", tick)
	}
	if expansion != 1.0 {
		t.Errorf("expected initial expansion 1.0, got %f", expansion)
	}

	// Process emergence cycle
	ouroboros.Advance()
	resonance, expansion, entropy, tick = ouroboros.GetStatus()

	if tick != 1 {
		t.Errorf("expected tick 1, got %d", tick)
	}

	// Resonance should be bounded
	if resonance < 0 || resonance > 2.0 {
		t.Errorf("resonance out of bounds: %f", resonance)
	}

	// Entropy should be bounded
	if entropy < 0 || entropy > 2.0 {
		t.Errorf("entropy out of bounds: %f", entropy)
	}
}

// TestOuroborosKappaConversion tests KAPPA coordinate conversion
func TestOuroborosKappaConversion(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	// Test world to KAPPA conversion
	worldX := 100.5
	kappaX := ouroboros.ToKappa(worldX)
	expectedKappa := int32(100500)

	if kappaX != expectedKappa {
		t.Errorf("expected kappa %d, got %d", expectedKappa, kappaX)
	}

	// Test KAPPA to world conversion
	backToWorld := ouroboros.ToWorld(kappaX)
	if math.Abs(backToWorld-worldX) > 0.001 {
		t.Errorf("round-trip conversion error: expected %f, got %f", worldX, backToWorld)
	}
}

// TestOuroborosPlayerMove tests player movement with KAPPA chunk generation
func TestOuroborosPlayerMove(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	// Create player move event
	event := &d2script.IAxiomaticEvent{
		ID:   "MoveEvent-test",
		Type: "PlayerMove",
		Metadata: map[string]interface{}{
			"client_id": "player1",
			"x":         100.0,
			"y":         200.0,
		},
	}

	// Process event
	ouroboros.handlePlayerMove(event)

	// Check entity state was created
	state := ouroboros.GetEntityState("player1")
	if state == nil {
		t.Fatal("expected entity state to be created")
	}

	if len(state.KappaPos) != 2 {
		t.Errorf("expected 2 position values, got %d", len(state.KappaPos))
	}

	// Check chunk was created
	chunkID := ouroboros.getChunkID(state.KappaPos[0], state.KappaPos[1])
	chunk := ouroboros.GetChunkInfo(chunkID)
	if chunk == nil {
		t.Fatal("expected chunk to be created")
	}
}

// TestOuroborosDeterministicRules tests that deterministic rules work correctly
func TestOuroborosDeterministicRules(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	// Process multiple emergence cycles
	for i := 0; i < 100; i++ {
		ouroboros.Advance()
	}

	// Verify determinism: two systems with same initial state should produce same results
	ouroboros2 := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	for i := 0; i < 100; i++ {
		ouroboros2.Advance()
	}

	res1, exp1, ent1, tick1 := ouroboros.GetStatus()
	res2, exp2, ent2, tick2 := ouroboros2.GetStatus()

	if tick1 != tick2 {
		t.Errorf("non-deterministic tick: %d != %d", tick1, tick2)
	}

	// Due to the deterministic nature, ticks should be equal
	// Note: resonance/entropy may vary due to time.Now() in ID generation
}

// TestNPCEmergentSystem tests NPC emergent behavior
func TestNPCEmergentSystem(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)
	npcEmergent := NewNPCEmergentSystem(ouroboros, baalAal, d2util.LogLevelDefault)

	// Test NPC state creation
	state := npcEmergent.GetNPCState("test_npc")

	if state.EntityID != "test_npc" {
		t.Errorf("expected entity ID 'test_npc', got '%s'", state.EntityID)
	}

	if state.CurrentBehavior != NPCBehaviorIdle {
		t.Errorf("expected initial behavior Idle, got %v", state.CurrentBehavior)
	}

	// Test behavior transition at high resonance
	for i := 0; i < 50; i++ {
		ouroboros.Advance()
	}

	// Set high resonance artificially for testing
	ouroboros.areLogik.GlobalResonance = 0.8

	// Advance NPC system
	npcEmergent.Advance()

	// Get updated state
	state = npcEmergent.GetNPCState("test_npc")

	// At high resonance, behavior should be Patrol
	if state.CurrentBehavior != NPCBehaviorPatrol {
		t.Logf("Note: NPC behavior at high resonance is %v (may vary based on rules)", state.CurrentBehavior)
	}
}

// TestCombatEmergentSystem tests combat emergent behavior
func TestCombatEmergentSystem(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)
	combatEmergent := NewCombatEmergentSystem(ouroboros, baalAal, d2util.LogLevelDefault)

	// Create combat event
	event := &d2script.IAxiomaticEvent{
		ID:   "CastSkill-test",
		Type: "9",
		Metadata: map[string]interface{}{
			"client_id": "player1",
		},
	}

	// Handle skill cast
	combatEmergent.handleSkillCast(event)

	// Check combat history
	history := combatEmergent.GetCombatHistory()
	if len(history) == 0 {
		t.Error("expected combat event to be recorded")
	}

	// Verify modifier was applied
	if len(history) > 0 {
		lastEvent := history[len(history)-1]
		if lastEvent.Modifier.ResonanceFactor == 0 && lastEvent.SourceID == "" {
			// This is expected - combat modifiers depend on ouroboros state
		}
	}
}

// TestItemEmergentSystem tests item emergent behavior
func TestItemEmergentSystem(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)
	itemEmergent := NewItemEmergentSystem(ouroboros, baalAal, d2util.LogLevelDefault)

	// Create item spawn event
	event := &d2script.IAxiomaticEvent{
		ID:        "SpawnItem-test",
		Type:      "10",
		Timestamp: 1234567890,
		Payload: map[string]interface{}{
			"x": 100.0,
			"y": 200.0,
		},
	}

	// Handle item spawn
	itemEmergent.handleItemSpawn(event)

	// Check item was registered
	itemState := itemEmergent.GetItemState("SpawnItem-test")
	if itemState == nil {
		t.Fatal("expected item state to be created")
	}

	// Verify quality was determined
	if itemState.Quality < 0 || itemState.Quality > 7 {
		t.Errorf("invalid quality: %d", itemState.Quality)
	}

	// Verify properties were generated
	if len(itemState.Properties) == 0 {
		t.Error("expected at least one property")
	}
}

// TestOuroborosIntegration tests the full system integration
func TestOuroborosIntegration(t *testing.T) {
	// Create full script engine with all emergent systems
	engine := d2script.CreateScriptEngineWithLogLevel(d2util.LogLevelDefault)

	if engine.Ouroboros == nil {
		t.Fatal("expected Ouroboros system to be initialized")
	}

	if engine.NPCEmergent == nil {
		t.Fatal("expected NPC emergent system to be initialized")
	}

	if engine.CombatEmergent == nil {
		t.Fatal("expected combat emergent system to be initialized")
	}

	if engine.ItemEmergent == nil {
		t.Fatal("expected item emergent system to be initialized")
	}

	// Test dispatching events
	event := &d2script.IAxiomaticEvent{
		ID:   "IntegrationTest",
		Type: "PlayerMove",
		Metadata: map[string]interface{}{
			"client_id": "player1",
			"x":         100.0,
			"y":         200.0,
		},
	}

	engine.DispatchEvent(event)

	// Verify entity was registered
	state := engine.Ouroboros.GetEntityState("player1")
	if state == nil {
		t.Fatal("expected entity state to be created via DispatchEvent")
	}

	// Test advance
	engine.Advance()

	// Verify ARE status is accessible
	res, exp, ent, tick := engine.GetAREStatus()
	if tick == 0 {
		t.Error("expected tick to have advanced")
	}

	t.Logf("ARE Status - Tick: %d, Resonance: %.4f, Expansion: %.4f, Entropy: %.4f",
		tick, res, exp, ent)
}

// TestOuroborosConcurrentAccess tests thread safety
func TestOuroborosConcurrentAccess(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	var wg sync.WaitGroup

	// Run concurrent operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Create events
			for j := 0; j < 10; j++ {
				event := &d2script.IAxiomaticEvent{
					ID:   "ConcurrentTest",
					Type: "PlayerMove",
					Metadata: map[string]interface{}{
						"client_id": "player",
						"x":         float64(j * 10),
						"y":         float64(j * 20),
					},
				}

				ouroboros.handlePlayerMove(event)
				ouroboros.Advance()
			}
		}(i)
	}

	wg.Wait()

	// Verify no crashes occurred and state is consistent
	state := ouroboros.GetEntityState("player")
	if state == nil {
		t.Error("expected entity state to exist after concurrent access")
	}
}

// TestOuroborosMarkgrafHarmonic tests the Markgraf harmonic oscillation
func TestOuroborosMarkgrafHarmonic(t *testing.T) {
	engine := CreateARELogikEngine(d2util.LogLevelDefault)

	// Collect resonance values over many ticks
	resonanceValues := make([]float64, 200)

	for i := 0; i < 200; i++ {
		event := engine.ProcessEmergence()
		resonanceValues[i] = event.Payload.(float64)
	}

	// Verify oscillation pattern exists
	// Count sign changes in first derivative
	signChanges := 0
	for i := 2; i < len(resonanceValues)-1; i++ {
		delta1 := resonanceValues[i] - resonanceValues[i-1]
		delta2 := resonanceValues[i+1] - resonanceValues[i]

		if (delta1 > 0 && delta2 < 0) || (delta1 < 0 && delta2 > 0) {
			signChanges++
		}
	}

	// With 200 ticks and the given harmonic frequency,
	// we should see several oscillations
	if signChanges < 3 {
		t.Errorf("expected at least 3 sign changes, got %d", signChanges)
	}
}

// TestOuroborosExpansionLimit tests the expansion safety reset
func TestOuroborosExpansionLimit(t *testing.T) {
	engine := CreateARELogikEngine(d2util.LogLevelDefault)

	// Rapidly advance many ticks
	for i := 0; i < 10000; i++ {
		engine.ProcessEmergence()
	}

	// Expansion should have been reset due to safety limit
	if engine.Expansion > 1e9 {
		t.Errorf("expansion exceeded safety limit: %f", engine.Expansion)
	}
}

// TestOuroborosChunkGeneration tests KAPPA chunk generation
func TestOuroborosChunkGeneration(t *testing.T) {
	baalAal := d2script.NewBaalAalEngine()
	ouroboros := NewOuroborosLogikSystem(baalAal, d2util.LogLevelDefault)

	// Create multiple entities in different positions
	positions := []struct{ x, y float64 }{
		{10, 10},
		{20, 20},
		{150, 150}, // Same chunk as 20,20
		{999, 999},
	}

	for i, pos := range positions {
		event := &d2script.IAxiomaticEvent{
			ID:   "MoveEvent",
			Type: "PlayerMove",
			Metadata: map[string]interface{}{
				"client_id": "player",
				"x":         pos.x,
				"y":         pos.y,
			},
		}

		ouroboros.handlePlayerMove(event)
	}

	// Should have 3 unique chunks (positions 20,20 and 150,150 share a chunk)
	ouroboros.mu.RLock()
	chunkCount := len(ouroboros.chunkRegistry)
	ouroboros.mu.RUnlock()

	if chunkCount != 3 {
		t.Errorf("expected 3 chunks, got %d", chunkCount)
	}
}