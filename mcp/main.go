package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Request struct {
}

func main() {
	s := server.NewMCPServer(
		"Quiz Generator",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	quizTool := mcp.NewTool("JSON Quiz Generator",
		mcp.WithDescription("Generate a quiz based on user input"),
		mcp.WithString("quiz_json_payload",
			mcp.Required(),
			mcp.Description(`A JSON string representing the full quiz content. The JSON object MUST conform to the following structure:
            {
              "title": "string",
              "description": "string",
              "questions": [
                {
                  "question_text": "string",
                  "type": "string",
                  "order_num": "integer",
                  "answers": [ { "answer_text": "string", "is_correct": "boolean" } ]
                }
              ]
      }`),
		),
	)

	// This is the handler function for your tool
	s.AddTool(quizTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 1. Get the value of the "quiz_json_payload" parameter from the request
		// The mcp.CallToolRequest object contains the parameters sent by Claude Desktop.
		quizJSONPayload, err := request.RequireString("quiz_json_payload")
		if err != nil {
			// If the required parameter is missing, log the error and return an error result to Claude.
			log.Printf("Error: Required parameter 'quiz_json_payload' is missing. %v", err)
			// The second argument to NewToolResultError can be an error code string, or empty.
			return mcp.NewToolResultError("Missing required parameter: quiz_json_payload. Details: " + err.Error()), nil
		}

		// 2. Log the received JSON payload
		// This will print the JSON string (that Claude Desktop sent) to your Go application's console.
		log.Printf("--- Received quiz_json_payload from Claude Desktop ---")
		log.Printf("\n%s\n", quizJSONPayload)
		log.Printf("--- End of quiz_json_payload ---")

		// For more direct output if your MCP server's stdout is what you're monitoring:
		// fmt.Printf("Received quiz_json_payload from Claude Desktop:\n%s\n", quizJSONPayload)

		// 3. For now, just acknowledge receipt and return the payload (or a success message).
		// In a real application, you would:
		//    a. Unmarshal quizJSONPayload into your Go structs (like GeneratedQuizData).
		//    b. Validate the data.
		//    c. Perform any transformations.
		//    d. Marshal the (potentially transformed) data back to a JSON string.
		//    e. Return that final JSON string using mcp.NewToolResultText().
		//
		// For this step, we'll just echo back the received payload.
		// This confirms to Claude Desktop that the tool executed and what it received.
		responseMessage := fmt.Sprintf("Successfully received and logged quiz_json_payload. Content received:\n%s", quizJSONPayload)
		return mcp.NewToolResultText(responseMessage), nil
		// Alternatively, just echo the payload if that's more useful for Claude's side:
		// return mcp.NewToolResultText(quizJSONPayload), nil
	})

	// You need to start the server for it to listen for tool calls from Claude Desktop
	log.Println("MCP Server starting and listening on stdio...")
	if err := server.ServeStdio(s); err != nil { // Or server.Serve("tcp", "localhost:PORT") for network
		log.Fatalf("Error starting MCP server: %v\n", err)
	}

}
