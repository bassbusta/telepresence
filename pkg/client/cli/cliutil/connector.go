package cliutil

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/datawire/dlib/dgroup"
	"github.com/telepresenceio/telepresence/rpc/v2/connector"
	"github.com/telepresenceio/telepresence/rpc/v2/manager"
	"github.com/telepresenceio/telepresence/v2/pkg/client"
	"github.com/telepresenceio/telepresence/v2/pkg/filelocation"
)

func launchConnector() error {
	args := []string{client.GetExe(), "connector-foreground"}

	cmd := exec.Command(args[0], args[1:]...)
	// Process must live in a process group of its own to prevent
	// getting affected by <ctrl-c> in the terminal
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s: %w", shellquote.Join(args...), err)
	}
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("%s: %w", shellquote.Join(args...), err)
	}

	return nil
}

// WithConnector ensures that the connector is running, establishes a connection to it, and runs the
// given function with that connection.  It streams to stdout any messages that the connector wants
// us to display to the user (which WithConnector listens for via the UserNotifications gRPC call).
// WithConnector does NOT make the "Connect" gRPC call or any other gRPC call but UserNotifications.
func WithConnector(ctx context.Context, fn func(context.Context, connector.ConnectorClient, manager.ManagerClient) error) error {
	if !client.SocketExists(client.ConnectorSocketName) {
		if err := launchConnector(); err != nil {
			return errors.Wrap(err, "failed to launch the connector service")
		}

		if err := client.WaitUntilSocketAppears("connector", client.ConnectorSocketName, 10*time.Second); err != nil {
			logDir, _ := filelocation.AppUserLogDir(ctx)
			return fmt.Errorf("connector service did not start (see %q for more info)", filepath.Join(logDir, "connector.log"))
		}
	}

	conn, err := client.DialSocket(ctx, client.ConnectorSocketName)
	if err != nil {
		return err
	}
	defer conn.Close()
	connectorClient := connector.NewConnectorClient(conn)
	managerClient := manager.NewManagerClient(conn)

	grp := dgroup.NewGroup(ctx, dgroup.GroupConfig{
		ShutdownOnNonError: true,
	})

	grp.Go("stdio", func(ctx context.Context) error {
		stream, err := connectorClient.UserNotifications(ctx, &empty.Empty{})
		if err != nil {
			return err
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			fmt.Println(msg.Message)
		}
	})
	grp.Go("main", func(ctx context.Context) error {
		return fn(ctx, connectorClient, managerClient)
	})

	return grp.Wait()
}
