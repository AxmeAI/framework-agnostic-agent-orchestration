// Framework-agnostic agent orchestration - Go example.
//
// Submit a document processing task. Any agent - LangGraph, CrewAI,
// AutoGen, or raw code - can handle it. No framework lock-in.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type":  "intent.orchestration.generic_task.v1",
		"to_agent":     "agent://myorg/production/generic-processor",
		"task_id":      "ORCH-2026-0015",
		"task_type":    "document_processing",
		"requested_by": "legal-pipeline",
		"input": map[string]any{
			"document_url": "s3://docs/contract-v3.pdf",
			"actions":      []string{"extract_clauses", "check_compliance", "generate_summary"},
		},
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Intent submitted: %s\n", intentID)

	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %v\n", result["status"])
}
