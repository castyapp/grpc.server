package tests

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/castyapp/grpc.server/config"
)

var defaultConfig = &config.Map{
	Debug:    false,
	Env:      "dev",
	Timezone: "America/California",
	Redis: config.RedisMap{
		Cluster:    false,
		MasterName: "casty",
		Addr:       "casty.redis:6379",
		Pass:       "super-secure-redis-password",
	},
	DB: config.DBMap{
		Name:       "casty",
		Host:       "casty.db",
		Port:       27017,
		User:       "gotest",
		Pass:       "super-secure-mongodb-password",
		AuthSource: "",
	},
	JWT: config.JWTMap{
		AccessToken: config.JWTToken{
			Secret: "random-secret",
			ExpiresAt: config.JWTExpiresAt{
				Type:  "days",
				Value: 1,
			},
		},
		RefreshToken: config.JWTToken{
			Secret: "random-secret",
			ExpiresAt: config.JWTExpiresAt{
				Type:  "weeks",
				Value: 1,
			},
		},
	},
	Oauth: config.OauthMap{
		RegistrationByOauth: true,
		Google: config.OauthClient{
			Enabled:      false,
			ClientID:     "",
			ClientSecret: "",
			AuthURI:      "https://accounts.google.com/o/oauth2/auth",
			TokenURI:     "https://oauth2.googleapis.com/token",
			RedirectURI:  "https://casty.ir/oauth/google/callback",
		},
		Spotify: config.OauthClient{
			Enabled:      false,
			ClientID:     "",
			ClientSecret: "",
			AuthURI:      "https://accounts.spotify.com/authorize",
			TokenURI:     "https://accounts.spotify.com/api/token",
			RedirectURI:  "https://casty.ir/oauth/spotify/callback",
		},
	},
	S3: config.S3Map{
		Endpoint:  "127.0.0.1:9000",
		AccessKey: "secret-access-key",
		SecretKey: "secret-key",
	},
	Sentry: config.SentryMap{
		Enabled: false,
		Dsn:     "sentry.dsn.here",
	},
	Recaptcha: config.RecaptchaMap{
		Enabled: false,
		Type:    "hcaptcha",
		Secret:  "hcaptcha-secret-token",
	},
}

func TestLoadConfig(t *testing.T) {
	configMap, err := config.LoadFile(filepath.Join(configFileName))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !reflect.DeepEqual(defaultConfig, configMap) {
		t.Fatalf("bad: %#v", configMap)
	}
}
