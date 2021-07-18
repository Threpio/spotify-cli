package queue

import (
	"spotify/internal"
	"spotify/internal/status"
	"strings"

	"github.com/brianstrauch/spotify"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "queue song",
		Aliases: []string{"q"},
		Short:   "Queue a specific song.",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := internal.Authenticate()
			if err != nil {
				return err
			}

			query := strings.Join(args, " ")

			output, err := Queue(api, query)
			if err != nil {
				return err
			}

			cmd.Print(output)
			return nil
		},
	}
}

func Queue(api internal.APIInterface, query string) (string, error) {
	track, err := internal.Search(api,"tracks", query)
	if err != nil {
		return "", err
	}

	if err := api.Queue(track.URI); err != nil {
		return "", err
	}

	return show(track), nil
}

func show(track *spotify.Track) string {
	output := status.PrefixLineWithEmoji("🎵", track.Name)
	output += status.PrefixLineWithEmoji("🎤", status.JoinArtists(track.Artists))
	return output
}
