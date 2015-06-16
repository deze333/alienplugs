package twilio

/*

example:

tc := sms.TwilioCfg{
    FromPhone: "+XXXX",
    AccountSID: "XXXX",
    AuthToken: "XXXX",
}

err := tc.SMS("+XXXXX", "message")

if err != nil {
    fmt.Print(err) // log error
} else {
    // sent success
}

*/

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiMsgUrl = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
)

//------------------------------------------------------------
// TwilioCfg
//------------------------------------------------------------

type TwilioCfg struct {
	FromPhone  string // outgoing phone
	AccountSID string // twilio accnt id
	AuthToken  string // twilio token
}

//------------------------------------------------------------
// TwilioCfg methods
//------------------------------------------------------------

//------------------------------------------------------------
// Send SMS
//
// silently ignore invalid numbers (numeric only)
// if twilio response status is not 201
//    put response body into err and return
// else return nil assuming success
//------------------------------------------------------------
func (tc *TwilioCfg) SMS(toPhone, body string) (err error) {

	var req *http.Request
	var resp *http.Response

	// Only allow [0...9] characters
	toPhone = cleanPhone(toPhone)

	// append +
	toPhone = "+" + toPhone

	// trim long msg
	if len(body) >= 160 {
		body = body[:155] + "..."
	}

	form := url.Values{
		"To": {
			toPhone,
		},
		"From": {
			tc.FromPhone,
		},
		"Body": {
			body,
		},
	}

	postUrl := fmt.Sprintf(apiMsgUrl, tc.AccountSID)
	req, err = http.NewRequest("POST", postUrl, bytes.NewBufferString(form.Encode()))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(tc.AccountSID, tc.AuthToken)

	client := &http.Client{}
	resp, err = client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		respBody, ioerr := ioutil.ReadAll(resp.Body)

		if ioerr != nil {
			return errors.New("twilio failed sending\nfailed to read http resp")
		} else {
			return errors.New("twilio failed sending: " + bytes.NewBuffer(respBody).String())
		}
	}

	return nil
}

func validPhone(phone string) bool {

	if len(phone) == 0 {
		return false
	}

	for _, k := range phone {
		if k < '0' || k > '9' {
			return false
		}
	}

	return true
}

// Removes all non-numeric characters from the phone.
func cleanPhone(phone string) string {

	var buf bytes.Buffer
	for _, d := range phone {
		if d < '0' || d > '9' {
			continue
		}
		buf.WriteRune(d)
	}

	return buf.String()
}
