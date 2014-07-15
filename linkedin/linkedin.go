package linkedin

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"
)

// Flow:
// 1. Redirect user to RedirectUri, they confirm linkedin. This process generates an AuthToken
// 2. The AuthToken is exchanged for an AccessToken. Posting to _LI_VALIDATE returns the AccessToken
// 3. Requests for profile information must contain AccessToken paramater

const (
    _LI_AUTH_URL     = "https://www.linkedin.com/uas/oauth2/authorization?"
    _LI_VALIDATE_URL = "https://www.linkedin.com/uas/oauth2/accessToken?grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s"
    _LI_PROFILE_URL  = "//api.linkedin.com/v1/people/~"
)

// Contains data for the linked in api
type LinkedIn struct {
    ApiKey    string
    ApiSecret string
    Redirect  string
    State     string
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
    }

    // the redirect_uri must be left unescaped

    return fmt.Sprint(_LI_AUTH_URL, v.Encode(), "&redirect_uri="+l.Redirect)
}

// Validate gets a (2 month) token that is used with Get to retreive
// profile information about the user. Entire response is returned, errors
// included (ie: linkedin reject is not an err)
func (l *LinkedIn) ValidateToken(authToken string) (data map[string]interface{}, err error) {

    postUrl := fmt.Sprintf(_LI_VALIDATE_URL, authToken, l.Redirect, l.ApiKey, l.ApiSecret)
    println(postUrl)
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
func Get(access_token string, fields ...string) (data map[string]interface{}, err error) {

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
        Opaque:   makeProfileQuery(fields...),
        RawQuery: v.Encode(),
    }
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
