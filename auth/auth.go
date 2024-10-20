package auth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/impersonate"
)

const (
	scope = "https://www.googleapis.com/auth/cloud-platform"
)

func NewClient(sa string) *http.Client {
	var err error
	var ts oauth2.TokenSource
	ctx := context.Background()
	if sa == "" {
		ts, err = defaultTokenSource(ctx)
	} else {
		ts, err = impersonatedTokenSource(ctx, sa)
	}
	if err != nil {
		return nil
	}
	return oauth2.NewClient(ctx, ts)
}

func defaultTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	return google.DefaultTokenSource(ctx, scope)
}

func impersonatedTokenSource(ctx context.Context, sa string) (oauth2.TokenSource, error) {
	cfg := impersonate.CredentialsConfig{
		TargetPrincipal: string(sa),
		Scopes:          []string{scope},
	}
	return impersonate.CredentialsTokenSource(ctx, cfg)
}
