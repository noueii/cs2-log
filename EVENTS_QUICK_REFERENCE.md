# CS2 Log Events - Quick Reference

## Event Types by Category

### 🎮 Match Flow
| Event | Description | Key Fields |
|-------|-------------|------------|
| `WorldMatchStart` | Match begins | map |
| `WorldRoundStart` | Round starts | - |
| `WorldRoundEnd` | Round ends | - |
| `GameOver` | Match ends | mode, map, score, duration |
| `TeamScored` | Team wins round | team, score, players |
| `TeamNotice` | Team notification | team, notice, scores |

### 👤 Player Connection
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerConnected` | Player connects | player, address |
| `PlayerDisconnected` | Player disconnects | player, reason |
| `PlayerEntered` | Player spawns first time | player |
| `PlayerValidated` | Steam ID validated | player |
| `PlayerSwitched` | Team change | player, from, to |
| `PlayerJoinedTeam` | Joins team | player, team |

### ⚔️ Combat
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerKill` | Player kills another | attacker, victim, weapon, headshot |
| `PlayerKillAssist` | Kill assist | attacker, victim |
| `PlayerFlashAssist` | Flash assist | attacker, victim |
| `PlayerAttack` | Damage dealt | attacker, victim, damage, hitgroup |
| `PlayerKilledBomb` | Killed by bomb | player, position |
| `PlayerKilledSuicide` | Suicide | player, with |
| `PlayerKilledOther` | Killed entity | attacker, victim, weapon |

### 💰 Economy
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerPurchase` | Item purchased | player, item |
| `PlayerMoneyChange` | Money change | player, equation, purchase |
| `PlayerPickedUp` | Item picked up | player, item |
| `PlayerDropped` | Item dropped | player, item |

### 💣 Bomb
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerBombGot` | Picks up bomb | player |
| `PlayerBombPlanted` | Plants bomb | player |
| `PlayerBombDropped` | Drops bomb | player |
| `PlayerBombBeginDefuse` | Starts defuse | player, kit |
| `PlayerBombDefused` | Defuses bomb | player |

### 🔫 Equipment
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerThrew` | Throws grenade | player, grenade, position |
| `PlayerBlinded` | Blinded by flash | victim, attacker, duration |
| `ProjectileSpawned` | Projectile created | position, velocity |
| `PlayerLeftBuyzone` | Exits buy zone | player, equipment |

### 💬 Communication
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerSay` | Chat message | player, text, team |
| `ChatCommand` | Chat command | player, command |
| `ServerSay` | Server message | message |

### ⚙️ Server
| Event | Description | Key Fields |
|-------|-------------|------------|
| `ServerCvar` | CVAR change | cvar, value |
| `CvarSet` | CVAR set | cvar, value |
| `Rcon` | RCON command | address, command, success |
| `LoadingMap` | Loading map | map |
| `StartedMap` | Map loaded | map, crc |

### 📊 Statistics
| Event | Description | Key Fields |
|-------|-------------|------------|
| `PlayerAccolade` | Player award | player, accolade, value |
| `RoundStats` | Round statistics | round, map, players |
| `MatchStatus` | Match status | score, map, rounds_played |

## Event Patterns (Regex)

### Most Common Patterns
```regex
# Player Kill
"(.+)<(\d+)><([\[\]\w:_]+)><(TERRORIST|CT)>" \[(-?\d+) (-?\d+) (-?\d+)\] killed "(.+)<(\d+)><([\[\]\w:_]+)><(TERRORIST|CT)>"

# Player Attack  
"(.+)<(\d+)><([^>]*)><(TERRORIST|CT)>" \[(-?\d+) (-?\d+) (-?\d+)\] attacked "(.+)<(\d+)><([^>]*)><(TERRORIST|CT)>"

# Player Say
"(.+)<(\d+)><([\[\]\w:_]+)><(TERRORIST|CT)>" say(_team)? "(.*)"

# Money Change
"(.+)<(\d+)><([\[\]\w:_]+)><(TERRORIST|CT)>" money change (\d+)([\+\-])(\d+) = \$(\d+)
```

## Parse Function Usage

```go
// Standard parsing
msg, err := cs2log.Parse(logLine)

// Enhanced parsing (includes custom events)
msg, err := cs2log.ParseEnhanced(logLine)

// Ordered parsing (respects pattern priority)
msg, err := cs2log.ParseOrdered(logLine)
```

## Event Type Detection

```go
switch msg.(type) {
case cs2log.PlayerKill:
    // Handle kill
case cs2log.PlayerAttack:
    // Handle attack
case cs2log.PlayerMoneyChange:
    // Handle economy
case cs2log.Unknown:
    // Handle unknown
}
```

## Session Markers

### Match Session
- Start: `WorldMatchStart`
- End: `GameOver`

### Round Session
- Start: `WorldRoundStart` or `FreezTimeStart`
- End: `WorldRoundEnd` or `RoundOfficiallyEnded`

### Player Session
- Start: `PlayerConnected` → `PlayerValidated` → `PlayerEntered`
- End: `PlayerDisconnected`

## Performance Tips

1. **Use ParseEnhanced** for full event support
2. **Cache parsed events** to avoid re-parsing
3. **Filter early** - check event type before processing
4. **Batch process** logs for better performance
5. **Handle Unknown events** gracefully

## Common Event Sequences

### Round Start
1. `FreezTimeStart`
2. `PlayerPurchase` (multiple)
3. `PlayerLeftBuyzone` (multiple)
4. `WorldRoundStart`

### Kill Sequence
1. `PlayerAttack` (multiple, optional)
2. `PlayerKill`
3. `PlayerKillAssist` (optional)
4. `PlayerFlashAssist` (optional)

### Bomb Plant/Defuse
1. `PlayerBombGot`
2. `PlayerBombPlanted`
3. `PlayerBombBeginDefuse`
4. `PlayerBombDefused` OR `TeamNotice` (bomb exploded)

### Round End
1. `WorldRoundEnd`
2. `TeamScored`
3. `TeamNotice`
4. `PlayerAccolade` (multiple)
5. `RoundStats`
6. `RoundOfficiallyEnded`

## Error Handling

```go
msg, err := cs2log.ParseEnhanced(line)
if err != nil {
    // Log line couldn't be parsed at all
    log.Printf("Parse error: %v", err)
    return
}

if unknown, ok := msg.(cs2log.Unknown); ok {
    // Recognized as log but unknown event type
    log.Printf("Unknown event: %s", unknown.Raw)
}
```