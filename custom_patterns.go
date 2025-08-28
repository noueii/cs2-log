package cs2log

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Custom event patterns
const (
	// Buy Zone Events
	PlayerLeftBuyzonePattern = `"(.+?)<(\d+)><(.+?)><(.*?)>" left buyzone with \[(.*?)\]`
	
	// Validation Events
	PlayerValidatedPattern = `"(.+?)<(\d+)><(.+?)><>" STEAM USERID validated`
	
	// Achievement/Award Events - handle tabs or commas as delimiters
	PlayerAccoladePattern = `ACCOLADE, (FINAL|ROUND): \{(.+?)\}[,\t]\s*(.+?)<(\d+)>[,\t]\s*VALUE: ([^,\t]+)`
	
	// Match Status Events
	MatchStatusScorePattern = `MatchStatus: Score: (\d+):(\d+) on map "(.+?)" RoundsPlayed: (-?\d+)`
	TeamPlayingPattern      = `Team playing "(TERRORIST|CT)": (.+)`
	MatchStatusTeamPattern  = `MatchStatus: Team playing "(TERRORIST|CT)": (.+)`
	
	// Pause Events
	MatchPauseEnabledPattern  = `Match pause is enabled`
	MatchPauseDisabledPattern = `Match pause is disabled`
	MatchUnpausePattern       = `Match unpaused`
	
	// Debug Events
	GrenadeThrowDebugPattern = `"(.+?)" (sv_throw_\w+) (-?\d+\.?\d*) (-?\d+\.?\d*) (-?\d+\.?\d*) (-?\d+\.?\d*) (-?\d+\.?\d*) (-?\d+\.?\d*)`
	
	// Server Configuration Events
	ServerCvarPattern = `server_cvar: "(.+?)" "(.+?)"`
	MpCvarPattern     = `"(mp_.+?)" = "(.+?)"`
	
	// RCON Events
	RconCommandPattern = `rcon from "(.+?)": command "(.+?)"`
	
	// Map Events
	LoadingMapPattern = `Loading map "(.+?)"`
	StartedMapPattern = `Started map "(.+?)"`
	
	// Log File Events
	LogFileStartedPattern = `Log file started \(file "(.+?)"\)`
	LogFileClosedPattern  = `Log file closed`
	
	// Extended Chat Events (commands)
	ChatCommandPattern = `"(.+?)<(\d+)><(.+?)><(.*?)>" say(?:_team)? "\.(\w+)\s*(.*)"`
	
	// Round Statistics Events
	RoundStatsFieldsPattern = `"fields"\s*:\s*"([^"]+)"`
	RoundStatsPlayerPattern = `"(player_\d+)"\s*:\s*"([^"]+)"`
	
	// JSON Round Stats markers
	JSONBeginPattern = `JSON_BEGIN\{`
	JSONEndPattern = `\}\}JSON_END`
	RoundStatsNamePattern = `"name"\s*:\s*"round_stats"`
	RoundStatsRoundPattern = `"round_number"\s*:\s*"(\d+)"`
	RoundStatsScoreTPattern = `"score_t"\s*:\s*"(\d+)"`
	RoundStatsScoreCTPattern = `"score_ct"\s*:\s*"(\d+)"`
	RoundStatsMapPattern = `"map"\s*:\s*"([^"]+)"`
	RoundStatsServerPattern = `"server"\s*:\s*"([^"]+)"`
	RoundStatsPlayersStartPattern = `"players"\s*:\s*\{`
	
	// Game Over with details
	GameOverDetailedPattern = `Game Over: (\w+) (.+?) score (\d+):(\d+) after (\d+) min`
	
	// Triggered Events
	TriggeredEventPattern = `(World|Team ".*?") triggered "(.+?)"`
	
	// Bomb Events (extended)
	BombBeginPlantPattern = `"(.+?)<(\d+)><(.+?)><(.*?)>" triggered "Bomb_Begin_Plant"`
	BombPlantedTriggerPattern = `"(.+?)<(\d+)><(.+?)><(.*?)>" triggered "Planted_The_Bomb"`
	BombDefusedTriggerPattern = `"(.+?)<(\d+)><(.+?)><(.*?)>" triggered "Defused_The_Bomb"`
	
	// Freeze Period
	FreezePeriodStartPattern = `Starting Freeze period`
	FreezePeriodEndPattern   = `World triggered "Round_Freeze_End"`
	
	// JSON Statistics
	StatsJSONStartPattern  = `JSON_BEGIN\{`
	StatsJSONEndPattern    = `\}\}JSON_END`
	StatsJSONPlayerPattern = `"player_\d+": \{`
	
	// NOTE: ServerMap and ServerName patterns were removed because 
	// "map" and "server" entries are part of JSON statistics blocks,
	// not standalone events. They should be parsed as StatsJSON.
	
	// Warmup Events
	WarmupStartPattern = `World triggered "Warmup_Start"`
	WarmupEndPattern = `World triggered "Warmup_End"`
)

