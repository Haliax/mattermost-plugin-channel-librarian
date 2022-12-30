package main

import "strings"

func ContainsValueCaseInsensitive(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}

	return false
}

func ContainsValue(arr []string, value string) bool {
	for _, element := range arr {
		if element == value {
			return true
		}
	}

	return false
}
