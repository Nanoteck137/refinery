package config

import (
	"fmt"

	"github.com/nanoteck137/refinery"
	"github.com/nanoteck137/validate"
	"github.com/spf13/viper"
)

type ConfigOidcProvider struct {
	Name         string `mapstructure:"name"`
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	IssuerUrl    string `mapstructure:"issuer_url"`
	RedirectUrl  string `mapstructure:"redirect_url"`
}

func (c ConfigOidcProvider) Validate() error {
	return validate.ValidateStruct(&c,
		validate.Field(&c.Name, validate.Required),
		validate.Field(&c.ClientId, validate.Required),
		validate.Field(&c.ClientSecret, validate.Required),
		validate.Field(&c.IssuerUrl, validate.Required),
		validate.Field(&c.RedirectUrl, validate.Required),
	)
}

type Config struct {
	RunMigrations bool   `mapstructure:"run_migrations"`
	ListenAddr    string `mapstructure:"listen_addr"`
	DataDir       string `mapstructure:"data_dir"`
	JwtSecret     string `mapstructure:"jwt_secret"`

	WebDir string `mapstructure:"web"`

	OidcProviders map[string]ConfigOidcProvider `mapstructure:"oidc_providers"`
}

func (c Config) Validate() error {
	return validate.ValidateStruct(&c,
		validate.Field(&c.ListenAddr, validate.Required),
		validate.Field(&c.DataDir, validate.Required),
		validate.Field(&c.JwtSecret, validate.Required),

		validate.Field(&c.OidcProviders, validate.Required, validate.Length(1, 0)),
	)
}

func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	// NOTE(patrik): Set default values here
	v.SetDefault("run_migrations", "true")
	v.SetDefault("listen_addr", ":3000")
	v.BindEnv("data_dir")
	v.BindEnv("jwt_secret")

	v.BindEnv("web")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	v.SetEnvPrefix(refinery.AppName)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	override := v.GetString("config_override")
	if override != "" {
		v.SetConfigFile(override)
		v.MergeInConfig()
	}

	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// configCopy := LoadedConfig
	// configCopy.OidcProviders = map[string]ConfigOidcProvider{}
	// for k, v := range LoadedConfig.OidcProviders {
	// 	configCopy.OidcProviders[k] = v
	// }
	//
	// configCopy.JwtSecret = "***"
	// configCopy.MeilisearchApiKey = "***"
	// for k, v := range configCopy.OidcProviders {
	// 	v.ClientSecret = "***"
	// 	configCopy.OidcProviders[k] = v
	// }
	//
	// slog.Info("loaded config", "config", configCopy)

	// TODO(patrik): I hate this
	oldTag := validate.ErrorTag
	validate.ErrorTag = "mapstructure"
	defer func() {
		validate.ErrorTag = oldTag
	}()

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
