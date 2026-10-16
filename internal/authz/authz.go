package authz

import (
	"errors"
	"strings"
)

const (
	RoleMajor    = "Major"
	RoleMinor    = "Minor"
	RoleObserver = "Observer"
)

type Subject struct {
	ID                   string   `json:"id"`
	Role                 string   `json:"role"`
	AssignedIncidents    []string `json:"assigned_incidents,omitempty"`
	GlobalIncidentAccess bool     `json:"global_incident_access,omitempty"`
	Compartments         []string `json:"compartments,omitempty"`
}

type Delegation struct {
	AgentID string `json:"agent_id"`
	HumanID string `json:"human_id"`
}

type Policy struct {
	Subjects                []Subject    `json:"subjects"`
	Delegations             []Delegation `json:"delegations,omitempty"`
	AllowSinglePersonWindow bool         `json:"allow_single_person_window,omitempty"`
}

type Item struct {
	IncidentID  string `json:"incident_id"`
	Compartment string `json:"compartment,omitempty"`
	AuthorID    string `json:"author_id,omitempty"`
}

type Request struct {
	ActorID            string `json:"actor_id"`
	HumanID            string `json:"human_id,omitempty"`
	Action             string `json:"action"`
	Item               Item   `json:"item"`
	SinglePersonWindow bool   `json:"single_person_window,omitempty"`
	SelfApprovalReason string `json:"self_approval_reason,omitempty"`
}

type Decision struct {
	Allowed      bool   `json:"allowed"`
	HumanID      string `json:"human_id,omitempty"`
	AgentID      string `json:"agent_id,omitempty"`
	SelfApproval bool   `json:"self_approval,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

func ValidatePolicy(p Policy) error {
	subjects := make(map[string]Subject, len(p.Subjects))
	for _, s := range p.Subjects {
		if strings.TrimSpace(s.ID) == "" {
			return errors.New("subject id is required")
		}
		if _, exists := subjects[s.ID]; exists {
			return errors.New("duplicate subject id")
		}
		if s.Role != RoleMajor && s.Role != RoleMinor && s.Role != RoleObserver {
			return errors.New("unknown subject role")
		}
		if hasDuplicateOrBlank(s.AssignedIncidents) {
			return errors.New("invalid assigned incidents")
		}
		if hasDuplicateOrBlank(s.Compartments) {
			return errors.New("invalid compartments")
		}
		subjects[s.ID] = s
	}

	agents := map[string]string{}
	for _, d := range p.Delegations {
		if strings.TrimSpace(d.AgentID) == "" || strings.TrimSpace(d.HumanID) == "" || d.AgentID == d.HumanID {
			return errors.New("invalid delegation")
		}
		if _, exists := subjects[d.AgentID]; exists {
			return errors.New("agent id cannot also be a human subject")
		}
		if _, exists := subjects[d.HumanID]; !exists {
			return errors.New("delegation references unknown human")
		}
		if prior, exists := agents[d.AgentID]; exists {
			if prior != d.HumanID {
				return errors.New("agent has ambiguous human delegation")
			}
			return errors.New("duplicate delegation")
		}
		agents[d.AgentID] = d.HumanID
	}
	return nil
}

func Evaluate(p Policy, r Request) (Decision, error) {
	if err := ValidatePolicy(p); err != nil {
		return Decision{}, err
	}
	if strings.TrimSpace(r.ActorID) == "" || strings.TrimSpace(r.Action) == "" || strings.TrimSpace(r.Item.IncidentID) == "" {
		return Decision{}, errors.New("actor, action and incident are required")
	}
	if r.HumanID == r.ActorID && r.HumanID != "" {
		return Decision{}, errors.New("human_id is only valid for an agent acting for another principal")
	}

	human := r.ActorID
	agent := ""
	if delegatedHuman, isAgent := agentPrincipal(p, r.ActorID); isAgent {
		if r.HumanID == "" {
			return Decision{AgentID: r.ActorID, Reason: "agent requires named human"}, nil
		}
		if r.HumanID != delegatedHuman {
			return Decision{AgentID: r.ActorID, Reason: "agent is not delegated by named human"}, nil
		}
		human = r.HumanID
		agent = r.ActorID
	} else if r.HumanID != "" {
		return Decision{Reason: "actor is not a delegated agent"}, nil
	}

	s, ok := subject(p, human)
	if !ok {
		return Decision{HumanID: human, AgentID: agent, Reason: "unknown human"}, nil
	}
	d := Decision{HumanID: human, AgentID: agent}
	if !incidentAllowed(s, r.Item.IncidentID) {
		d.Reason = "not assigned to incident"
		return d, nil
	}
	if r.Item.Compartment != "" && !has(s.Compartments, r.Item.Compartment) {
		d.Reason = "compartment clearance required"
		return d, nil
	}

	switch r.Action {
	case "read":
		d.Allowed = true
	case "comment":
		if s.Role == RoleMajor || s.Role == RoleMinor {
			d.Allowed = true
		} else {
			d.Reason = "role cannot comment"
		}
	case "propose":
		if s.Role == RoleMajor {
			d.Allowed = true
		} else {
			d.Reason = "role cannot propose"
		}
	case "confirm":
		if s.Role != RoleMajor {
			d.Reason = "role cannot confirm"
			break
		}
		if r.Item.AuthorID == human {
			if !p.AllowSinglePersonWindow || !r.SinglePersonWindow || strings.TrimSpace(r.SelfApprovalReason) == "" {
				d.Reason = "author cannot confirm own version"
				break
			}
			d.SelfApproval = true
		}
		d.Allowed = true
	default:
		return Decision{}, errors.New("unknown action")
	}
	return d, nil
}

func subject(p Policy, id string) (Subject, bool) {
	for _, s := range p.Subjects {
		if s.ID == id {
			return s, true
		}
	}
	return Subject{}, false
}

func agentPrincipal(p Policy, agent string) (string, bool) {
	for _, d := range p.Delegations {
		if d.AgentID == agent {
			return d.HumanID, true
		}
	}
	return "", false
}

func incidentAllowed(s Subject, id string) bool {
	return s.GlobalIncidentAccess || has(s.AssignedIncidents, id)
}

func has(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func hasDuplicateOrBlank(xs []string) bool {
	seen := map[string]bool{}
	for _, x := range xs {
		if strings.TrimSpace(x) == "" || seen[x] {
			return true
		}
		seen[x] = true
	}
	return false
}
