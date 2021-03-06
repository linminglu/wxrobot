package logic

import (
	"strings"
)

func ExecCheckEqualFunc(f, value string) bool {
	value = strings.Replace(value, " ", "", -1)
	if strings.HasPrefix(f, NOTINCLUDE) {
		v := strings.Replace(f, NOTINCLUDE, "", -1)
		vSlice := strings.Split(v, ",")
		for _, vStr := range vSlice {
			if value == vStr {
				return false
			}
		}
	} else if strings.HasPrefix(f, INCLUDE) {
		v := strings.Replace(f, INCLUDE, "", -1)
		vSlice := strings.Split(v, ",")
		for _, vStr := range vSlice {
			if value == vStr {
				return true
			}
		}
		return false
	} else if strings.HasPrefix(f, EQUAL) {
		v := strings.Replace(f, EQUAL, "", -1)
		return value == v
	}

	return true
}

func ExecCheckFunc(f, value string) bool {
	value = strings.Replace(value, " ", "", -1)
	if strings.HasPrefix(f, NOTINCLUDE) {
		v := strings.Replace(f, NOTINCLUDE, "", -1)
		vSlice := strings.Split(v, ",")
		for _, vStr := range vSlice {
			if strings.Contains(value, vStr) {
				return false
			}
		}
	} else if strings.HasPrefix(f, INCLUDE) {
		v := strings.Replace(f, INCLUDE, "", -1)
		vSlice := strings.Split(v, ",")
		for _, vStr := range vSlice {
			if strings.Contains(value, vStr) {
				return true
			}
		}
		return false
	} else if strings.HasPrefix(f, EQUAL) {
		v := strings.Replace(f, EQUAL, "", -1)
		vSlice := strings.Split(v, ",")
		for _, vStr := range vSlice {
			if value == vStr {
				return true
			}
		}
		return false
	}

	return true
}

func ExecGetArgvFunc(f string) string {
	var v string
	if strings.HasPrefix(f, NOTINCLUDE) {
		v = strings.Replace(f, NOTINCLUDE, "", -1)
	} else if strings.HasPrefix(f, INCLUDE) {
		v = strings.Replace(f, INCLUDE, "", -1)
	} else if strings.HasPrefix(f, EQUAL) {
		v = strings.Replace(f, EQUAL, "", -1)
	}

	return v
}
