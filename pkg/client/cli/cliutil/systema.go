package cliutil

import (
	"context"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/telepresenceio/telepresence/rpc/v2/connector"
	"github.com/telepresenceio/telepresence/v2/pkg/client/cache"
)

// EnsureLoggedIn ensures that the user is logged in to Ambassador Cloud.  An error is returned if
// login fails.  The result code will indicate if this is a new login or if it resued an existing
// login.
func EnsureLoggedIn(ctx context.Context) (connector.LoginResult_Code, error) {
	var resp *connector.LoginResult
	err := WithConnector(ctx, func(ctx context.Context, connectorClient connector.ConnectorClient) error {
		var err error
		resp, err = connectorClient.Login(ctx, &empty.Empty{})
		return err
	})
	return resp.GetCode(), err
}

// HasLoggedIn returns true if either the user has an active login session or an expired login
// session, and returns false the user has never logged in or has explicitly logged out.
func HasLoggedIn(ctx context.Context) bool {
	var dat interface{}
	err := cache.LoadFromUserCache(ctx, &dat, "tokens.json")
	return err == nil
}
