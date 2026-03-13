package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultTimeout = 300 * time.Second
	launchTimeout  = 600 * time.Second
)

// runMultipass executes a multipass CLI command and returns stdout.
// Returns an error with stderr content on non-zero exit.
func runMultipass(ctx context.Context, timeout time.Duration, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("multipass: no command specified")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "multipass", args...)
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("multipass %s timed out after %v", args[0], timeout)
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("multipass %s failed (exit %d): %s",
				strings.Join(args, " "), exitErr.ExitCode(), strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("multipass %s: %w", args[0], err)
	}
	return strings.TrimSpace(string(out)), nil
}

// runMultipassJSON executes a multipass CLI command with --format json and parses the result.
func runMultipassJSON(ctx context.Context, timeout time.Duration, args ...string) (json.RawMessage, error) {
	fullArgs := append(append([]string{}, args...), "--format", "json")
	result, err := runMultipass(ctx, timeout, fullArgs...)
	if err != nil {
		return nil, err
	}
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(result), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}
	return raw, nil
}
