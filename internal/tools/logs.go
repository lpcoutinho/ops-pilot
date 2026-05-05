package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type AnalyzeLogsTool struct{}

type LogsReport struct {
	FilePath string   `json:"file_path"`
	Lines    []string `json:"lines"`
	OS       string   `json:"os"`
}

func (t *AnalyzeLogsTool) Name() string {
	return "analyze_logs"
}

func (t *AnalyzeLogsTool) Description() string {
	return "Reads the last 50 lines of the system log file to help identify recent errors or system events."
}

func (t *AnalyzeLogsTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *AnalyzeLogsTool) Execute(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var logPath string
	switch runtime.GOOS {
	case "linux":
		logPath = "/var/log/syslog"
		// Some distros use /var/log/messages
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			logPath = "/var/log/messages"
		}
	case "darwin":
		logPath = "/var/log/system.log"
	default:
		return nil, fmt.Errorf("log analysis not supported on %s", runtime.GOOS)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w (try running with higher privileges or check file existence)", logPath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > 50 {
			lines = lines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return LogsReport{
		FilePath: logPath,
		Lines:    lines,
		OS:       runtime.GOOS,
	}, nil
}
