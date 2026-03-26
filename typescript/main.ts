/**
 * Framework-agnostic agent orchestration - TypeScript example.
 *
 * Submit a document processing task. Any agent - LangGraph, CrewAI,
 * AutoGen, or raw code - can handle it. No framework lock-in.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
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
  console.log(`Intent submitted: ${intentId}`);

  const result = await client.waitFor(intentId);
  console.log(`Final status: ${result.status}`);
}

main().catch(console.error);
