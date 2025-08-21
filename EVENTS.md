# CS2 Log Events Documentation

This document describes all event types that can be parsed from CS2 (Counter-Strike 2) server logs.

## Table of Contents

- [Core Events](#core-events)
  - [Server Events](#server-events)
  - [Match Events](#match-events)
  - [Team Events](#team-events)
- [Player Events](#player-events)
  - [Connection Events](#connection-events)
  - [Combat Events](#combat-events)
  - [Economy Events](#economy-events)
  - [Equipment Events](#equipment-events)
  - [Bomb Events](#bomb-events)
  - [Communication Events](#communication-events)
- [Custom Events](#custom-events)
  - [Player State Events](#player-state-events)
  - [Match Management Events](#match-management-events)
  - [Server Management Events](#server-management-events)
  - [Statistics Events](#statistics-events)

---

## Core Events

### Server Events

#### ServerMessage
General server messages and console output.
```
server_message: "mp_roundtime" = "1.92"
```

#### FreezTimeStart
Triggered when the freeze time period begins at the start of a round.
```
Starting Freeze period
```

### Match Events

#### WorldMatchStart
Triggered when a match officially begins.
```
World triggered "Match_Start" on "de_dust2"
```

#### WorldRoundStart
Triggered when a new round starts.
```
World triggered "Round_Start"
```

#### WorldRoundRestart
Triggered when a round is restarted (typically during warmup or practice).
```
World triggered "Restart_Round_(1_second)"
```

#### WorldRoundEnd
Triggered when a round ends.
```
World triggered "Round_End"
```

#### WorldGameCommencing
Triggered when the game is about to start.
```
World triggered "Game_Commencing"
```

#### GameOver
Triggered when the match ends with final score.
```
Game Over: competitive mg_active de_dust2 score 16:8 after 42 min
```

### Team Events

#### TeamScored
Records when a team scores a round.
```
Team "CT" scored "1" with "5" players
```

#### TeamNotice
Various team-related notifications (e.g., team winning, bomb defused/exploded).
```
Team "CT" triggered "SFUI_Notice_CTs_Win" (CT "1") (T "0")
```

---

## Player Events

### Connection Events

#### PlayerConnected
When a player connects to the server.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]"
  },
  "address": "192.168.1.100"
}
```

#### PlayerDisconnected
When a player disconnects from the server.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "reason": "Disconnect"
}
```

#### PlayerEntered
When a player enters the game (spawns for the first time).
```
"Player1<2><[U:1:123456789]><>" entered the game
```

#### PlayerBanned
When a player is banned from the server.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]"
  },
  "duration": "5 minutes",
  "by": "Console"
}
```

#### PlayerSwitched
When a player switches teams.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]"
  },
  "from": "Spectator",
  "to": "CT"
}
```

### Combat Events

#### PlayerKill
When a player kills another player.
```json
{
  "attacker": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "attacker_pos": {"x": 100, "y": 200, "z": 50},
  "victim": {
    "name": "Player2",
    "id": 3,
    "steam_id": "[U:1:987654321]",
    "side": "TERRORIST"
  },
  "victim_pos": {"x": 150, "y": 250, "z": 50},
  "weapon": "ak47",
  "headshot": true,
  "penetrated": false
}
```

#### PlayerKillAssist
When a player assists in killing another player.
```json
{
  "attacker": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "victim": {
    "name": "Player2",
    "id": 3,
    "steam_id": "[U:1:987654321]",
    "side": "TERRORIST"
  }
}
```

#### PlayerFlashAssist
When a player gets a flash assist for a kill.
```json
{
  "attacker": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "victim": {
    "name": "Player2",
    "id": 3,
    "steam_id": "[U:1:987654321]",
    "side": "TERRORIST"
  }
}
```

#### PlayerAttack
When a player damages another player.
```json
{
  "attacker": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "attacker_pos": {"x": 100, "y": 200, "z": 50},
  "victim": {
    "name": "Player2",
    "id": 3,
    "steam_id": "[U:1:987654321]",
    "side": "TERRORIST"
  },
  "victim_pos": {"x": 150, "y": 250, "z": 50},
  "weapon": "ak47",
  "damage": 27,
  "damage_armor": 3,
  "health": 73,
  "armor": 97,
  "hitgroup": "chest"
}
```

#### PlayerKilledBomb
When a player is killed by the bomb explosion.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "position": {"x": 100, "y": 200, "z": 50}
}
```

#### PlayerKilledSuicide
When a player commits suicide.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "position": {"x": 100, "y": 200, "z": 50},
  "with": "world"
}
```

#### PlayerKilledOther
When a player kills a non-player entity (e.g., chicken).
```json
{
  "attacker": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "victim": "chicken",
  "weapon": "knife"
}
```

### Economy Events

#### PlayerPurchase
When a player purchases an item.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "item": "ak47"
}
```

#### PlayerMoneyChange
When a player's money changes.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "equation": {
    "start": 2000,
    "delta": -2700,
    "end": 4700,
    "operation": "+"
  },
  "purchase": "ak47"
}
```

### Equipment Events

#### PlayerPickedUp
When a player picks up an item.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "item": "ak47"
}
```

#### PlayerDropped
When a player drops an item.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "item": "ak47"
}
```

#### PlayerThrew
When a player throws a grenade.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "position": {"x": 100, "y": 200, "z": 50},
  "grenade": "flashbang",
  "entindex": 234
}
```

#### PlayerBlinded
When a player is blinded by a flashbang.
```json
{
  "victim": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "attacker": {
    "name": "Player2",
    "id": 3,
    "steam_id": "[U:1:987654321]",
    "side": "TERRORIST"
  },
  "for": 3.45,
  "entindex": 234
}
```

#### ProjectileSpawned
When a projectile (e.g., molotov) is spawned.
```json
{
  "position": {"x": 100.5, "y": 200.3, "z": 50.2},
  "velocity": {"x": 500.1, "y": 100.2, "z": 300.3}
}
```

### Bomb Events

#### PlayerBombGot
When a player picks up the bomb.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "TERRORIST"
  }
}
```

#### PlayerBombPlanted
When a player plants the bomb.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "TERRORIST"
  }
}
```

#### PlayerBombDropped
When a player drops the bomb.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "TERRORIST"
  }
}
```

#### PlayerBombBeginDefuse
When a player begins defusing the bomb.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "kit": true
}
```

#### PlayerBombDefused
When a player successfully defuses the bomb.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  }
}
```

### Communication Events

#### PlayerSay
When a player sends a chat message.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "text": "Good luck, have fun!",
  "team": false
}
```

---

## Custom Events

These are additional events added to support more comprehensive CS2 log parsing.

### Player State Events

#### PlayerLeftBuyzone
When a player leaves the buy zone.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "equipment": ["m4a1", "deagle", "kevlar", "defuser", "flashbang", "hegrenade"]
}
```

#### PlayerValidated
When a player's Steam ID is validated by the server.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]"
  }
}
```

#### PlayerJoinedTeam
When a player joins a specific team (different from switching teams).
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]"
  },
  "team": "CT"
}
```

