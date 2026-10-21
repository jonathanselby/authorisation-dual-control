# Authorisation and Dual Control — Engineering Design Document

**Version:** 0.1

## 1. Purpose

A small deterministic policy engine for incident access, compartments, role permissions, delegated agents and maker-checker confirmation.

> **One access model. One item lifecycle. No author approves their own version unless policy explicitly permits the declared exception.**

## 2. Non-goals

No authentication, persistence, workflow scheduling, user scoring/monitoring, legal privilege determination, application item schema or policy-administration workflow.

## 3. Runtime model

```text
policy + trusted request/context
             |
             v
        ValidatePolicy
             |
             v
     resolve human principal
             |
             v
       incident scope
             |
             v
      compartment scope
             |
             v
        role matrix
             |
             v
 confirm? -> maker/checker + controlled exception
             |
        +----+----+
        v         v
      ALLOW      DENY
```

The CLI is a thin JSON-file edge around `internal/authz`.

## 4. Policy model

`Subject` is a human principal with one role, incident scope and compartments. `Delegation` binds one agent ID to one human ID. `Policy.AllowSinglePersonWindow` controls whether the exceptional maker-checker path exists at all.

Agent and human namespaces are deliberately disjoint. An agent ID cannot also be a `Subject`.

## 5. Policy validation

`ValidatePolicy` rejects:

- blank or duplicate subject IDs;
- roles outside Major/Minor/Observer;
- blank or duplicate incident/compartment entries;
- blank/self delegation;
- an agent ID that is also a human subject;
- delegation to an unknown human;
- duplicate or ambiguous delegation for an agent.

A policy error is not a governed denial; it is invalid configuration.

## 6. Principal resolution

If `ActorID` is a configured agent, `HumanID` is mandatory and must equal the agent's configured human. The evaluator then uses only the named human's `Subject` for all access decisions and returns both IDs in `Decision`.

This inheritance is intentional: V0.1 does not define a second role/scope object for an agent. The agent is an attributed delegate of the named human, not an independently authorised principal.

If an ordinary human supplies a different `HumanID`, the request is denied rather than allowing the human to manufacture agent context.

## 7. Incident and compartment gates

Incident access is granted by explicit assignment or `GlobalIncidentAccess`. This is checked before role permission.

A non-empty item compartment creates an additional gate; the human subject must hold that compartment.

## 8. Role permission matrix

| Action | Major | Minor | Observer |
|---|:---:|:---:|:---:|
| read | yes | yes | yes |
| comment | yes | yes | no |
| propose | yes | no | no |
| confirm | yes | no | no |

Unknown actions are malformed input errors, not policy denials.

## 9. Maker-checker

For `confirm`, `Item.AuthorID` is mandatory. A missing or whitespace-only author is malformed input because maker-checker cannot be evaluated without an authoritative maker identity.

After principal resolution, the resolved human principal is compared with `Item.AuthorID`.

If equal, self-approval is denied unless:

```text
policy.AllowSinglePersonWindow
AND request.SinglePersonWindow
AND strings.TrimSpace(request.SelfApprovalReason) != ""
```

A successful exception sets `Decision.SelfApproval=true`.

The policy flag and runtime flag intentionally come from different logical authorities. The consuming application must ensure the runtime flag is trusted workflow context rather than actor-controlled form data.

## 10. Decision contract

`Decision` returns `Allowed`, resolved `HumanID`, optional `AgentID`, `SelfApproval` and a deterministic denial `Reason`.

The component does not write an approval record; the caller owns persistence/audit.

## 11. CLI and exits

```text
authzcheck -policy <policy.json> -request <request.json>
```

`-version` prints `authzcheck 0.1.0`.

```text
0 allowed
1 deterministic policy denial
2 invalid CLI usage
3 malformed policy/request, file I/O or execution error
```

## 12. State and persistence

No persistence, cache, database or global mutable state. Policy, request, item and exceptional runtime context are explicit inputs.

## 13. Test strategy

Named tests cover incident/compartment gates, role matrix, required confirmation author identity, policy-controlled maker-checker exception, agent principal rules, policy validation, malformed request and unknown actions.

CI additionally enforces `gofmt`, `go vet`, race tests, native build, exact version, Linux ARM64 build and Bafflewick acceptance.

## 14. Security boundary

The engine assumes actor IDs, human IDs, item author IDs and the runtime single-person-window signal are authenticated/authoritative inputs from the caller. It does not itself authenticate them.

For confirmation, absence of the authoritative item author fails closed as malformed input rather than being interpreted as “not self”.

The durable policy permission prevents an actor-controlled request flag from being sufficient to enable self-approval.

## 15. Operational behaviour

One Go executable; standard library only; deterministic and side-effect free core; structured JSON result; no hidden telemetry or background service.

ARM64 evidence is buildability, not a constrained-device performance claim.

## 16. Known limits

- fixed V0.1 role vocabulary;
- no field/action ABAC beyond the explicit model;
- no agent-specific scope independent of the named human principal;
- no policy persistence/versioning/approval;
- no approval ledger;
- no scheduler deciding when a single-person window is active;
- no proof that an agent or human ID was authenticated;
- no legal interpretation of a compartment label.

## 17. Acceptance-to-evidence map

| Requirement | Evidence |
|---|---|
| incident + compartment scope | `TestCompartmentAndIncidentScope` |
| role matrix | `TestRoles` |
| confirm requires authoritative author | `TestConfirmRequiresAuthorIdentity` |
| maker-checker + policy-gated exception | `TestMakerCheckerRequiresPolicyControlledException` |
| named-human delegation | `TestAgentInheritsNamedHumanOnly` |
| invalid/ambiguous policy | `TestValidatePolicy` |
| malformed request/action | `TestMalformedRequestAndUnknownActionAreErrors` |
| format/vet/race/build/version/ARM64/Bafflewick | `.github/workflows/test.yml` |