// Constructor functions for custom events

func NewPlayerLeftBuyzone(ti time.Time, r []string) Message {
	equipment := strings.Fields(r[5])
	return PlayerLeftBuyzone{
		Meta:      NewMeta(ti, "PlayerLeftBuyzone"),
		Player:    NewPlayer(r[1], r[2], r[3], r[4]),
		Equipment: equipment,
	}
}

func NewPlayerValidated(ti time.Time, r []string) Message {
	return PlayerValidated{
		Meta:   NewMeta(ti, "PlayerValidated"),
		Player: NewPlayer(r[1], r[2], r[3], ""),
	}
}

func NewPlayerAccolade(ti time.Time, r []string) Message {
	value, _ := strconv.ParseFloat(r[5], 64)
	playerName := r[3]
	playerID := r[4]
	
	return PlayerAccolade{
		Meta:    NewMeta(ti, "PlayerAccolade"),
		Type:    r[2],
		Player:  NewPlayer(playerName, playerID, "", ""),
		Value:   value,
		IsFinal: r[1] == "FINAL",
	}
}

func NewMatchStatus(ti time.Time, r []string) Message {
	scoreCT, _ := strconv.Atoi(r[1])
	scoreT, _ := strconv.Atoi(r[2])
	rounds, _ := strconv.Atoi(r[4])
	
	return MatchStatus{
		Meta:         NewMeta(ti, "MatchStatus"),
		ScoreCT:      scoreCT,
		ScoreT:       scoreT,
		Map:          r[3],
		RoundsPlayed: rounds,
	}
}

func NewTeamPlaying(ti time.Time, r []string) Message {
	return TeamPlaying{
		Meta:     NewMeta(ti, "TeamPlaying"),
		Side:     r[1],
		TeamName: r[2],
	}
}

func NewMatchPause(ti time.Time, action string, reason string) Message {
	return MatchPause{
		Meta:   NewMeta(ti, "MatchPause"),
		Action: action,
		Reason: reason,
	}
}

func NewMatchPauseEnabled(ti time.Time, r []string) Message {
	return NewMatchPause(ti, "enabled", "")
}

func NewMatchPauseDisabled(ti time.Time, r []string) Message {
	return NewMatchPause(ti, "disabled", "")
}

func NewMatchUnpause(ti time.Time, r []string) Message {
	return NewMatchPause(ti, "unpaused", "")
}

func NewGrenadeThrowDebug(ti time.Time, r []string) Message {
	// Extract grenade type from sv_throw_xxx command
	grenadeType := strings.TrimPrefix(r[2], "sv_throw_")
	
	posX, _ := strconv.ParseFloat(r[3], 32)
	posY, _ := strconv.ParseFloat(r[4], 32)
	posZ, _ := strconv.ParseFloat(r[5], 32)
	
	velX, _ := strconv.ParseFloat(r[6], 32)
	velY, _ := strconv.ParseFloat(r[7], 32)
	velZ, _ := strconv.ParseFloat(r[8], 32)
	
	return GrenadeThrowDebug{
		Meta:        NewMeta(ti, "GrenadeThrowDebug"),
		Player:      NewPlayer(r[1], "", "", ""),
		GrenadeType: grenadeType,
		Position: PositionFloat{
			X: float32(posX),
			Y: float32(posY),
			Z: float32(posZ),
		},
		Velocity: Velocity{
			X: float32(velX),
			Y: float32(velY),
			Z: float32(velZ),
		},
		DebugCommand: strings.Join(r[2:], " "),
	}
}

