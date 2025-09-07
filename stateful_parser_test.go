package cs2log

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestParseStateful_SingleLineEvent(t *testing.T) {
	state := NewParserState()
	line := `08/29/2025 - 10:26:49.000: "ragga<6><[U:1:109933575]><TERRORIST>" purchased "item_assaultsuit"`
	
	msg, err := ParseStateful(line, state)
	if err != nil {
		t.Fatalf("Failed to parse single line event: %v", err)
	}
	
	if msg == nil {
		t.Fatal("Expected message, got nil")
	}
	
	if msg.GetType() != "PlayerPurchase" {
		t.Errorf("Expected PlayerPurchase, got %s", msg.GetType())
	}
	
	// State should be reset after single line event
	if state.InJSONBlock {
		t.Error("State should not be in JSON block after single line event")
	}
}

func TestParseStateful_JSONBlock(t *testing.T) {
	state := NewParserState()
	
	// Sample JSON block lines
	lines := []string{
		`08/31/2025 - 16:30:18.000: JSON_BEGIN{`,
		`08/31/2025 - 16:30:18.000: "name": "round_stats",`,
		`08/31/2025 - 16:30:18.000: "round_number" : "33",`,
		`08/31/2025 - 16:30:18.000: "score_t" : "16",`,
		`08/31/2025 - 16:30:18.000: "score_ct" : "15",`,
		`08/31/2025 - 16:30:18.000: "map" : "de_dust2",`,
		`08/31/2025 - 16:30:18.000: "server" : "Test Server",`,
		`08/31/2025 - 16:30:18.000: "fields" : "accountid,team,money,kills,deaths,assists,dmg,hsp,kdr,adr,mvp,ef,ud,3k,4k,5k,clutchk,firstk,pistolk,sniperk,blindk,bombk,firedmg,uniquek,dinks,chickenk"`,
		`08/31/2025 - 16:30:18.000: "players" : {`,
		`08/31/2025 - 16:30:18.000: "player_0" : "208135644,2,10250,19,23,9,2649,57.89,0.83,83,4,11,131,2,0,0,4,3,5,0,0,4,47,84,5,0"`,
		`08/31/2025 - 16:30:18.000: "player_1" : "1014228401,2,10050,23,26,4,2537,65.22,0.88,79,2,8,7,0,0,0,0,3,1,0,0,5,0,7,4,1"`,
		`08/31/2025 - 16:30:18.000: }}JSON_END`,
	}
	
	var msg Message
	var err error
	
	// Process all lines except the last one
	for i, line := range lines[:len(lines)-1] {
		msg, err = ParseStateful(line, state)
		if err != nil {
			t.Fatalf("Error parsing line %d: %v", i, err)
		}
		
		if msg != nil {
			t.Fatalf("Expected nil message at line %d (buffering), got %v", i, msg)
		}
		
		if !state.InJSONBlock {
			t.Fatalf("Expected to be in JSON block at line %d", i)
		}
	}
	
	// Process the last line - should return complete message
	msg, err = ParseStateful(lines[len(lines)-1], state)
	if err != nil {
		t.Fatalf("Error parsing final line: %v", err)
	}
	
	if msg == nil {
		t.Fatal("Expected complete message after JSON_END")
	}
	
	// Check the message type
	if msg.GetType() != "JSONStatistics" {
		t.Errorf("Expected JSONStatistics type, got %s", msg.GetType())
	}
	
	// Check parsed content
	jsonStats, ok := msg.(JSONStatistics)
	if !ok {
		t.Fatalf("Failed to cast to JSONStatistics")
	}
	
	// Debug: print raw JSON
	t.Logf("Raw JSON: %s", jsonStats.RawJSON)
	
	// Try parsing the raw JSON to debug
	var testData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStats.RawJSON), &testData); err != nil {
		t.Logf("JSON parse error: %v", err)
	} else {
		t.Logf("Parsed JSON data: %+v", testData)
	}
	
	if jsonStats.Name != "round_stats" {
		t.Errorf("Expected name 'round_stats', got '%s'", jsonStats.Name)
	}
	
	if jsonStats.RoundNumber != 33 {
		t.Errorf("Expected round number 33, got %d", jsonStats.RoundNumber)
	}
	
	if jsonStats.ScoreT != 16 {
		t.Errorf("Expected score_t 16, got %d", jsonStats.ScoreT)
	}
	
	if jsonStats.ScoreCT != 15 {
		t.Errorf("Expected score_ct 15, got %d", jsonStats.ScoreCT)
	}
	
	if jsonStats.Map != "de_dust2" {
		t.Errorf("Expected map 'de_dust2', got '%s'", jsonStats.Map)
	}
	
	// Check player data
	if len(jsonStats.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(jsonStats.Players))
	}
	
	if player0, ok := jsonStats.Players["player_0"]; ok {
		if player0.AccountID != 208135644 {
			t.Errorf("Expected player_0 accountid 208135644, got %d", player0.AccountID)
		}
		if player0.Kills != 19 {
			t.Errorf("Expected player_0 kills 19, got %d", player0.Kills)
		}
		if player0.Deaths != 23 {
			t.Errorf("Expected player_0 deaths 23, got %d", player0.Deaths)
		}
	} else {
		t.Error("player_0 not found in players map")
	}
	
	// State should be reset after complete JSON block
	if state.InJSONBlock {
		t.Error("State should not be in JSON block after complete JSON")
	}
}

