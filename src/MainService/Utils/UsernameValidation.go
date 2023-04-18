package Utils

import "regexp"

func UsernameValidation(username string) int {
	regexes := []string{
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		`^+98 9[0-9]{2} [0-9]{3} [0-9]{4}$`,
		`^[+989]+[0-9]{2}[0-9]{3}[0-9]{4}$`,
		`^0 9[0-9]{2} [0-9]{3} [0-9]{4}$`,
		`^09[0-9]{2}[0-9]{3}[0-9]{4}$`,
	}
	for i, regex := range regexes {
		if match, _ := regexp.MatchString(regex, username); match {
			return i
		}
	}
	return -1
}
