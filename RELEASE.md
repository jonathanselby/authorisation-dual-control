# Release checklist

## Source and version

- [ ] `VERSION`, CLI `-version` and release tag agree.
- [ ] `go.mod` uses the Garage Go baseline.
- [ ] `CHANGELOG.md` records the release.

## Policy/security gate

- [ ] `ValidatePolicy` rejects unknown roles and ambiguous principal/delegation configuration.
- [ ] an agent ID cannot also be a human subject.
- [ ] agent evaluation always resolves through its named human.
- [ ] incident and compartment gates precede role permission.
- [ ] request context alone cannot enable maker-checker self-approval.
- [ ] a self-approval exception requires policy permission, active trusted window and a reason.

## Engineering gate

- [ ] `gofmt` produces no changes.
- [ ] `go vet ./...` passes.
- [ ] `go test -race -count=1 ./...` passes.
- [ ] native and `linux/arm64` builds pass.
- [ ] exact CLI version check passes.
- [ ] Bafflewick acceptance passes.

## Evidence

- [ ] every PRD criterion maps to executable evidence.
- [ ] `generated/EDD.md` is recovered from the exact source commit.
- [ ] `generated/COMPARISON.md` records matches, limits and unknowns.

Tag format: `v0.1.0`, `v0.2.0`, and so on.