### Match Management Events

#### MatchStatus
Provides current match status information.
```json
{
  "score": {
    "ct": 8,
    "t": 7
  },
  "map": "de_dust2",
  "rounds_played": 15
}
```

#### RoundStart (Custom)
Enhanced round start event with more details.
```json
{
  "timelimit": 115,
  "fraglimit": 0,
  "objective": "de_dust2"
}
```

#### RoundEnd (Custom)
Enhanced round end event with winner and reason.
```json
{
  "winner": "CT",
  "reason": "Bomb_Defused",
  "message": "Counter-Terrorists Win"
}
```

#### RoundOfficiallyEnded
When a round officially ends (after round end delay).
```
World triggered "Round_Officially_Ended"
```

#### BeginNewMatchReady
When players signal they are ready for a new match.
```json
{
  "ready_players": 10,
  "needed_players": 10
}
```

#### GameCommencing
When the game is commencing (transitioning from warmup).
```
World triggered "Game_Commencing"
```

### Server Management Events

#### ServerCvar
When a server console variable is changed.
```json
{
  "cvar": "mp_roundtime",
  "value": "1.92"
}
```

#### ServerSay
When the server sends a message.
```json
{
  "message": "Match will start when all players are ready"
}
```

#### LoadingMap
When the server is loading a new map.
```json
{
  "map": "de_mirage"
}
```

#### StartedMap
When the server has finished loading a map.
```json
{
  "map": "de_mirage",
  "crc": "1234567890"
}
```

#### Rcon
When an RCON command is executed.
```json
{
  "address": "192.168.1.100:27015",
  "command": "status",
  "success": true
}
```

#### CvarSet
When a CVAR value is set.
```json
{
  "cvar": "sv_cheats",
  "value": "0"
}
```

### Statistics Events

