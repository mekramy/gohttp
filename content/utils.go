package content

import (
	"strings"
)

func isValidContent(c string, valids ...string) bool {
	c = strings.ToLower(strings.TrimSpace(c))
	for _, v := range valids {
		if strings.HasPrefix(c, strings.ToLower(v)) {
			return true
		}
	}

	return false
}
