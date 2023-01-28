package config

type DatabaseConfig struct {
	Dialect string `hcl:"dialect,label"`
	URL     string `hcl:"url,label"`
}

type AppConfig struct {
	TokenExpiryInSec   int64 `hcl:"token_expiry_seconds,label"`
	RefreshExpiryInSec int64 `hcl:"refresh_expiry_seconds,label"`
}
type Config struct {
	*DatabaseConfig `hcl:"database,block"`
	*AppConfig      `hcl:"app,block"`
}
