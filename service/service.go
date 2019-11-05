package service

var countryCodes []string

func init() {
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