func NewServerCvar(ti time.Time, r []string) Message {
	return ServerCvar{
		Meta:  NewMeta(ti, "ServerCvar"),
		Name:  r[1],
		Value: r[2],
	}
}

func NewRconCommand(ti time.Time, r []string) Message {
	return RconCommand{
		Meta:    NewMeta(ti, "RconCommand"),
		Source:  r[1],
		Command: r[2],
	}
}

func NewLoadingMap(ti time.Time, r []string) Message {
	return LoadingMap{
		Meta: NewMeta(ti, "LoadingMap"),
		Map:  r[1],
	}
}

func NewStartedMap(ti time.Time, r []string) Message {
	return StartedMap{
		Meta: NewMeta(ti, "StartedMap"),
		Map:  r[1],
	}
}

func NewLogFileStarted(ti time.Time, r []string) Message {
	return LogFile{
		Meta:     NewMeta(ti, "LogFile"),
		Action:   "started",
		Filename: r[1],
	}
}

func NewLogFileClosed(ti time.Time, r []string) Message {
	return LogFile{
		Meta:   NewMeta(ti, "LogFile"),
		Action: "closed",
	}
}

func NewChatCommand(ti time.Time, r []string) Message {
	return ChatCommand{
		Meta:    NewMeta(ti, "ChatCommand"),
		Player:  NewPlayer(r[1], r[2], r[3], r[4]),
		Command: r[5],
		Args:    r[6],
		Text:    "." + r[5] + " " + r[6],
	}
}

func NewGameOverDetailed(ti time.Time, r []string) Message {
	scoreCT, _ := strconv.Atoi(r[3])
	scoreT, _ := strconv.Atoi(r[4])
	duration, _ := strconv.Atoi(r[5])
	
	return GameOverDetailed{
		Meta:     NewMeta(ti, "GameOverDetailed"),
		Mode:     r[1],
		Map:      r[2],
		ScoreCT:  scoreCT,
		ScoreT:   scoreT,
		Duration: duration,
	}
}

func NewTriggeredEvent(ti time.Time, r []string) Message {
	data := make(map[string]string)
	// Parse additional data if present
	if len(r) > 2 {
		// Store any additional captured groups
		for i := 3; i < len(r); i++ {
			data[strconv.Itoa(i-2)] = r[i]
		}
	}
	
	return TriggeredEvent{
		Meta:  NewMeta(ti, "TriggeredEvent"),
		Event: r[2],
		Data:  data,
	}
}

func NewBombBeginPlant(ti time.Time, r []string) Message {
	return BombEvent{
		Meta:   NewMeta(ti, "BombEvent"),
		Player: NewPlayer(r[1], r[2], r[3], r[4]),
		Action: "begin_plant",
	}
}

func NewFreezePeriodStart(ti time.Time, r []string) Message {
	return FreezePeriod{
		Meta:   NewMeta(ti, "FreezePeriod"),
		Action: "start",
	}
}

func NewFreezePeriodEnd(ti time.Time, r []string) Message {
	return FreezePeriod{
		Meta:   NewMeta(ti, "FreezePeriod"),
		Action: "end",
	}
}

func NewStatsJSON(ti time.Time, statsType string, data string) Message {
	return StatsJSON{
		Meta: NewMeta(ti, "StatsJSON"),
		Type: statsType,
		Data: data,
	}
}

func NewJSONBegin(ti time.Time, r []string) Message {
	return StatsJSON{
		Meta: NewMeta(ti, "StatsJSON"),
		Type: "begin",
		Data: "JSON_BEGIN{",
	}
}

