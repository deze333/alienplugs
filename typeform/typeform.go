package typeform

import (
	"fmt"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
)

//------------------------------------------------------------
// Typeform
//------------------------------------------------------------

type Typeform struct {
	AccountKey string // account key
}

//------------------------------------------------------------
// Typeform API Endpoints
//------------------------------------------------------------

const (
	API_RESPONSES = "https://api.typeform.com/forms/%v/responses"
	API_FORM      = "https://api.typeform.com/forms/%v"
)

//------------------------------------------------------------
// API
//------------------------------------------------------------

func NewTypeform(key string) Typeform {
	return Typeform{AccountKey: key}
}

// Retrieves form responses since given date until (optional) date.
func (tf *Typeform) GetFormResponses(formId string, completed bool, since time.Time, until *time.Time) (responses *Responses, err error) {

	// API URL
	var apiUrl *url.URL
	apiUrl, err = url.Parse(fmt.Sprintf(API_RESPONSES, formId))
	if err != nil {
		return
	}

	// Add URL parameters
	parameters := url.Values{}
	parameters.Add("completed", fmt.Sprintf("%v", completed))
	parameters.Add("since", since.Format("2006-01-02T15:04:05"))
	if until != nil {
		parameters.Add("until", until.Format("2006-01-02T15:04:05"))
	} else {
		//parameters.Add("until", time.Now().Format("2006-01-02T15:04:05"))
	}
	apiUrl.RawQuery = parameters.Encode()

	// Send request
	var data []byte
	data, err = tf.sendRequest("GET", apiUrl.String())
	if err != nil {
		return
	}

	responses = &Responses{}
	err = json.Unmarshal(data, responses)
	if err != nil {
		s := fmt.Sprintf("Error unmarshalling JSON: %v", err)
		err = errors.New(s)
		return
	}

	return
}

// Retrieves form data.
func (tf *Typeform) GetForm(formId string) (form *Form, err error) {

	// API URL
	var apiUrl *url.URL
	apiUrl, err = url.Parse(fmt.Sprintf(API_FORM, formId))
	if err != nil {
		return
	}

	// Send request
	var data []byte
	data, err = tf.sendRequest("GET", apiUrl.String())
	if err != nil {
		return
	}

	form = &Form{}
	err = json.Unmarshal(data, form)
	if err != nil {
		s := fmt.Sprintf("Error unmarshalling JSON: %v", err)
		err = errors.New(s)
		return
	}

	return
}


//------------------------------------------------------------
// Private methods
//------------------------------------------------------------

// Sends request to Typeform.
func (tf *Typeform) sendRequest(method, url string) (data []byte, err error) {

	client := http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	req.Header.Add("authorization", fmt.Sprintf("bearer %v", tf.AccountKey))

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// Status OK?
	if resp.StatusCode != http.StatusOK {
		s := fmt.Sprintf("TypeForm server error: Response status: %v / %v", resp.StatusCode, resp.Status)
		err = errors.New(s)
		return
	}

	// Empty response?
	if resp.ContentLength == 0 {
		s := fmt.Sprintf("TypeForm server error: Empty response")
		err = errors.New(s)
		return
	}

	// Read response
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s := fmt.Sprintf("TypeForm response reading error: %v", err)
		err = errors.New(s)
	}

	return
}