func TestParseStateful_IncompleteJSONBlock(t *testing.T) {
	state := NewParserState()
	
	// Start JSON block
	msg, err := ParseStateful(`08/31/2025 - 16:30:18.000: JSON_BEGIN{`, state)
	if err != nil {
		t.Fatalf("Error starting JSON block: %v", err)
	}
	if msg != nil {
		t.Fatal("Expected nil message while buffering")
	}
	
	// Add a line
	msg, err = ParseStateful(`08/31/2025 - 16:30:18.000: "name": "round_stats",`, state)
	if err != nil {
		t.Fatalf("Error adding to JSON block: %v", err)
	}
	
	// Now a line with different timestamp - should error and reset
	msg, err = ParseStateful(`08/31/2025 - 16:30:19.000: World triggered "Round_Start"`, state)
	if err == nil {
		t.Fatal("Expected error for interrupted JSON block")
	}
	
	if !strings.Contains(err.Error(), "different timestamp") {
		t.Errorf("Expected timestamp error, got: %v", err)
	}
	
	// State should be reset
	if state.InJSONBlock {
		t.Error("State should be reset after interrupted JSON block")
	}
}

func TestParseStateful_MixedEvents(t *testing.T) {
	state := NewParserState()
	
	// Regular event
	msg, err := ParseStateful(`08/31/2025 - 16:30:17.000: World triggered "Round_Start"`, state)
	if err != nil {
		t.Fatalf("Error parsing round start: %v", err)
	}
	if msg.GetType() != "WorldRoundStart" {
		t.Errorf("Expected WorldRoundStart, got %s", msg.GetType())
	}
	
	// Start JSON block
	msg, err = ParseStateful(`08/31/2025 - 16:30:18.000: JSON_BEGIN{`, state)
	if err != nil {
		t.Fatalf("Error starting JSON: %v", err)
	}
	if msg != nil {
		t.Fatal("Expected nil while buffering JSON")
	}
	
	// Continue JSON
	msg, err = ParseStateful(`08/31/2025 - 16:30:18.000: "name": "test",`, state)
	if err != nil {
		t.Fatalf("Error continuing JSON: %v", err)
	}
	
	// End JSON
	msg, err = ParseStateful(`08/31/2025 - 16:30:18.000: }}JSON_END`, state)
	if err != nil {
		t.Fatalf("Error ending JSON: %v", err)
	}
	if msg == nil {
		t.Fatal("Expected complete JSON message")
	}
	if msg.GetType() != "JSONStatistics" {
		t.Errorf("Expected JSONStatistics, got %s", msg.GetType())
	}
	
	// Another regular event
	msg, err = ParseStateful(`08/31/2025 - 16:30:19.000: World triggered "Round_End"`, state)
	if err != nil {
		t.Fatalf("Error parsing round end: %v", err)
	}
	if msg.GetType() != "WorldRoundEnd" {
		t.Errorf("Expected WorldRoundEnd, got %s", msg.GetType())
	}
}

