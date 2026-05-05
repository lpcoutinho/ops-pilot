package tools

import (
	"context"
	"encoding/json"
)

// SystemTool defines the interface for all diagnostic and administrative tools.
type SystemTool interface {
	// Name returns the unique identifier of the tool.
	Name() string
	
	// Description returns a natural language description of what the tool does.
	Description() string
	
	// Schema returns the JSON schema representing the tool's input parameters.
	Schema() map[string]interface{}
	
	// Execute performs the tool's action with the given arguments.
	Execute(ctx context.Context, args json.RawMessage) (interface{}, error)
}
