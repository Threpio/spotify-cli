package update

import (
	"errors"
	"spotify/internal"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const repo = "brianstrauch/spotify-cli"

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update CLI to the latest version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			isUpdated, err := IsUpdated(cmd)
			if err != nil {
				return err
			}
			if isUpdated {
				return errors.New(internal.ErrAlreadyUpToDate)
			}

			current, err := parseVersion(cmd)
			if err != nil {
				return err
			}

			latest, err := selfupdate.UpdateSelf(current, repo)
			if err != nil {
				return err
			}

			cmd.Printf("🎁 v%s\n", latest.Version.String())
			return nil
		},
	}
}

func IsUpdated(cmd *cobra.Command) (bool, error) {
	current, err := parseVersion(cmd)
	if err != nil {
		return true, err
	}

	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		return true, err
	}

	isUpdated := !found || current.Equals(latest.Version)
	return isUpdated, nil
}

func parseVersion(cmd *cobra.Command) (semver.Version, error) {
	version := strings.TrimPrefix(cmd.Root().Version, "v")
	return semver.Parse(version)
}
