package service

import (
	"fmt"
	"os"
	"time"
)

var countryCodes []string

// Location contains the London timezone
var Location *time.Location

// TimeDisplayFormat is the standard format for output.
const TimeDisplayFormat = "2006-01-02 15:04"

func init() {
	var err error
	Location, err = time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.LoadLocation(%q) failed: %+v",
			"Europe/London", err.Error())
		return
	}

	countryCodes = []string{
		"UK - United Kingdom",
		"AT - Austria",
		"BG - Bulgaria",
		"CZ - Czechia",
		"DK - Denmark",
		"EE - Estonia",
		"FI - Finland",
		"FR - France",
		"DE - Germany",
		"HU - Hungary",
		"IT - Italy",
		"IE - Ireland",
		"LT - Lithuania",
		"LU - Luxembourg",
		"NL - Netherlands",
		"NO - Norway",
		"PL - Poland",
		"PT - Portugal",
		"RO - Romania",
		"SK - Slovakia",
		"SI - Slovenia",
		"ES - Spain (including Balearic islands)",
		"SE - Sweden",
		"CH - Switzerland",
		"US - United States",
	}
}

// CountryCodes returns a slice of country codes
func CountryCodes() []string {
	return countryCodes
}

// TitleFromCountryCode accepts the 2 character country code
// and returns the full country title.
func TitleFromCountryCode(match string) string {
	for _, c := range countryCodes {
		if match == c[0:2] {
			return c
		}
	}
	return ""
}