#### PlayerAccolade
Player receives an accolade (award) at the end of a round.
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "accolade": "mvp",
  "value": 1,
  "position": 0.82,
  "score": 98,
  "nemesis_kills": 0,
  "ff_kills": 0,
  "ff_deaths": 0,
  "kills": 3
}
```

#### RoundStats
Comprehensive round statistics (JSON format).
```json
{
  "name": "round_stats",
  "round": 15,
  "map": "de_dust2",
  "players": {
    "player_0": {
      "accountid": 123456789,
      "team": 1,
      "money": 10000,
      "kills": 2,
      "deaths": 1,
      "assists": 0,
      "damage": 186,
      "headshot_percentage": 50.0,
      "kill_death_ratio": 2.0,
      "average_damage_per_round": 93,
      "mvp": 1,
      "enemies_flashed": 2,
      "utility_damage": 45,
      "triple_kills": 0,
      "quad_kills": 0,
      "ace": 0
    }
  }
}
```

#### ChatCommand
When a player uses a chat command (e.g., !ready, .gg).
```json
{
  "player": {
    "name": "Player1",
    "id": 2,
    "steam_id": "[U:1:123456789]",
    "side": "CT"
  },
  "command": ".ready",
  "team": false
}
```

---

## Event Categories

### By Frequency (Typical Match)
1. **Very Common**: PlayerAttack, PlayerSay, PlayerMoneyChange
2. **Common**: PlayerKill, PlayerPurchase, PlayerThrew, PlayerPickedUp
3. **Round Events**: WorldRoundStart, WorldRoundEnd, TeamScored
4. **Occasional**: PlayerSwitched, PlayerBlinded, PlayerKillAssist
5. **Rare**: PlayerBanned, PlayerKilledBomb, GameOver

### By Importance
1. **Critical**: PlayerKill, WorldRoundEnd, PlayerBombPlanted, PlayerBombDefused
2. **Important**: PlayerAttack, TeamScored, PlayerPurchase
3. **Informational**: PlayerSay, ServerCvar, PlayerMoneyChange
4. **Debug**: ProjectileSpawned, PlayerAccolade, RoundStats

### By Game Phase
- **Pre-match**: PlayerConnected, PlayerValidated, ServerCvar
- **Warmup**: PlayerSwitched, LoadingMap, BeginNewMatchReady
- **Freeze Time**: PlayerPurchase, PlayerLeftBuyzone
- **Round Active**: PlayerAttack, PlayerKill, PlayerThrew
- **Round End**: TeamScored, PlayerAccolade, RoundStats
- **Post-match**: GameOver, PlayerDisconnected

---

## Usage Examples

### Parsing Events
```go
import cs2log "github.com/noueii/cs2-log"

// Parse a log line
msg, err := cs2log.ParseEnhanced(logLine)
if err != nil {
    // Handle parsing error
}

// Type assertion to specific event
switch event := msg.(type) {
case cs2log.PlayerKill:
    fmt.Printf("%s killed %s with %s\n", 
        event.Attacker.Name, 
        event.Victim.Name, 
        event.Weapon)
case cs2log.PlayerMoneyChange:
    fmt.Printf("%s money: $%d %s $%d = $%d\n",
        event.Player.Name,
        event.Equation.Start,
        event.Equation.Operation,
        event.Equation.Delta,
        event.Equation.End)
}
```

### Event Filtering
```go
// Filter for combat events only
combatEvents := []string{
    "PlayerKill",
    "PlayerAttack", 
    "PlayerKillAssist",
    "PlayerFlashAssist",
    "PlayerBlinded",
}

if contains(combatEvents, event.GetType()) {
    // Process combat event
}
```

### Session Detection
```go
// Detect match start/end
if _, ok := msg.(cs2log.WorldMatchStart); ok {
    // Start new match session
}

if _, ok := msg.(cs2log.GameOver); ok {
    // End match session
}

// Detect round boundaries
if _, ok := msg.(cs2log.WorldRoundStart); ok {
    // Start new round
}

if _, ok := msg.(cs2log.WorldRoundEnd); ok {
    // End current round
}
```

---

## Notes

- All player events include Steam ID in CS2 format: `[U:1:xxxxxxxxx]`
- Position coordinates are in game units
- Money values are in dollars (not cents)
- Damage values are in HP units
- Time values are in seconds (float for precise values)
- Team sides are: "CT", "TERRORIST", "Spectator", "Unassigned"
- All events include a timestamp from when they occurred

## Unknown Events

Some events may still parse as `Unknown`. Common patterns include:
- JSON-formatted round statistics (partially supported via RoundStats)
- Server configuration messages
- Plugin-specific messages
- Corrupted or malformed log entries

For unknown events, the raw log line is preserved in the `Unknown.Raw` field for manual processing.