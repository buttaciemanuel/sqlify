package prompt

import (
	"fmt"
	"strings"

	"github.com/buttaciemanuel/sqlify/context"
)

type PromptSection interface {
	Dump(stream *strings.Builder)
}

type TextSection struct {
	Title string
	Body  string
	Items []string
}

func (section *TextSection) Dump(stream *strings.Builder) {
	title := strings.Join(
		strings.Split(strings.ToLower(section.Title), " "),
		"_",
	)

	stream.WriteString("<")
	stream.WriteString(title)
	stream.WriteString(">")

	if len(section.Body) > 0 {
		stream.WriteString("\n  ")
		stream.WriteString(section.Body)
	}

	for index, item := range section.Items {
		stream.WriteString("\n  <item>\n    ")
		stream.WriteString(fmt.Sprintf("%d. ", index+1))
		stream.WriteString(item)
		stream.WriteString("\n  </item>")
	}

	stream.WriteString("\n</")
	stream.WriteString(title)
	stream.WriteString(">\n")
}

type CodeSection struct {
	Title   string
	Body    string
	Samples []context.Document
}

func (section *CodeSection) Dump(stream *strings.Builder) {
	title := strings.Join(
		strings.Split(strings.ToLower(section.Title), " "),
		"_",
	)

	stream.WriteString("<")
	stream.WriteString(title)
	stream.WriteString(">")

	if len(section.Body) > 0 {
		stream.WriteString("\n  ")
		stream.WriteString(section.Body)
	}

	for _, sample := range section.Samples {
		stream.WriteString("\n  <sql>\n    ")
		stream.WriteString(sample.Value)
		stream.WriteString("\n  </sql>")
	}

	stream.WriteString("\n</")
	stream.WriteString(title)
	stream.WriteString(">\n")
}

type Prompt struct {
	Sections []TextSection
	Query    TextSection
	Schemas  CodeSection
	Examples CodeSection
}

func (prompt *Prompt) Build(
	query string,
	schemas, examples []context.Document,
) string {
	var stream strings.Builder

	for _, section := range prompt.Sections {
		section.Dump(&stream)
	}

	prompt.Schemas.Samples = schemas
	prompt.Examples.Samples = examples

	prompt.Schemas.Dump(&stream)
	prompt.Examples.Dump(&stream)

	prompt.Query.Body = query
	prompt.Query.Items = []string{}

	prompt.Query.Dump(&stream)

	return stream.String()
}
