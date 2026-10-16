# Authorisation and Dual Control — Recovered Engineering Design

**Evidence commit:** `938cbad9760bd23211f8a2a487b3ccf71900a00f`  
**Recovery method:** source, tests, CLI and CI inspection  
**Classification:** KNOWN / INFERRED / UNKNOWN

## 1. Observed purpose

**KNOWN.** The repository contains one deterministic authorisation evaluator. `internal/authz.Evaluate` combines incident scope, compartment scope, role permission, delegated-agent attribution and maker-checker into one decision.

**KNOWN.** The package does not authenticate identities, persist policy/decisions, schedule workflows or score users.

## 2. Core data model

**KNOWN.** `Subject` represents a human principal with `ID`, `Role`, incident assignments, optional global incident access and compartments.

**KNOWN.** `Delegation` maps `AgentID` to one `HumanID`.

**KNOWN.** `Policy` contains subjects, delegations and `AllowSinglePersonWindow`.

**KNOWN.** `Request` supplies actor, optional named human, action, item, runtime single-person-window state and optional self-approval reason.

**KNOWN.** `Decision` exposes allow/deny, resolved human, optional agent, self-approval flag and reason.

## 3. Policy-validation gate

**KNOWN.** `Evaluate` calls `ValidatePolicy` before evaluating a request.

**KNOWN.** Validation rejects blank/duplicate subject IDs, roles outside Major/Minor/Observer, blank/duplicate incident or compartment entries, self/blank delegation, unknown delegation humans, duplicate/ambiguous agent delegation, and any agent ID that is also a human subject.

**INFERRED.** This creates disjoint human and agent namespaces and prevents an agent from acquiring a second authority path through its own subject record.

## 4. Human/agent resolution

**KNOWN.** A configured agent must supply `HumanID`; it must exactly match the human in the delegation. The evaluator then uses only that human's `Subject` for access checks.

**KNOWN.** A configured agent without `HumanID` is denied. An agent naming another human is denied. An ordinary human supplying another `HumanID` is denied.

**KNOWN.** Successful delegated decisions preserve both `HumanID` and `AgentID` in the result.

## 5. Incident scope

**KNOWN.** `incidentAllowed` grants scope when the resolved human has `GlobalIncidentAccess` or the item incident appears in `AssignedIncidents`.

**KNOWN.** Incident scope is evaluated before role permission.

## 6. Compartment scope

**KNOWN.** If `Item.Compartment` is non-empty, the resolved human must contain the same string in `Compartments`.

**KNOWN.** Failure is a deterministic denial with `compartment clearance required`.

**UNKNOWN.** The repository does not define the legal/business meaning, hierarchy or lifecycle of a compartment string.

## 7. Role matrix

**KNOWN.** Action permission is hard-coded as:

```text
             read  comment  propose  confirm
Major         yes     yes       yes      yes
Minor         yes     yes       no       no
Observer      yes     no        no       no
```

**KNOWN.** Unknown actions are malformed-input errors, not governed denials.

## 8. Maker-checker

**KNOWN.** Confirmation is available only to a Major.

**KNOWN.** The resolved human principal, not the raw agent ID, is compared with `Item.AuthorID`.

**KNOWN.** If the same human authored the item, confirmation is denied unless all of the following are true:

```text
Policy.AllowSinglePersonWindow == true
Request.SinglePersonWindow == true
strings.TrimSpace(Request.SelfApprovalReason) != ""
```

**KNOWN.** A successful exception sets `Decision.SelfApproval=true`.

**INFERRED.** Durable policy and runtime context act as independent gates. The source prevents the request's window flag alone from enabling self-approval.

**UNKNOWN.** The package cannot prove that the runtime window flag came from a trusted scheduler; the caller must enforce that trust boundary.

## 9. Request validity

**KNOWN.** Actor, action and incident ID are mandatory. Supplying `HumanID == ActorID` is treated as malformed input because `HumanID` is reserved for delegated-agent attribution.

**KNOWN.** Missing mandatory fields or an unknown action returns an error rather than a policy denial.

## 10. Decision/failure semantics

**KNOWN.** Unknown humans, delegation mismatch, missing incident access, missing compartment clearance, insufficient role and maker-checker refusal return structured `Allowed=false` decisions.

**KNOWN.** Invalid policy and malformed request/action return errors.

**INFERRED.** This cleanly separates expected governed refusal from broken policy/input.

## 11. CLI contract

**KNOWN.** `cmd/authzcheck` reads policy and request JSON files, invokes `Evaluate`, writes a JSON decision and uses:

```text
0 allowed
1 deterministic policy denial
2 CLI usage error
3 malformed policy/request, file I/O or execution error
```

**KNOWN.** `authzcheck -version` prints `authzcheck 0.1.0`.

## 12. State and side effects

**KNOWN.** The core has no package-global mutable state, persistence, network access, background routine or database/cache.

**KNOWN.** Policy and item state are caller-owned inputs.

**INFERRED.** Core evaluation is suitable for concurrent calls when caller-owned inputs are not concurrently mutated.

## 13. Dependencies and portability

**KNOWN.** `go.mod` declares Go 1.26 and no third-party modules.

**KNOWN.** CI builds both the native executable and `linux/arm64`.

**INFERRED.** ARM64 compilation proves portability/buildability, not constrained-device runtime performance.

## 14. Test evidence

**KNOWN.** `TestCompartmentAndIncidentScope` proves incident and compartment gates.

**KNOWN.** `TestRoles` proves the role matrix.

**KNOWN.** `TestMakerCheckerRequiresPolicyControlledException` proves normal self-approval refusal, incomplete exception refusal, valid three-part exception, and refusal when policy does not permit the exception.

**KNOWN.** `TestAgentInheritsNamedHumanOnly` proves valid attribution, mandatory human naming, inability to switch principals and refusal of fake agent context by a human.

**KNOWN.** `TestValidatePolicy` exercises unknown role, agent/human namespace collision, unknown human and ambiguous agent delegation.

**KNOWN.** `TestMalformedRequestAndUnknownActionAreErrors` proves input/error semantics.

## 15. Security boundary

**KNOWN.** This repository does not authenticate `ActorID`, `HumanID`, `AuthorID` or the runtime single-person-window signal.

**KNOWN.** It does ensure a request flag alone cannot enable self-approval because durable policy permission is separately required.

**UNKNOWN.** Authentication strength, policy administration, runtime-window authority, decision audit and downstream enforcement belong to the consuming application.

## 16. Observed limits

**KNOWN.** V0.1 has a fixed three-role vocabulary and four actions. It does not implement general ABAC, role inheritance, policy persistence/versioning, approval persistence, legal privilege logic, workflow scheduling or user monitoring.

**KNOWN.** There is no constrained-hardware runtime benchmark in this repository.

## 17. Evidence conclusion

**KNOWN.** Source commit `938cbad9760bd23211f8a2a487b3ccf71900a00f` implements one small authorisation decision engine with explicit human/agent separation and a two-authority maker-checker exception gate.

**INFERRED.** The strongest V0.1 property is that delegation and exceptional self-approval do not become alternate hidden privilege systems: both remain inside the same deterministic policy decision.

**UNKNOWN.** Production assurance still depends on authenticated input identities, trusted runtime context, policy governance and durable recording/enforcement outside this package.
