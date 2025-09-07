package cs2log

import "time"

// Custom event types for CS2 logs that are not in the original parser

// PlayerLeftBuyzone is received when a player leaves the buy zone
type PlayerLeftBuyzone struct {
	Meta
	Player    Player   `json:"player"`
	Equipment []string `json:"equipment"`
}

// PlayerValidated is received when a player's Steam ID is validated
type PlayerValidated struct {
	Meta
	Player Player `json:"player"`
}

// PlayerAccolade is received when a player gets an achievement/award
type PlayerAccolade struct {
	Meta
	Type    string  `json:"type"`     // "3k", "4k", "5k", "mvp", etc.
	Player  Player  `json:"player"`
	Value   float64 `json:"value"`
	IsFinal bool    `json:"is_final"` // FINAL vs ROUND
}

// MatchStatus is received for match score updates
type MatchStatus struct {
	Meta
	ScoreCT      int    `json:"score_ct"`
	ScoreT       int    `json:"score_t"`
	Map          string `json:"map"`
	RoundsPlayed int    `json:"rounds_played"`
}

// TeamPlaying is received when team names are set
type TeamPlaying struct {
	Meta
	Side     string `json:"side"`      // "CT" or "TERRORIST"
	TeamName string `json:"team_name"` // Actual team name
}

// MatchPause is received when match is paused/unpaused
type MatchPause struct {
	Meta
	Action string `json:"action"` // "enabled", "disabled", "unpaused"
	Reason string `json:"reason,omitempty"`
}

// GrenadeThrowDebug is received for grenade trajectory debug data
type GrenadeThrowDebug struct {
	Meta
	Player       Player        `json:"player"`
	GrenadeType  string        `json:"grenade_type"` // "molotov", "smokegrenade", "flashgrenade", "hegrenade"
	Position     PositionFloat `json:"position"`
	Velocity     Velocity      `json:"velocity"`
	DebugCommand string        `json:"debug_command"` // Full sv_throw command
}

