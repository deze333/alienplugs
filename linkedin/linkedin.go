package linkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	//"io/ioutil"
)

// Flow:
// 1. Redirect user to RedirectUri, they confirm linkedin. This process generates an AuthToken
// 2. The AuthToken is exchanged for an AccessToken. Posting to _LI_VALIDATE returns the AccessToken
// 3. Requests for profile information must contain AccessToken paramater

const (
	_LI_AUTH_URL     = "https://www.linkedin.com/oauth/v2/authorization?"
	_LI_VALIDATE_URL = "https://www.linkedin.com/oauth/v2/accessToken?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s"
	_LI_PROFILE_URL  = "https://api.linkedin.com/v2/me"
)

// Contains data for the linked in api
type LinkedIn struct {
	ApiKey    string
	ApiSecret string
	Redirect  string
	State     string
}

// LinkedIn member profile data
type MemberProfile struct {
	Id string
	FirstName string
	LastName string
	Photos []PhotoDescriptor
}

// Describes a single profile photo
type PhotoDescriptor struct {
	Width  int
	Height int
	Url    string
}

// Finds photo with size closest to given size.
// Otherwise will return the one with size slightly under.
// Otherwise will return nil.
func (mp *MemberProfile) PhotoLargerThan(width int) (descr *PhotoDescriptor) {

	// No width meants first photo, if available
	if width <= 0 && len(mp.Photos) > 0 {
		photo := mp.Photos[0]
		return &photo
	}
	
	// Find two sizes that are close to given width
	var smallestPositiveDelta int = 10000
	var smallestNegativeDelta int = -10000
	var smallestPositiveIdx int = -1
	var smallestNegativeIdx int = -1

	for i, photo := range mp.Photos {
		delta := photo.Width - width
		if delta == 0 {
			return &photo
		}
		if delta > 0 && delta < smallestPositiveDelta {
			smallestPositiveDelta = delta
			smallestPositiveIdx = i
		} else if delta > smallestNegativeDelta {
			smallestNegativeDelta = delta
			smallestNegativeIdx = i
		}
	}

	if 0 <= smallestPositiveIdx && smallestPositiveIdx < len(mp.Photos) {
		photo := mp.Photos[smallestPositiveIdx]
		return &photo
	} else if 0 <= smallestNegativeIdx && smallestNegativeIdx < len(mp.Photos) {
		photo := mp.Photos[smallestNegativeIdx]
		return &photo
	}

	return
}

// String function.
func (mp *MemberProfile) String() string {

	ss := []string{}

	ss = append(ss, fmt.Sprintf("ID: %v", mp.Id))
	ss = append(ss, fmt.Sprintf("FirstName: %v", mp.FirstName))
	ss = append(ss, fmt.Sprintf("LastName: %v", mp.LastName))
	ss = append(ss, fmt.Sprintf("LastName: %v", mp.LastName))
	ss = append(ss, fmt.Sprintf("Photos: %v", len(mp.Photos)))

	return strings.Join(ss, ", ")
}

// Generate a an AuthUri - the user must be redirected here and accept
// the linkedin request. This process produces an AuthToken which must be
// sent to Validate quickly (under a minute) before linkedin expires it
func (l *LinkedIn) AuthUri(scope ...string) string {

	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", l.ApiKey)
	v.Add("state", l.State)

	if len(scope) > 0 {
		v.Add("scope", strings.Join(scope, " "))
	} else {
		v.Add("scope", "r_liteprofile")
		//v.Add("scope", "r_basicprofile")
		//v.Add("scope", "r_fullprofile")
	}

	// The redirect_uri must be left unescaped

	return fmt.Sprint(_LI_AUTH_URL, v.Encode(), "&redirect_uri="+l.Redirect)
}

// Validate gets a (2 month) token that is used with Get to retreive
// profile information about the user. Entire response is returned, errors
// included (ie: linkedin reject is not an err)
func (l *LinkedIn) ValidateToken(authToken string) (data map[string]interface{}, err error) {

	postUrl := fmt.Sprintf(_LI_VALIDATE_URL, authToken, l.Redirect, l.ApiKey, l.ApiSecret)
	//resp, err := http.PostForm(postUrl, url.Values{})

	req, _ := http.NewRequest("POST", postUrl, nil)
	req.Header.Add("x-li-format", "json")

	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	data = map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)

	return
}

// Get a user profile values from linkedin. If fields is blank the default
// linkedin response contains name and linkedinUri. Otherwise the selected
// fields are requested from linkedin.
func getUserProfile(access_token string) (data map[string]interface{}, err error) {

	v := url.Values{}
	v.Add("oauth2_access_token", access_token)

	// Since linkedin parameters are going to get escaped (http does not like parenthesis
	// in url) its best to make the request with no URL and create one manually.
	// The Opaque value below makes sure that the Get request works with linkedin (otherwise
	// go url-ecapes the parethesies and the other silly characters)
	req, _ := http.NewRequest("GET", "", nil)
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "api.linkedin.com",
		//Opaque:   makeProfileQuery(fields...),
		Opaque:    _LI_PROFILE_URL,
		RawQuery: v.Encode(),
	}
	//req.Header.Add("x-li-format", "json")
	req.Header.Add("X-RestLi-Protocol-Version", "2.0.0")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// DEBUG:
	/*
	fmt.Println()
	fmt.Println("--- LI HTTP REQUEST---")
	fmt.Println(req.URL)
	fmt.Println("--- LI HTTP RESPONSE---")
	fmt.Println("Error    =", err)
	fmt.Println("Response =", resp)
	fmt.Println()
	// If body has been read then it'll become emptied and 
	// JSON decode below will fail
	//fmt.Println("Response String =")
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println(string(bodyBytes))
	*/

	// Parse JSON

	data = map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)

	return
}

