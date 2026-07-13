package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/cfgtree"
	"github.com/nveeser/vyconfigure/section"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultMappingFile = "vysync.yaml"
	appVersion         = "development"
)

// NewRoot creates an instance of the root command.
func NewRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vysync",
		Short:   "Declarative configuration for VyOS.",
		Version: appVersion,
	}

	pf := rootCmd.PersistentFlags()
	pf.String("host", "", "The hostname of the VyOS HTTP API.")
	viper.BindEnv("host", "VYSYNC_HOST")

	pf.StringP("api-key", "k", "", "API key for the HTTP API.")
	viper.BindEnv("api-key", "VYSYNC_API_KEY")

	pf.StringP("config-dir", "d", ".", "Directory where config is stored.")
	viper.BindEnv("config-dir", "VYSYNC_CONFIG_DIR")

	pf.BoolP("insecure", "i", false, "Whether to skip verifying the SSL certificate.")
	viper.BindEnv("insecure", "VYSYNC_INSECURE")

	pf.BoolP("debug", "D", false, "Enable Debug mode.")
	viper.BindEnv("debug", "VYSYNC_DEBUG")
	viper.BindPFlags(pf)

	rootCmd.AddCommand(
		newVersionCmd(),
		newApplyCmd(),
		newSyncCmd(),
		newPlanCmd(),
		newShowCmd(),
	)
	return rootCmd
}

type SyncConfig struct {
	ConfigDir string                 `yaml:"ConfigDir"`
	Host      string                 `yaml:"Host"`
	APIKey    string                 `yaml:"APIKey"`
	Insecure  bool                   `yaml:"Insecure"`
	Debug     bool                   `yaml:"Debug"`
	Mappings  []section.MappingEntry `yaml:"Mappings"`
}

func newRepo() (*cfgtree.Repo, error) {
	var mapper *section.Mapper
	for _, filename := range []string{
		path.Join(viper.GetString("config-dir"), defaultMappingFile),
	} {
		var err error
		mapper, err = section.NewMapper(filename)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}
	return &cfgtree.Repo{
		ConfigDirectory: viper.GetString("config-dir"),
		SectionMapper:   mapper,
		Ignore:          []string{strings.TrimSuffix(defaultMappingFile, ".yaml")},
	}, nil
}

func newRepoAndClient(ctx context.Context) (*cfgtree.Repo, *cfgtree.Client, error) {
	repo, err := newRepo()
	if err != nil {
		return nil, nil, err
	}
	opts := []vyos.Option{
		vyos.Token(viper.GetString("api-key")),
		vyos.Timeout(0),
	}
	if viper.GetBool("insecure") {
		opts = append(opts, vyos.Insecure())
	}
	if viper.GetBool("debug") {
		opts = append(opts, vyos.DebugLogging())
	}
	client, err := vyos.NewClient(viper.GetString("host"), opts...)
	if err != nil {
		return nil, nil, err
	}
	cc := &cfgtree.Client{
		Client:        client,
		SectionMapper: repo.SectionMapper,
	}
	ctx, done := context.WithTimeout(ctx, 5*time.Second)
	defer done()
	r, err := cc.OpMode().Info(ctx, vyos.InfoRequest{Version: true, Hostname: true})
	if err != nil {
		return nil, nil, err
	}
	color.Cyan("Host: %s (%s): %s", r.Hostname, r.Version, r.Banner)
	return repo, cc, nil
}
