package cli

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/telepresenceio/telepresence/rpc/v2/connector"
	"github.com/telepresenceio/telepresence/v2/pkg/client"
	"github.com/telepresenceio/telepresence/v2/pkg/client/cache"
)

func dashboardCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "dashboard",
		Args: cobra.NoArgs,

		Short: "Open the dashboard in a web page",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := client.LoadEnv(cmd.Context())
			if err != nil {
				return err
			}

			// Ensure we're logged in
			userinfo, err := EnsureLoggedIn(cmd.Context())
			if err != nil {
				return err
			}

			if true /* TODO */ {
				// The LoginFlow actually takes the user to the dashboard. Hence the else here.
				err := browser.OpenURL(fmt.Sprintf("https://%s/cloud/preview", env.SystemAHost))
				if err != nil {
					return err
				}
			}

			return nil
		}}
}
