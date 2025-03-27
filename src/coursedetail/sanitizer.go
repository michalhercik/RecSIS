package coursedetail

import (
	"regexp"
	"strings"

	"log"

	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	goldmarkRenderer "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
)

func (d *Description) SanitizeContent(keepNewLines bool) string {
	content := d.Content
	// convert "\n" to "\n\n" for hard line breaks, if needed
	hardNewLines := makeHardLineBreaks(content, keepNewLines)
	// convert markdown to html
	htmlContent := markdownToHTML(hardNewLines)
	// clean html elements
	cleanedContent := cleanHTML(htmlContent)
	// sanitize html
	sanitizedContent := sanitizeHTML(cleanedContent)

	// ================== FOR TESTING PURPOSES ==================
	removedTags := getRemovedTags(cleanedContent, sanitizedContent)
	if len(removedTags) > 0 {
		log.Println("Removed tags: ", removedTags)
	}
	// ==========================================================

	return sanitizedContent
}

func makeHardLineBreaks(content string, keepOriginal bool) string {
	if keepOriginal {
		return content
	}

	// Split the content into lines.
	lines := strings.Split(content, "\n")
	var builder strings.Builder

	for i, line := range lines {
		builder.WriteString(line)
		// Determine if we should add an extra newline.
		// Check if we're not at the last line.
		if i < len(lines)-1 {
			// If either the current or the next line is a list item,
			// insert only a single newline.
			if isListItem(line) || isListItem(lines[i+1]) {
				builder.WriteString("\n")
			} else {
				builder.WriteString("\n\n")
			}
		}
	}

	newContent := builder.String()
	return newContent
}

// isListItem checks if a line starts with a common Markdown list marker.
func isListItem(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ")
}

// convert markdown to HTML using goldmark
func markdownToHTML(content string) string {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			goldmarkRenderer.WithUnsafe(), // Allow raw HTML
		),
	)
	// convert markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		panic(err)
	}
	return buf.String()
}

// for correct rendering, it is necessary to replace invalid tags with <p>
// known invalid tags: <h7>, <tema>
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
