package gplay

const BaseURL = "https://play.google.com"

type Category string

type Collection string

type Sort int

type Age string

type PermissionGroup int

const (
	CategoryApplication       Category = "APPLICATION"
	CategoryAndroidWear       Category = "ANDROID_WEAR"
	CategoryArtAndDesign      Category = "ART_AND_DESIGN"
	CategoryAutoAndVehicles   Category = "AUTO_AND_VEHICLES"
	CategoryBeauty            Category = "BEAUTY"
	CategoryBooksAndReference Category = "BOOKS_AND_REFERENCE"
	CategoryBusiness          Category = "BUSINESS"
	CategoryComics            Category = "COMICS"
	CategoryCommunication     Category = "COMMUNICATION"
	CategoryDating            Category = "DATING"
	CategoryEducation         Category = "EDUCATION"
	CategoryEntertainment     Category = "ENTERTAINMENT"
	CategoryEvents            Category = "EVENTS"
	CategoryFinance           Category = "FINANCE"
	CategoryFoodAndDrink      Category = "FOOD_AND_DRINK"
	CategoryHealthAndFitness  Category = "HEALTH_AND_FITNESS"
	CategoryHouseAndHome      Category = "HOUSE_AND_HOME"
	CategoryLibrariesAndDemo  Category = "LIBRARIES_AND_DEMO"
	CategoryLifestyle         Category = "LIFESTYLE"
	CategoryMapsAndNavigation Category = "MAPS_AND_NAVIGATION"
	CategoryMedical           Category = "MEDICAL"
	CategoryMusicAndAudio     Category = "MUSIC_AND_AUDIO"
	CategoryNewsAndMagazines  Category = "NEWS_AND_MAGAZINES"
	CategoryParenting         Category = "PARENTING"
	CategoryPersonalization   Category = "PERSONALIZATION"
	CategoryPhotography       Category = "PHOTOGRAPHY"
	CategoryProductivity      Category = "PRODUCTIVITY"
	CategoryShopping          Category = "SHOPPING"
	CategorySocial            Category = "SOCIAL"
	CategorySports            Category = "SPORTS"
	CategoryTools             Category = "TOOLS"
	CategoryTravelAndLocal    Category = "TRAVEL_AND_LOCAL"
	CategoryVideoPlayers      Category = "VIDEO_PLAYERS"
	CategoryWatchFace         Category = "WATCH_FACE"
	CategoryWeather           Category = "WEATHER"
	CategoryGame              Category = "GAME"
	CategoryGameAction        Category = "GAME_ACTION"
	CategoryGameAdventure     Category = "GAME_ADVENTURE"
	CategoryGameArcade        Category = "GAME_ARCADE"
	CategoryGameBoard         Category = "GAME_BOARD"
	CategoryGameCard          Category = "GAME_CARD"
	CategoryGameCasino        Category = "GAME_CASINO"
	CategoryGameCasual        Category = "GAME_CASUAL"
	CategoryGameEducational   Category = "GAME_EDUCATIONAL"
	CategoryGameMusic         Category = "GAME_MUSIC"
	CategoryGamePuzzle        Category = "GAME_PUZZLE"
	CategoryGameRacing        Category = "GAME_RACING"
	CategoryGameRolePlaying   Category = "GAME_ROLE_PLAYING"
	CategoryGameSimulation    Category = "GAME_SIMULATION"
	CategoryGameSports        Category = "GAME_SPORTS"
	CategoryGameStrategy      Category = "GAME_STRATEGY"
	CategoryGameTrivia        Category = "GAME_TRIVIA"
	CategoryGameWord          Category = "GAME_WORD"
	CategoryFamily            Category = "FAMILY"
)

const (
	CollectionTopFree  Collection = "TOP_FREE"
	CollectionTopPaid  Collection = "TOP_PAID"
	CollectionGrossing Collection = "GROSSING"
)

const (
	SortHelpfulness Sort = 1
	SortNewest      Sort = 2
	SortRating      Sort = 3
)

const (
	AgeFiveUnder Age = "AGE_RANGE1"
	AgeSixEight  Age = "AGE_RANGE2"
	AgeNineUp    Age = "AGE_RANGE3"
)

const (
	PermissionGroupCommon PermissionGroup = 0
	PermissionGroupOther  PermissionGroup = 1
)
