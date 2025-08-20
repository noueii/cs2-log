package cs2log

import (
	"regexp"
	"time"
)

// OrderedPattern represents a pattern with its handler function
type OrderedPattern struct {
	Pattern *regexp.Regexp
	Handler MessageFunc
}

// GetOrderedPatterns returns all patterns in the correct priority order
// More specific patterns come before general patterns
func GetOrderedPatterns() []OrderedPattern {
	patterns := []OrderedPattern{
		// Original specific patterns first (these should match before general ones)
		{regexp.MustCompile(ServerMessagePattern), NewServerMessage},
		{regexp.MustCompile(FreezTimeStartPattern), NewFreezTimeStart},
		{regexp.MustCompile(WorldMatchStartPattern), NewWorldMatchStart},
		{regexp.MustCompile(WorldRoundStartPattern), NewWorldRoundStart},
		{regexp.MustCompile(WorldRoundRestartPattern), NewWorldRoundRestart},
		{regexp.MustCompile(WorldRoundEndPattern), NewWorldRoundEnd},
		{regexp.MustCompile(WorldGameCommencingPattern), NewWorldGameCommencing},
		{regexp.MustCompile(TeamScoredPattern), NewTeamScored},
		{regexp.MustCompile(TeamNoticePattern), NewTeamNotice},
		{regexp.MustCompile(PlayerConnectedPattern), NewPlayerConnected},
		{regexp.MustCompile(PlayerDisconnectedPattern), NewPlayerDisconnected},
		{regexp.MustCompile(PlayerEnteredPattern), NewPlayerEntered},
		{regexp.MustCompile(PlayerBannedPattern), NewPlayerBanned},
		{regexp.MustCompile(PlayerSwitchedPattern), NewPlayerSwitched},
		
		// Chat command MUST come before PlayerSay
		{regexp.MustCompile(ChatCommandPattern), NewChatCommand},
		
		{regexp.MustCompile(PlayerSayPattern), NewPlayerSay},
		{regexp.MustCompile(PlayerPurchasePattern), NewPlayerPurchase},
		{regexp.MustCompile(PlayerKillPattern), NewPlayerKill},
		{regexp.MustCompile(PlayerKillAssistPattern), NewPlayerKillAssist},
		{regexp.MustCompile(PlayerAttackPattern), NewPlayerAttack},
		{regexp.MustCompile(PlayerKilledBombPattern), NewPlayerKilledBomb},
		{regexp.MustCompile(PlayerKilledSuicidePattern), NewPlayerKilledSuicide},
		{regexp.MustCompile(PlayerPickedUpPattern), NewPlayerPickedUp},
		{regexp.MustCompile(PlayerDroppedPattern), NewPlayerDropped},
		{regexp.MustCompile(PlayerMoneyChangePattern), NewPlayerMoneyChange},
		{regexp.MustCompile(PlayerBombGotPattern), NewPlayerBombGot},
		{regexp.MustCompile(PlayerBombPlantedPattern), NewPlayerBombPlanted},
		{regexp.MustCompile(PlayerBombDroppedPattern), NewPlayerBombDropped},
		{regexp.MustCompile(PlayerBombBeginDefusePattern), NewPlayerBombBeginDefuse},
		{regexp.MustCompile(PlayerBombDefusedPattern), NewPlayerBombDefused},
		{regexp.MustCompile(PlayerThrewPattern), NewPlayerThrew},
		{regexp.MustCompile(PlayerBlindedPattern), NewPlayerBlinded},
		{regexp.MustCompile(ProjectileSpawnedPattern), NewProjectileSpawned},
		{regexp.MustCompile(GameOverPattern), NewGameOver},
		
		// Custom specific patterns
		{regexp.MustCompile(PlayerLeftBuyzonePattern), NewPlayerLeftBuyzone},
		{regexp.MustCompile(PlayerValidatedPattern), NewPlayerValidated},
		{regexp.MustCompile(PlayerAccoladePattern), NewPlayerAccolade},
		{regexp.MustCompile(MatchStatusScorePattern), NewMatchStatus},
		{regexp.MustCompile(TeamPlayingPattern), NewTeamPlaying},
		{regexp.MustCompile(MatchStatusTeamPattern), NewTeamPlaying},
		{regexp.MustCompile(MatchPauseEnabledPattern), NewMatchPauseEnabled},
		{regexp.MustCompile(MatchPauseDisabledPattern), NewMatchPauseDisabled},
		{regexp.MustCompile(MatchUnpausePattern), NewMatchUnpause},
		{regexp.MustCompile(GrenadeThrowDebugPattern), NewGrenadeThrowDebug},
		{regexp.MustCompile(ServerCvarPattern), NewServerCvar},
		{regexp.MustCompile(MpCvarPattern), NewServerCvar},
		{regexp.MustCompile(RconCommandPattern), NewRconCommand},
		{regexp.MustCompile(LoadingMapPattern), NewLoadingMap},
		{regexp.MustCompile(StartedMapPattern), NewStartedMap},
		{regexp.MustCompile(LogFileStartedPattern), NewLogFileStarted},
		{regexp.MustCompile(LogFileClosedPattern), NewLogFileClosed},
		{regexp.MustCompile(GameOverDetailedPattern), NewGameOverDetailed},
		{regexp.MustCompile(BombBeginPlantPattern), NewBombBeginPlant},
		{regexp.MustCompile(BombPlantedTriggerPattern), NewBombBeginPlant},
		{regexp.MustCompile(BombDefusedTriggerPattern), NewBombBeginPlant},
		{regexp.MustCompile(FreezePeriodStartPattern), NewFreezePeriodStart},
		{regexp.MustCompile(FreezePeriodEndPattern), NewFreezePeriodEnd},
		{regexp.MustCompile(StatsJSONStartPattern), NewStatsJSONStart},
		{regexp.MustCompile(StatsJSONEndPattern), NewStatsJSONEnd},
		
		// TriggeredEvent MUST be last as it's very general
		{regexp.MustCompile(TriggeredEventPattern), NewTriggeredEvent},
	}
	
	return patterns
}

// ParseOrdered parses using ordered patterns for correct priority
func ParseOrdered(line string) (Message, error) {
	// pattern for date, beginning of a log message
	result := LogLinePattern.FindStringSubmatch(line)
	
	// if result set is empty, parsing failed, return error
	if result == nil {
		return nil, ErrorNoMatch
	}
	
	// parse time
	ti, err := time.Parse("01/02/2006 - 15:04:05", result[1])
	
	// if parsing the date failed, return error
	if err != nil {
		return nil, err
	}
	
	// Check patterns in order
	patterns := GetOrderedPatterns()
	for _, p := range patterns {
		if matches := p.Pattern.FindStringSubmatch(result[2]); matches != nil {
			return p.Handler(ti, matches), nil
		}
	}
	
	// if there was no match above but format of the log message was correct
	// it's a valid logline but pattern is not defined, return unknown type
	return NewUnknown(ti, result[1:]), nil
}