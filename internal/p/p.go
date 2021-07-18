package p

import (
	"errors"
	"spotify/internal"
	"spotify/internal/pause"
	"spotify/internal/play"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "p [song]",
		// Keep hidden, since this command is an alias.
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := internal.Authenticate()
			if err != nil {
				return err
			}

			query := strings.Join(args, " ")

			deviceID, err := cmd.Flags().GetString("device-id")
			if err != nil {
				return err
			}
			playlistQuery, err := cmd.Flags().GetString("playlist")
			if err != nil {
				return err
			}

			status, err := p(api, query, playlistQuery, deviceID)
			if err != nil {
				return err
			}

			cmd.Print(status)
			return nil
		},
	}

	cmd.Flags().String("device-id", "", "Device ID from 'spotify device list'.")
	cmd.Flags().String("playlist", "", "Playlist name from 'spotify playlist list'")

	return cmd
}

func p(api internal.APIInterface, query, contextQuery, deviceID string) (string, error) {
	if len(query) > 0 {
		return play.Play(api, query, contextQuery, deviceID)
	}

	playback, err := api.GetPlayback()
	if err != nil {
		return "", err
	}

	if playback == nil {
		return "", errors.New(internal.ErrNoActiveDevice)
	}

	if playback.IsPlaying {
		return pause.Pause(api, deviceID)
	} else {
		return play.Play(api, "", "", deviceID)
	}
}
