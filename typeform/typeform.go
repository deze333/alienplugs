package typeform

import (
	"fmt"
	"encoding/json"
	"strings"
)

//------------------------------------------------------------
// TypeForm data structures
//------------------------------------------------------------

type Output struct {
	Questions []*Question `json:"questions"`
	Responses []*Response `json:"responses"`
}

type Question struct {
	Id      string `json:"id"`
	Text    string `json:"question"`
	FieldId int    `json:"field_id"`
}

type Response struct {
	Completed string            `json:"completed"`
	Metadata  map[string]string `json:"metadata"`
	Token     string            `json:"token"`
	Hidden    map[string]string `json:"hidden"`
	Answers   map[string]string `json:"answers"`
}

//------------------------------------------------------------
// Structure methods
//------------------------------------------------------------

func (o *Output) QuestionsString() string {

	var ss []string
	for i, val := range o.Questions {
		s := fmt.Sprintf("%v: %v", i, val)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

/// Produces Go code in the form of switch statements.
func (o *Output) QuestionsAsCodeString() string {

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

func (o *Output) ResponsesString() string {

	var ss []string
	for i, val := range o.Responses {
		s := fmt.Sprintf("%v: %v", i, val)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

func (q *Question) String() string {
	return fmt.Sprintf("id: %v\nfield_id: %v\ntext: %v\n", q.Id, q.FieldId, q.Text)
}

func (r *Response) String() string {

	var ss []string
	for key, val := range r.Answers {
		s := fmt.Sprintf("%v: %v", key, val)
		ss = append(ss, s)
	}

	return strings.Join(ss, "\n")
}

//------------------------------------------------------------
// Helpers
//------------------------------------------------------------

func AddToArray(target map[string]string, arrayName string, val string) {

	var ss []string
	if target[arrayName] != "" {
		err := json.Unmarshal([]byte(target[arrayName]), &ss)
		if err != nil {
			fmt.Println("ERROR WHILE PARSING EXISTING VALUES")
			return
		}
	}

	ss = append(ss, val)

	result, err := json.Marshal(ss)
	if err != nil {
		fmt.Println("ERROR WHILE CONVERTING TO JSON")
		return
	}

	target[arrayName] = string(result)
	return
}
