# Debug mode
debug = false

# Metrics api enabled?
matrics = false

# Application environment
env = "dev"

# Timezone
timezone = "America/California"

# gRPC TCP listener config
listener {
  host = "0.0.0.0"
  port = 8000
}

# Redis configurations
redis {
  # if you wish to use redis cluster, set this value to true
  # If cluster is true, sentinels is required
  # If cluster is false, addr is required
  cluster     = false
  master_name = "casty"
  addr        = "127.0.0.1:26379"
  sentinels   = [
    "127.0.0.1:26379"
  ]
  pass = "super-secure-password"
  sentinel_pass = "super-secure-sentinels-password"
}

# Database (mongodb) config
db {
  name = "casty"
  host = "localhost"
  port = 27017
  user = "service"
  pass = "super-secure-password"
}

# JWT secrets
jwt {
  access_token {
    # make sure to use a strong secret key
    secret = "random-secret"
    # If you wish to change valid duration of a access_token, change this value
    expires_at {
      type  = "days" # can be [seconds|minutes|hours|days]
      value = 1
    }
  }
  refresh_token {
    # make sure to use a strong secret key
    secret = "random-secret"
    # If you wish to change valid duration of a refresh_token, change this value
    expires_at {
      type  = "days" # can be [seconds|minutes|hours|days]
      value = 7
    }
  }
}

# oauth details
oauth {

  # Let user to register with oauth
  registration_by_oauth = true

  # Google config
  google {
    enabled       = false
    client_id     = ""
    client_secret = ""
    auth_uri      = "https://accounts.google.com/o/oauth2/auth"
    token_uri     = "https://oauth2.googleapis.com/token"
    redirect_uri = "https://casty.ir/oauth/google/callback"
  }
 
  # Spotify config
  spotify {
    enabled       = false
    client_id     = ""
    client_secret = ""
    auth_uri      = "https://accounts.spotify.com/authorize"
    token_uri     = "https://accounts.spotify.com/api/token"
    redirect_uri  = "https://casty.ir/oauth/spotify/callback"
  }

}

# S3 bucket config
s3 {
  endpoint = "127.0.0.1:9000"
  access_key = "secret-access-key"
  secret_key = "secret-key"
}

# Sentry config
sentry {
  enabled = false
  dsn     = "sentry.dsn.here"
}

# Recaptcha config, it can be google or hcaptcha
recaptcha {
  enabled = false
  type    = "hcaptcha"
  secret  = "hcaptcha-secret-token"
}
