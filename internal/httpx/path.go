package httpx

import (
	"fmt"
	"strconv"
	"strings"
)

// SplitPath trims leading/trailing slashes and splits a URL path into parts.
func SplitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func ParseID(value string, fieldName string) (int, error){
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", fieldName, err)
	}
	return id, nil
}
