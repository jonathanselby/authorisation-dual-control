# Authorisation and Dual Control — Product Requirements Document

**Version:** 0.1

## 1. Problem

Role checks, need-to-know, agent delegation and maker-checker are often implemented in separate paths until those paths disagree. A reusable component needs to answer the decision once without becoming an identity platform, workflow engine or user-monitoring system.

A self-approval exception is particularly dangerous if the request itself can enable it. Durable policy must permit the exception independently from the runtime statement that the exceptional window is active.

## 2. Product objective

Provide a small deterministic Go component that decides `read`, `comment`, `propose` and `confirm` using role, incident assignment, compartment clearance, named-human delegation and maker-checker rules.

> **One access model. One item lifecycle. No author approves their own version unless policy explicitly permits the declared exception.**

## 3. Users

Applications with high-trust human teams that still require explicit need-to-know, delegated-agent attribution and dual control.

## 4. In scope

- Major, Minor and Observer roles;
- incident assignment and optional global incident access;
- named compartment clearances;
- structurally separate agent IDs delegated to one named human;
- read, comment, propose and confirm decisions;
- policy-gated single-person-window self-approval with mandatory reason;
- policy validation;
- deterministic decision result and stable CLI exit codes;
- Go 1.26 and `linux/arm64` build proof.

## 5. Out of scope

Authentication, identity lifecycle, approval persistence, workflow scheduling, single-person-window scheduling, user monitoring/scoring, legal determination of privilege, application-specific schemas and policy administration/maker-checker.

## 6. Functional requirements

### FR1 — Valid policy first
Unknown roles, duplicate subjects, blank/duplicate scopes, invalid delegation, delegation to an unknown human, ambiguous agent delegation, or an agent ID also defined as a human subject are configuration errors.

### FR2 — Need-to-know before role
A human requires incident assignment or explicit global incident access. A compartmented item additionally requires the named compartment.

### FR3 — Role matrix
Observers may read. Minors may read/comment. Majors may read/comment/propose/confirm.

### FR4 — Agent principal
A configured agent is not a human subject. Every agent request must name the human principal, the delegation must match, and the decision is evaluated entirely from that human's subject record.

### FR5 — Maker-checker
A confirmation requires the authoritative item author ID. A Major cannot confirm an item version authored by the same human principal under normal operation. A confirmation with no author identity is malformed input because the maker-checker rule cannot be evaluated.

### FR6 — Controlled single-person exception
Self-approval requires all three:

1. durable policy `allow_single_person_window=true`;
2. trusted runtime context `single_person_window=true`;
3. non-blank `self_approval_reason`.

Successful self-approval is explicitly flagged in the decision.

### FR7 — Denial vs malformed input
Unknown human, missing scope, insufficient role and delegation mismatch are deterministic denials. Invalid policy, missing required request fields (including item author for `confirm`) and unknown actions are execution/input errors.

## 7. Acceptance criteria

V0.1 is accepted when automated evidence shows:

1. an unassigned Minor cannot read another incident;
2. a compartment requires both incident access and clearance;
3. the role matrix is enforced;
4. a valid agent inherits its named human's rights and attribution;
5. an agent without a named human is denied;
6. an agent cannot switch to another human;
7. a normal human cannot spoof agent context;
8. invalid/ambiguous policy is rejected;
9. a confirmation without an authoritative author ID is rejected as malformed input;
10. author self-confirmation is denied normally;
11. request-only single-person activation cannot enable self-approval;
12. policy + active window + reason can permit self-approval and the decision flags it;
13. malformed requests and unknown actions are errors;
14. format, vet, race tests, native build, exact version and `linux/arm64` build pass.

## 8. Definition of done

All criteria have executable evidence; the trust boundary for actor/runtime context is explicit; Go 1.26 CI is green; and generated design evidence independently recovers the release-candidate implementation.
