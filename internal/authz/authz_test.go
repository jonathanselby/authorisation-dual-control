package authz

import "testing"

func policy() Policy {
	return Policy{
		Subjects: []Subject{
			{ID: "major", Role: RoleMajor, GlobalIncidentAccess: true, Compartments: []string{"legal"}},
			{ID: "minor", Role: RoleMinor, AssignedIncidents: []string{"I"}},
			{ID: "observer", Role: RoleObserver, AssignedIncidents: []string{"I"}},
		},
		Delegations:             []Delegation{{AgentID: "agent", HumanID: "major"}},
		AllowSinglePersonWindow: true,
	}
}

func TestCompartmentAndIncidentScope(t *testing.T) {
	p := policy()
	d, _ := Evaluate(p, Request{ActorID: "minor", Action: "read", Item: Item{IncidentID: "I", Compartment: "legal"}})
	if d.Allowed || d.Reason != "compartment clearance required" {
		t.Fatalf("compartment leak: %+v", d)
	}
	d, _ = Evaluate(p, Request{ActorID: "major", Action: "read", Item: Item{IncidentID: "I", Compartment: "legal"}})
	if !d.Allowed {
		t.Fatal(d.Reason)
	}
	d, _ = Evaluate(p, Request{ActorID: "minor", Action: "read", Item: Item{IncidentID: "OTHER"}})
	if d.Allowed || d.Reason != "not assigned to incident" {
		t.Fatalf("incident scope failed: %+v", d)
	}
}

func TestRoles(t *testing.T) {
	p := policy()
	cases := []struct {
		actor, action string
		want          bool
	}{
		{"observer", "read", true},
		{"observer", "comment", false},
		{"minor", "comment", true},
		{"minor", "propose", false},
		{"minor", "confirm", false},
		{"major", "propose", true},
		{"major", "confirm", true},
	}
	for _, tc := range cases {
		d, err := Evaluate(p, Request{ActorID: tc.actor, Action: tc.action, Item: Item{IncidentID: "I", AuthorID: "someone-else"}})
		if err != nil || d.Allowed != tc.want {
			t.Fatalf("%s %s: %+v %v", tc.actor, tc.action, d, err)
		}
	}
}

func TestMakerCheckerRequiresPolicyControlledException(t *testing.T) {
	p := policy()
	base := Request{ActorID: "major", Action: "confirm", Item: Item{IncidentID: "I", AuthorID: "major"}}
	for _, r := range []Request{
		base,
		func() Request { x := base; x.SinglePersonWindow = true; return x }(),
		func() Request { x := base; x.SelfApprovalReason = "out of hours"; return x }(),
	} {
		d, err := Evaluate(p, r)
		if err != nil || d.Allowed {
			t.Fatalf("self approval escaped gate: %+v %v", d, err)
		}
	}

	r := base
	r.SinglePersonWindow = true
	r.SelfApprovalReason = "out of hours"
	d, err := Evaluate(p, r)
	if err != nil || !d.Allowed || !d.SelfApproval {
		t.Fatalf("declared exception rejected: %+v %v", d, err)
	}

	p.AllowSinglePersonWindow = false
	d, err = Evaluate(p, r)
	if err != nil || d.Allowed {
		t.Fatalf("request enabled exception without policy permission: %+v %v", d, err)
	}
}

func TestAgentInheritsNamedHumanOnly(t *testing.T) {
	p := policy()
	d, err := Evaluate(p, Request{ActorID: "agent", HumanID: "major", Action: "propose", Item: Item{IncidentID: "I"}})
	if err != nil || !d.Allowed || d.HumanID != "major" || d.AgentID != "agent" {
		t.Fatalf("valid delegation: %+v %v", d, err)
	}

	d, err = Evaluate(p, Request{ActorID: "agent", Action: "propose", Item: Item{IncidentID: "I"}})
	if err != nil || d.Allowed || d.Reason != "agent requires named human" {
		t.Fatalf("agent without human: %+v %v", d, err)
	}

	d, err = Evaluate(p, Request{ActorID: "agent", HumanID: "minor", Action: "comment", Item: Item{IncidentID: "I"}})
	if err != nil || d.Allowed || d.Reason != "agent is not delegated by named human" {
		t.Fatalf("agent changed principal: %+v %v", d, err)
	}

	d, err = Evaluate(p, Request{ActorID: "major", HumanID: "minor", Action: "read", Item: Item{IncidentID: "I"}})
	if err != nil || d.Allowed || d.Reason != "actor is not a delegated agent" {
		t.Fatalf("human spoofed agent: %+v %v", d, err)
	}
}

func TestValidatePolicy(t *testing.T) {
	p := policy()
	if err := ValidatePolicy(p); err != nil {
		t.Fatal(err)
	}

	bad := p
	bad.Subjects = append([]Subject{}, p.Subjects...)
	bad.Subjects[0].Role = "Administrator"
	if err := ValidatePolicy(bad); err == nil {
		t.Fatal("unknown role accepted")
	}

	bad = p
	bad.Subjects = append([]Subject{}, p.Subjects...)
	bad.Subjects = append(bad.Subjects, Subject{ID: "agent", Role: RoleMajor})
	if err := ValidatePolicy(bad); err == nil {
		t.Fatal("agent id also accepted as human subject")
	}

	bad = p
	bad.Delegations = []Delegation{{AgentID: "agent", HumanID: "missing"}}
	if err := ValidatePolicy(bad); err == nil {
		t.Fatal("unknown delegation human accepted")
	}

	bad = p
	bad.Delegations = []Delegation{{AgentID: "agent", HumanID: "major"}, {AgentID: "agent", HumanID: "minor"}}
	if err := ValidatePolicy(bad); err == nil {
		t.Fatal("ambiguous agent accepted")
	}
}

func TestMalformedRequestAndUnknownActionAreErrors(t *testing.T) {
	p := policy()
	if _, err := Evaluate(p, Request{ActorID: "major", Action: "read"}); err == nil {
		t.Fatal("missing incident accepted")
	}
	if _, err := Evaluate(p, Request{ActorID: "major", Action: "delete", Item: Item{IncidentID: "I"}}); err == nil {
		t.Fatal("unknown action accepted")
	}
}
