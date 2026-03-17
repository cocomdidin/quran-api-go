package validator

import (
	"testing"
)

func TestIDParamValidator(t *testing.T) {
	// Check if id is empty return error value
	id := ""
	id, err := ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id does not contain any")
		t.Fail()
	}

	// Check, Validate Param must return error if id is not a number
	id = "Get"
	id, err = ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id is not number")
		t.Fail()
	}

	// Check, Validate Param must return errror if id is negative value
	id = "-2"
	_, err = ValidateIDParam(id)
	if err == nil {
		t.Logf("expect validation id param to return error if the id is negative")
		t.Fail()
	}

	// Check, Validate Param must not return error if id is not negative value
	id = "1"
	id, err = ValidateIDParam(id)
	if err != nil {
		t.Logf("expect validation id param to accept id value")
		t.Fail()
	}
}

func TestRangeParamValidator(t *testing.T) {
	// Check if ValidateRangeParam return error if from and to is empty string
	err := ValidateRangeParam("", "")
	if err == nil {
		t.Fatal("expect validation \"from\", \"to\" to return error if the from, to does not contain any")
	}

	// Check if ValidateRangeParam return error if from and to is contain not a numberV
	err = ValidateRangeParam("Get", "Post")
	if err == nil {
		t.Fatal("expect validation \"from\", \"to\" to return error if the from, to does not not a number")
	}

	// Check if ValidateRangeParam return error if from and to is negative number
	err = ValidateRangeParam("-2", "5")
	if err == nil {
		t.Fatal("expect validation \"from\", \"to\" to return error if the from, to contain negative value")
	}

	// Check if ValidateRangeParam return error if from value is bigger than to
	err = ValidateRangeParam("10", "5")
	if err == nil {
		t.Fatal("expect validation \"from\", \"to\" to return error if the from value is bigger than to")
	}
}
