package typeform

import (
	"fmt"
	"strings"
)

//------------------------------------------------------------
// TypeForm data structures: Create API
//------------------------------------------------------------

type Form struct {
	Id        string            `json:"id"`
	Title     string            `json:"title"`
	Language  string            `json:"language"`
	Fields    []*Form_Field     `json:"fields"`
	Hidden    []string          `json:"hidden"`

	// Not yet implemented...
	//WelcomeScreens  []Form_WelcomeScreen    `json:"welcome_screens"`
	//ThankyouScreens []Form_ThankyouScreen   `json:"thankyou_screens"`
	//Logic           []Form_Logic            `json:"logic"`
	//Theme           map[string]interface{}  `json:"theme"`
	//Workspace       map[string]interface{}  `json:"workspace"`
	//Links           map[string]interface{}  `json:"_links"`
	//Settings        map[string]interface{}  `json:"settings"`
}

type Form_Field struct {
	Id       string            `json:"id"`
	Ref      string            `json:"ref"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`

	// Not yet implemented...
	//Properties     map[string]interface{}  `json:"properties"`
}

//------------------------------------------------------------
// Code producer
//------------------------------------------------------------

// Produces Go code copy/paste switch code
// to simplify form parsing.
// Note that fields are referred by Id,
// alternatively they can be referred by Ref.
func (f *Form) GoCodeString() string {

	var ss []string

	s := ` 
	switch answer.Field.Id {
`
	ss = append(ss, s)

	fs := `
	case "%v": // %v
		// %v
		// %v
`

	for _, field := range f.Fields {
		s := fmt.Sprintf(fs, field.Id, field.Ref, field.Title, field.Type)
		ss = append(ss, s)
	}

	s  = ` 
	default:
		// TODO: Handle error
	    break
	}
`
	ss = append(ss, s)

	return strings.Join(ss, "")
}
