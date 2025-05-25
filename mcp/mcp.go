package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Quiz Generator",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	quizTool := mcp.NewTool("quiz-generator/create",
		mcp.WithDescription("A simple calculator tool"),
	)

	s.AddTool(quizTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, fmt.Errorf("this is an error from the quiz generator tool")
	})
}
