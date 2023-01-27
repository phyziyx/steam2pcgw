package main

const (
	APP_NAME = "Steam 2 PCGW Converter"
	VERSION  = "v0.0.74-RC7"
	API_LINK = "https://store.steampowered.com/api/appdetails?appids="
	LOCALE   = "&l=english"
	GH_LINK  = "https://github.com/phyziyx/steam2pcgw"
)

type GenreId int

const (
	Action                GenreId = 1  // Action,
	Strategy              GenreId = 2  // Strategy,
	RPG                   GenreId = 3  // RPG,
	Casual                GenreId = 4  // Casual,
	Racing                GenreId = 9  // Racing,
	Sports                GenreId = 18 // Sports,
	Indie                 GenreId = 23 // Indie,
	Adventure             GenreId = 25 // Adventure,
	Simulation            GenreId = 28 // Simulation,
	MassivelyMultiplayer  GenreId = 29 // Massively Multiplayer,
	FreeToPlay            GenreId = 37 // Free to Play,
	Accounting            GenreId = 50 // Accounting,
	AnimationAndModeling  GenreId = 51 // Animation & Modeling,
	AudioProduction       GenreId = 52 // Audio Production,
	DesignAndIllustration GenreId = 53 // Design & Illustration,
	Education             GenreId = 54 // Education,
	PhotoEditing          GenreId = 55 // Photo Editing,
	SoftwareTraining      GenreId = 56 // Software Training,
	Utilities             GenreId = 57 // Utilities,
	VideoProduction       GenreId = 58 // Video Production,
	WebPublishing         GenreId = 59 // Web Publishing,
	GameDevelopment       GenreId = 60 // Game Development,
	EarlyAccess           GenreId = 70 // Early Access,
	SexualContent         GenreId = 71 // Sexual Content,
	Nudity                GenreId = 72 // Nudity,
	Violent               GenreId = 73 // Violent,
	Gore                  GenreId = 74 // Gore,
	Documentary           GenreId = 81 // Documentary,
	Tutorial              GenreId = 84 // Tutorial
)

type CategoryId int

const (
	Multiplayer              CategoryId = 1  // Multi-player,
	Singleplayer             CategoryId = 2  // Single-player,
	HL2Mods                  CategoryId = 6  // Mods (require HL2),
	VAC                      CategoryId = 8  // Valve Anti-Cheat enabled,
	CoOp                     CategoryId = 9  // Co-op,
	Captions                 CategoryId = 13 // Captions available,
	Commentary               CategoryId = 14 // Commentary available,
	Stats                    CategoryId = 15 // Stats,
	SourceSDK                CategoryId = 16 // Includes Source SDK,
	LevelEditor              CategoryId = 17 // Includes level editor,
	PartialControllerSupport CategoryId = 18 // Partial Controller Support,
	Mods                     CategoryId = 19 // Mods,
	MMO                      CategoryId = 20 // MMO,
	SteamAchievements        CategoryId = 22 // Steam Achievements,
	SteamCloud               CategoryId = 23 // Steam Cloud,
	SharedOrSplitScreen      CategoryId = 24 // Shared/Split Screen,
	SteamLeaderboards        CategoryId = 25 // Steam Leaderboards,
	CrossPlatformMultiplayer CategoryId = 27 // Cross-Platform Multiplayer,
	FullControllerSupport    CategoryId = 28 // Full controller support,
	TradingCards             CategoryId = 29 // Steam Trading Cards,
	Workshop                 CategoryId = 30 // Steam Workshop,
	VRSupport                CategoryId = 31 // VR Support,
	TurnNotifications        CategoryId = 32 // Steam Turn Notifications,
	InAppPurchases           CategoryId = 35 // In-App Purchases,
	OnlineMultiPlayer        CategoryId = 36 // Online Multi-Player,
	LocalMultiPlayer         CategoryId = 37 // Local Multi-Player,
	OnlineCoOp               CategoryId = 38 // Online Co-op,
	LocalCoOp                CategoryId = 39 // Local Co-op,
	SteamVRCollectibles      CategoryId = 40 // SteamVR Collectibles
)
