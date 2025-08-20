package cs2log

import (
	"regexp"
)

// CombinedPatterns merges DefaultPatterns with ExtendedPatterns
// Returns patterns in priority order (specific patterns before general ones)
func CombinedPatterns() map[*regexp.Regexp]MessageFunc {
	// We need to use a slice to maintain order since maps don't preserve order
	// But for now, we'll return a map and rely on ParseWithOrderedPatterns
	combined := make(map[*regexp.Regexp]MessageFunc)
	
	// First add default patterns
	for pattern, fn := range DefaultPatterns {
		combined[pattern] = fn
	}
	
	// Then add extended patterns (will override defaults if same pattern)
	for pattern, fn := range ExtendedPatterns {
		combined[pattern] = fn
	}
	
	return combined
}

// ParseEnhanced parses a log line using both default and extended patterns
// This is the main function to use for parsing CS2 logs with custom events
func ParseEnhanced(line string) (Message, error) {
	// Use ordered parsing to ensure correct pattern priority
	return ParseOrdered(line)
}

// ParseExtendedOnly parses using only the extended patterns
// Useful for testing or when you only want custom events
func ParseExtendedOnly(line string) (Message, error) {
	return ParseWithPatterns(line, ExtendedPatterns)
}