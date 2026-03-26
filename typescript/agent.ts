/**
 * Generic processor agent - TypeScript example.
 *
 * Listens for intents via SSE, processes document tasks,
 * resumes with structured results.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "generic-processor-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const taskId = payload.task_id ?? "unknown";
  const taskType = payload.task_type ?? "unknown";
  const input = payload.input ?? {};
  const documentUrl = input.document_url ?? "unknown";
  const actions = input.actions ?? [];

  console.log(`  Task ID: ${taskId}`);
  console.log(`  Task type: ${taskType}`);
  console.log(`  Document: ${documentUrl}`);
  console.log(`  Actions: ${actions.join(", ")}`);

  console.log(`  Processing document...`);
  await new Promise((r) => setTimeout(r, 1000));

  console.log(`  Extracting clauses...`);
  await new Promise((r) => setTimeout(r, 1000));

  console.log(`  Checking compliance...`);
  await new Promise((r) => setTimeout(r, 1000));

  console.log(`  Generating summary...`);
  await new Promise((r) => setTimeout(r, 1000));

  const result = {
    action: "complete",
    task_id: taskId,
    results: {
      clauses_extracted: 24,
      compliance_status: "passed",
      summary:
        "Standard SaaS agreement with mutual indemnification, 12-month term, auto-renewal clause, and standard limitation of liability.",
    },
    processed_at: new Date().toISOString(),
  };

  await client.resumeIntent(intentId, result, { ownerAgent: AGENT_ADDRESS });
  console.log(
    `  Processing complete. Clauses: ${result.results.clauses_extracted}, Compliance: ${result.results.compliance_status}`
  );
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;

    if (!intentId) continue;

    if (["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error processing intent: ${e}`);
      }
    }
  }
}

main().catch(console.error);
