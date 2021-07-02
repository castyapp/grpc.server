package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/castyapp/grpc.server/core"
	"github.com/hashicorp/hcl"
)

type Map struct {
	Debug     bool         `hcl:"debug"`
	Env       string       `hcl:"env"`
	Timezone  string       `hcl:"timezone"`
	Redis     RedisMap     `hcl:"redis,block"`
	DB        DBMap        `hcl:"db,block"`
	Oauth     OauthMap     `hcl:"oauth,block"`
	S3        S3Map        `hcl:"s3,block"`
	Sentry    SentryMap    `hcl:"sentry,block"`
	JWT       JWTMap       `hcl:"jwt,block"`
	Recaptcha RecaptchaMap `hcl:"recaptcha,block"`
}

type RedisMap struct {
	Cluster      bool     `hcl:"cluster"`
	MasterName   string   `hcl:"master_name"`
	Addr         string   `hcl:"addr"`
	Sentinels    []string `hcl:"sentinels"`
	Pass         string   `hcl:"pass"`
	SentinelPass string   `hcl:"sentinel_pass"`
}

type DBMap struct {
	Name       string `hcl:"name"`
	Host       string `hcl:"host"`
	Port       int    `hcl:"port"`
	User       string `hcl:"user"`
	Pass       string `hcl:"pass"`
	AuthSource string `hcl:"auth_source"`
}

type OauthClient struct {
	Enabled      bool   `hcl:"enabled"`
	ClientID     string `hcl:"client_id"`
	ClientSecret string `hcl:"client_secret"`
	AuthURI      string `hcl:"auth_uri"`
	TokenURI     string `hcl:"token_uri"`
	RedirectURI  string `hcl:"redirect_uri"`
}

type OauthMap struct {
	RegistrationByOauth bool        `hcl:"registration_by_oauth"`
	Google              OauthClient `hcl:"google,block"`
	Spotify             OauthClient `hcl:"spotify,block"`
}

type S3Map struct {
	Endpoint  string `hcl:"endpoint"`
	AccessKey string `hcl:"access_key"`
	SecretKey string `hcl:"secret_key"`
}

type SentryMap struct {
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

type JWTMap struct {
	AccessToken  JWTToken `hcl:"access_token,block"`
	RefreshToken JWTToken `hcl:"refresh_token,block"`
}

type RecaptchaMap struct {
	Enabled bool   `hcl:"enabled"`
	Type    string `hcl:"type"`
	Secret  string `hcl:"secret"`
}

func LoadFile(filename string) (c *Map, err error) {

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
	return ctx.Set("config.map", configMap)
}
