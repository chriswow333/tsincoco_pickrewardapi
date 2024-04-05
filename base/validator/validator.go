package validator

import (
	"fmt"
	"regexp"
	"strings"

	phone "github.com/nyaruka/phonenumbers"
)

var (
	// ErrCountryCallingCode ...
	ErrCountryCallingCode = fmt.Errorf("Invalid countryCallingCode")
	// ErrLocalPhoneNumber ...
	ErrLocalPhoneNumber = fmt.Errorf("Invalid localPhoneNumber")
	// ErrTwoDigitISO ...
	ErrTwoDigitISO = fmt.Errorf("Invalid twoDigitISO")
	// ErrPhoneNumberFormat ...
	ErrPhoneNumberFormat = fmt.Errorf("Phone number wrong format")
)

// IsValidCountryCallingCode returns if countryCallingCode format is valid
// Example: 886, 81
func IsValidCountryCallingCode(countryCallingCode string) bool {
	return regexp.MustCompile(`^[1-9][0-9]{0,2}$`).MatchString(countryCallingCode)
}

// IsValidLocalPhoneNumber returns if localPhoneNumber format is valid
// Example: 0987654321, 987654321
func IsValidLocalPhoneNumber(localPhoneNumber string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(localPhoneNumber)
}

// IsValidTwoDigitISO returns if twoDigitISO format is valid
// Example: tw, JP
func IsValidTwoDigitISO(twoDigitISO string) bool {
	return regexp.MustCompile(`^[a-zA-Z]{2}$`).MatchString(twoDigitISO)
}

// ValidatePhoneNumber ...
func ValidatePhoneNumber(twoDigitISO, countryCallingCode, localPhoneNumber string) error {
	if !IsValidCountryCallingCode(countryCallingCode) {
		return ErrCountryCallingCode
	}
	if !IsValidLocalPhoneNumber(localPhoneNumber) {
		return ErrLocalPhoneNumber
	}
	if !IsValidTwoDigitISO(twoDigitISO) {
		return ErrTwoDigitISO
	}

	// copy from sms.IsSMSReceiver
	formattedPhoneNumber, err := phone.Parse("+"+countryCallingCode+localPhoneNumber, "")
	if err != nil {
		return err
	}

	numberType := phone.GetNumberType(formattedPhoneNumber)
	switch numberType {
	case phone.FIXED_LINE, phone.MOBILE, phone.FIXED_LINE_OR_MOBILE, phone.PERSONAL_NUMBER:
		return nil
	}
	return ErrPhoneNumberFormat
}

// IsValidOpenID returns if openID format is valid
func IsValidOpenID(openID string) bool {
	openID = strings.ToLower(openID)
	return regexp.MustCompile(`^[\p{Han}\x{3041}-\x{3096}\x{30A1}-\x{30FC}\w.]{2,20}$`).MatchString(openID) && len(regexp.MustCompile(`[\p{Han}\x{3041}-\x{3096}\x{30A1}-\x{30FC}]`).FindAllStringIndex(openID, -1)) <= 4 && !regexp.MustCompile(`^17`).MatchString(openID)
}

// IsInStringSlice validates target exists in list or not
func IsInStringSlice(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

// IsInIntSlice validates target exists in list or not
func IsInIntSlice(intlist []int, target int) bool {
	for _, s := range intlist {
		if s == target {
			return true
		}
	}
	return false
}

// IsInInt32Slice validates target exists in list or not
func IsInInt32Slice(intlist []int32, target int32) bool {
	for _, s := range intlist {
		if s == target {
			return true
		}
	}
	return false
}

// IsInInt64Slice validates target exists in list or not
func IsInInt64Slice(int64list []int64, target int64) bool {
	for _, s := range int64list {
		if s == target {
			return true
		}
	}
	return false
}

type StringSlice []string

// Intersect returns intersection between self and another string slice
func (s StringSlice) Intersect(s2 []string) StringSlice {
	hash := map[string]int{}
	for _, elem := range s {
		hash[elem]++
	}

	result := StringSlice{}
	for _, target := range s2 {
		if count, ok := hash[target]; ok && count > 0 {
			result = append(result, target)
			hash[target]--
		}
	}

	return result
}

// IsEmpty checks if this slice is empty
func (s StringSlice) IsEmpty() bool {
	return len(s) == 0
}
