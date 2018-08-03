package util

import (
	"bytes"
	"fmt"
)

//------------------------------------------------------------
// Utils: Stringify any object
//------------------------------------------------------------

func DumpToString(name string, kvals []interface{}) string {
	var buf bytes.Buffer

	buf.WriteString("\n")
	buf.WriteString("-----------------------------------\n")
	buf.WriteString(name)
	buf.WriteString("\n")
	buf.WriteString("-----------------------------------\n")

	for i := 0; i < len(kvals); i += 2 {
		buf.WriteString(fmt.Sprint(kvals[i]))
		buf.WriteString(": ")
		buf.WriteString(fmt.Sprint(kvals[i+1]))
		buf.WriteString("\n")
	}

	return buf.String()
}

func DumpToStringLine(name string, kvals []interface{}) string {
	var buf bytes.Buffer

	buf.WriteString(name)
	buf.WriteString(" = {")

	for i := 0; i < len(kvals); i += 2 {
		buf.WriteString(fmt.Sprint(kvals[i]))
		if i+1 == len(kvals) {
			break
		}
		buf.WriteString(": ")
		buf.WriteString(fmt.Sprint(kvals[i+1]))
		buf.WriteString(", ")
	}

	buf.WriteString(" }")
	return buf.String()
}

