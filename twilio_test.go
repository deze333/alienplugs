package alienplugs

import (
	"fmt"
	"os"
	"testing"

	"github.com/deze333/alienplugs/twilio"
	"github.com/deze333/skini"
)

type TwilioParams struct {
	AccountSID string
	AuthToken  string
	FromPhone  string
	ToPhone    string
}

func (t *TwilioParams) ok() bool {
	return t.AccountSID != "" && t.AuthToken != "" && t.FromPhone != "" && t.ToPhone != ""
}

func TestTwilio(t *testing.T) {
	var err error

	// Load parameters
	var tp TwilioParams
	// Try private data first, if not use public
	fname := "private_twilio.ini"
	if _, err := os.Stat(fname); err == nil {
		err = skini.ParseFile(&tp, fname)
	} else {
		fname = "public_twilio.ini"
		if _, err := os.Stat(fname); err == nil {
			err = skini.ParseFile(&tp, fname)
		} else {
			t.Fatal(fmt.Sprint("No suitable parameters found, exiting"))
		}
	}

	if err != nil {
		t.Fatal(fmt.Sprintf("Error parsing %v: %v", fname, err))
	}

	if !tp.ok() {
		t.Fatal(fmt.Sprintf("Error parsing %v: missing parameters for test: \n%+v", fname, tp))
	}

	tcfg := twilio.TwilioCfg{
		FromPhone:  tp.FromPhone,
		AccountSID: tp.AccountSID,
		AuthToken:  tp.AuthToken,
	}

	err = tcfg.SMS(tp.ToPhone, "unit test message")

	if err != nil {
		t.Fatal(fmt.Sprintf("Failed sending twilio message:\n%s", err))
	}
}
