package validator

import "testing"

func TestLangValidator(t *testing.T) {

	// Check For Default Value
	lang := ""
	lang, err := ValidateLang(lang)
	if err != nil && lang != "id" {
		t.Logf("expect default value is id, but got %s", lang)
		t.Fail()
	}

	// Check if the lang is still the same when throw the value into function
	lang = "en"
	lang, err = ValidateLang(lang)
	if err != nil && lang != "en" {
		t.Logf("expect value is id, but got %s", lang)
		t.Fail()
	}

	// Check if the function return error if the chosed lang is not id or en
	lang = "jp"
	_, err = ValidateLang(lang)
	if err == nil {
		t.Logf("expect to be error when chosing language other than id or en")
		t.Fail()
	}
}
