// Package logpath provides helpers for resolving log file paths.
package logpath

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var logPrefixRe = regexp.MustCompile(`^(\d{3})_`)

// ResolveRecordPath returns a log file path for recording.
func ResolveRecordPath(path string, listenAddr string, target *url.URL) (string, error) {
	// Returns a log file path. If path is a directory or has no extension,
	// it creates a new numbered log file name under that directory.
	// Used for rwnd proxy log file creation
	if isDirPath(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
		next, err := nextLogNumber(path)
		if err != nil {
			return "", err
		}
		name := buildLogFilename(next, listenAddr, target)
		return filepath.Join(path, name), nil
	}

	return path, nil
}

// ResolveReplayPath returns the log file path to replay.
func ResolveReplayPath(path string) (string, error) {
	// Returns a log file path. If path is a directory or has no extension,
	// it selects the highest numbered log file from that directory.
	// Used for rwnd replay log file selection
	if isDirPath(path) {
		latest, err := latestLogFile(path)
		if err != nil {
			return "", err
		}
		return filepath.Join(path, latest), nil
	}

	return path, nil
}

func buildLogFilename(seq int, listenAddr string, target *url.URL) string {
	// Assembles the log filename with sequence, time, and metadata.
	stamp := time.Now().UTC().Format("20060102T150405Z")
	listen := sanitizeFilenamePart(listenAddr)
	targetStr := ""
	if target != nil {
		targetStr = sanitizeFilenamePart(target.Host)
		if targetStr == "" {
			targetStr = sanitizeFilenamePart(target.String())
		}
	}
	if targetStr != "" {
		return fmt.Sprintf("%03d_%s_listen-%s_target-%s.jsonl", seq, stamp, listen, targetStr)
	}
	return fmt.Sprintf("%03d_%s_listen-%s.jsonl", seq, stamp, listen)
}

func isDirPath(path string) bool {
	// Determines if the path should be treated as a directory.
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return true
	}
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return filepath.Ext(path) == ""
}

func nextLogNumber(dir string) (int, error) {
	// Returns the next sequential log number for a directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	max := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		m := logPrefixRe.FindStringSubmatch(entry.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}

	return max + 1, nil
}

func latestLogFile(dir string) (string, error) {
	// Returns the highest-numbered log file in a directory.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	max := -1
	name := ""
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		m := logPrefixRe.FindStringSubmatch(entry.Name())
		if len(m) != 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > max {
			max = n
			name = entry.Name()
		}
	}

	if name == "" {
		return "", fmt.Errorf("No log files found in %s", dir)
	}
	return name, nil
}

func sanitizeFilenamePart(value string) string {
	// Converts a value into a safe filename segment.
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(value))
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	return out
}
