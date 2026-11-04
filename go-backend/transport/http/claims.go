package transporthttp

import (
	"context"
	"fmt"
	"slices"
)

type Claims struct {
	Subject     string `json:"sub"`
	Typ         string `json:"typ"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAcesss map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

func (c *Claims) DomainRoles() []string {
	if roles, ok := c.ResourceAcesss["backend-api"]; ok {
		return roles.Roles
	}

	return nil
}

func (c *Claims) MainRole() (string, error) {
	roles := c.DomainRoles()

	if slices.Contains(roles, "expert_verifier") {
		return "expert_verifier", nil
	} else if slices.Contains(roles, "ordinary_verifier") {
		return "ordinary_verifier", nil
	}

	return "", fmt.Errorf("The user has no assigned roles.")
}

type ctxKey struct{}

func ClaimsFrom(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(Claims)
	return c, ok
}
