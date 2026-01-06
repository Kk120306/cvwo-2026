package helpers

import (
	"regexp"
	"strings"
)

// converts title into a slug for topic
func GenerateSlug(name string) string {
	// lowercase
	slug := strings.ToLower(name)
	// keep only alphabetics and - and space and replace the space with -
	reg, _ := regexp.Compile(`[^\w\s-]`)
	slug = reg.ReplaceAllString(slug, "")
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-") // remove any continousou -
	slug = strings.Trim(slug, "-")                              //  remove trailing dash

	return slug
}
