package alienplugs

import (
	"fmt"
	"os"
	"testing"

	"github.com/deze333/alienplugs/mandrill"
	"github.com/deze333/skini"
)

//------------------------------------------------------------
// Mandrill
//------------------------------------------------------------

type MandrillParams struct {
	ApiKey    string
	TplKeys   []string
	Sender    map[string]string
	Recipient map[string]string
	Support   map[string]string

	Tpl      map[string]map[string]string
	VarsGlob map[string]map[string]string
	Vars     map[string]map[string]string
}

// Mandrill testing
func TestMandrill(t *testing.T) {
	var err error

	// Load parameters
	var mp MandrillParams
	// Try private data first, if not use public
	fname := "private_mandrill.ini"
	if _, err := os.Stat(fname); err == nil {
		err = skini.ParseFile(&mp, fname)
	} else {
		fname = "public_mandrill.ini"
		if _, err := os.Stat(fname); err == nil {
			err = skini.ParseFile(&mp, fname)
		} else {
			t.Error(fmt.Sprint("No suitable parameters found, exiting"))
		}
	}

	if err != nil {
		t.Error(fmt.Sprintf("Error parsing %v: %v", fname, err))
	}

	// Test ping-pong
	//_, err = mandrill.New(API_KEY)
	//if err != nil {
	//    t.Error(fmt.Sprintf("Error creating Mandrill connection: %v", err))
	//}

	// User info
	//resp, err := mdrl.UserInfo()
	//fmt.Println(resp)

	// Send email via template
	sender := map[string]string{
		"email":    mp.Sender["email"],
		"identity": mp.Sender["identity"]}

	recipient := map[string]string{
		"email":    mp.Recipient["email"],
		"identity": mp.Recipient["identity"]}

	bcc := map[string]string{
		"email":    mp.Support["email"],
		"identity": mp.Support["identity"]}

	// Get first template
	tpl := mp.TplKeys[0]

	// Option A:
	// Send via template
	//mm := mandrill.NewEmail(mp.Tpl[tpl]["id"], mp.Tpl[tpl]["subj"])

	// Option B:
	// Send via provided HTML template
	mm := mandrill.NewEmail_Templateless("<p>*|name|*</p><p>Inline HTML text.</p>", mp.Tpl[tpl]["subj"])

	mm.SetSender(sender)
	mm.AddTo(recipient)
	mm.SetReplyTo(recipient)

	if mp.Support["active"] == "true" {
		mm.SetBcc(bcc)
	}

	// Replaces the whole body section with HTML code
	//mm.AddTplContent("brief", mp.Vars[tpl]["brief"])

	// Template merge vars
	for k, v := range mp.Vars[tpl] {
		mm.AddVar(recipient, k, v)
	}

	err = mm.Send(mp.ApiKey)
	if err != nil {
		t.Error(fmt.Sprintf("Error sending Mandrill email: %v", err))
	}

}
