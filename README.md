# cs2-log

Go package for parsing cs2 server logfiles. It exports types for cs2 logfiles, their regular expressions, a function for
parsing and a function for converting to non-html-escaped JSON.

## Fork Features

This is an enhanced fork of the original [janstuemmel/cs2-log](https://github.com/janstuemmel/cs2-log) with:

- **30+ additional event types** for comprehensive CS2 log parsing
- **Enhanced parsing functions** (`ParseEnhanced`, `ParseOrdered`)
- **Improved pattern matching** for CS2's `[U:1:xxx]` Steam ID format
- **Support for corrupted/malformed log entries**
- **Custom event categories** (match management, statistics, etc.)

## Documentation

- 📚 **[Full Event Documentation](./EVENTS.md)** - Detailed documentation of all event types
- 🚀 **[Quick Reference Guide](./EVENTS_QUICK_REFERENCE.md)** - Quick lookup for events and patterns

## Usage

For more examples look at the [tests](./cs2log_test.go) and the command-line utility in [examples folder](./example).
Have also a look at [godoc](http://godoc.org/github.com/janstuemmel/cs2-log).

### Basic Usage

```go
package main

import (
	"fmt"

	cs2log "github.com/noueii/cs2-log"
)

func main() {

	var msg cs2log.Message

	// a line from a server logfile
	line := `L 11/05/2018 - 15:44:36: "Player<12><[U:1:29384012]><CT>" purchased "m4a1"`

	// parse into Message (standard parsing)
	msg, err := cs2log.Parse(line)
	
	// OR use enhanced parsing for custom events
	msg, err = cs2log.ParseEnhanced(line)
	
	// OR use ordered parsing for priority matching
	msg, err = cs2log.ParseOrdered(line)

	if err != nil {
		panic(err)
	}

	fmt.Println(msg.GetType(), msg.GetTime().String())

	// cast Message interface to PlayerPurchase type
	playerPurchase, ok := msg.(cs2log.PlayerPurchase)

	if ok != true {
		panic("casting failed")
	}

	fmt.Println(playerPurchase.Player.SteamID, playerPurchase.Item)

	// get json non-htmlescaped
	jsn := cs2log.ToJSON(msg)

	fmt.Println(jsn)
}
```

Example JSON output:

```json
{
  "time": "2018-11-05T15:44:36Z",
  "type": "PlayerPurchase",
  "player": {
    "name": "Player",
    "id": 12,
    "steam_id": "[U:1:29384012]",
    "side": "CT"
  },
  "item": "m4a1"
}
```

### Parsing Functions

The library provides multiple parsing functions:

#### Single-Line Parsing

##### `Parse(line string) (Message, error)`
Standard parsing using default CS2 event patterns. Use this for basic CS2 log parsing.

##### `ParseEnhanced(line string) (Message, error)`
Enhanced parsing that includes both default patterns and 30+ custom event types. Recommended for comprehensive log analysis.

##### `ParseOrdered(line string) (Message, error)`
Ordered parsing that respects pattern priority (e.g., chat commands before regular chat). Use when pattern matching order matters.

#### Batch Parsing (Recommended)

##### `ParseLines(lines []string) ([]Message, []error)`
Batch parsing that takes multiple log lines and returns all parsed events. Automatically handles multi-line JSON statistics blocks. This is the recommended approach for parsing log files.

##### `ParseLinesOrdered(lines []string) ([]Message, []error)`
Batch parsing with ordered pattern matching for priority-sensitive parsing.

##### `ParseLinesEnhanced(lines []string) ([]Message, []error)`
Batch parsing with enhanced pattern support for all custom event types.

### Custom Events

This fork adds support for many additional events:

- **Player Events**: `PlayerLeftBuyzone`, `PlayerValidated`, `PlayerJoinedTeam`, `PlayerAccolade`
- **Match Events**: `MatchStatus`, `RoundOfficiallyEnded`, `BeginNewMatchReady`
- **Server Events**: `ServerCvar`, `ServerSay`, `LoadingMap`, `StartedMap`, `Rcon`
- **Combat Events**: `PlayerFlashAssist`, `PlayerKilledOther`
- **Statistics**: `RoundStats` (JSON format), `PlayerAccolade`
- **Chat**: `ChatCommand` (for commands like `.ready`, `!gg`)

See [EVENTS.md](./EVENTS.md) for complete documentation of all supported events.