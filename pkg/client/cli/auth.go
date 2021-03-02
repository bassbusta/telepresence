package cli

import (
	"errors"
	"os"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/telepresenceio/telepresence/rpc/v2/connector"
	"github.com/telepresenceio/telepresence/rpc/v2/manager"
	"github.com/telepresenceio/telepresence/v2/pkg/client"
	"github.com/telepresenceio/telepresence/v2/pkg/client/cache"
)

func EnsureLoggedIn(ctx context.Context) (manager.UserInfo, error) {
	var userinfo *manager.UserInfo
	err := WithConnector(ctx, func(
		ctx context.Context,
		connectorClient connector.ConnectorClient,
		_ manager.ManagerClient,
	) error {
		var err error
		user, err = connectorClient.Login(ctx, &empty.Empty{})
		return err
	})
	return userinfo, err
}

// Command returns the telepresence sub-command "auth"
func LoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "login",
		Args: cobra.NoArgs,

		Short: "Authenticate to Ambassador Cloud",
		Long:  "Authenticate to Ambassador Cloud",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return EnsureLoggedIn(cmd.Context())
		},
	}
}

func LogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "logout",
		Args: cobra.NoArgs,

		Short: "Logout from Ambassador Cloud",
		Long:  "Logout from Ambassador Cloud",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return WithConnector(cmd.Context(), func(
				ctx context.Context,
				connectorClient connector.ConnectorClient,
				_ manager.ManagerClient,
			) error {
				_, err := connectorClient.Logout(ctx, &empty.Empty{})
				return err
			})
		},
	}
}
