package coursedetail

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func (d *description) sanitizeContent(squash bool) string {
	content := d.content
	htmlContent := parseToHTML(content, squash)
	cleanedContent := cleanHTML(htmlContent)
	sanitizedContent := sanitizeHTML(cleanedContent)

	return sanitizedContent
}

/*
Replace SIS formatting rules with HTML tags.

Known rules:
  - '* [line]' ->  <b> [line] </b>
  - '-[line]'  ->  <ul><li> [line] </li></ul>
  - '[line]'   ->  <p> [line] </p>

Is used in text but sis formatter does not support it:
  - '<[link]>' -> <a href="[link]"> [link] </a>

annotation seems to be just squashed text
*/
func parseToHTML(content string, squash bool) string {
	var builder strings.Builder

	lines := strings.Split(content, "\n")
	if squash {
		// squash all lines into one <p> tag
		builder.WriteString("<p>")
		for _, line := range lines {
			builder.WriteString(strings.TrimSpace(line) + " ")
		}
		builder.WriteString("</p>")
	} else {
		// parse each line according to the rules
		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.Contains(line, "<") && strings.Contains(line, ">") {
				// find all substrings enclosed in angle brackets
				re := regexp.MustCompile(`<([^>]+)>`)
				matches := re.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					trimmed := strings.TrimSpace(match[1]) // trim content inside the angle brackets
					if strings.HasPrefix(trimmed, "http") {
						// convert to hyperlink
						line = strings.Replace(line, match[0], "<a href=\""+trimmed+"\">"+trimmed+"</a>", 1)
					} else {
						// leave the original content unchanged
						line = strings.Replace(line, match[0], "<"+trimmed+">", 1)
					}
				}
			}

			if strings.HasPrefix(line, "* ") {
				// bold text
				builder.WriteString("<b>" + strings.TrimPrefix(line, "* ") + "</b>")
			} else if strings.HasPrefix(line, "-") {
				// unordered list item
				builder.WriteString("<ul><li>" + strings.TrimPrefix(line, "-") + "</li></ul>")
			} else {
				// regular paragraph
				builder.WriteString("<p>" + line + "</p>")
			}
		}
	}

	return builder.String()
}

/*
For correct rendering, it is necessary to replace invalid tags with <p>

Known invalid tags:
  - <h7>
  - <tema>
*/
func cleanHTML(content string) string {
	re := regexp.MustCompile(`(?i)</?(h7|tema)>`)
	output := re.ReplaceAllStringFunc(content, func(match string) string {
		if strings.HasPrefix(match, "</") {
			return "</p>"
		}
		return "<p>"
	})
	return output
}

func sanitizeHTML(content string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements("a", "b", "i", "strong", "em", "p", "small", "br", "br/", "h1", "h2", "h3", "h4", "h5", "h6", "span", "var", "sub", "sup")
	p.AllowAttrs("href").OnElements("a")
	p.AllowStandardURLs()
	p.AllowElements("font")                 // for backward compatibility
	p.AllowAttrs("face").OnElements("font") // for backward compatibility
	p.AllowStandardAttributes()
	p.AllowAttrs("hidden").Globally()
	p.AllowImages()
	p.AllowLists()
	p.AllowTables()
	p.AllowAttrs("border").OnElements("table")
	return p.Sanitize(content)
}
