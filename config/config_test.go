package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

var defaultConfig = &ConfigMap{
	Debug:    false,
	Metrics:  false,
	Env:      "dev",
	Timezone: "America/California",
	Listener: GrpcListener{
		Host: "0.0.0.0",
		Port: 8000,
	},
	Redis: RedisConfig{
		MasterName:   "casty",
		Sentinels:    []string{"127.0.0.1:26379"},
		Pass:         "super-secure-password",
		SentinelPass: "super-secure-sentinels-password",
	},
	DB: DBConfig{
		Name: "casty",
		Host: "localhost",
		Port: 27017,
		User: "service",
		Pass: "super-secure-password",
	},
	JWT: JWTConfig{
		AccessToken: JWTToken{
			Secret: "random-secret",
			ExpiresAt: JWTExpiresAt{
				Type:  "days",
				Value: 1,
			},
		},
		RefreshToken: JWTToken{
			Secret: "random-secret",
			ExpiresAt: JWTExpiresAt{
				Type:  "days",
				Value: 7,
			},
		},
	},
	Oauth: OauthConfig{
		RegistrationByOauth: true,
		Google: OauthClient{
			Enabled:      false,
			ClientID:     "",
			ClientSecret: "",
			AuthUri:      "https://accounts.google.com/o/oauth2/auth",
			TokenUri:     "https://oauth2.googleapis.com/token",
			RedirectUri:  "https://casty.ir/oauth/google/callback",
		},
		Spotify: OauthClient{
			Enabled:      false,
			ClientID:     "",
			ClientSecret: "",
			AuthUri:      "https://accounts.spotify.com/authorize",
			TokenUri:     "https://accounts.spotify.com/api/token",
			RedirectUri:  "https://casty.ir/oauth/spotify/callback",
		},
	},
	S3: S3Config{
		Endpoint:  "127.0.0.1:9000",
		AccessKey: "secret-access-key",
		SecretKey: "secret-key",
	},
	Sentry: SentryConfig{
		Enabled: false,
		Dsn:     "sentry.dsn.here",
	},
	Recaptcha: RecaptchaConfig{
		Enabled: false,
		Type:    "hcaptcha",
		Secret:  "hcaptcha-secret-token",
	},
}

func TestLoadConfig(t *testing.T) {
	configMap, err := LoadFile(filepath.Join("config_test.hcl"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !reflect.DeepEqual(defaultConfig, configMap) {
		t.Fatalf("bad: %#v", configMap)
	}
}
