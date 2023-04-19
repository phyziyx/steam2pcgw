package types

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
