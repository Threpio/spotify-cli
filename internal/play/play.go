package play

import (
	"errors"
	"github.com/brianstrauch/spotify"
	"spotify/internal"
	"spotify/internal/status"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play [song]",
		Short: "Play current song, or a specific song.",
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

			status, err := Play(api, query, playlistQuery, deviceID)
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

func Play(api internal.APIInterface, query, contextQuery, deviceID string) (string, error) {
	playback, err := api.GetPlayback()
	if err != nil {
		return "", err
	}

	if playback == nil {
		return "", errors.New(internal.ErrNoActiveDevice)
	}

	if len(query) > 0 {
		track, err := internal.Search(api, "track", query)
		if err != nil {
			return "", err
		}

		if err := api.Play(deviceID, contextQuery, track.URI); err != nil {
			return "", err
		}
	} else if len(contextQuery) > 0 {

		// Return a different API interface required for the playlist commands?
		api, err := internal.Authenticate()
		if err != nil {
			return "", err
		}

		playlists, err := api.GetPlaylists()
		if err != nil {
			return "", err
		}

		for _, playlist := range playlists {
			if strings.EqualFold(playlist.Name, contextQuery) {
				if err := api.Play(deviceID, playlist.URI); err != nil {
					return "", err
				}
			}
		}
	} else {
		if err := api.Play(deviceID, "", ""); err != nil {
			return "", err
		}
	}

	playback, err = internal.WaitForUpdatedPlayback(api, func(playback *spotify.Playback) bool {
		// The first check safeguards against empty playback objects
		return len(playback.Item.ID) > 0 && playback.IsPlaying
	})
	if err != nil {
		return "", err
	}

	return status.Show(playback), nil
}
