package internal

import "github.com/Powerdfi-com/Backend/internal/models"

const KYC_PLATFORM = "shufti"

const (
	HEDERA = "hedera"
	RIPPLE = "ripple"
)

var (
	USDC models.ListingCurrency = "usdc"
)

var ChainNameMap = map[models.Chain]string{
	HEDERA: "hedera",
	RIPPLE: "ripple",
}
var ChainLogoMap = map[models.Chain]string{
	HEDERA: "https://res.cloudinary.com/dxxyljdfd/image/upload/v1721041609/blockchains/hpzx51xbwgacnvprzc5h.png",
	RIPPLE: "https://res.cloudinary.com/dxxyljdfd/image/upload/v1721041607/blockchains/m9hvgguuyoerjcrcicpj.png",
}

var CurrencyNameMap = map[models.ListingCurrency]string{
	USDC: "USDC",
}

var ChainCurrencyListMap = map[models.Chain][]models.ListingCurrency{
	HEDERA: {USDC},
}

var ChainDisabledMap = map[models.Chain]bool{
	HEDERA: true,
}