func NewJSONEnd(ti time.Time, r []string) Message {
	return StatsJSON{
		Meta: NewMeta(ti, "StatsJSON"),
		Type: "end",
		Data: "}}JSON_END",
	}
}

func NewRoundStatsName(ti time.Time, r []string) Message {
	return StatsJSON{
		Meta: NewMeta(ti, "StatsJSON"),
		Type: "round_stats_name",
		Data: "round_stats",
	}
}

func NewRoundStatsMetadata(ti time.Time, r []string) Message {
	return StatsJSON{
		Meta: NewMeta(ti, "StatsJSON"),
		Type: "round_stats_metadata",
		Data: r[0],
	}
}

func NewRoundStatsFields(ti time.Time, r []string) Message {
	// Split the fields by comma and trim whitespace
	fieldsStr := r[1]
	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	
	return RoundStatsFields{
		Meta:   NewMeta(ti, "RoundStatsFields"),
		Fields: fields,
	}
}

func NewRoundStatsPlayer(ti time.Time, r []string) Message {
	playerID := r[1]
	statsStr := r[2]
	
	// Split the stats by comma and trim whitespace
	stats := strings.Split(statsStr, ",")
	for i := range stats {
		stats[i] = strings.TrimSpace(stats[i])
	}
	
	// Parse all the statistics
	player := RoundStatsPlayer{
		Meta:     NewMeta(ti, "RoundStatsPlayer"),
		PlayerID: playerID,
	}
	
	// Parse each field based on position
	// Expected order: accountid, team, money, kills, deaths, assists, dmg, hsp, kdr, adr, mvp, ef, ud, 3k, 4k, 5k, clutchk, firstk, pistolk, sniperk, blindk, bombk, firedmg, uniquek, dinks, chickenk
	if len(stats) >= 26 {
		player.AccountID, _ = strconv.Atoi(stats[0])
		player.Team, _ = strconv.Atoi(stats[1])
		player.Money, _ = strconv.Atoi(stats[2])
		player.Kills, _ = strconv.Atoi(stats[3])
		player.Deaths, _ = strconv.Atoi(stats[4])
		player.Assists, _ = strconv.Atoi(stats[5])
		player.Damage, _ = strconv.Atoi(stats[6])
		player.HeadshotPct, _ = strconv.ParseFloat(stats[7], 64)
		player.KDR, _ = strconv.ParseFloat(stats[8], 64)
		player.ADR, _ = strconv.Atoi(stats[9])
		player.MVP, _ = strconv.Atoi(stats[10])
		player.EnemiesFlashed, _ = strconv.Atoi(stats[11])
		player.UtilityDamage, _ = strconv.Atoi(stats[12])
		player.TripleKills, _ = strconv.Atoi(stats[13])
		player.QuadKills, _ = strconv.Atoi(stats[14])
		player.AceKills, _ = strconv.Atoi(stats[15])
		player.ClutchKills, _ = strconv.Atoi(stats[16])
		player.FirstKills, _ = strconv.Atoi(stats[17])
		player.PistolKills, _ = strconv.Atoi(stats[18])
		player.SniperKills, _ = strconv.Atoi(stats[19])
		player.BlindKills, _ = strconv.Atoi(stats[20])
		player.BombKills, _ = strconv.Atoi(stats[21])
		player.FireDamage, _ = strconv.Atoi(stats[22])
		player.UniqueKills, _ = strconv.Atoi(stats[23])
		player.Dinks, _ = strconv.Atoi(stats[24])
		player.ChickenKills, _ = strconv.Atoi(stats[25])
	}
	
	return player
}

func NewStatsJSONStart(ti time.Time, r []string) Message {
	return NewStatsJSON(ti, "start", r[0])
}

func NewStatsJSONEnd(ti time.Time, r []string) Message {
	return NewStatsJSON(ti, "end", r[0])
}

// NOTE: NewServerMap and NewServerName removed because these events
// are part of JSON statistics blocks, not standalone events.

