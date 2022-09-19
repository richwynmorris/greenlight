package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedJsonValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	stringsSlice := strings.Split(unquotedJsonValue, " ")

	// Sanity Check that the parts are as expected.
	if len(stringsSlice) != 2 || stringsSlice[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Parse the string into int32
	runtimeInt, err := strconv.ParseInt(stringsSlice[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(runtimeInt)

	return nil
}
