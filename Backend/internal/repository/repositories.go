package repository

type Repositories struct {
	Nonce        NonceRepository
	User         UserRepository
	Asset        AssetRepository
	Listing      ListingRepository
	Stats        StatsRepository
	UserKyc      UserKycRepository
	Category     CategoryRepository
	AssetOwner   AssetOwnerRepository
	Order        OrderRepository
	Trade        TradeRepository
	Activity     ActivityRepository
	Admin        AdminRepository
	Notification NotificationRepository
	UserKYB      UserKYBRepository
}