func NewWarmupStart(ti time.Time, r []string) Message {
	return WarmupStart{
		Meta: NewMeta(ti, "WarmupStart"),
	}
}

func NewWarmupEnd(ti time.Time, r []string) Message {
	return WarmupEnd{
		Meta: NewMeta(ti, "WarmupEnd"),
	}
}

// Helper function to create a Player struct
func NewPlayer(name, id, steamID, side string) Player {
	idInt, _ := strconv.Atoi(id)
	return Player{
		Name:    name,
		ID:      idInt,
		SteamID: steamID,
		Side:    side,
	}
}

// ExtendedPatterns contains all custom patterns
var ExtendedPatterns = map[*regexp.Regexp]MessageFunc{
	// Buy Zone
	regexp.MustCompile(PlayerLeftBuyzonePattern): NewPlayerLeftBuyzone,
	
	// Validation
	regexp.MustCompile(PlayerValidatedPattern): NewPlayerValidated,
	
	// Achievements
	regexp.MustCompile(PlayerAccoladePattern): NewPlayerAccolade,
	
	// Match Status
	regexp.MustCompile(MatchStatusScorePattern): NewMatchStatus,
	regexp.MustCompile(TeamPlayingPattern):      NewTeamPlaying,
	regexp.MustCompile(MatchStatusTeamPattern):  NewTeamPlaying,
	
	// Pause Events
	regexp.MustCompile(MatchPauseEnabledPattern):  NewMatchPauseEnabled,
	regexp.MustCompile(MatchPauseDisabledPattern): NewMatchPauseDisabled,
	regexp.MustCompile(MatchUnpausePattern):       NewMatchUnpause,
	
	// Debug
	regexp.MustCompile(GrenadeThrowDebugPattern): NewGrenadeThrowDebug,
	
	// Server Config
	regexp.MustCompile(ServerCvarPattern): NewServerCvar,
	regexp.MustCompile(MpCvarPattern):     NewServerCvar,
	
	// RCON
	regexp.MustCompile(RconCommandPattern): NewRconCommand,
	
	// Map Events
	regexp.MustCompile(LoadingMapPattern): NewLoadingMap,
	regexp.MustCompile(StartedMapPattern): NewStartedMap,
	
	// Log File
	regexp.MustCompile(LogFileStartedPattern): NewLogFileStarted,
	regexp.MustCompile(LogFileClosedPattern):  NewLogFileClosed,
	
	// Round Statistics
	regexp.MustCompile(RoundStatsFieldsPattern): NewRoundStatsFields,
	regexp.MustCompile(RoundStatsPlayerPattern): NewRoundStatsPlayer,
	
	// Chat Commands
	regexp.MustCompile(ChatCommandPattern): NewChatCommand,
	
	// Game Over Detailed
	regexp.MustCompile(GameOverDetailedPattern): NewGameOverDetailed,
	
	// Triggered Events
	regexp.MustCompile(TriggeredEventPattern): NewTriggeredEvent,
	
	// Bomb Events
	regexp.MustCompile(BombBeginPlantPattern):     NewBombBeginPlant,
	regexp.MustCompile(BombPlantedTriggerPattern): NewBombBeginPlant,
	regexp.MustCompile(BombDefusedTriggerPattern): NewBombBeginPlant,
	
	// Freeze Period
	regexp.MustCompile(FreezePeriodStartPattern): NewFreezePeriodStart,
	regexp.MustCompile(FreezePeriodEndPattern):   NewFreezePeriodEnd,
	
	// JSON Stats
	regexp.MustCompile(StatsJSONStartPattern): NewStatsJSONStart,
	regexp.MustCompile(StatsJSONEndPattern):   NewStatsJSONEnd,
	
	// NOTE: ServerMap and ServerName patterns removed - they're part of JSON blocks
	
	// Warmup Events
	regexp.MustCompile(WarmupStartPattern): NewWarmupStart,
	regexp.MustCompile(WarmupEndPattern):   NewWarmupEnd,
}