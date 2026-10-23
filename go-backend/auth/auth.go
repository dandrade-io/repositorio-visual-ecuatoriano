package auth

import (
	"context"
	"repository/config"

	"github.com/coreos/go-oidc/v3/oidc"
)

var Verfier *oidc.IDTokenVerifier
var AppClient = "backend-api"

func init() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.KCBaseURL+"/realms/"+config.KCRealm)
	if err != nil {
		panic(err)
	}

	Verfier = provider.VerifierContext(ctx, &oidc.Config{
		ClientID:             AppClient,
		SupportedSigningAlgs: []string{oidc.RS256},
	})
}
