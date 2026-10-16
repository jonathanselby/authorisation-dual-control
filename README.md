# Authorisation and Dual Control

[![test](https://github.com/jonathanselby/authorisation-dual-control/actions/workflows/test.yml/badge.svg)](https://github.com/jonathanselby/authorisation-dual-control/actions/workflows/test.yml)

**One access model. One item lifecycle. No author approves their own version unless policy explicitly permits the declared exception.**

`authorisation-dual-control` is a small deterministic Go policy component for role access, incident need-to-know, privilege compartments, delegated agents and maker-checker confirmation. It is independently usable; CRAB is one consumer, not part of the runtime contract.

## What V0.1 proves

- Major, Minor and Observer roles have one explicit action matrix;
- incident assignment is checked before role permission unless the human has explicit global incident access;
- compartment access is an additional need-to-know gate;
- agent IDs are structurally separate from human subjects;
- each agent resolves to exactly one named human and inherits only that human's rights;
- an agent must name its human principal on every request;
- an author cannot confirm their own version unless policy permits single-person exceptions, the trusted workflow says the window is active, and a non-blank reason is supplied;
- policy configuration is validated before a decision;
- deterministic denials are distinct from malformed policy/request errors.

## Role matrix

| Action | Major | Minor | Observer |
|---|:---:|:---:|:---:|
| read | yes | yes | yes |
| comment | yes | yes | no |
| propose | yes | no | no |
| confirm | yes | no | no |

Incident and compartment gates apply before this matrix.

## Install

Requires Go 1.26+.

```bash
go install github.com/jonathanselby/authorisation-dual-control/cmd/authzcheck@main
authzcheck -version
```

Expected:

```text
authzcheck 0.1.0
```

## CLI

```bash
authzcheck -policy policy.json -request request.json
```

## Exit codes

| Code | Meaning |
|---:|---|
| `0` | allowed |
| `1` | deterministic policy denial |
| `2` | invalid CLI usage |
| `3` | malformed policy/request, file I/O or execution error |

## Trust boundary

The engine **does not authenticate** IDs. The caller must provide authentic actor, human, item and runtime-context data.

`allow_single_person_window` is durable policy. `single_person_window` is runtime context supplied by a trusted workflow/scheduler. Self-approval requires both. If either can be actor-controlled, the exception is not safe.

The engine does not persist approvals, schedule windows, determine legal privilege, score people, monitor behaviour or own application item schemas.

## Repository map

```text
authorisation-dual-control/
├── .github/workflows/test.yml
├── PRD.md
├── EDD.md
├── README.md
├── SKILL.md
├── RELEASE.md
├── VERSION
├── CHANGELOG.md
├── LICENSE
├── architecture.txt
├── go.mod
├── cmd/authzcheck/main.go
├── internal/authz/
├── examples/bafflewick/
└── generated/
    ├── EDD.md
    └── COMPARISON.md
```

## Engineering evidence

- [`PRD.md`](PRD.md) — boundary and acceptance criteria;
- [`EDD.md`](EDD.md) — intended policy design;
- [`architecture.txt`](architecture.txt) — compact decision flow;
- [`internal/authz`](internal/authz) — implementation and behavioural tests;
- [`examples/bafflewick`](examples/bafflewick) — fictional acceptance case;
- [`generated/EDD.md`](generated/EDD.md) — code-derived design evidence pinned to a source commit;
- [`generated/COMPARISON.md`](generated/COMPARISON.md) — intended design versus recovered implementation;
- [`RELEASE.md`](RELEASE.md) — release gate.

## Design rule

> **One access model. One item lifecycle. No author approves their own version unless policy explicitly permits the declared exception.**

## Licence

MIT. See [`LICENSE`](LICENSE).
