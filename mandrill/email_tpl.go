package mandrill

import (
	"encoding/json"
	"fmt"
)

//------------------------------------------------------------
// Model - template based email
//------------------------------------------------------------

type Email struct {
	Key        string   `json:"key"`
	TplName    string   `json:"template_name,omitempty"`
	TplContent []KeyVal `json:"template_content"`
	Message    Message  `json:"message"`
}

type KeyVal struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Person struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Message struct {
	Html               string              `json:"html"`
	Text               string              `json:"text"`
	Subject            string              `json:"subject"`
	FromEmail          string              `json:"from_email"`
	FromName           string              `json:"from_name"`
	To                 []Person            `json:"to"`
	Headers            map[string]string   `json:"headers,omitempty"`
	Bcc                string              `json:"bcc_address,omitempty"`
	TrackOpens         bool                `json:"track_opens"`
	TrackClicks        bool                `json:"track_clicks"`
	AutoText           bool                `json:"auto_text"`
	PreserveRecipients bool                `json:"preserve_recipients"`
	VarsGlob           []KeyVal            `json:"global_merge_vars,omitempty"`
	Vars               []RcptVars          `json:"merge_vars,omitempty"`
	Attachments        []map[string]string `json:"attachments,omitempty"`
}

type RcptVars struct {
	Rcpt string   `json:"rcpt"`
	Vars []KeyVal `json:"vars"`
}

//------------------------------------------------------------
// Templated Mail
//------------------------------------------------------------

// Create email based on template key.
func NewEmail(tpl, subj string) *Email {
	return &Email{
		TplName: tpl,
		Message: Message{Subject: subj, AutoText: true},
	}
}

// Create email based on provided HTML template.
func NewEmail_Templateless(html, subj string) *Email {
	return &Email{
		Message: Message{Html: html, Subject: subj, AutoText: true},
	}
}

//------------------------------------------------------------
// Sending methods
//------------------------------------------------------------

func (m *Email) SetSender(params map[string]string) {
	m.Message.FromEmail = params["email"]
	m.Message.FromName = params["identity"]
}

func (m *Email) AddTo(params map[string]string) {
	m.Message.To = append(
		m.Message.To,
		Person{params["email"], params["identity"]})
}

func (m *Email) ClearTo() {
	m.Message.To = []Person{}
}

func (m *Email) SetReplyTo(params map[string]string) {
	if m.Message.Headers == nil {
		m.Message.Headers = map[string]string{}
	}
	m.Message.Headers["Reply-To"] = params["email"]
}

func (m *Email) SetBcc(params map[string]string) {
	m.Message.Bcc = params["email"]
}

//------------------------------------------------------------
// Template methods
//------------------------------------------------------------

func (m *Email) AddTplContent(key, val string) {
	m.TplContent = append(m.TplContent, KeyVal{key, val})
}

func (m *Email) AddGlobalVar(key, val string) {
	m.Message.VarsGlob = append(
		m.Message.VarsGlob,
		KeyVal{key, val})
}

func (m *Email) AddVar(rcpt map[string]string, key, val string) {
	email := rcpt["email"]
	// Check if this recipient's values already exist
	for i, v := range m.Message.Vars {
		if v.Rcpt == email {
			m.Message.Vars[i].Vars = append(v.Vars, KeyVal{key, val})
			return
		}
	}
	// Add new recipient
	m.Message.Vars = append(
		m.Message.Vars,
		RcptVars{
			Rcpt: email,
			Vars: []KeyVal{KeyVal{key, val}}})
}

func (m *Email) AddAttachment(mimeType, name, content string) {
	m.Message.Attachments = append(m.Message.Attachments, map[string]string{
		"type":    mimeType,
		"name":    name,
		"content": content,
	})
}

//------------------------------------------------------------
// Sending
//------------------------------------------------------------

// Sends email.
func (m *Email) Send(apikey string) (err error) {
	m.Key = apikey

	// DEBUG
	d, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(d))
	fmt.Println("------ Sending ------")

	var resp interface{}
	if m.TplName != "" {
		resp, err = post(MNDRL_MESSAGES_TEMPLATE, m)
	} else {
		resp, err = post(MNDRL_MESSAGES_TEMPLATELESS, m)
	}

	fmt.Println(resp)
	return
}
