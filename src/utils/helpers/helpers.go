package helpers

import (
	"strings"
	"time"
)

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDBLayout   = "2006-01-02 15:04:05"
)

type newWhereT struct {
	fields             []whereField
	customFields       []whereField
	wherePrefix        string
	appendWhereAtStart bool
}

type whereField struct {
	query string
	args  []interface{}
}

func NewWhere() *newWhereT {
	return &newWhereT{
		wherePrefix:        "WHERE ",
		appendWhereAtStart: false,
	}
}

func (i *newWhereT) AppendWhereAtStart() *newWhereT {
	i.appendWhereAtStart = true
	return i
}

func (i *newWhereT) ChangePrefix(prefix string) *newWhereT {
	i.wherePrefix = prefix
	return i
}

func (i *newWhereT) Where(query string, args ...interface{}) *newWhereT {
	i.fields = append(i.fields, whereField{
		query: query,
		args:  args,
	})
	return i
}

func (i *newWhereT) CustomWhere(query string, args ...interface{}) *newWhereT {
	i.customFields = append(i.customFields, whereField{
		query: query,
		args:  args,
	})
	return i
}

func (i *newWhereT) String() (string, []interface{}) {
	var where []string
	var args []interface{}

	for index := range i.fields {
		el := i.fields[index]

		where = append(where, el.query)
		args = append(args, el.args...)
	}

	str := strings.Join(where, " AND ")
	if len(i.customFields) > 0 {
		str += " "
	}

	// Add custom where
	for index := range i.customFields {
		el := i.customFields[index]

		str += el.query
		args = append(args, el.args...)
	}

	if len(where) > 0 && i.appendWhereAtStart {
		return i.wherePrefix + str, args
	}
	return str, args
}

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDBLayout)
}
