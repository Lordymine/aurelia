package persona

import (
	"context"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

func (s *CanonicalIdentityService) ResolveIdentity(ctx context.Context, userID string) (*Persona, CanonicalIdentity, error) {
	p, err := LoadPersona(s.identityPath, s.soulPath, s.userPath)
	if err != nil {
		return nil, CanonicalIdentity{}, err
	}

	identity := p.CanonicalIdentity
	if s.memory == nil {
		return p, identity, nil
	}

	identity.AgentName = s.resolveFact(ctx, "agent", "default", "agent.name", identity.AgentName)
	identity.AgentRole = s.resolveFact(ctx, "agent", "default", "agent.role", identity.AgentRole)
	if strings.TrimSpace(userID) != "" {
		identity.UserName = s.resolveFact(ctx, "user", userID, "user.name", identity.UserName)
	}

	return p, identity, nil
}

func (s *CanonicalIdentityService) SeedIdentityFacts(ctx context.Context, userID string, identity CanonicalIdentity, source string) error {
	if s.memory == nil {
		return nil
	}

	facts := map[string]memory.Fact{
		"agent.name": {
			Scope:    "agent",
			EntityID: "default",
			Key:      "agent.name",
			Value:    identity.AgentName,
			Source:   source,
		},
		"agent.role": {
			Scope:    "agent",
			EntityID: "default",
			Key:      "agent.role",
			Value:    identity.AgentRole,
			Source:   source,
		},
	}
	if strings.TrimSpace(userID) != "" && strings.TrimSpace(identity.UserName) != "" && identity.UserName != "nao definido" {
		facts["user.name"] = memory.Fact{
			Scope:    "user",
			EntityID: userID,
			Key:      "user.name",
			Value:    identity.UserName,
			Source:   source,
		}
	}

	return s.ApplyFacts(ctx, userID, facts)
}

func (s *CanonicalIdentityService) ApplyFacts(ctx context.Context, userID string, facts map[string]memory.Fact) error {
	if s.memory == nil {
		return nil
	}

	for _, fact := range facts {
		if strings.TrimSpace(fact.Value) == "" {
			continue
		}
		if fact.Scope == "user" && strings.TrimSpace(fact.EntityID) == "" {
			fact.EntityID = userID
		}
		if err := s.upsertFactWithPriority(ctx, fact); err != nil {
			return err
		}
	}

	return nil
}

func (s *CanonicalIdentityService) ApplyFactsAndSync(ctx context.Context, userID string, facts map[string]memory.Fact) error {
	if err := s.ApplyFacts(ctx, userID, facts); err != nil {
		return err
	}
	return s.SyncFactsToFiles(facts)
}

func (s *CanonicalIdentityService) resolveFact(ctx context.Context, scope, entityID, key, fallback string) string {
	if fact, ok, err := s.memory.GetFact(ctx, scope, entityID, key); err == nil && ok && strings.TrimSpace(fact.Value) != "" {
		return fact.Value
	}
	return fallback
}

func (s *CanonicalIdentityService) upsertFactWithPriority(ctx context.Context, fact memory.Fact) error {
	existing, ok, err := s.memory.GetFact(ctx, fact.Scope, fact.EntityID, fact.Key)
	if err != nil {
		return err
	}
	if ok && factSourcePriority(existing.Source) > factSourcePriority(fact.Source) {
		return nil
	}
	return s.memory.UpsertFact(ctx, fact)
}

func factSourcePriority(source string) int {
	switch strings.TrimSpace(strings.ToLower(source)) {
	case "correction":
		return 40
	case "bootstrap":
		return 30
	case "persona":
		return 20
	case "conversation":
		return 10
	default:
		return 0
	}
}


