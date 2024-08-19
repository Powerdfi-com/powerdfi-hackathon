package internal

import (
	"github.com/Powerdfi-com/Backend/config"
	"github.com/Powerdfi-com/Backend/external/hederaUtils"
	"github.com/Powerdfi-com/Backend/external/shufti"

	repositories "github.com/Powerdfi-com/Backend/internal/repository"
)

type Application struct {
	Config       config.Config
	Repositories repositories.Repositories
	HederaClient *hederaUtils.Client
	ShuftiClient *shufti.Client
}
