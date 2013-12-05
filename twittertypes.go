package httpstream

import (
	//"encoding/json"
	"net/url"
	"regexp"
)

type User struct {
	ID                        *Int64Nullable
	IDStr                     StringNullable `json:"id_str"` // "id_str":"608729011",
	Name                      string
	ScreenName                string         `json:"screen_name"`
	ContributorsEnabled       bool           `json:"contributors_enabled"`
	CreatedAt                 string         `json:"created_at"`
	Description               StringNullable `json:"description"`
	FavouritesCount           int            `json:"favourites_count"`
	Followerscount            int            `json:"followers_count"`
	Following                 *BoolNullable  // "following":null,
	Friendscount              int            `json:"friends_count"`
	GeoEnabled                bool
	Lang                      string
	Location                  StringNullable
	ListedCount               int            `json:"listed_count"`
	Notifications             StringNullable //"notifications":null,
	ProfileTextColor          string
	ProfileLinkColor          string
	ProfileBackgroundImageURL string
	ProfileBackgroundColor    string
	ProfileSidebarFillColor   string
	ProfileImageURL           string
	ProfileSidebarBorderColor string
	ProfileBackgroundTile     bool
	Protected                 bool
	StatusesCount             int `json:"statuses_count"`
	TimeZone                  StringNullable
	URL                       StringNullable // "url":null
	UtcOffset                 *IntNullable   // "utc_offset":null,
	Verified                  bool
	ShowAllInlineMedia        *BoolNullable `json:"show_all_inline_media"`
	RawBytes                  []byte
	//"default_profile":false,
	//"follow_request_sent":null,
	//"is_translator":false,
	//"profile_use_background_image":true,
	//"default_profile_image":false,
}

type Tweet struct {
	Text                string
	Entities            Entity
	Favorited           bool
	Source              string
	Contributors        []Contributor
	Coordinates         *Coordinate
	InReplyToScreenName StringNullable
	InReplyToStatusID   *Int64Nullable
	InReplyToUserID     *Int64Nullable
	ID                  *Int64Nullable
	IDStr               string
	CreatedAt           string
	RetweetCount        int32
	Retweeted           *BoolNullable
	PossiblySensitive   *BoolNullable
	User                *User
	RawBytes            []byte
	Truncated           *BoolNullable
	Place               *Place // "place":null,
	//Geo                     string   // deprecated
	//RetweetedStatus         Tweet `json:"retweeted_status"`
}

func (t *Tweet) URLs() []string {
	if len(t.Entities.URLs) > 0 {
		urls := make([]string, 0)
		for _, u := range t.Entities.URLs {
			if len(string(u.ExpandedURL)) > 0 {
				if eu, err := url.QueryUnescape(string(u.ExpandedURL)); err == nil {
					urls = append(urls, eu)
				}
			}
		}
		return urls
	}
	return nil
}

func (t *Tweet) Hashes() []string {
	if len(t.Entities.Hashtags) > 0 {
		tags := make([]string, 0)
		for _, t := range t.Entities.Hashtags {
			tags = append(tags, t.Text)
		}
		return tags
	}
	return nil
}

// Return a list of usernames found in the tweet entity mentions
func (t *Tweet) Mentions() []string {
	if len(t.Entities.UserMentions) > 0 {
		users := make([]string, 0)
		for _, m := range t.Entities.UserMentions {
			users = append(users, m.ScreenName)
		}
		return users
	}
	return nil
}

// Create a nullable coordinates, as the data comes across like so:
//    "coordinates":null,
type Coordinate struct {
	Coordinates []float64
	Type        string
}

type Place struct {
	Attributes  interface{}
	Bounding    BoundingBox `json:"bounding_box"`
	Country     string      `json:"country"`
	CountryCode string      `json:"country_code"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	PlaceType   string      `json:"place_type"`
	URL         string      `json:"url"`
}

// Location bounding box of coordinates
type BoundingBox struct {
	Coordinates [][][]float64
	Type        string // "Polygon"
}

/*
func (c *Coordinate) UnmarshalJSON(data []byte) error {
	// do we need this, can't we just use pointer?
	if len(data) > 0 {
		m := make(map[string]interface{})
		if err := json.Unmarshal(data, m); err == nil {
			if co, ok := m["coordinates"]; ok {
				if cof, ok := co.([]float64); ok {
					c.Coordinates = cof
				}
			}
			if ty, ok := m["type"]; ok {
				if tys, ok := ty.(string); ok {
					c.Type = tys
				}
			}
		}
	}
	return nil
}
*/
type Contributor struct {
	ID         int64
	IDStr      string
	ScreenName string
}

type SiteStreamMessage struct {
	ForUser int64
	Message Tweet
}

type Event struct {
	Target    User
	Source    User
	CreatedAt string
	Event     string
}

type Entity struct {
	Hashtags     []Hashtag
	URLs         []TwitterURL
	UserMentions []Mention
	Media        []Media
}

type Hashtag struct {
	Text    string
	Indices []int
}

// A twitter url
//  "urls":[{"indices":[123,136],"url":"http:\/\/t.co\/a","display_url":null,"expanded_url":null}]
type TwitterURL struct {
	URL         string
	ExpandedURL StringNullable // may be null
	DisplayURL  StringNullable // may be null if it gets chopped off after t.co because of shortenring
	Indices     []int
}
type Mention struct {
	ScreenName string
	Name       StringNullable // No idea why this could be null, if a username gets mentioned that doesn't exist?
	ID         *Int64Nullable
	IDStr      string
	Indices    []int
}

type Media struct {
	ID             int64
	IDStr          string
	DisplayURL     string
	ExpandedURL    string
	Indices        []int
	MediaURL       string
	MediaURL_https string
	URL            string
	Type           string
	ScreenName     string
	Sizes          Sizes
}

type Sizes struct {
	Large  Dimensions
	Medium Dimensions
	Small  Dimensions
	Thumb  Dimensions
}

type Dimensions struct {
	W      int
	Resize string
	H      int
}

type FriendList struct {
	Friends []int64
}

/*
The twitter stream contains non-tweets (deletes)

{"delete":{"status":{"user_id_str":"36484472","id_str":"191029491823423488","user_id":36484472,"id":191029491823423488}}}
{"delete":{"status":{"id_str":"191184618165256194","id":191184618165256194,"user_id":355665960,"user_id_str":"355665960"}}}
{"delete":{"status":{"id_str":"172129790210482176","id":172129790210482176,"user_id_str":"499324766","user_id":499324766}}}
{"delete":{"status":{"user_id_str":"366839894","user_id":366839894,"id_str":"116974717763719168","id":116974717763719168}}}
{"delete":{"status":{"user_id_str":"382739413","id":191184546841112579,"user_id":382739413,"id_str":"191184546841112579"}}}
{"delete":{"status":{"user_id_str":"388738304","id_str":"123723751366987776","id":123723751366987776,"user_id":388738304}}}
{"delete":{"status":{"user_id_str":"156157535","id_str":"190608148829179907","id":190608148829179907,"user_id":156157535}}}

*/
// a function to filter out the delete messages
func OnlyTweetsFilter(handler func([]byte)) func([]byte) {
	delTw := regexp.MustCompile(`"delete"`)
	return func(line []byte) {
		if delTw.Find(line) != nil {
			handler(line)
		}
	}
}
