"""
Generic processor agent - processes tasks from any framework orchestration.

Listens for intents via SSE. Simulates document processing: clause extraction,
compliance checking, and summary generation. Resumes with structured results.

Usage:
    export AXME_API_KEY="<agent-key>"
    python agent.py
"""

import os
import sys
import time

sys.stdout.reconfigure(line_buffering=True)

from axme import AxmeClient, AxmeClientConfig


AGENT_ADDRESS = "generic-processor-demo"


def handle_intent(client, intent_id):
    """Process a document task and resume with results."""
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    task_id = payload.get("task_id", "unknown")
    task_type = payload.get("task_type", "unknown")
    input_data = payload.get("input", {})
    document_url = input_data.get("document_url", "unknown")
    actions = input_data.get("actions", [])

    print(f"  Task ID: {task_id}")
    print(f"  Task type: {task_type}")
    print(f"  Document: {document_url}")
    print(f"  Actions: {', '.join(actions)}")

    print(f"  Processing document...")
    time.sleep(1)

    print(f"  Extracting clauses...")
    time.sleep(1)

    print(f"  Checking compliance...")
    time.sleep(1)

    print(f"  Generating summary...")
    time.sleep(1)

    result = {
        "action": "complete",
        "task_id": task_id,
        "results": {
            "clauses_extracted": 24,
            "compliance_status": "passed",
            "summary": "Standard SaaS agreement with mutual indemnification, 12-month term, auto-renewal clause, and standard limitation of liability.",
        },
        "processed_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }

    client.resume_intent(intent_id, result)
    print(f"  Processing complete. Clauses: {result['results']['clauses_extracted']}, Compliance: {result['results']['compliance_status']}")


def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set.")
        print("Run the scenario first: axme scenarios apply scenario.json")
        print("Then get the agent key from ~/.config/axme/scenario-agents.json")
        sys.exit(1)

    client = AxmeClient(AxmeClientConfig(api_key=api_key))

    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")

    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")

        if not intent_id:
            continue

        if status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error processing intent: {e}")


if __name__ == "__main__":
    main()