// Get a user profile associated with given access token.
func GetUserProfile(access_token string) (profile MemberProfile, err error) {

	// Get member profile

	var profileData map[string]interface{}
	profileData, err = getUserProfile(access_token)
	if err != nil {
		return
	}
	if profileData == nil {
		err = errors.New("LinkedIn returned empty data")
		return
	}

	// FUNC: Gets string value from map, or empty string.
	getStringValue := func(m map[string]interface{}, key string) (s string) {
		if v := m[key]; v != nil {
			return fmt.Sprint(v)
		}
		return
	}

	// FUNC: Gets lite profile field using default localization.
	// map[localized:map[en_US:SOME_VALUE] preferredLocale:map[country:US language:en]]
	getValue := func(key string) (s string) {
		if amap, ok := profileData[key].(map[string]interface{}); ok {
			if amap1, ok1 := amap["preferredLocale"].(map[string]interface{}); ok1 {
				country := getStringValue(amap1, "country")
				language := getStringValue(amap1, "language")
				localization := fmt.Sprintf("%v_%v", language, country)

				if amap2, ok2 := amap["localized"].(map[string]interface{}); ok2 {
					return getStringValue(amap2, localization)
				}
			}
		}

		return
	}

	// Parse member basic information

	if s := getValue("firstName"); s != "" {
		profile.FirstName = s
	}
	if s := getValue("lastName"); s != "" {
		profile.LastName = s
	}

	// Member ID

	var profileId string
	if v := profileData["id"]; v != nil {
		profileId = fmt.Sprint(v)
		profile.Id = profileId
	} else {
		return
	}

	// Get member photo URL

	var photos []PhotoDescriptor
	photos, err = GetProfilePhotos(access_token, profileId)
	if err != nil {
		return
	}

	profile.Photos = photos

	// Positions ?
	GetProfilePositions(access_token, profileId)

	return
}

// Get a user profile photo from linkedin for given ID.
// Will not work for r_liteprofile
func GetProfilePositions(access_token string, id string) (photoUrl string, err error) {

	v := url.Values{}
	v.Add("oauth2_access_token", access_token)
	v.Add("projection", "(id,positions,profilePicture)") // profilePicture added just for testing

	// Since linkedin parameters are going to get escaped (http does not like parenthesis
	// in url) its best to make the request with no URL and create one manually.
	// The Opaque value below makes sure that the Get request works with linkedin (otherwise
	// go url-ecapes the parethesies and the other silly characters)
	req, _ := http.NewRequest("GET", "", nil)
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "api.linkedin.com",
		Opaque:    _LI_PROFILE_URL,
		RawQuery: v.Encode(),
	}
	req.Header.Add("X-RestLi-Protocol-Version", "2.0.0")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// DEBUG:
	/*
	fmt.Println()
	fmt.Println("--- LI HTTP REQUEST---")
	fmt.Println(req.URL)
	fmt.Println("--- LI HTTP RESPONSE---")
	fmt.Println("Error    =", err)
	fmt.Println("Response =", resp)
	fmt.Println()
	// If body has been read then it'll become emptied and 
	// JSON decode below will fail
	//fmt.Println("Response String =")
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println(string(bodyBytes))
	*/

	data := map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)

	return
}

