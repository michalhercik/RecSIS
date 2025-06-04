package coursedetail

import (
	"regexp"
	"strings"

	"log"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

func (d *Description) SanitizeContent(squash bool) string {
	content := d.Content
	// parse content to html
	htmlContent := parseToHTML(content, squash)
	// clean html elements
	cleanedContent := cleanHTML(htmlContent)
	// sanitize html
	sanitizedContent := sanitizeHTML(cleanedContent)

	// ================== FOR TESTING PURPOSES ==================
	// TODO: remove
	removedTags := getRemovedTags(cleanedContent, sanitizedContent)
	if len(removedTags) > 0 {
		log.Println("Removed tags: ", removedTags)
	}
	// ==========================================================

	return sanitizedContent
}

// creating own parser for SIS content
// ==================================
// *** KNOWN RULES ***
//
// content		translation
// ----------------------------------
// '* line'		<b> line </b>
// '-line'		<ul><li> line </li></ul>
// 'line'		<p> line </p>
//
// extra
// ----------------------------------
// '<link>'		<a href="link"> link </a>
//
// ==================================
// annotation seems to be just squashed text
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

// for correct rendering, it is necessary to replace invalid tags with <p>
// ==================================
// *** KNOWN INVALID TAGS ***
//
// <h7>
// <tema>
//
// ==================================
func cleanHTML(content string) string {
	// regex to match <h7> and <tema> tags case-insensitively
	re := regexp.MustCompile(`(?i)</?(h7|tema)>`)
	// replace with <p> for opening tags and </p> for closing tags
	output := re.ReplaceAllStringFunc(content, func(match string) string {
		if strings.HasPrefix(match, "</") {
			return "</p>"
		}
		return "<p>"
	})
	return output
}

// sanitizing HTML using bluemonday
func sanitizeHTML(content string) string {
	p := bluemonday.NewPolicy()
	// allow only the most basic HTML elements
	p.AllowElements("a", "b", "i", "strong", "em", "p", "small", "br", "br/", "h1", "h2", "h3", "h4", "h5", "h6", "span", "var", "sub", "sup")
	// allow links
	p.AllowAttrs("href").OnElements("a")
	p.AllowStandardURLs()
	// allow font for backward compatibility
	p.AllowElements("font")
	p.AllowAttrs("face").OnElements("font")
	// allow the most basic attributes
	p.AllowStandardAttributes()
	p.AllowAttrs("hidden").Globally()
	// allow images, lists, and tables
	p.AllowImages()
	p.AllowLists()
	p.AllowTables()
	p.AllowAttrs("border").OnElements("table")
	// sanitize the content
	return p.Sanitize(content)
}

// =============== FOR TESTING PURPOSES ===============
func getRemovedTags(original string, sanitized string) []string {
	// Parse the original and sanitized HTML into Go HTML nodes
	originalTokens := tokenizeHTML(original)
	sanitizedTokens := tokenizeHTML(sanitized)

	// Find tags that were removed by comparing the tokens
	var removedTags []string
	for _, originalTag := range originalTokens {
		if !contains(sanitizedTokens, originalTag) {
			removedTags = append(removedTags, originalTag)
		}
	}
	return removedTags
}

func tokenizeHTML(htmlContent string) []string {
	var tags []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return tags // end of document
		case html.StartTagToken, html.SelfClosingTagToken:
			tagName, _ := tokenizer.TagName()
			tags = append(tags, string(tagName))
		}
	}
}

func contains(tokens []string, token string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}
