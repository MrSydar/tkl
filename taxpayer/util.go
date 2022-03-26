package taxpayer

import (
	"fmt"
	"regexp"
)

func getFirstSubgroupMatch(text string, re *regexp.Regexp) (string, error) {
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return "", fmt.Errorf("no submatch found")
	} else {
		return match[1], nil
	}
}
