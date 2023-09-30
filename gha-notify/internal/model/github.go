package model

import (
	"github.com/shogo82148/goat/jwt"
)

// GitHub's custom claims
// https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect
type ActionsIDToken struct {
	*jwt.Claims
	Environment          string `jwt:"environment"`
	Ref                  string `jwt:"ref"`
	SHA                  string `jwt:"sha"`
	Repository           string `jwt:"repository"`
	RepositoryOwner      string `jwt:"repository_owner"`
	ActorID              string `jwt:"actor_id"`
	RepositoryVisibility string `jwt:"repository_visibility"`
	RepositoryID         string `jwt:"repository_id"`
	RepositoryOwnerID    string `jwt:"repository_owner_id"`
	RunID                string `jwt:"run_id"`
	RunNumber            string `jwt:"run_number"`
	RunAttempt           string `jwt:"run_attempt"`
	Actor                string `jwt:"actor"`
	Workflow             string `jwt:"workflow"`
	HeadRef              string `jwt:"head_ref"`
	BaseRef              string `jwt:"base_ref"`
	EventName            string `jwt:"event_name"`
	EventType            string `jwt:"branch"`
	RefType              string `jwt:"ref_type"`
	JobWorkflowRef       string `jwt:"job_workflow_ref"`
}
