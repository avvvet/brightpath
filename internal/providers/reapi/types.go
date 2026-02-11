package reapi

// SearchRequest for Property Search API (ids_only mode)
type SearchRequest struct {
	Size             int    `json:"size"`
	Start            int    `json:"start,omitempty"`
	State            string `json:"state"`
	Auction          bool   `json:"auction"`
	AuctionDateMin   string `json:"auction_date_min"`
	AuctionDateMax   string `json:"auction_date_max"`
	EquityPercentMin int    `json:"equity_percent_min"`
	CorporateOwned   bool   `json:"corporate_owned"`
	PropertyType     string `json:"property_type"`
	IDsOnly          bool   `json:"ids_only"`
}

// SearchResponse for Property Search API (ids_only mode)
type SearchResponse struct {
	Data        []int64 `json:"data"`
	ResultCount int     `json:"resultCount"`
	StatusCode  int     `json:"statusCode"`
}

// PropertyDetailRequest for Property Detail API
type PropertyDetailRequest struct {
	ID int64 `json:"id"`
}

// PropertyDetailResponse for Property Detail API
type PropertyDetailResponse struct {
	Data       PropertyDetailData `json:"data"`
	StatusCode int                `json:"statusCode"`
}

// PropertyDetailData contains the full property details
type PropertyDetailData struct {
	ID                  int64  `json:"id"`
	EstimatedValue      int    `json:"estimatedValue"`
	EstimatedEquity     int    `json:"estimatedEquity"`
	EquityPercent       int    `json:"equityPercent"`
	OpenMortgageBalance int    `json:"openMortgageBalance"`
	PropertyType        string `json:"propertyType"`
	PreForeclosure      bool   `json:"preForeclosure"`
	Auction             bool   `json:"auction"`
	NoticeType          string `json:"noticeType"`
	OwnerOccupied       bool   `json:"ownerOccupied"`
	InvestorBuyer       bool   `json:"investorBuyer"`
	CorporateOwned      bool   `json:"corporateOwned"`

	OwnerInfo    OwnerInfo    `json:"ownerInfo"`
	PropertyInfo PropertyInfo `json:"propertyInfo"`
	AuctionInfo  AuctionInfo  `json:"auctionInfo"`
	TaxInfo      TaxInfo      `json:"taxInfo"`
	Demographics Demographics `json:"demographics"`

	CurrentMortgages []Mortgage `json:"currentMortgages"`
}

// OwnerInfo contains owner details
type OwnerInfo struct {
	Owner1FirstName string      `json:"owner1FirstName"`
	Owner1LastName  string      `json:"owner1LastName"`
	Owner1FullName  string      `json:"owner1FullName"`
	Owner1Type      string      `json:"owner1Type"`
	Owner2FirstName string      `json:"owner2FirstName"`
	Owner2LastName  string      `json:"owner2LastName"`
	Owner2Type      string      `json:"owner2Type"`
	OwnerOccupied   bool        `json:"ownerOccupied"`
	OwnershipLength int         `json:"ownershipLength"`
	MailAddress     MailAddress `json:"mailAddress"`
}

// MailAddress contains mailing address
type MailAddress struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

// PropertyInfo contains property details
type PropertyInfo struct {
	Address          Address `json:"address"`
	Bedrooms         int     `json:"bedrooms"`
	Bathrooms        int     `json:"bathrooms"`
	LivingSquareFeet int     `json:"livingSquareFeet"`
	YearBuilt        int     `json:"yearBuilt"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	HOA              bool    `json:"hoa"`
	Pool             bool    `json:"pool"`
	Fireplace        bool    `json:"fireplace"`
}

// Address contains property address
type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	County  string `json:"county"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

// AuctionInfo contains foreclosure auction details
type AuctionInfo struct {
	AuctionDate          string `json:"auctionDate"`
	AuctionTime          string `json:"auctionTime"`
	AuctionStreetAddress string `json:"auctionStreetAddress"`
	DocumentType         string `json:"documentType"`
	RecordingDate        string `json:"recordingDate"`
	EstimatedBankValue   string `json:"estimatedBankValue"`
	LenderName           string `json:"lenderName"`
	TrusteeFullName      string `json:"trusteeFullName"`
	TrusteePhone         string `json:"trusteePhone"`
	TrusteeAddress       string `json:"trusteeAddress"`
}

// TaxInfo contains tax details
type TaxInfo struct {
	AssessedValue     int    `json:"assessedValue"`
	TaxAmount         string `json:"taxAmount"`
	TaxDelinquentYear string `json:"taxDelinquentYear"`
}

// Demographics contains area demographics
type Demographics struct {
	MedianIncome  string `json:"medianIncome"`
	SuggestedRent string `json:"suggestedRent"`
}

// Mortgage contains mortgage details
type Mortgage struct {
	Amount       int    `json:"amount"`
	LenderName   string `json:"lenderName"`
	LoanType     string `json:"loanType"`
	MaturityDate string `json:"maturityDate"`
	Position     string `json:"position"`
	Term         string `json:"term"`
}

// SkipTraceRequest for Skip Trace API
type SkipTraceRequest struct {
	FirstName         string            `json:"first_name"`
	LastName          string            `json:"last_name"`
	MailAddress       string            `json:"mail_address"`
	MailCity          string            `json:"mail_city"`
	MailState         string            `json:"mail_state"`
	MailZip           string            `json:"mail_zip"`
	MatchRequirements MatchRequirements `json:"match_requirements"`
}

// MatchRequirements specifies what data is required
type MatchRequirements struct {
	Phones bool `json:"phones"`
}

// SkipTraceResponse for Skip Trace API
type SkipTraceResponse struct {
	Output     SkipTraceOutput `json:"output"`
	Match      bool            `json:"match"`
	Credits    int             `json:"credits"`
	StatusCode int             `json:"statusCode"`
}

// SkipTraceOutput contains identity and demographics
type SkipTraceOutput struct {
	Identity     Identity              `json:"identity"`
	Demographics SkipTraceDemographics `json:"demographics"`
}

// Identity contains phones, emails, names
type Identity struct {
	Phones []Phone `json:"phones"`
	Emails []Email `json:"emails"`
	Names  []Name  `json:"names"`
}

// Phone contains phone details
type Phone struct {
	PersonID     string `json:"personId"`
	Phone        string `json:"phone"`
	PhoneDisplay string `json:"phoneDisplay"`
	IsConnected  bool   `json:"isConnected"`
	DoNotCall    bool   `json:"doNotCall"`
	PhoneType    string `json:"phoneType"`
	LastSeen     string `json:"lastSeen"`
}

// Email contains email details
type Email struct {
	PersonID  string `json:"personId"`
	Email     string `json:"email"`
	EmailType string `json:"emailType"`
}

// Name contains person name details
type Name struct {
	PersonID  string `json:"personId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// SkipTraceDemographics contains age and gender
type SkipTraceDemographics struct {
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}