func TestParseStateful_RealWorldExample(t *testing.T) {
	state := NewParserState()
	
	// Real example from the provided logs
	realExample := `08/31/2025 - 16:30:18.000: JSON_BEGIN{
08/31/2025 - 16:30:18.000: "name": "round_stats",
08/31/2025 - 16:30:18.000: "round_number" : "33",
08/31/2025 - 16:30:18.000: "score_t" : "16",
08/31/2025 - 16:30:18.000: "score_ct" : "16",
08/31/2025 - 16:30:18.000: "map" : "de_nuke",
08/31/2025 - 16:30:18.000: "server" : "DraculaN | ENCE vs Monte",
08/31/2025 - 16:30:18.000: "fields" : "accountid,team,money,kills,deaths,assists,dmg,hsp,kdr,adr,mvp,ef,ud,3k,4k,5k,clutchk,firstk,pistolk,sniperk,blindk,bombk,firedmg,uniquek,dinks,chickenk"
08/31/2025 - 16:30:18.000: "players" : {
08/31/2025 - 16:30:18.000: "player_0" : "208135644,2,10250,19,23,9,2649,57.89,0.83,83,4,11,131,2,0,0,4,3,5,0,0,4,47,84,5,0"
08/31/2025 - 16:30:18.000: "player_1" : "1014228401,2,10050,23,26,4,2537,65.22,0.88,79,2,8,7,0,0,0,0,3,1,0,0,5,0,7,4,1"
08/31/2025 - 16:30:18.000: }}JSON_END`

	lines := strings.Split(realExample, "\n")
	
	var msg Message
	var err error
	
	// Process all lines
	for i, line := range lines {
		msg, err = ParseStateful(line, state)
		if err != nil {
			t.Fatalf("Error parsing line %d: %v", i, err)
		}
		
		// Only the last line should return a message
		if i < len(lines)-1 {
			if msg != nil {
				t.Fatalf("Expected nil message at line %d, got %v", i, msg)
			}
		}
	}
	
	// Final message should be complete
	if msg == nil {
		t.Fatal("Expected complete message after processing all lines")
	}
	
	jsonStats, ok := msg.(JSONStatistics)
	if !ok {
		t.Fatal("Failed to cast to JSONStatistics")
	}
	
	// Verify key fields
	if jsonStats.Name != "round_stats" {
		t.Errorf("Expected name 'round_stats', got '%s'", jsonStats.Name)
	}
	
	if jsonStats.Map != "de_nuke" {
		t.Errorf("Expected map 'de_nuke', got '%s'", jsonStats.Map)
	}
	
	if jsonStats.Server != "DraculaN | ENCE vs Monte" {
		t.Errorf("Expected server 'DraculaN | ENCE vs Monte', got '%s'", jsonStats.Server)
	}
	
	// Check timestamp
	expectedTime := time.Date(2025, 8, 31, 16, 30, 18, 0, time.UTC)
	if !jsonStats.GetTime().Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, jsonStats.GetTime())
	}
}

func TestParseStateful_EmptyJSONBlock(t *testing.T) {
	state := NewParserState()
	
	lines := []string{
		`08/31/2025 - 16:30:18.000: JSON_BEGIN{`,
		`08/31/2025 - 16:30:18.000: }}JSON_END`,
	}
	
	// Process first line
	msg, err := ParseStateful(lines[0], state)
	if err != nil {
		t.Fatalf("Error starting JSON: %v", err)
	}
	if msg != nil {
		t.Fatal("Expected nil while buffering")
	}
	
	// Process end line
	msg, err = ParseStateful(lines[1], state)
	if err != nil {
		t.Fatalf("Error ending JSON: %v", err)
	}
	
	if msg == nil {
		t.Fatal("Expected message for empty JSON block")
	}
	
	jsonStats, ok := msg.(JSONStatistics)
	if !ok {
		t.Fatal("Failed to cast to JSONStatistics")
	}
	
	// Should have empty/default values
	if jsonStats.Name != "" {
		t.Errorf("Expected empty name, got '%s'", jsonStats.Name)
	}
	
	if len(jsonStats.Players) != 0 {
		t.Errorf("Expected no players, got %d", len(jsonStats.Players))
	}
}