package intercom

import (
	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"io/ioutil"
	"time"

	"github.com/google/go-querystring/query"
)

//------------------------------------------------------------
// Intercom
//------------------------------------------------------------

type Intercom struct {
	AccountKey string // account key
}

//------------------------------------------------------------
// Intercom API Endpoints
//------------------------------------------------------------

const (
	API_USERS  = "https://api.intercom.io/users"
	API_EVENTS = "https://api.intercom.io/events"
	API_CONTACTS = "https://api.intercom.io/contacts"
)

//------------------------------------------------------------
// API
//------------------------------------------------------------

func NewIntercom(key string) Intercom {
	return Intercom{AccountKey: key}
}

// Adds an event to a user.
func (ic *Intercom) CreateUserEvent(email, eventName string, metadata map[string]string) (err error) {

	now := time.Now()
	req := map[string]interface{}{
		"event_name": eventName,
		"email":      email,
		"created_at": now.Unix(),
		"metadata":   metadata,
	}

	_, err = ic.sendRequest("POST", API_EVENTS, nil, req)
	return
}

// Creates or updates a user.
func (ic *Intercom) UpsertUser(email, name, typ string) (err error) {

	req := map[string]interface{}{
		"email": email,
		"name":  name,
		"custom_attributes": map[string]string{
			"user_type": typ,
		},
	}

	_, err = ic.sendRequest("POST", API_USERS, nil, req)
	return
}

// Lists users:
// page - page number
// order - "asc", "desc"
// sort - which field to sort by: 
//        created_at, last_request_at, signed_up_at, updated_at
func (ic *Intercom) ListUsers(page int64, order, sort string) (userList UserList, err error) {

	params := UserListRequestParams{
		Page: page,
		Order: order,
		Sort: sort,
	}

	var data []byte
	data, err = ic.sendRequest("GET", API_USERS, params, nil)
	if err != nil {
		return
	}

	//s := string(data)
	//fmt.Println(s)

	err = json.Unmarshal(data, &userList)
	if err != nil {
		err = errors.New(string(data))
	}

	return
}

// Archives user.
func (ic *Intercom) ArchiveUser(user User) (user1 User, err error) {

	url := fmt.Sprintf("%s/%s", API_USERS, user.ID)

	var data []byte
	data, err = ic.sendRequest("DELETE", url, nil, nil)
	if err != nil {
		return
	}

	//s := string(data)
	//fmt.Println(s)

	err = json.Unmarshal(data, &user1)
	if err != nil {
		err = errors.New(string(data))
	}

	return
}

// Lists contacts:
// page - page number
// order - "asc", "desc"
// sort - which field to sort by: 
//        created_at, last_request_at, signed_up_at, updated_at
func (ic *Intercom) ListContacts(page int64, order, sort string) (contactList ContactList, err error) {

	params := UserListRequestParams{
		Page: page,
		Order: order,
		Sort: sort,
	}

	var data []byte
	data, err = ic.sendRequest("GET", API_CONTACTS, params, nil)
	if err != nil {
		return
	}

	//s := string(data)
	//fmt.Println(s)

	err = json.Unmarshal(data, &contactList)
	if err != nil {
		err = errors.New(string(data))
	}

	return
}


// Archives contact.
func (ic *Intercom) ArchiveContact(contact Contact) (contact1 Contact, err error) {

	url := fmt.Sprintf("%s/%s", API_CONTACTS, contact.ID)

	var data []byte
	data, err = ic.sendRequest("DELETE", url, nil, nil)
	if err != nil {
		return
	}

	//s := string(data)
	//fmt.Println(s)

	err = json.Unmarshal(data, &contact1)
	if err != nil {
		err = errors.New(string(data))
	}

	return
}

//------------------------------------------------------------
// Private methods
//------------------------------------------------------------

// Sends request to Intercom.
func (ic *Intercom) sendRequest(method, url string, queryParams interface{}, payload map[string]interface{}) (data []byte, err error) {

	client := http.Client{}

	var req *http.Request
	var resp *http.Response

	if payload != nil {
		// With JSON payload
		var buf bytes.Buffer

		if err = json.NewEncoder(&buf).Encode(payload); err != nil {
			return
		}

		if req, err = http.NewRequest(method, url, &buf); err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")

	} else {
		// Without JSON payload
		if req, err = http.NewRequest(method, url, nil); err != nil {
			return
		}
	}

	req.SetBasicAuth(ic.AccountKey, "")
	req.Header.Set("Accept", "application/json")

	// Optional query parameters
	if queryParams != nil {
		ic.addQueryParams(req, queryParams)
	}

	//fmt.Println(req.Method, req.URL, req.Body, req.ContentLength)

	// Send request
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	// Read response
	var err2 error
	data, err2 = ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return data, err2
	}

	// Error returned?
	if resp.StatusCode >= 400 {
		s := string(data)
		err = errors.New(s)
	}

	return
}

func (ic *Intercom) addQueryParams(req *http.Request, params interface{}) {
	v, _ := query.Values(params)
	req.URL.RawQuery = v.Encode()
}

