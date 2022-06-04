package main

import (
	"fmt"
	"strings"
	"yamlfix/misc"
)

// TagDecorate Add HTML open/close tag to a phrase, and
// accept any attributes as well such as:
// TagDecorate("Huh?", "b", "style=\"color: hotpink;\"", "id=\"+GetUnusedId()+"\"")
// would produce:
// <b style="color: hotpink; id="0000000F_spizfuzzle">Huh?</b>
func TagDecorate(base string, tags ...string) (decoratedString string) {
	var sb strings.Builder
	if 0 >= len(tags) {
		xLog.Printf("huh? Attempting to HTML-ize %s without an HTML tag?", base)
	}
	misc.WriteSB(&sb, "<", tags[0])
	for ix := 1; ix < len(tags); ix++ {
		if "" != tags[ix] {
			misc.WriteSB(&sb, " ", tags[ix])
		}
	}
	misc.WriteSB(&sb, ">", MarkdownToHtml(base), "</", tags[0], ">")
	return sb.String()
}

// MarkdownToHtml converts simple Markdown formatting
// into HTML-tagged formatted text, naively.
func MarkdownToHtml(markdown string) (html string) {
	const ThreeStars = "***"
	const TwoStars = "**"
	const OneStar = "*"
	const ThreeUnder = "___"
	const TwoUnder = "__"
	// const OneUnder = "_"
	const Sneech = "`" // https://bit.ly/3qSumos

	if !IsStringSet(&markdown) {
		return ""
	}
	// MIGHT be a better way to do it?  MUST use 3-2-1 order!
	html = MarkdownSeparatorToHtmlTags(markdown, ThreeStars, "b", "i")
	html = MarkdownSeparatorToHtmlTags(html, TwoStars, "b")
	html = MarkdownSeparatorToHtmlTags(html, OneStar, "i")
	html = MarkdownSeparatorToHtmlTags(html, ThreeUnder, "i", "b")
	html = MarkdownSeparatorToHtmlTags(html, TwoUnder, "b")
	// html = MarkdownSeparatorToHtmlTags(html, OneUnder, "i")
	return MarkdownSeparatorToHtmlTags(html, Sneech, "code")
}

// MarkdownSeparatorToHtmlTags do the actual transformation
// of A MARKDOWN FORMAT MODE to HTML-TAGGED TEXT
func MarkdownSeparatorToHtmlTags(target string, sep string, tags ...string) (rep string) {
	var replace strings.Builder
	substr := strings.SplitAfter(target, sep)

	if nil == substr || 0 >= len(substr) || 1 != len(substr)%2 {
		if FlagDebug {
			replace.WriteString(
				fmt.Sprintf("huh?\n\tDESCRIPTION: %s\n\nSEPARATOR: %s\n\t TAGS: %q\n\t",
					target, sep, tags))
			if nil == substr {
				replace.WriteString("-- no generated substrings --")
			} else {
				replace.WriteString(
					fmt.Sprintf("unexpected (%d is even) number of substrings:\n\t[ %q ]",
						len(substr)%2, substr))
			}
			xLog.Println(replace.String())
			replace.Reset()
		}
		return target
	}

	if 2 < len(substr) { // any separations happen?
		for ix01 := 0; ix01 < len(substr)-1; ix01 += 2 {
			markedText := strings.ReplaceAll(substr[ix01+1], sep, "")
			for _, val := range tags {
				markedText = TagDecorate(markedText, val)
			}
			misc.WriteSB(&replace,
				strings.ReplaceAll(substr[ix01], sep, ""),
				markedText)
		}
		misc.WriteSB(&replace, substr[len(substr)-1])
	} else {
		return target
	}
	return replace.String()
}

func IsStringSet(s *string) (isSet bool) {
	if nil != s && "" != *s {
		return true
	}
	return false
}
