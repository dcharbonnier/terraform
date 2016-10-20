package scaleway

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// Bool returns a pointer to of the bool value passed in.
func Bool(val bool) *bool {
	return &val
}

// String returns a pointer to of the string value passed in.
func String(val string) *string {
	return &val
}

// NOTE copied from github.com/scaleway/scaleway-cli/pkg/api/helpers.go
// the helpers.go file pulls in quite a lot dependencies, and they're just convenience wrappers anyway

func deleteServerSafe(s *api.ScalewayAPI, serverID string) error {
	server, err := s.GetServer(serverID)
	if err != nil {
		return err
	}

	if server.State != "stopped" {
		if err := s.PostServerAction(serverID, "poweroff"); err != nil {
			return err
		}
		if err := waitForServerState(s, serverID, "stopped"); err != nil {
			return err
		}
	}

	if err := s.DeleteServer(serverID); err != nil {
		return err
	}
	if rootVolume, ok := server.Volumes["0"]; ok {
		if err := s.DeleteVolume(rootVolume.Identifier); err != nil {
			return err
		}
	}

	return nil
}

// NOTE copied from github.com/scaleway/scaleway-cli/pkg/api/helpers.go
// the helpers.go file pulls in quite a lot dependencies, and they're just convenience wrappers anyway

func waitForServerState(scaleway *api.ScalewayAPI, serverID string, targetState string) error {
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		s, err := scaleway.GetServer(serverID)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		if s.State != targetState {
			return resource.RetryableError(fmt.Errorf("Waiting for server to enter %q state", targetState))
		}

		return nil
	})
}
