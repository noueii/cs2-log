package cs2log

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParserState maintains state for parsing multi-line events
type ParserState struct {
	InJSONBlock   bool
	JSONBuffer    []string
	JSONStartTime time.Time
	LastTimestamp string
}

// NewParserState creates a new parser state
func NewParserState() *ParserState {
	return &ParserState{
		InJSONBlock: false,
		JSONBuffer:  make([]string, 0),
	}
}

// Reset clears the parser state
func (ps *ParserState) Reset() {
	ps.InJSONBlock = false
	ps.JSONBuffer = ps.JSONBuffer[:0]
	ps.JSONStartTime = time.Time{}
	ps.LastTimestamp = ""
}

// ParseLines takes multiple log lines and returns all parsed events
// This handles both single-line events and multi-line JSON blocks correctly
func ParseLines(lines []string) ([]Message, []error) {
	state := NewParserState()
	var messages []Message
	var errors []error
	
	for _, line := range lines {
		msg, err := ParseStateful(line, state)
		
		if err != nil {
			errors = append(errors, err)
			continue
		}
		
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	
	// Check if there's an incomplete JSON block at the end
	if state.InJSONBlock {
		errors = append(errors, fmt.Errorf("incomplete JSON block at end of input"))
	}
	
	return messages, errors
}

// ParseLinesOrdered takes multiple log lines and uses ParseOrdered for single-line events
// This provides enhanced parsing with correct pattern priority
func ParseLinesOrdered(lines []string) ([]Message, []error) {
	state := NewParserState()
	var messages []Message
	var errors []error
	
	for _, line := range lines {
		msg, err := parseStatefulWithParser(line, state, ParseOrdered)
		
		if err != nil {
			errors = append(errors, err)
			continue
		}
		
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	
	// Check if there's an incomplete JSON block at the end
	if state.InJSONBlock {
		errors = append(errors, fmt.Errorf("incomplete JSON block at end of input"))
	}
	
	return messages, errors
}

// ParseLinesEnhanced takes multiple log lines and uses ParseEnhanced for single-line events
// This provides access to all custom event types
func ParseLinesEnhanced(lines []string) ([]Message, []error) {
	state := NewParserState()
	var messages []Message
	var errors []error
	
	for _, line := range lines {
		msg, err := parseStatefulWithParser(line, state, ParseEnhanced)
		
		if err != nil {
			errors = append(errors, err)
			continue
		}
		
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	
	// Check if there's an incomplete JSON block at the end
	if state.InJSONBlock {
		errors = append(errors, fmt.Errorf("incomplete JSON block at end of input"))
	}
	
	return messages, errors
}

// ParseStateful parses a log line with state management for multi-line events
// Returns nil, nil if the line is part of an incomplete multi-line event
// Returns the complete Message when a multi-line event is completed
func ParseStateful(line string, state *ParserState) (Message, error) {
	return parseStatefulWithParser(line, state, Parse)
}

// parseStatefulWithParser is the internal implementation that accepts a parser function
func parseStatefulWithParser(line string, state *ParserState, parser func(string) (Message, error)) (Message, error) {
	// First extract timestamp and content
	result := LogLinePattern.FindStringSubmatch(line)
	if result == nil {
		// If we're in a JSON block and get an invalid line, treat it as an error
		if state.InJSONBlock {
			state.Reset()
			return nil, fmt.Errorf("invalid line format during JSON block: %s", line)
		}
		return nil, ErrorNoMatch
	}

	timestamp := result[1]
	content := result[2]

	// Check if this is the start of a JSON block
	if strings.HasPrefix(content, "JSON_BEGIN{") {
		// Parse the timestamp
		ti, err := time.Parse("01/02/2006 - 15:04:05.000", timestamp)
		if err != nil {
			return nil, err
		}
		
		state.InJSONBlock = true
		state.JSONStartTime = ti
		state.LastTimestamp = timestamp
		state.JSONBuffer = append(state.JSONBuffer, content)
		return nil, nil // Buffer is building, no complete message yet
	}

	// If we're in a JSON block, accumulate lines
	if state.InJSONBlock {
		// Only accept lines with the same timestamp
		if timestamp != state.LastTimestamp {
			// Different timestamp means the JSON block was incomplete
			state.Reset()
			return nil, fmt.Errorf("JSON block interrupted by different timestamp")
		}

		state.JSONBuffer = append(state.JSONBuffer, content)

		// Check if this is the end of the JSON block
		if strings.HasSuffix(content, "}}JSON_END") {
			// Parse the complete JSON block
			msg := parseJSONBlock(state.JSONStartTime, state.JSONBuffer)
			state.Reset()
			return msg, nil
		}

		return nil, nil // Still building the buffer
	}

	// Not in a JSON block, parse as a regular single-line event
	return parser(line)
}

// parseJSONBlock parses a complete JSON statistics block
func parseJSONBlock(timestamp time.Time, lines []string) Message {
	// Join all lines to form the JSON content
	var jsonLines []string
	
	for i, line := range lines {
		if i == 0 {
			// Remove JSON_BEGIN{ prefix
			line = strings.TrimPrefix(line, "JSON_BEGIN{")
			if line != "" {
				jsonLines = append(jsonLines, line)
			}
		} else if i == len(lines)-1 {
			// Remove }}JSON_END suffix but keep the closing }
			if strings.HasSuffix(line, "}}JSON_END") {
				line = strings.TrimSuffix(line, "}JSON_END") // Keep one }
			}
			if line != "" {
				jsonLines = append(jsonLines, line)
			}
		} else {
			jsonLines = append(jsonLines, line)
		}
	}
	
	// Join lines and add commas where needed
	var jsonContent strings.Builder
	jsonContent.WriteString("{")
	
	for i, line := range jsonLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		jsonContent.WriteString("\n  ")
		jsonContent.WriteString(line)
		
		// Add comma if not last line and doesn't already end with comma or open brace
		if i < len(jsonLines)-1 {
			trimmed := strings.TrimSpace(line)
			if !strings.HasSuffix(trimmed, ",") && 
			   !strings.HasSuffix(trimmed, "{") && 
			   !strings.HasPrefix(strings.TrimSpace(jsonLines[i+1]), "}") {
				jsonContent.WriteString(",")
			}
		}
	}
	
	jsonContent.WriteString("\n}")
	fullJSON := jsonContent.String()

	// Try to parse as structured JSON statistics
	var stats JSONStatistics
	stats.Meta = NewMeta(timestamp, "JSONStatistics")
	stats.RawJSON = fullJSON

	// Parse the JSON to extract key fields
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(fullJSON), &jsonData)
	if err == nil {
		// Extract known fields
		if name, ok := jsonData["name"].(string); ok {
			stats.Name = name
		}
		
		if roundStr, ok := jsonData["round_number"].(string); ok {
			stats.RoundNumber, _ = strconv.Atoi(roundStr)
		}
		
		if scoreT, ok := jsonData["score_t"].(string); ok {
			stats.ScoreT, _ = strconv.Atoi(scoreT)
		}
		
		if scoreCT, ok := jsonData["score_ct"].(string); ok {
			stats.ScoreCT, _ = strconv.Atoi(scoreCT)
		}
		
		if mapName, ok := jsonData["map"].(string); ok {
			stats.Map = mapName
		}
		
		if server, ok := jsonData["server"].(string); ok {
			stats.Server = server
		}
		
		if fields, ok := jsonData["fields"].(string); ok {
			// Split fields by comma and trim
			fieldList := strings.Split(fields, ",")
			for i := range fieldList {
				fieldList[i] = strings.TrimSpace(fieldList[i])
			}
			stats.Fields = fieldList
		}
		
		// Parse players data
		if players, ok := jsonData["players"].(map[string]interface{}); ok {
			stats.Players = make(map[string]PlayerStatistics)
			
			for playerID, playerDataStr := range players {
				if dataStr, ok := playerDataStr.(string); ok {
					// Parse the comma-separated values
					values := strings.Split(dataStr, ",")
					playerStats := parsePlayerStatistics(values, stats.Fields)
					stats.Players[playerID] = playerStats
				}
			}
		}
	}

	return stats
}

// parsePlayerStatistics parses player statistics from comma-separated values
func parsePlayerStatistics(values []string, fields []string) PlayerStatistics {
	stats := PlayerStatistics{}
	
	// Trim all values
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	
	// Map known fields by position
	// Expected order: accountid, team, money, kills, deaths, assists, dmg, hsp, kdr, adr, mvp, ef, ud, 3k, 4k, 5k, clutchk, firstk, pistolk, sniperk, blindk, bombk, firedmg, uniquek, dinks, chickenk
	if len(values) >= 26 {
		stats.AccountID, _ = strconv.Atoi(values[0])
		stats.Team, _ = strconv.Atoi(values[1])
		stats.Money, _ = strconv.Atoi(values[2])
		stats.Kills, _ = strconv.Atoi(values[3])
		stats.Deaths, _ = strconv.Atoi(values[4])
		stats.Assists, _ = strconv.Atoi(values[5])
		stats.Damage, _ = strconv.Atoi(values[6])
		stats.HeadshotPct, _ = strconv.ParseFloat(values[7], 64)
		stats.KDR, _ = strconv.ParseFloat(values[8], 64)
		stats.ADR, _ = strconv.Atoi(values[9])
		stats.MVP, _ = strconv.Atoi(values[10])
		stats.EnemiesFlashed, _ = strconv.Atoi(values[11])
		stats.UtilityDamage, _ = strconv.Atoi(values[12])
		stats.TripleKills, _ = strconv.Atoi(values[13])
		stats.QuadKills, _ = strconv.Atoi(values[14])
		stats.AceKills, _ = strconv.Atoi(values[15])
		stats.ClutchKills, _ = strconv.Atoi(values[16])
		stats.FirstKills, _ = strconv.Atoi(values[17])
		stats.PistolKills, _ = strconv.Atoi(values[18])
		stats.SniperKills, _ = strconv.Atoi(values[19])
		stats.BlindKills, _ = strconv.Atoi(values[20])
		stats.BombKills, _ = strconv.Atoi(values[21])
		stats.FireDamage, _ = strconv.Atoi(values[22])
		stats.UniqueKills, _ = strconv.Atoi(values[23])
		stats.Dinks, _ = strconv.Atoi(values[24])
		stats.ChickenKills, _ = strconv.Atoi(values[25])
	}
	
	return stats
}