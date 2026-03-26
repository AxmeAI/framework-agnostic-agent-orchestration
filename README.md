# Framework-Agnostic Agent Orchestration

Orchestrate LangGraph, CrewAI, AutoGen, and raw Python agents in one workflow. No framework lock-in. One ScenarioBundle JSON.

You picked a framework. LangGraph, CrewAI, AutoGen - each one has its own orchestration layer, its own agent lifecycle, its own way of routing tasks. Now you need two of them to work together. Or you need to swap one out. Or you need a raw Python script to participate in the same workflow. Every framework says "use our orchestrator." None of them interoperate.

> **Alpha** - Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) - [hello@axme.ai](mailto:hello@axme.ai)

---

## The Problem

Each agent framework has its own orchestration. None of them talk to each other.

- **LangGraph** orchestrates LangGraph agents. Cannot route to CrewAI.
- **CrewAI** orchestrates CrewAI crews. Cannot route to AutoGen.
- **AutoGen** orchestrates AutoGen agents. Cannot route to LangGraph.
- **Raw Python** has no orchestration at all.

Want a LangGraph agent to hand off to a CrewAI crew? Build a custom bridge. Want to swap AutoGen for raw Python? Rewrite the orchestration. Want all four in one workflow? Good luck.

### Framework-Locked Orchestration (one framework at a time)

```python
# LangGraph - locked to LangGraph agents
from langgraph.graph import StateGraph
graph = StateGraph(AgentState)
graph.add_node("process", langgraph_agent)   # LangGraph only
graph.add_node("review", langgraph_agent_2)  # LangGraph only
# Cannot add a CrewAI crew or AutoGen agent here

# CrewAI - locked to CrewAI agents
from crewai import Crew, Task
crew = Crew(agents=[crewai_agent], tasks=[task])  # CrewAI only
# Cannot add a LangGraph node or raw Python here

# AutoGen - locked to AutoGen agents
from autogen import GroupChat
chat = GroupChat(agents=[autogen_agent])  # AutoGen only
# Cannot add a CrewAI crew or LangGraph node here
```

### AXME (framework-agnostic, 4 lines)

```python
intent_id = client.send_intent({
    "intent_type": "intent.orchestration.generic_task.v1",
    "to_agent": "agent://myorg/production/generic-processor",
    "payload": {"task_id": "ORCH-2026-0015", "task_type": "document_processing"},
})
result = client.wait_for(intent_id)
```

The agent behind `generic-processor` can be LangGraph, CrewAI, AutoGen, or raw Python. The caller does not know and does not care. Swap frameworks without changing a single line of orchestration code.

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
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
})

print(f"Submitted: {intent_id}")
result = client.wait_for(intent_id)
print(f"Done: {result['status']}")
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "intent.orchestration.generic_task.v1",
  toAgent: "agent://myorg/production/generic-processor",
  payload: {
    taskId: "ORCH-2026-0015",
    taskType: "document_processing",
    input: {
      documentUrl: "s3://docs/contract-v3.pdf",
      actions: ["extract_clauses", "check_compliance", "generate_summary"],
    },
    requestedBy: "legal-pipeline",
  },
});

console.log(`Submitted: ${intentId}`);
const result = await client.waitFor(intentId);
console.log(`Done: ${result.status}`);
```

### Go

```bash
go get github.com/AxmeAI/axme-sdk-go
```

```go
client, _ := axme.NewClient(axme.ClientConfig{APIKey: os.Getenv("AXME_API_KEY")})

intentID, _ := client.SendIntent(ctx, map[string]any{
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

result, _ := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
fmt.Printf("Final status: %v\n", result["status"])
```

---

## Before / After

### Before: Framework Lock-In

```
+-----------+     +------------+     +-----------+
| LangGraph |---->| LangGraph  |---->| LangGraph |
|  Agent A  |     | Orchestr.  |     |  Agent B  |
+-----------+     +------------+     +-----------+
                       X
                       X  (cannot cross)
                       X
+-----------+     +------------+     +-----------+
|  CrewAI   |---->|   CrewAI   |---->|  CrewAI   |
|  Agent C  |     | Orchestr.  |     |  Agent D  |
+-----------+     +------------+     +-----------+
```

Each framework is a silo. Agents cannot participate in workflows outside their framework.

### After: Framework-Agnostic with AXME

```
+-----------+                        +-----------+
| LangGraph |--+                 +-->| LangGraph |
|  Agent A  |  |  +-----------+  |   |  Agent B  |
+-----------+  +->|           |--+   +-----------+
                  | AXME Cloud|
+-----------+  +->| (intents) |--+   +-----------+
|  CrewAI   |--+  |           |  +-->| AutoGen   |
|  Agent C  |     +-----------+      |  Agent E  |
+-----------+          |             +-----------+
                       v
               +-----------+
               | Raw Python|
               |  Agent F  |
               +-----------+
```

One orchestration layer. Any framework. Swap agents without changing callers.

---

## Works With

| Framework | How It Connects | Lock-In |
|-----------|----------------|---------|
| **LangGraph** | Wrap graph invocation in an AXME agent listener | None - LangGraph runs inside, AXME orchestrates outside |
| **CrewAI** | Wrap crew.kickoff() in an AXME agent listener | None - CrewAI runs inside, AXME orchestrates outside |
| **AutoGen** | Wrap group chat in an AXME agent listener | None - AutoGen runs inside, AXME orchestrates outside |
| **Raw Python** | Use AXME SDK directly | None - no framework required |
| **Custom / Internal** | Any code that can make HTTP calls | None - language-agnostic protocol |

The pattern is always the same: your framework runs **inside** the agent. AXME runs **outside** as the orchestration layer. The agent receives an intent, does its work using whatever framework it wants, and resumes with results.

---

## More Languages

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |

---

## Run the Full Example

### Prerequisites

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
# Open a new terminal, or run the "source" command shown by the installer

# Log in
axme login

# Install Python SDK
pip install axme
```

### Terminal 1 - submit the scenario

```bash
axme scenarios apply scenario.json
# Note the intent_id in the output
```

### Terminal 2 - start the agent

Get the agent key after scenario apply:

```bash
# macOS
cat ~/Library/Application\ Support/axme/scenario-agents.json | grep -A2 generic-processor-demo

# Linux
cat ~/.config/axme/scenario-agents.json | grep -A2 generic-processor-demo
```

Run in your language of choice:

```bash
# Python
AXME_API_KEY=<agent-key> python agent.py

# TypeScript (requires Node 20+)
cd typescript && npm install
AXME_API_KEY=<agent-key> npx tsx agent.ts

# Go
cd go && go run ./cmd/agent/
```

### Verify

```bash
axme intents get <intent_id>
# lifecycle_status: COMPLETED
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) - project overview
- [AXP Spec](https://github.com/AxmeAI/axme-spec) - open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) - 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) - manage intents, agents, scenarios from the terminal
- [Durable Execution with Human Approval](https://github.com/AxmeAI/durable-execution-with-human-approval) - durable workflows with human gates

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
