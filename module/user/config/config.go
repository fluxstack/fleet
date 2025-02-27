package config

type Config struct {
	Auth struct {
		JWTSecret string `json:"jwt_secret" yaml:"jwt_secret"`
	} `json:"auth"`
}
