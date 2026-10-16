# Bafflewick acceptance example

Bafflewick is fictional.

`major-1` has global incident access and the `legal` compartment. `agent-1` is delegated only to `major-1`. The example asks the agent to propose a legal-compartment item while naming that human principal.

```bash
go build -o ./bin/authzcheck ./cmd/authzcheck
./bin/authzcheck \
  -policy examples/bafflewick/policy.json \
  -request examples/bafflewick/request.json
```

Expected: `allowed=true`, `human_id=major-1`, `agent_id=agent-1`.

The package tests additionally prove that an agent cannot omit or switch its human, a human cannot spoof agent context, and self-approval cannot be enabled by request context unless durable policy also permits the exception.
