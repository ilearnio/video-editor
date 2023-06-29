package shotcutProjectBuilder

import (
	"fmt"
	"strings"

	"videoeditor/src/services/shotcutProjectBuilder/helpers"
)

type XMLNodeAttr struct {
	Name  string
	Value string
}

type XMLNode struct {
	TagName     string
	Attrs       []XMLNodeAttr
	ChildNodes  []XMLNode
	TextContent string
}

func NodeToXML(node XMLNode, depth int) string {
	indent := strings.Repeat("  ", depth)
	var attrs string
	if node.Attrs != nil {
		attrStrings := []string{}
		for _, attr := range node.Attrs {
			if attr.Value == "" {
				continue
			}
			attrStrings = append(
				attrStrings,
				fmt.Sprintf(`%s="%s"`, attr.Name, helpers.EscapeAttr(attr.Value)),
			)
		}
		attrs = " " + strings.Join(attrStrings, " ")
	}

	if node.TextContent != "" {
		return fmt.Sprintf(
			`%s<%s%s>%s</%s>`,
			indent, node.TagName, attrs, helpers.EscapeHTML(node.TextContent), node.TagName,
		)
	}

	var childrenXml string
	if node.ChildNodes != nil && len(node.ChildNodes) > 0 {
		childStrings := make([]string, len(node.ChildNodes))
		for i, child := range node.ChildNodes {
			childStrings[i] = NodeToXML(child, depth+1)
		}
		childrenXml = strings.Join(childStrings, "\n")
	}

	if childrenXml != "" {
		return fmt.Sprintf(
			"%s<%s%s>\n%s\n%s</%s>",
			indent, node.TagName, attrs, childrenXml, indent, node.TagName,
		)
	}

	return fmt.Sprintf(
		`%s<%s%s/>`,
		indent, node.TagName, attrs,
	)
}
