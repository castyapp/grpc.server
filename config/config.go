package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/castyapp/grpc.server/core"
	"github.com/hashicorp/hcl"
)

type ConfigMap struct {
	Debug     bool            `hcl:"debug"`
	Env       string          `hcl:"env"`
	Metrics   bool            `hcl:"metrics"`
	Timezone  string          `hcl:"timezone"`
	Redis     RedisConfig     `hcl:"redis,block"`
	DB        DBConfig        `hcl:"db,block"`
	Oauth     OauthConfig     `hcl:"oauth,block"`
	S3        S3Config        `hcl:"s3,block"`
	Sentry    SentryConfig    `hcl:"sentry,block"`
	JWT       JWTConfig       `hcl:"jwt,block"`
	Recaptcha RecaptchaConfig `hcl:"recaptcha,block"`
}

type RedisConfig struct {
	Cluster      bool     `hcl:"cluster"`
	MasterName   string   `hcl:"master_name"`
	Addr         string   `hcl:"addr"`
	Sentinels    []string `hcl:"sentinels"`
	Pass         string   `hcl:"pass"`
	SentinelPass string   `hcl:"sentinel_pass"`
}

type DBConfig struct {
	Name string `hcl:"name"`
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
	User string `hcl:"user"`
	Pass string `hcl:"pass"`
}

type OauthClient struct {
	Enabled      bool   `hcl:"enabled"`
	ClientID     string `hcl:"client_id"`
	ClientSecret string `hcl:"client_secret"`
	AuthUri      string `hcl:"auth_uri"`
	TokenUri     string `hcl:"token_uri"`
	RedirectUri  string `hcl:"redirect_uri"`
}

type OauthConfig struct {
	RegistrationByOauth bool        `hcl:"registration_by_oauth"`
	Google              OauthClient `hcl:"google,block"`
	Spotify             OauthClient `hcl:"spotify,block"`
}

type S3Config struct {
	Endpoint  string `hcl:"endpoint"`
	AccessKey string `hcl:"access_key"`
	SecretKey string `hcl:"secret_key"`
}

type SentryConfig struct {
	Enabled bool   `hcl:"enabled"`
	Dsn     string `hcl:"dsn"`
}

type JWTExpiresAt struct {
	Type  string `hcl:"type"`
	Value int    `hcl:"value"`
}

type JWTToken struct {
	Secret    string       `hcl:"secret"`
	ExpiresAt JWTExpiresAt `hcl:"expires_at,block"`
}

func (t JWTToken) GetSecretAtBytes() []byte {
	return []byte(t.Secret)
}

func (t JWTToken) GetExpireDuration() time.Duration {
	switch t.ExpiresAt.Type {
	case "days":
		return (time.Hour * 24) * time.Duration(t.ExpiresAt.Value)
	case "weeks":
		aweek := (time.Hour * 24) * 7
		return aweek * time.Duration(t.ExpiresAt.Value)
	case "minutes":
		return time.Minute * time.Duration(t.ExpiresAt.Value)
	case "seconds":
		return time.Second * time.Duration(t.ExpiresAt.Value)
	case "hours":
		return time.Hour * time.Duration(t.ExpiresAt.Value)
	}
	return 0
}

type JWTConfig struct {
	AccessToken  JWTToken `hcl:"access_token,block"`
	RefreshToken JWTToken `hcl:"refresh_token,block"`
}

type RecaptchaConfig struct {
	Enabled bool   `hcl:"enabled"`
	Type    string `hcl:"type"`
	Secret  string `hcl:"secret"`
}

func LoadFile(filename string) (c *ConfigMap, err error) {

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, err
	}

	// Build up the result
	if err := hcl.DecodeObject(&c, obj); err != nil {
		return nil, err
	}

	return
}

func Provider(ctx *core.Context) error {
	configFilePath := ctx.MustGetString("config.filepath")
	configMap, err := LoadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}
	ctx.Set("config.map", configMap)
	return nil
}
