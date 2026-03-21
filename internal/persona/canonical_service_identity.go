package persona

import (
	"context"
)

// ResolveIdentity loads the persona from files and returns the canonical identity.
func (s *CanonicalIdentityService) ResolveIdentity(ctx context.Context) (*Persona, CanonicalIdentity, error) {
	_ = ctx
	p, err := LoadPersona(s.identityPath, s.soulPath, s.userPath)
	if err != nil {
		return nil, CanonicalIdentity{}, err
	}

	return p, p.CanonicalIdentity, nil
}
