package hubspot

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

//------------------------------------------------------------
// Constants
//------------------------------------------------------------

const (
    hubFormsUrl = "https://forms.hubspot.com/uploads/form/v2/%s/%s"
)

//------------------------------------------------------------
// Constructor
//------------------------------------------------------------

func Submit(portalId, formId string, form map[string]string) (err error) {

    v := toValues(form)
    url := buildFormsUrl(portalId, formId)

    resp, err := http.PostForm(url, v)

    if err != nil {
        return
    }

    if resp.StatusCode != 204 {
        err = fmt.Errorf("Error submitting HubSpot form to %s. StatusCode: %d, expected 204", url, resp.StatusCode)
    }

    return
}

// Convert a map to url.Values
func toValues(m map[string]string) (vs url.Values) {

    vs = url.Values{}
    for k, v := range m {
        vs.Set(k, v)
    }

    return
}

// Build hubspot url using portalId and formId
func buildFormsUrl(portalId, formId string) string {
    return fmt.Sprintf(hubFormsUrl, portalId, formId)
}

// Convenience function to make a proper HubSpot request.
// hubspotuk is taken from the request cookies
// The resulting map should be filled with the other parameters
// for the form.
func BuildSubmit(pageName string, r *http.Request) (m map[string]string) {

    m = map[string]string{}

    // get context cookie
    hubspotutk, err := r.Cookie("hubspotutk")
    if err != nil {
        return
    }

    // build context
    hubCtx := map[string]string{
        "hutk":      hubspotutk.Value,
        "ipAddress": r.RemoteAddr,
        "pageUrl":   r.URL.Host + r.URL.Path,
        "pageName":  pageName,
    }

    // encode context to json
    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    enc.Encode(hubCtx)

    // place context into result map
    m = map[string]string{
        "hs_context": buf.String(),
    }

    return
}
