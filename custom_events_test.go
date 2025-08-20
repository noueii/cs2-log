package cs2log

import (
	"testing"
)

func TestPlayerLeftBuyzone(t *testing.T) {
	logLine := `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" left buyzone with [ weapon_knife weapon_usp_silencer kevlar(100) weapon_awp ]`
	
	msg, err := ParseEnhanced(logLine)
	if err != nil {
		t.Fatalf("Failed to parse PlayerLeftBuyzone: %v", err)
	}
	
	leftBuyzone, ok := msg.(PlayerLeftBuyzone)
	if !ok {
		t.Fatalf("Expected PlayerLeftBuyzone, got %T", msg)
	}
	
	if leftBuyzone.Player.Name != "Magixx" {
		t.Errorf("Expected player name 'Magixx', got '%s'", leftBuyzone.Player.Name)
	}
	
	if len(leftBuyzone.Equipment) != 4 {
		t.Errorf("Expected 4 equipment items, got %d", len(leftBuyzone.Equipment))
	}
}

func TestPlayerValidated(t *testing.T) {
	logLine := `L 08/19/2025 - 15:12:44: "sh1ro<456><STEAM_1:0:654321><>" STEAM USERID validated`
	
	msg, err := ParseEnhanced(logLine)
	if err != nil {
		t.Fatalf("Failed to parse PlayerValidated: %v", err)
	}
	
	validated, ok := msg.(PlayerValidated)
	if !ok {
		t.Fatalf("Expected PlayerValidated, got %T", msg)
	}
	
	if validated.Player.Name != "sh1ro" {
		t.Errorf("Expected player name 'sh1ro', got '%s'", validated.Player.Name)
	}
}

func TestPlayerAccolade(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected struct {
			Type    string
			IsFinal bool
			Value   float64
		}
	}{
		{
			name:    "Final 3k",
			logLine: `L 08/19/2025 - 15:12:44: ACCOLADE, FINAL: {3k}, sh1ro<456>, VALUE: 2.000000`,
			expected: struct {
				Type    string
				IsFinal bool
				Value   float64
			}{Type: "3k", IsFinal: true, Value: 2.0},
		},
		{
			name:    "Round MVP",
			logLine: `L 08/19/2025 - 15:12:44: ACCOLADE, ROUND: {mvp}, HooXi<789>, VALUE: 5.000000`,
			expected: struct {
				Type    string
				IsFinal bool
				Value   float64
			}{Type: "mvp", IsFinal: false, Value: 5.0},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseEnhanced(tt.logLine)
			if err != nil {
				t.Fatalf("Failed to parse PlayerAccolade: %v", err)
			}
			
			accolade, ok := msg.(PlayerAccolade)
			if !ok {
				t.Fatalf("Expected PlayerAccolade, got %T", msg)
			}
			
			if accolade.Type != tt.expected.Type {
				t.Errorf("Expected type '%s', got '%s'", tt.expected.Type, accolade.Type)
			}
			
			if accolade.IsFinal != tt.expected.IsFinal {
				t.Errorf("Expected IsFinal %v, got %v", tt.expected.IsFinal, accolade.IsFinal)
			}
			
			if accolade.Value != tt.expected.Value {
				t.Errorf("Expected value %f, got %f", tt.expected.Value, accolade.Value)
			}
		})
	}
}

func TestMatchStatus(t *testing.T) {
	logLine := `L 08/19/2025 - 15:12:44: MatchStatus: Score: 17:19 on map "de_dust2" RoundsPlayed: 36`
	
	msg, err := ParseEnhanced(logLine)
	if err != nil {
		t.Fatalf("Failed to parse MatchStatus: %v", err)
	}
	
	status, ok := msg.(MatchStatus)
	if !ok {
		t.Fatalf("Expected MatchStatus, got %T", msg)
	}
	
	if status.ScoreCT != 17 {
		t.Errorf("Expected CT score 17, got %d", status.ScoreCT)
	}
	
	if status.ScoreT != 19 {
		t.Errorf("Expected T score 19, got %d", status.ScoreT)
	}
	
	if status.Map != "de_dust2" {
		t.Errorf("Expected map 'de_dust2', got '%s'", status.Map)
	}
	
	if status.RoundsPlayed != 36 {
		t.Errorf("Expected 36 rounds played, got %d", status.RoundsPlayed)
	}
}

func TestMatchPause(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected string
	}{
		{
			name:     "Pause enabled",
			logLine:  `L 08/19/2025 - 15:12:44: Match pause is enabled`,
			expected: "enabled",
		},
		{
			name:     "Pause disabled",
			logLine:  `L 08/19/2025 - 15:12:44: Match pause is disabled`,
			expected: "disabled",
		},
		{
			name:     "Match unpaused",
			logLine:  `L 08/19/2025 - 15:12:44: Match unpaused`,
			expected: "unpaused",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseEnhanced(tt.logLine)
			if err != nil {
				t.Fatalf("Failed to parse MatchPause: %v", err)
			}
			
			pause, ok := msg.(MatchPause)
			if !ok {
				t.Fatalf("Expected MatchPause, got %T", msg)
			}
			
			if pause.Action != tt.expected {
				t.Errorf("Expected action '%s', got '%s'", tt.expected, pause.Action)
			}
		})
	}
}

