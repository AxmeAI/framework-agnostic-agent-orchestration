// Generic processor agent - Go example.
//
// Listens for intents via SSE, processes document tasks,
// resumes with structured results.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run ./cmd/agent/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "generic-processor-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	taskID, _ := payload["task_id"].(string)
	if taskID == "" {
		taskID = "unknown"
	}
	taskType, _ := payload["task_type"].(string)
	if taskType == "" {
		taskType = "unknown"
	}
	input, _ := payload["input"].(map[string]any)
	if input == nil {
		input = map[string]any{}
	}
	documentURL, _ := input["document_url"].(string)
	if documentURL == "" {
		documentURL = "unknown"
	}
	actionsRaw, _ := input["actions"].([]any)
	actions := make([]string, 0, len(actionsRaw))
	for _, a := range actionsRaw {
		if s, ok := a.(string); ok {
			actions = append(actions, s)
		}
	}

	fmt.Printf("  Task ID: %s\n", taskID)
	fmt.Printf("  Task type: %s\n", taskType)
	fmt.Printf("  Document: %s\n", documentURL)
	fmt.Printf("  Actions: %s\n", strings.Join(actions, ", "))

	fmt.Println("  Processing document...")
	time.Sleep(1 * time.Second)

	fmt.Println("  Extracting clauses...")
	time.Sleep(1 * time.Second)

	fmt.Println("  Checking compliance...")
	time.Sleep(1 * time.Second)

	fmt.Println("  Generating summary...")
	time.Sleep(1 * time.Second)

	result := map[string]any{
		"action":  "complete",
		"task_id": taskID,
		"results": map[string]any{
			"clauses_extracted":  24,
			"compliance_status":  "passed",
			"summary":           "Standard SaaS agreement with mutual indemnification, 12-month term, auto-renewal clause, and standard limitation of liability.",
		},
		"processed_at": time.Now().UTC().Format(time.RFC3339),
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}
	fmt.Printf("  Processing complete. Clauses: 24, Compliance: passed\n")
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)

		if intentID == "" {
			continue
		}

		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error processing intent: %v\n", err)
			}
		}
	}
}
