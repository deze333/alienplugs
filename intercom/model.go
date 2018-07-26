package intercom

//------------------------------------------------------------
// Intercom data structures: API generic
//------------------------------------------------------------

// PageParams determine paging information to and from the API
type PageParams struct {
	Page       int64 `json:"page" url:"page,omitempty"`
	PerPage    int64 `json:"per_page" url:"per_page,omitempty"`
	TotalPages int64 `json:"total_pages" url:"-"`
}

//------------------------------------------------------------
// Intercom data structures: User List
//------------------------------------------------------------

// Request parameters: user list
type UserListRequestParams struct {
    Page       int64 `json:"page" url:"page,omitempty"`
    PerPage    int64 `json:"per_page" url:"per_page,omitempty"`
    TotalPages int64 `json:"total_pages" url:"-"`
    Order      string `json:"order" url:"order,omitempty"`
    Sort       string `json:"sort" url:"sort,omitempty"`
    CreatedSince       string `json:"created_since" url:"created_since,omitempty"`
}

// UserList holds a list of Users and paging information
type UserList struct {
	Pages PageParams
	Users []User
	ScrollParam string `json:"scroll_param,omitempty"`
}

// User represents a User within Intercom.
// Not all of the fields are writeable to the API, non-writeable fields are
// stripped out from the request. Please see the API documentation for details.
type User struct {
	ID                     string                 `json:"id,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Anonymous              *bool                  `json:"anonymous,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Pseudonym              string                 `json:"pseudonym,omitempty"`
	//Avatar                 *UserAvatar            `json:"avatar,omitempty"`
	//LocationData           *LocationData          `json:"location_data,omitempty"`
	SignedUpAt             int64                  `json:"signed_up_at,omitempty"`
	RemoteCreatedAt        int64                  `json:"remote_created_at,omitempty"`
	LastRequestAt          int64                  `json:"last_request_at,omitempty"`
	CreatedAt              int64                  `json:"created_at,omitempty"`
	UpdatedAt              int64                  `json:"updated_at,omitempty"`
	SessionCount           int64                  `json:"session_count,omitempty"`
	LastSeenIP             string                 `json:"last_seen_ip,omitempty"`
	//SocialProfiles         *SocialProfileList     `json:"social_profiles,omitempty"`
	UnsubscribedFromEmails *bool                  `json:"unsubscribed_from_emails,omitempty"`
	UserAgentData          string                 `json:"user_agent_data,omitempty"`
	//Tags                   *TagList               `json:"tags,omitempty"`
	//Segments               *SegmentList           `json:"segments,omitempty"`
	//Companies              *CompanyList           `json:"companies,omitempty"`
	CustomAttributes       map[string]interface{} `json:"custom_attributes,omitempty"`
	UpdateLastRequestAt    *bool                  `json:"update_last_request_at,omitempty"`
	NewSession             *bool                  `json:"new_session,omitempty"`
	LastSeenUserAgent      string                 `json:"last_seen_user_agent,omitempty"`
}

//------------------------------------------------------------
// Intercom data structures: Contact List
//------------------------------------------------------------

// ContactList holds a list of Contacts and paging information
type ContactList struct {
	Pages    PageParams
	Contacts []Contact
	ScrollParam string `json:"scroll_param,omitempty"`
}

// Contact represents a Contact within Intercom.
// Not all of the fields are writeable to the API, non-writeable fields are
// stripped out from the request. Please see the API documentation for details.
type Contact struct {
	ID                     string                 `json:"id,omitempty"`
	Email                  string                 `json:"email,omitempty"`
	Phone                  string                 `json:"phone,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	//Avatar                 *UserAvatar            `json:"avatar,omitempty"`
	//LocationData           *LocationData          `json:"location_data,omitempty"`
	LastRequestAt          int64                  `json:"last_request_at,omitempty"`
	CreatedAt              int64                  `json:"created_at,omitempty"`
	UpdatedAt              int64                  `json:"updated_at,omitempty"`
	SessionCount           int64                  `json:"session_count,omitempty"`
	LastSeenIP             string                 `json:"last_seen_ip,omitempty"`
	//SocialProfiles         *SocialProfileList     `json:"social_profiles,omitempty"`
	UnsubscribedFromEmails *bool                  `json:"unsubscribed_from_emails,omitempty"`
	UserAgentData          string                 `json:"user_agent_data,omitempty"`
	//Tags                   *TagList               `json:"tags,omitempty"`
	//Segments               *SegmentList           `json:"segments,omitempty"`
	//Companies              *CompanyList           `json:"companies,omitempty"`
	CustomAttributes       map[string]interface{} `json:"custom_attributes,omitempty"`
	UpdateLastRequestAt    *bool                  `json:"update_last_request_at,omitempty"`
	NewSession             *bool                  `json:"new_session,omitempty"`
}

type contactListParams struct {
	PageParams
	SegmentID string `url:"segment_id,omitempty"`
	TagID     string `url:"tag_id,omitempty"`
	Email     string `url:"email,omitempty"`
}