func TestChatCommand(t *testing.T) {
	tests := []struct {
		name     string
		logLine  string
		expected struct {
			Command string
			Args    string
		}
	}{
		{
			name:    "Pause command",
			logLine: `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" say ".pause"`,
			expected: struct {
				Command string
				Args    string
			}{Command: "pause", Args: ""},
		},
		{
			name:    "Restore with args",
			logLine: `L 08/19/2025 - 15:12:44: "sh1ro<456><STEAM_1:0:654321><T>" say ".restore 35"`,
			expected: struct {
				Command string
				Args    string
			}{Command: "restore", Args: "35"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseEnhanced(tt.logLine)
			if err != nil {
				t.Fatalf("Failed to parse ChatCommand: %v", err)
			}
			
			cmd, ok := msg.(ChatCommand)
			if !ok {
				t.Fatalf("Expected ChatCommand, got %T", msg)
			}
			
			if cmd.Command != tt.expected.Command {
				t.Errorf("Expected command '%s', got '%s'", tt.expected.Command, cmd.Command)
			}
			
			if cmd.Args != tt.expected.Args {
				t.Errorf("Expected args '%s', got '%s'", tt.expected.Args, cmd.Args)
			}
		})
	}
}

func TestServerCvar(t *testing.T) {
	tests := []struct {
		name    string
		logLine string
		expected struct {
			Name  string
			Value string
		}
	}{
		{
			name:    "server_cvar format",
			logLine: `L 08/19/2025 - 15:12:44: server_cvar: "mp_freezetime" "20"`,
			expected: struct {
				Name  string
				Value string
			}{Name: "mp_freezetime", Value: "20"},
		},
		{
			name:    "mp_ format",
			logLine: `L 08/19/2025 - 15:12:44: "mp_maxrounds" = "24"`,
			expected: struct {
				Name  string
				Value string
			}{Name: "mp_maxrounds", Value: "24"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseEnhanced(tt.logLine)
			if err != nil {
				t.Fatalf("Failed to parse ServerCvar: %v", err)
			}
			
			cvar, ok := msg.(ServerCvar)
			if !ok {
				t.Fatalf("Expected ServerCvar, got %T", msg)
			}
			
			if cvar.Name != tt.expected.Name {
				t.Errorf("Expected name '%s', got '%s'", tt.expected.Name, cvar.Name)
			}
			
			if cvar.Value != tt.expected.Value {
				t.Errorf("Expected value '%s', got '%s'", tt.expected.Value, cvar.Value)
			}
		})
	}
}

func TestRconCommand(t *testing.T) {
	logLine := `L 08/19/2025 - 15:12:44: rcon from "192.168.1.100:12345": command "mp_pause_match 1"`
	
	msg, err := ParseEnhanced(logLine)
	if err != nil {
		t.Fatalf("Failed to parse RconCommand: %v", err)
	}
	
	rcon, ok := msg.(RconCommand)
	if !ok {
		t.Fatalf("Expected RconCommand, got %T", msg)
	}
	
	if rcon.Source != "192.168.1.100:12345" {
		t.Errorf("Expected source '192.168.1.100:12345', got '%s'", rcon.Source)
	}
	
	if rcon.Command != "mp_pause_match 1" {
		t.Errorf("Expected command 'mp_pause_match 1', got '%s'", rcon.Command)
	}
}

func TestGrenadeThrowDebug(t *testing.T) {
	logLine := `L 08/19/2025 - 15:12:44: "Magixx" sv_throw_molotov -1943.109 1620.291 94.267 -123.456 456.789 789.012`
	
	msg, err := ParseEnhanced(logLine)
	if err != nil {
		t.Fatalf("Failed to parse GrenadeThrowDebug: %v", err)
	}
	
	debug, ok := msg.(GrenadeThrowDebug)
	if !ok {
		t.Fatalf("Expected GrenadeThrowDebug, got %T", msg)
	}
	
	if debug.Player.Name != "Magixx" {
		t.Errorf("Expected player 'Magixx', got '%s'", debug.Player.Name)
	}
	
	if debug.GrenadeType != "molotov" {
		t.Errorf("Expected grenade type 'molotov', got '%s'", debug.GrenadeType)
	}
	
	// Check position values (approximate due to float conversion)
	if debug.Position.X < -1944 || debug.Position.X > -1943 {
		t.Errorf("Expected X position around -1943, got %f", debug.Position.X)
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that original events still work
	tests := []struct {
		name    string
		logLine string
		msgType string
	}{
		{
			name:    "Player kill",
			logLine: `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" [1 2 3] killed "sh1ro<456><STEAM_1:0:654321><TERRORIST>" [4 5 6] with "ak47"`,
			msgType: "PlayerKill",
		},
		{
			name:    "Round start",
			logLine: `L 08/19/2025 - 15:12:44: World triggered "Round_Start"`,
			msgType: "WorldRoundStart",
		},
		{
			name:    "Player say",
			logLine: `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" say "nice shot"`,
			msgType: "PlayerSay",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseEnhanced(tt.logLine)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.msgType, err)
			}
			
			if msg.GetType() != tt.msgType {
				t.Errorf("Expected type '%s', got '%s'", tt.msgType, msg.GetType())
			}
		})
	}
}

// Benchmark to compare performance
func BenchmarkParseOriginal(b *testing.B) {
	logLine := `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" [1 2 3] killed "sh1ro<456><STEAM_1:0:654321><T>" [4 5 6] with "ak47"`
	
	for i := 0; i < b.N; i++ {
		Parse(logLine)
	}
}

func BenchmarkParseEnhanced(b *testing.B) {
	logLine := `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" [1 2 3] killed "sh1ro<456><STEAM_1:0:654321><T>" [4 5 6] with "ak47"`
	
	for i := 0; i < b.N; i++ {
		ParseEnhanced(logLine)
	}
}

func BenchmarkParseCustomEvent(b *testing.B) {
	logLine := `L 08/19/2025 - 15:12:44: "Magixx<123><STEAM_1:0:123456><CT>" left buyzone with [ weapon_knife weapon_usp_silencer kevlar(100) weapon_awp ]`
	
	for i := 0; i < b.N; i++ {
		ParseEnhanced(logLine)
	}
}