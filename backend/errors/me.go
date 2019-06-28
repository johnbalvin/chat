package errors

import (
	"errors"
)

//NotTheSame used when parameters are not equals
var NotTheSame = errors.New("Error: Fail to compare, parameters are not equals")

//NotTheSame trying to create new data with an id already in database
var Duplicated = errors.New("Error: Fail to compare, parameters are not equals")

//NotAllow when user is not allow
var NotAllow = errors.New("Error: Not alow to apply")

//TimesOut when time is out
var TimesOut = errors.New("Error: Time is over")

//Empty used when there paramater is empty
var Empty = errors.New("Error: parameter are empty")

//Seguridad used when user broke some rules already in frontend
var Security = errors.New("Error: seguridad")
