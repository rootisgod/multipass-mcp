package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterExecTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("multipass_run_script",
			mcp.WithDescription("Run a multi-line script inside a Multipass instance. The script is transferred to the instance, executed with the specified interpreter, and cleaned up automatically. Use this for complex multi-step operations instead of chaining exec_command calls."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithString("script", mcp.Required(), mcp.Description("The script content to execute (multi-line string).")),
			mcp.WithString("interpreter", mcp.Description("Interpreter to run the script with (default \"bash\"). Examples: \"bash\", \"python3\", \"sh\".")),
			mcp.WithString("working_directory", mcp.Description("Working directory inside the instance.")),
			mcp.WithNumber("timeout", mcp.Description("Timeout in seconds (default 300). Increase for long-running scripts.")),
			mcp.WithTitleAnnotation("Multipass: Run Script"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleRunScript,
	)

	s.AddTool(
		mcp.NewTool("multipass_exec_command",
			mcp.WithDescription("Execute a command inside a Multipass instance."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Instance name.")),
			mcp.WithArray("command",
				mcp.Required(),
				mcp.Description("Command and arguments as a list (e.g. [\"ls\", \"-la\", \"/tmp\"])."),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithString("working_directory", mcp.Description("Working directory inside the instance.")),
			mcp.WithTitleAnnotation("Multipass: Execute Command"),
			mcp.WithReadOnlyHintAnnotation(false),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
		),
		handleExecCommand,
	)
}

func handleRunScript(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	script, err := req.RequireString("script")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	interpreter := req.GetString("interpreter", "bash")

	// Write script to a temp file on host
	tmpFile, err := os.CreateTemp("", "mcp-script-*")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create temp file: %v", err)), nil
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(script); err != nil {
		tmpFile.Close()
		return mcp.NewToolResultError(fmt.Sprintf("failed to write script: %v", err)), nil
	}
	tmpFile.Close()

	timeout := defaultTimeout
	if t := req.GetInt("timeout", 0); t > 0 {
		timeout = time.Duration(t) * time.Second
	}

	// Transfer to instance
	remotePath := fmt.Sprintf("/tmp/%s", filepath.Base(tmpPath))
	dest := fmt.Sprintf("%s:%s", name, remotePath)
	if _, err := runMultipass(ctx, defaultTimeout, "transfer", tmpPath, dest); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to transfer script: %v", err)), nil
	}

	// Make executable and run
	execArgs := []string{"exec", name}
	if wd := req.GetString("working_directory", ""); wd != "" {
		execArgs = append(execArgs, "--working-directory", wd)
	}
	execArgs = append(execArgs, "--", interpreter, remotePath)

	result, execErr := runMultipass(ctx, timeout, execArgs...)

	// Clean up remote script (best effort, use background context in case parent was cancelled)
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cleanupCancel()
	runMultipass(cleanupCtx, 10*time.Second, "exec", name, "--", "rm", "-f", remotePath)

	if execErr != nil {
		return mcp.NewToolResultError(execErr.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleExecCommand(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	rawArgs := req.GetArguments()
	commandRaw, ok := rawArgs["command"]
	if !ok {
		return mcp.NewToolResultError("missing required parameter: command"), nil
	}

	commandSlice, ok := commandRaw.([]any)
	if !ok {
		return mcp.NewToolResultError("command must be a list of strings"), nil
	}
	if len(commandSlice) == 0 {
		return mcp.NewToolResultError("command must contain at least one element"), nil
	}

	cmdArgs := []string{"exec", name}

	if wd := req.GetString("working_directory", ""); wd != "" {
		cmdArgs = append(cmdArgs, "--working-directory", wd)
	}

	cmdArgs = append(cmdArgs, "--")
	for _, item := range commandSlice {
		s, ok := item.(string)
		if !ok {
			return mcp.NewToolResultError("each command element must be a string"), nil
		}
		cmdArgs = append(cmdArgs, s)
	}

	result, err := runMultipass(ctx, defaultTimeout, cmdArgs...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}