// Get a user profile photo from linkedin for given ID.
// Returns empty array on failure.
func GetProfilePhotos(access_token string, id string) (photos []PhotoDescriptor, err error) {

	v := url.Values{}
	v.Add("oauth2_access_token", access_token)
	v.Add("projection", "(id,profilePicture(displayImage~:playableStreams))")

	// Since linkedin parameters are going to get escaped (http does not like parenthesis
	// in url) its best to make the request with no URL and create one manually.
	// The Opaque value below makes sure that the Get request works with linkedin (otherwise
	// go url-ecapes the parethesies and the other silly characters)
	req, _ := http.NewRequest("GET", "", nil)
	req.URL = &url.URL{
		Scheme:   "https",
		Host:     "api.linkedin.com",
		Opaque:    _LI_PROFILE_URL,
		RawQuery: v.Encode(),
	}
	req.Header.Add("X-RestLi-Protocol-Version", "2.0.0")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// DEBUG:
	/*
	fmt.Println()
	fmt.Println("--- LI HTTP REQUEST---")
	fmt.Println(req.URL)
	fmt.Println("--- LI HTTP RESPONSE---")
	fmt.Println("Error    =", err)
	fmt.Println("Response =", resp)
	fmt.Println()
	*/

	data := map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&data)

	// FUNC: Gets string value from map, or empty string.
	getStringValue := func(m map[string]interface{}, key string) (s string) {
		if v := m[key]; v != nil {
			return fmt.Sprint(v)
		}
		return
	}

	// FUNC: Gets int value from map, or 0.
	getIntValue := func(m map[string]interface{}, key string) (n int) {
		if v, ok := m[key]; ok && v != nil {
			if n, ok = v.(int); ok {
				return
			}
			if f, ok := v.(float64); ok {
				n = int(f)
				return
			}
			n, _ = strconv.Atoi(fmt.Sprint(v))
			return
		}
		return
	}

	// FUNC: Gets map value from map, or nil.
	getMapValue := func(m map[string]interface{}, key string) (m2 map[string]interface{}) {
		if v, ok := m[key].(map[string]interface{}); ok {
			return v
		}
		return
	}

	// FUNC: Gets array value from map, or nil.
	getArrayValue := func(m map[string]interface{}, key string) (arr []interface{}) {
		if v, ok := m[key].([]interface{}); ok {
			return v
		}
		return
	}

	// Parse photo data

	// Build array of all images from:
	// data.profilePicture["displayImage~"].elements[]

	dataProfilePicture := getMapValue(data, "profilePicture")
	if dataProfilePicture == nil {
		return
	}

	dataDisplayImage := getMapValue(dataProfilePicture, "displayImage~")
	if dataDisplayImage == nil {
		return
	}

	dataElements := getArrayValue(dataDisplayImage, "elements")
	if dataElements == nil {
		return
	}

	// For each element get size
	// elements[0].data["com.linkedin.digitalmedia.mediaartifact.StillImage"].displaySize.width
	// and URL
	// elements[0].identifiers[0].identifier

	//photos := []PhotoDescriptors{}

	for _, d := range dataElements {
		if dataElement, ok := d.(map[string]interface{}); ok {
			dataElementData := getMapValue(dataElement, "data")
			if dataElementData == nil {
				continue
			}

			dataStillImage := getMapValue(dataElementData, "com.linkedin.digitalmedia.mediaartifact.StillImage")
			if dataStillImage == nil {
				continue
			}

			dataDisplaySize := getMapValue(dataStillImage, "displaySize")
			if dataDisplaySize == nil {
				continue
			}

			w := getIntValue(dataDisplaySize, "width")
			h := getIntValue(dataDisplaySize, "height")

			dataIdentifiers := getArrayValue(dataElement, "identifiers")
			if dataIdentifiers == nil {
				continue
			}

			var url string
			for _, d1 := range dataIdentifiers {
				if dataIdentifier := d1.(map[string]interface{}); ok {
					if s := getStringValue(dataIdentifier, "identifier"); s != "" {
						url = s
						break
					}
				}
			}

			if url == "" {
				continue
			}

			ph := PhotoDescriptor{ Width: w, Height: h, Url: url }
			photos = append(photos, ph)
		}
	}

	return
}

// Get all the companies listed as current in the linkedin response. For
// proper results the passed argument should be a linkedin response of
// containing "positions" and "headline".
//
// Its up to the implementation to decide which to use.
func GetCurrentCompanies(lresp map[string]interface{}) (curPositions []map[string]string, headline string) {

	var positionsMap map[string]interface{}
	var ok bool

	if headline, ok = lresp["headline"].(string); !ok {
		headline = ""
	}

	if positionsMap, ok = lresp["positions"].(map[string]interface{}); !ok {
		return
	}

	for _, pItem := range positionsMap {
		positions, ok := pItem.([]interface{})

		if !ok {
			continue
		}

		for _, positionRaw := range positions {
			if position, ok := positionRaw.(map[string]interface{}); ok {
				if isCurrent, ok := position["isCurrent"].(bool); ok && isCurrent {
					data := position["company"].(map[string]interface{})
					pos := map[string]string{}

					// Get company name
					if val, ok := data["name"]; ok && val != nil {
						pos["company"] = fmt.Sprint(val)
					}

					// Get title
					if val, ok := position["title"]; ok && val != nil {
						pos["position"] = fmt.Sprint(val)
					}

					curPositions = append(curPositions, pos)

					/* BUGGY previous version:
					   curPositions = append(curPositions, map[string]string{
					       "company":  data["name"].(string),
					       "position": position["title"].(string),
					   })
					*/
				}
			}
		}
	}

	return
}

// Create LinkedIn-formatted profile request. LinkedIn format looks
// like [api-url]/~:(id, picture-url, etc)
func makeProfileQuery(params ...string) string {

	if len(params) == 0 {
		return _LI_PROFILE_URL
	}

	buf := bytes.NewBufferString(_LI_PROFILE_URL + ":(")

	last := len(params) - 1
	comma := []byte(",")
	for _, param := range params[:last] {
		buf.Write([]byte(param))
		buf.Write(comma)
	}

	buf.Write([]byte(params[last]))
	buf.Write([]byte(")"))

	return buf.String()
}

