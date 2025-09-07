package cs2log

import (
	"fmt"
	"strings"
	"testing"
)

func TestJSONStatistics_SingleComprehensiveEvent(t *testing.T) {
	// Your actual JSON block with multiple players
	jsonBlock := `08/31/2025 - 16:30:18.000: JSON_BEGIN{
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

	lines := strings.Split(jsonBlock, "\n")
	
	// Parse using ParseLines - should return exactly ONE event
	messages, errors := ParseLines(lines)
	
	// Should have no errors
	if len(errors) != 0 {
		t.Fatalf("Expected no errors, got %d: %v", len(errors), errors)
	}
	
	// Should return exactly ONE event (not multiple)
	if len(messages) != 1 {
		t.Fatalf("Expected exactly 1 message for JSON block, got %d messages", len(messages))
	}
	
	// The single event should be JSONStatistics
	msg := messages[0]
	if msg.GetType() != "JSONStatistics" {
		t.Fatalf("Expected JSONStatistics type, got %s", msg.GetType())
	}
	
	// Cast to JSONStatistics and verify it contains ALL data
	stats, ok := msg.(JSONStatistics)
	if !ok {
		t.Fatal("Failed to cast to JSONStatistics")
	}
	
	// Verify all fields are present in the single event
	fmt.Printf("Single Event Contents:\n")
	fmt.Printf("  Name: %s\n", stats.Name)
	fmt.Printf("  Round: %d\n", stats.RoundNumber)
	fmt.Printf("  Score: T=%d CT=%d\n", stats.ScoreT, stats.ScoreCT)
	fmt.Printf("  Map: %s\n", stats.Map)
	fmt.Printf("  Server: %s\n", stats.Server)
	fmt.Printf("  Fields: %d field names\n", len(stats.Fields))
	fmt.Printf("  Players: %d players\n", len(stats.Players))
	
	// Verify specific values
	if stats.Name != "round_stats" {
		t.Errorf("Expected name 'round_stats', got '%s'", stats.Name)
	}
	
	if stats.RoundNumber != 33 {
		t.Errorf("Expected round 33, got %d", stats.RoundNumber)
	}
	
	if len(stats.Players) != 2 {
		t.Errorf("Expected 2 players in the single event, got %d", len(stats.Players))
	}
	
	// Check that player data is embedded in the single event
	if player0, exists := stats.Players["player_0"]; exists {
		fmt.Printf("  Player 0 data: AccountID=%d, Kills=%d, Deaths=%d\n", 
			player0.AccountID, player0.Kills, player0.Deaths)
		
		if player0.AccountID != 208135644 {
			t.Errorf("Player 0 AccountID mismatch")
		}
	} else {
		t.Error("Player 0 not found in the single event")
	}
	
	fmt.Printf("\n✅ JSON block successfully parsed as a SINGLE comprehensive event!\n")
}