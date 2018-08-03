package typeform

import (
	"fmt"
	//"encoding/json"
	"strings"
	"time"

	"github.com/deze333/alienplugs/util"
)

//------------------------------------------------------------
// TypeForm data structures: Responses API
//------------------------------------------------------------

type Responses struct {
	TotalItems      int `json:"total_items"`
	PageCount       int `json:"page_count"`
	Items           []*Item `json:"items"`
	Hidden          map[string]string `json:"hidden,omitempty"`
	Calculated      map[string]interface{} `json:"calculated,omitempty"`
}

type Item struct {
	LandingId    string            `json:"landing_id"`
	Token        string            `json:"token"`
	LandedAt     time.Time         `json:"landed_at"`
	SubmittedAt  time.Time         `json:"submitted_at"`
	Metadata     map[string]string `json:"metadata"`
	Answers      []Answer          `json:"answers"`
}

type Answer struct {
	Field  AnswerField       `json:"field"`
	Type   string            `json:"type"`

	/* 
	Following field will depend on 'type':
	text
	email
	number
	boolean
	date
	choice
	choices
	*/

	Text     string              `json:"text"`
	Email    string              `json:"email"`
	Number   int                 `json:"number"`
	Boolean  bool                `json:"boolean"`
	Date     time.Time           `json:"date"`
	Choice   AnswerValue_Choice  `json:"choice"`
	Choices  AnswerValue_Choices `json:"choices"`
}

type AnswerField struct {
	Id     string            `json:"id"`
	Type   string            `json:"type"`
	Ref    string            `json:"ref"`
}

type AnswerValue_Choice struct {
	Label     string            `json:"label"`
}

type AnswerValue_Choices struct {
	Labels     []string            `json:"labels"`
}

//------------------------------------------------------------
// String: Responses
//------------------------------------------------------------

func (o *Responses) String() string {

	var ss []string
	s := util.DumpToString("Responses", o.Dump())
	ss = append(ss, s)

	for _, item := range o.Items {
		ss = append(ss, item.String())
	}

	return strings.Join(ss, "\n")
}

func (o *Responses) Dump() []interface{} {
	return []interface{}{
		"TotalItems", o.TotalItems,
		"PageCount", o.PageCount,
		"Hidden", o.Hidden,
		"Calculated", o.Calculated,
	}
}

//------------------------------------------------------------
// String: Item
//------------------------------------------------------------

func (o *Item) String() string {
	return util.DumpToString("Item", o.Dump())
}

func (o *Item) Dump() []interface{} {

	ds := []interface{}{
		"LandingId", o.LandingId,
		"Token", o.Token,
		"LandedAt", o.LandedAt,
		"SubmittedAt", o.SubmittedAt,
		"Metadata", o.Metadata,
	}

	for _, answer := range o.Answers {
		for _, ad := range answer.Dump() {
			ds = append(ds, ad)
		}
	}

	return ds
}

//------------------------------------------------------------
// String: Answer
//------------------------------------------------------------

func (o *Answer) String() string {
	return util.DumpToString("Answer", o.Dump())
}

func (o *Answer) Dump() []interface{} {

	ds := []interface{}{
		"Field", o.Field.String(),
		"Type", o.Type,
	}

	switch o.Type {
	case "text": 
		ds = append(ds, "Text")
		ds = append(ds, o.Text)
	case "email": 
		ds = append(ds, "Email")
		ds = append(ds, o.Email)
	case "number": 
		ds = append(ds, "Number")
		ds = append(ds, o.Number)
	case "boolean": 
		ds = append(ds, "Boolean")
		ds = append(ds, o.Boolean)
	case "date": 
		ds = append(ds, "Date")
		ds = append(ds, o.Date)
	case "choice": 
		ds = append(ds, "Choice")
		ds = append(ds, o.Choice)
	case "choices": 
		ds = append(ds, "Choices")
		ds = append(ds, o.Choices)
	default:
		ds = append(ds, "???")
		ds = append(ds, o)

		s := fmt.Sprintf("Not supported answer type. Answer data: %v", o)
		panic(s)
	}

	return ds
}

//------------------------------------------------------------
// String: AnswerField
//------------------------------------------------------------

func (o *AnswerField) String() string {
	return util.DumpToStringLine("Field", o.Dump())
}

func (o *AnswerField) Dump() []interface{} {
	return []interface{}{
		"Id", o.Id,
		"Type", o.Type,
		"Ref", o.Ref,
	}
}

//------------------------------------------------------------
// Code producer
//------------------------------------------------------------

/*
func (o *Responses) GoCodeString() string {

	var ss []string

	s := ` 
	switch question.Id {
`
	ss = append(ss, s)

	fs := `
	case "%v":
		// %v
`

	for _, val := range o.Questions {
		s := fmt.Sprintf(fs, val.Id, val.Text)
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
*/
