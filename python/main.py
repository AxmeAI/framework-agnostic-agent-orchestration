"""
Framework-agnostic agent orchestration - Python example.

Submit a document processing task. Any agent - LangGraph, CrewAI, AutoGen,
or raw Python - can pick it up. No framework lock-in.

Usage:
    pip install axme
    export AXME_API_KEY="your-key"
    python main.py
"""

import os
from axme import AxmeClient, AxmeClientConfig


def main():
    client = AxmeClient(
        AxmeClientConfig(api_key=os.environ["AXME_API_KEY"])
    )

    # Submit a document processing task - any framework can handle it
    intent_id = client.send_intent(
        {
            "intent_type": "intent.orchestration.generic_task.v1",
            "to_agent": "agent://myorg/production/generic-processor",
            "payload": {
                "task_id": "ORCH-2026-0015",
                "task_type": "document_processing",
                "input": {
                    "document_url": "s3://docs/contract-v3.pdf",
                    "actions": ["extract_clauses", "check_compliance", "generate_summary"],
                },
                "requested_by": "legal-pipeline",
            },
        }
    )
    print(f"Intent submitted: {intent_id}")

    # Wait for completion
    print("Watching lifecycle...")
    for event in client.observe(intent_id):
        status = event.get("status", "")
        print(f"  [{status}] {event.get('event_type', '')}")
        if status in ("COMPLETED", "FAILED", "TIMED_OUT", "CANCELLED"):
            break

    intent = client.get_intent(intent_id)
    print(f"\nFinal status: {intent['intent']['lifecycle_status']}")


if __name__ == "__main__":
    main()
