# Authorisation and Dual Control — Skill

## Use this when

Use this component when an application needs one deterministic decision for role permissions, incident/compartment need-to-know, delegated agents and maker-checker confirmation.

## Core rule

> **One access model. One item lifecycle. No author approves their own version unless policy explicitly permits the declared exception.**

## Read order

1. `PRD.md` — boundary and acceptance criteria.
2. `EDD.md` — principal, maker-checker and failure semantics.
3. `architecture.txt` — compact decision order.
4. `examples/bafflewick/` — executable fictional case.

## Change discipline

1. Validate policy before evaluating a request.
2. Keep agent IDs structurally separate from human subjects.
3. Never evaluate an agent from an agent-owned role; always resolve the named human.
4. Check incident and compartment scope before the action matrix.
5. Never allow a request flag by itself to enable self-approval.
6. Keep deterministic denials separate from malformed-input/configuration errors.
7. Do not add authentication, persistence, scheduling or user scoring here.
8. Run format, vet, race, native and ARM64 gates after changes.
9. Regenerate implementation-derived evidence after code changes.
10. Do not tag until `RELEASE.md` is satisfied.