// ServerCvar is received when a server variable changes
type ServerCvar struct {
	Meta
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RconCommand is received when an RCON command is executed
type RconCommand struct {
	Meta
	Source  string `json:"source"`  // IP:Port of RCON client
	Command string `json:"command"` // Command executed
}

// LoadingMap is received when server starts loading a map
type LoadingMap struct {
	Meta
	Map string `json:"map"`
}

// StartedMap is received when map is fully loaded
type StartedMap struct {
	Meta
	Map string `json:"map"`
}

// LogFile is received for log file events
type LogFile struct {
	Meta
	Action   string `json:"action"`   // "started", "closed"
	Filename string `json:"filename,omitempty"`
}

// MatchStatusTeam is received when team assignments are shown
type MatchStatusTeam struct {
	Meta
	Side     string `json:"side"`
	TeamName string `json:"team_name"`
}

// TriggeredEvent is a generic event for World triggered events
type TriggeredEvent struct {
	Meta
	Event string            `json:"event"`
	Data  map[string]string `json:"data,omitempty"`
}

// ChatCommand is received when a player uses a chat command
type ChatCommand struct {
	Meta
	Player  Player `json:"player"`
	Command string `json:"command"` // The command without the dot (e.g., "pause", "ready")
	Args    string `json:"args,omitempty"`
	Text    string `json:"text"` // Full text including the command
}

// GameOverDetailed provides more detail about game ending
type GameOverDetailed struct {
	Meta
	Mode     string `json:"mode"`      // "competitive", "casual", etc.
	Map      string `json:"map"`
	ScoreCT  int    `json:"score_ct"`
	ScoreT   int    `json:"score_t"`
	Duration int    `json:"duration"` // in minutes
}

// StatsJSON represents JSON statistics dump events
type StatsJSON struct {
	Meta
	Type string `json:"stats_type"` // "start", "end", "player_data", etc.
	Data string `json:"data"`       // Raw JSON data
}

// RoundStatsFields defines the field names for round statistics
type RoundStatsFields struct {
	Meta
	Fields []string `json:"fields"` // Field names in order
}

// RoundStatsJSON represents complete round statistics in JSON format
type RoundStatsJSON struct {
	Meta
	Name        string                       `json:"name"`         // "round_stats"
	RoundNumber int                          `json:"round_number"`
	ScoreT      int                          `json:"score_t"`
	ScoreCT     int                          `json:"score_ct"`
	Map         string                       `json:"map"`
	Server      string                       `json:"server"`
	Fields      []string                     `json:"fields"`       // Field names
	Players     map[string]RoundStatsPlayerData `json:"players"`      // player_0, player_1, etc.
}

// RoundStatsPlayerData represents player data within RoundStatsJSON
type RoundStatsPlayerData struct {
	AccountID    int     `json:"accountid"`
	Team         int     `json:"team"`           // 1=T, 2=CT
	Money        int     `json:"money"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	Damage       int     `json:"damage"`
	HeadshotPct  float64 `json:"headshot_pct"`   // HSP percentage
	KDR          float64 `json:"kdr"`            // Kill/Death ratio
	ADR          int     `json:"adr"`            // Average Damage per Round
	MVP          int     `json:"mvp"`
	EnemiesFlashed int  `json:"enemies_flashed"` // EF
	UtilityDamage int   `json:"utility_damage"`  // UD
	TripleKills  int    `json:"triple_kills"`    // 3K
	QuadKills    int    `json:"quad_kills"`      // 4K
	AceKills     int    `json:"ace_kills"`       // 5K
	ClutchKills  int    `json:"clutch_kills"`    // clutchk
	FirstKills   int    `json:"first_kills"`     // firstk
	PistolKills  int    `json:"pistol_kills"`    // pistolk
	SniperKills  int    `json:"sniper_kills"`    // sniperk
	BlindKills   int    `json:"blind_kills"`     // blindk
	BombKills    int    `json:"bomb_kills"`      // bombk
	FireDamage   int    `json:"fire_damage"`     // firedmg
	UniqueKills  int    `json:"unique_kills"`    // uniquek
	Dinks        int    `json:"dinks"`           // Headshot dinks
	ChickenKills int    `json:"chicken_kills"`   // chickenk
}

// RoundStatsPlayer represents a single player's round statistics
type RoundStatsPlayer struct {
	Meta
	PlayerID     string `json:"player_id"`      // "player_0", "player_1", etc.
	AccountID    int    `json:"accountid"`
	Team         int    `json:"team"`           // 1=T, 2=CT
	Money        int    `json:"money"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	Damage       int    `json:"damage"`
	HeadshotPct  float64 `json:"headshot_pct"`  // HSP percentage
	KDR          float64 `json:"kdr"`           // Kill/Death ratio
	ADR          int    `json:"adr"`            // Average Damage per Round
	MVP          int    `json:"mvp"`
	EnemiesFlashed int  `json:"enemies_flashed"` // EF
	UtilityDamage int   `json:"utility_damage"`  // UD
	TripleKills  int    `json:"triple_kills"`    // 3K
	QuadKills    int    `json:"quad_kills"`      // 4K
	AceKills     int    `json:"ace_kills"`       // 5K
	ClutchKills  int    `json:"clutch_kills"`    // clutchk
	FirstKills   int    `json:"first_kills"`     // firstk
	PistolKills  int    `json:"pistol_kills"`    // pistolk
	SniperKills  int    `json:"sniper_kills"`    // sniperk
	BlindKills   int    `json:"blind_kills"`     // blindk
	BombKills    int    `json:"bomb_kills"`      // bombk
	FireDamage   int    `json:"fire_damage"`     // firedmg
	UniqueKills  int    `json:"unique_kills"`    // uniquek
	Dinks        int    `json:"dinks"`           // Headshot dinks
	ChickenKills int    `json:"chicken_kills"`   // chickenk
}

// BombEvent for additional bomb-related triggers
type BombEvent struct {
	Meta
	Player   Player   `json:"player"`
	Action   string   `json:"action"` // "begin_plant", "abort_plant", etc.
	Site     string   `json:"site,omitempty"`
	Position Position `json:"position,omitempty"`
}

// FreezePeriod for freeze period events
type FreezePeriod struct {
	Meta
	Action string `json:"action"` // "start", "end"
}

// NOTE: ServerMap and ServerName events removed because "map" and "server" 
// entries are part of JSON statistics blocks, not standalone events.

// WarmupStart is received when warmup period starts
type WarmupStart struct {
	Meta
}

// WarmupEnd is received when warmup period ends
type WarmupEnd struct {
	Meta
}

// PlayerStatistics represents detailed player statistics
type PlayerStatistics struct {
	AccountID      int     `json:"accountid"`
	Team           int     `json:"team"`            // 1=T, 2=CT, 3=Spectator
	Money          int     `json:"money"`
	Kills          int     `json:"kills"`
	Deaths         int     `json:"deaths"`
	Assists        int     `json:"assists"`
	Damage         int     `json:"damage"`
	HeadshotPct    float64 `json:"headshot_pct"`    // HSP percentage
	KDR            float64 `json:"kdr"`             // Kill/Death ratio
	ADR            int     `json:"adr"`             // Average Damage per Round
	MVP            int     `json:"mvp"`
	EnemiesFlashed int     `json:"enemies_flashed"` // EF
	UtilityDamage  int     `json:"utility_damage"`  // UD
	TripleKills    int     `json:"triple_kills"`    // 3K
	QuadKills      int     `json:"quad_kills"`      // 4K
	AceKills       int     `json:"ace_kills"`       // 5K
	ClutchKills    int     `json:"clutch_kills"`
	FirstKills     int     `json:"first_kills"`
	PistolKills    int     `json:"pistol_kills"`
	SniperKills    int     `json:"sniper_kills"`
	BlindKills     int     `json:"blind_kills"`
	BombKills      int     `json:"bomb_kills"`
	FireDamage     int     `json:"fire_damage"`
	UniqueKills    int     `json:"unique_kills"`
	Dinks          int     `json:"dinks"`
	ChickenKills   int     `json:"chicken_kills"`
}

// JSONStatistics represents a complete JSON statistics block from the logs
type JSONStatistics struct {
	Meta
	Name        string                      `json:"name"`         // e.g., "round_stats"
	RoundNumber int                         `json:"round_number"`
	ScoreT      int                         `json:"score_t"`
	ScoreCT     int                         `json:"score_ct"`
	Map         string                      `json:"map"`
	Server      string                      `json:"server"`
	Fields      []string                    `json:"fields"`       // List of field names
	Players     map[string]PlayerStatistics `json:"players"`      // Map of player_id to stats
	RawJSON     string                      `json:"raw_json"`     // Complete raw JSON for custom processing
}

// GetType returns the type of the JSONStatistics message
func (j JSONStatistics) GetType() string {
	return "JSONStatistics"
}

// GetTime returns the timestamp of the JSONStatistics message
func (j JSONStatistics) GetTime() time.Time {
	return j.Meta.Time
}