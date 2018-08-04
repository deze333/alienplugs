package typeform

import (
	"fmt"
	"strings"
	"time"
	"encoding/json"

	"github.com/deze333/alienplugs/util"
)

//------------------------------------------------------------
// TypeForm data structures: Responses API
//------------------------------------------------------------

type Responses struct {
	TotalItems      int                     `json:"total_items"`
	PageCount       int                     `json:"page_count"`
	Items           []*Responses_Item       `json:"items"`
	Hidden          map[string]string       `json:"hidden,omitempty"`
	Calculated      map[string]interface{}  `json:"calculated,omitempty"`
}

type Responses_Item struct {
	LandingId    string              `json:"landing_id"`
	Token        string              `json:"token"`
	LandedAt     *time.Time          `json:"landed_at"`
	SubmittedAt  *time.Time          `json:"submitted_at"`
	Metadata     map[string]string   `json:"metadata"`
	Answers      []Responses_Answer  `json:"answers"`
	Hidden       map[string]string   `json:"hidden,omitempty"`
	//Calculated   map[string]interface{}  `json:"hidden,omitempty"`
}

type Responses_Answer struct {
	Field  Responses_FormField       `json:"field"`
	Type   string                    `json:"type"`

	/* 
	Following fields content will depend on 'type':
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

type Responses_FormField struct {
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
// Methods: Responses_Answer
//------------------------------------------------------------

// Returns either single choice or multiple choices
// as a JSON string. Will return text value is such is present.
func (o *Responses_Answer) ChoiceOrChoicesAsJsonString() string {

	switch o.Type {
	case "text":
		return o.Text
	case "choice":
		return o.Choice.Label
	case "choises":
		return o.Choices.JsonString()
	}

	return ""
}

//------------------------------------------------------------
// Methods: AnswerValue_Choices
//------------------------------------------------------------

// Checks if choices contain given choice.
func (o *AnswerValue_Choices) Contains(s string) bool {

	for _, choice := range o.Labels {
		if choice == s {
			return true
		}
	}

	return false
}

// Converts labels to JSON string, ie ["one", "two", "three"].
func (o *AnswerValue_Choices) JsonString() string {
	
	// Not using JSON Marshal since don't want to encode chars like '<'
	if false {
		var s string
		b, err := json.Marshal(o.Labels)
		if err != nil {
			s = fmt.Sprintf("%v", o.Labels)
		} else {
			s = string(b)
		}

		return s
	}

	// Simple serialization
	var ss []string

	for _, v := range o.Labels {
		ss = append(ss, fmt.Sprintf("\"%v\"", v))
	}

	s := strings.Join(ss, ", ")
	return fmt.Sprintf("[%v]", s)
}

// Converts labels to JSON string, ie ["one", "two", "three"],
// by also applying optional transform like "one" --> "OPTION_A".
func (o *AnswerValue_Choices) JsonStringByMapping(transform map[string]string) string {
	
	var ss []string

	for _, v := range o.Labels {
		if v1, ok := transform[v]; ok {
			v = v1
		}

		ss = append(ss, fmt.Sprintf("\"%v\"", v))
	}

	s := strings.Join(ss, ", ")
	return fmt.Sprintf("[%v]", s)
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
// String: Responses_Item
//------------------------------------------------------------

func (o *Responses_Item) String() string {
	return util.DumpToString("Item", o.Dump())
}

func (o *Responses_Item) Dump() []interface{} {

	ds := []interface{}{
		"LandingId", o.LandingId,
		"Token", o.Token,
		"LandedAt", o.LandedAt,
		"SubmittedAt", o.SubmittedAt,
		"Metadata", o.Metadata,
		"Hidden", o.Hidden,
	}

	for _, answer := range o.Answers {
		for _, ad := range answer.Dump() {
			ds = append(ds, ad)
		}
	}

	return ds
}

//------------------------------------------------------------
// String: Responses_Answer
//------------------------------------------------------------

func (o *Responses_Answer) String() string {
	return util.DumpToString("Answer", o.Dump())
}

func (o *Responses_Answer) Dump() []interface{} {

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
// String: Responses_FormField
//------------------------------------------------------------

func (o *Responses_FormField) String() string {
	return util.DumpToStringLine("Field", o.Dump())
}

func (o *Responses_FormField) Dump() []interface{} {
	return []interface{}{
		"Id", o.Id,
		"Type", o.Type,
		"Ref", o.Ref,
	}
}
