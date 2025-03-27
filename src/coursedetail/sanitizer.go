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

func (d *Description) SanitizeContent(softNewLines bool) string {
	content := d.Content
	if softNewLines {
		// convert "\n" to "  \n" for hard line breaks
		content = strings.Replace(content, "\n", "\n\n", -1)
	}
	// convert markdown to html
	htmlContent := markdownToHTML(content)
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
