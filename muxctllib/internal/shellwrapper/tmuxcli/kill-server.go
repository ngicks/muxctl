package tmuxcli

import "context"

// kill-server
//         Kill the tmux server and clients and destroy all sessions.

// KillServer runs kill-server.
func (s *ShellOut) KillServer(ctx context.Context) error {
	_, err := s.run(ctx, "kill-server")
	return err
}
