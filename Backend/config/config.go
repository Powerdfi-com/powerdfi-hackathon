package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Port int `env:"PORT,required"`
	Jwt  struct {
		Access       string `env:"JWT_ACCESS,required"`
		AdminAccess  string `env:"JWT_ADMIN_ACCESS,required"`
		Refresh      string `env:"JWT_REFRESH,required"`
		AdminRefresh string `env:"JWT_ADMIN_REFRESH,required"`
	}
	DbUri     string `env:"DB_URI,required"`
	Env       string `env:"ENV,required"`
	MasterKey string `env:"ENCRYPTION_KEY,required"`
	Hedera    struct {
		TokenIdUSDC        string `env:"HEDERA_USDC_TOKEN_ID,required"`
		TokenId            string `env:"HEDERA_TOKEN_ID,required"`
		TreasuryAccountId  string `env:"HEDERA_TREASURY_ACCOUNT_ID,required"`
		TreasuryPrivateKey string `env:"HEDERA_TREASURY_PRIVATE_KEY,required"`
	}
	Shufti struct {
		BaseUrl   string `env:"SHUFTI_API_URL,required"`
		ClientId  string `env:"SHUFTI_CLIENT_ID,required"`
		SecretKey string `env:"SHUFTI_SECRET_KEY,required"`
		JourneyId string `env:"SHUFTI_JOURNEY_ID,required"`
	}
}

func LoadConfig() (Config, error) {

	cfg := Config{} // 👈 new instance of `Config`
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	err = env.Parse(&cfg) // 👈 Parse environment variables into `Config`
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
