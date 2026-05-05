package validator

import (
	"fmt"
	"strings"
)

// CommandValidator handles security checks for system commands.
type CommandValidator struct {
	DangerousMode bool
}

var forbiddenCommands = []string{
	"sudo",
	"rm",
	"mv",
	"dd",
	"mkfs",
	"fdisk",
	"parted",
	"chmod",
	"chown",
	"reboot",
	"shutdown",
	"iptables",
	"ufw",
}

// Validate checks if a command string is safe to execute.
func (v *CommandValidator) Validate(command string) error {
	if v.DangerousMode {
		return nil
	}

	parts := strings.Fields(strings.ToLower(command))
	if len(parts) == 0 {
		return nil
	}

	for _, forbidden := range forbiddenCommands {
		if parts[0] == forbidden {
			return fmt.Errorf("command '%s' is restricted. use --dangerous-mode to override", forbidden)
		}
	}

	// Check for dangerous pipes or redirections
	dangerousTokens := []string{">", ">>", "|"}
	for _, part := range parts {
		for _, token := range dangerousTokens {
			if part == token {
				return fmt.Errorf("redirections and pipes are restricted. use --dangerous-mode to override")
			}
		}
	}

	return nil
}
