package types

import (
	"encoding/json"
	"github.com/deadblue/elevengo/internal/util"
	"strconv"
)

// "IntString" uses for JSON field which's type may be int or string, and store it as string.
type IntString string

func (v *IntString) UnmarshalJSON(b []byte) (err error) {
	var s string
	if b[0] == '"' {
		err = json.Unmarshal(b, &s)
	} else {
		var n int
		if err := json.Unmarshal(b, &n); err == nil {
			s = strconv.Itoa(n)
		}
	}
	if err == nil {
		*v = IntString(s)
	}
	return nil
}

// "IntString" uses for JSON field which's type may be string or int, and store it as int.
type StringInt int

func (v *StringInt) UnmarshalJSON(b []byte) (err error) {
	var i int
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			i = util.MustAtoi(s)
		}
	} else {
		err = json.Unmarshal(b, &i)
	}
	if err == nil {
		*v = StringInt(i)
	}
	return nil
}

// "StringInt64" uses for JSON field which's type may be int64 or string, and store it as int64.
type StringInt64 int64

func (v *StringInt64) UnmarshalJSON(b []byte) (err error) {
	var i int64
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			i = util.MustParseInt(s)
		}
	} else {
		err = json.Unmarshal(b, &i)
	}
	if err == nil {
		*v = StringInt64(i)
	}
	return
}

type StringFloat64 float64

func (v *StringFloat64) UnmarshalJSON(b []byte) (err error) {
	var f float64
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			f = util.MustParseFloat(s)
		}
	} else {
		err = json.Unmarshal(b, &f)
	}
	if err == nil {
		*v = StringFloat64(f)
	}
	return
}
