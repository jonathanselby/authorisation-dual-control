# Authorisation and Dual Control — Intended vs Recovered Design

**Implementation evidence:** `938cbad9760bd23211f8a2a487b3ccf71900a00f`

| Area | Result | Recovered evidence |
|---|---|---|
| One deterministic access model | MATCH | One `Evaluate` path owns incident, compartment, role, delegation and maker-checker. |
| Go 1.26 / standard library | MATCH | `go.mod` declares Go 1.26 and no third-party modules. |
| Valid-policy gate | MATCH | `ValidatePolicy` runs before request evaluation and rejects invalid/ambiguous principals, roles, scopes and delegations. |
| Human/agent separation | MATCH | Agent ID cannot also be a `Subject`; a delegated agent must name its configured human. |
| Agent inherits human only | MATCH | Evaluation uses the resolved human's subject; agent has no independent role/scope. |
| Incident need-to-know | MATCH | Assignment/global access is checked before role matrix. |
| Compartment need-to-know | MATCH | Non-empty item compartment requires matching human clearance. |
| Major/Minor/Observer matrix | MATCH | Read/comment/propose/confirm behaviour matches PRD/EDD. |
| Maker-checker | MATCH | Resolved human cannot confirm their own authored item under normal operation. |
| Single-person exception cannot be request-only | MATCH | Requires policy permission + active runtime window + non-blank reason. |
| Self-approval visibility | MATCH | Successful exception sets `SelfApproval=true`. |
| Governed denial vs invalid input | MATCH | Permission/refusal paths return decisions; invalid policy/malformed request/action return errors. |
| Stable CLI and version | MATCH | `authzcheck` maps allowed/denied/usage/error to 0/1/2/3 and exposes 0.1.0. |
| ARM64 evidence | MATCH | CI cross-builds `linux/arm64`. |
| Authentication | DELIBERATE GAP | IDs are trusted caller inputs; no authentication is implemented. |
| Runtime-window scheduling | DELIBERATE GAP | Policy permits the exception, but another trusted component decides whether the window is active. |
| Decision/policy persistence | DELIBERATE GAP | No database or audit store exists here. |
| General ABAC / policy language | DELIBERATE GAP | V0.1 intentionally uses a fixed role/action/scope model. |
| Trustworthiness of supplied runtime context | UNKNOWN | Cannot be proven within this stateless evaluator. |
| Legal meaning of compartments | UNKNOWN | Strings are enforced but not interpreted. |

## Acceptance traceability

The source contains named tests for incident/compartment scope, role permissions, controlled self-approval, delegated-agent attribution, policy validation and malformed request/action handling. CI adds format, vet, race, native/version, ARM64 and Bafflewick checks.

## Material correction made before publication

The initial draft allowed `Request.SinglePersonWindow` plus a reason to enable author self-confirmation. The release-candidate implementation does **not** retain that design. It adds `Policy.AllowSinglePersonWindow`, so self-approval requires both durable policy permission and trusted runtime context.

The initial draft also did not validate the policy model. The release candidate validates roles, duplicate scopes, delegation targets/ambiguity and disjoint human/agent IDs before evaluating access.

These are security corrections, not scope expansion: the component still performs one authorisation decision and remains stateless.

## Conclusion

**MATCH for the hardened V0.1 design.** The recovered implementation supports the intended one-job boundary, with stronger controls around agent identity and the maker-checker exception than the original draft.
