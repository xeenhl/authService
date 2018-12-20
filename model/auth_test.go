package model

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestAuthError_Error(t *testing.T) {

	err := AuthError{
		ErrorCode: "Error",
		Reason:    "Reasons",
	}

	expected := `{"error": "Error", "reason": "Reasons"}`
	actual := err.Error()

	if !IsEqualJson(actual, expected) {
		t.Errorf("Expected %v json for error but got %v value", expected, actual)
	}

}

func TestAuthError_ToBytes(t *testing.T) {

	err := AuthError{
		ErrorCode: "Error",
		Reason:    "Reasons",
	}

	s := `{"error": "Error", "reason": "Reasons"}`
	expected := []byte(s)
	actual := err.ToBytes()

	if bytes.Equal(actual, expected) {
		t.Errorf(" Bites slice for %v was not as expected", s)
	}
}

func IsEqualJson(s1, s2 string) bool {

	if s1 == s2 {
		return true
	}

	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(s1), &o2)

	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
