package gplay

type AppCategory struct {
	Name string  `json:"name"`
	ID   *string `json:"id"`
}

type App struct {
	AppID string `json:"appId"`
	URL   string `json:"url"`

	Title   string `json:"title"`
	Summary string `json:"summary"`

	Developer   string `json:"developer"`
	DeveloperID string `json:"developerId"`

	Icon      string   `json:"icon"`
	Score     *float64 `json:"score"`
	ScoreText *string  `json:"scoreText"`

	PriceText *string  `json:"priceText"`
	Free      *bool    `json:"free"`
	Currency  *string  `json:"currency"`
	Price     *float64 `json:"price"`

	Description     *string `json:"description"`
	DescriptionHTML *string `json:"descriptionHTML"`

	Installs    *string `json:"installs"`
	MinInstalls *int64  `json:"minInstalls"`
	MaxInstalls *int64  `json:"maxInstalls"`

	Ratings *int64 `json:"ratings"`
	Reviews *int64 `json:"reviews"`

	Histogram map[string]int64 `json:"histogram"`

	OriginalPrice   *float64 `json:"originalPrice"`
	DiscountEndDate *string  `json:"discountEndDate"`

	Available *bool   `json:"available"`
	OffersIAP *bool   `json:"offersIAP"`
	IAPRange  *string `json:"IAPRange"`
	Size      *string `json:"size"`

	AndroidVersion     *string `json:"androidVersion"`
	AndroidVersionText *string `json:"androidVersionText"`
	AndroidMaxVersion  *string `json:"androidMaxVersion"`

	DeveloperInternalID       *string `json:"developerInternalID"`
	DeveloperEmail            *string `json:"developerEmail"`
	DeveloperWebsite          *string `json:"developerWebsite"`
	DeveloperAddress          *string `json:"developerAddress"`
	DeveloperLegalName        *string `json:"developerLegalName"`
	DeveloperLegalEmail       *string `json:"developerLegalEmail"`
	DeveloperLegalAddress     *string `json:"developerLegalAddress"`
	DeveloperLegalPhoneNumber *string `json:"developerLegalPhoneNumber"`
	PrivacyPolicy             *string `json:"privacyPolicy"`

	Genre   *string `json:"genre"`
	GenreID *string `json:"genreId"`

	Categories []AppCategory `json:"categories"`

	HeaderImage  *string  `json:"headerImage"`
	Screenshots  []string `json:"screenshots"`
	Video        *string  `json:"video"`
	VideoImage   *string  `json:"videoImage"`
	PreviewVideo *string  `json:"previewVideo"`

	ContentRating            *string `json:"contentRating"`
	ContentRatingDescription *string `json:"contentRatingDescription"`
	AdSupported              *bool   `json:"adSupported"`

	Released *string `json:"released"`
	Updated  *int64  `json:"updated"`
	Version  *string `json:"version"`

	RecentChanges *string  `json:"recentChanges"`
	Comments      []string `json:"comments"`

	Preregister           *bool `json:"preregister"`
	EarlyAccessEnabled    *bool `json:"earlyAccessEnabled"`
	IsAvailableInPlayPass *bool `json:"isAvailableInPlayPass"`
}

type Review struct {
	ID        string           `json:"id"`
	UserName  string           `json:"userName"`
	UserImage string           `json:"userImage"`
	Date      string           `json:"date"`
	Score     int64            `json:"score"`
	ScoreText string           `json:"scoreText"`
	URL       string           `json:"url"`
	Title     *string          `json:"title"`
	Text      string           `json:"text"`
	ReplyDate *string          `json:"replyDate"`
	ReplyText *string          `json:"replyText"`
	Version   *string          `json:"version"`
	ThumbsUp  *int64           `json:"thumbsUp"`
	Criterias []ReviewCriteria `json:"criterias"`
}

type ReviewCriteria struct {
	Criteria string `json:"criteria"`
	Rating   *int64 `json:"rating"`
}

type ReviewsResult struct {
	Data                []Review `json:"data"`
	NextPaginationToken *string  `json:"nextPaginationToken"`
}

type PermissionItem struct {
	Permission string `json:"permission"`
	Type       string `json:"type"`
}

type PermissionsResult struct {
	Short bool             `json:"short"`
	Items []PermissionItem `json:"items"`
	Names []string         `json:"names"`
}

type DataSafetyEntry struct {
	Data     string `json:"data"`
	Optional bool   `json:"optional"`
	Purpose  string `json:"purpose"`
	Type     string `json:"type"`
}

type SecurityPractice struct {
	Practice    string `json:"practice"`
	Description string `json:"description"`
}

type DataSafetyResult struct {
	SharedData        []DataSafetyEntry  `json:"dataShared"`
	CollectedData     []DataSafetyEntry  `json:"dataCollected"`
	SecurityPractices []SecurityPractice `json:"securityPractices"`
	PrivacyPolicyURL  *string            `json:"privacyPolicyUrl"`
}
