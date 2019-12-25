// Package figures is a extension for the goldmark
// (http://github.com/yuin/goldmark).
//
// This extension adds html render for paragraph with image output as figure.
//
// An image with nonempty alt text, occurring by itself in a paragraph, will be
// rendered as a figure with a caption. The image’s alt text will be used as the
// caption.
//  ![**Figure:** [description](/link)](image.png "title")
// This syntax is borrowed from [Pandoc](https://pandoc.org/MANUAL.html#images).
//
// If you just want a regular inline image, just make sure it is not the only
// thing in the paragraph. One way to do this is to insert a nonbreaking space
// after the image:
//  ![This image won't be a figure](/url/of/image.png)\
package figures

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type render struct {
	html.Config
}

// NewRenderer return new renderer for image figure.
func NewRenderer(opts ...html.Option) renderer.NodeRenderer {
	var r = &render{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

var defaultRenderer = NewRenderer()

// RegisterFuncs implement renderer.NodeRenderer interface.
func (r *render) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *render) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	var (
		name     = "p"
		filter   = html.ParagraphAttributeFilter
		isFigure = isFigure(n)
	)
	if isFigure {
		name = "figure"
		filter = html.GlobalAttributeFilter
	}
	if entering {
		w.WriteByte('<')
		w.WriteString(name)
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, filter)
		}
		w.WriteByte('>')
		if isFigure {
			w.WriteByte('\n')
		}
	} else {
		w.WriteString("</")
		w.WriteString(name)
		w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

func (r *render) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// special case for figure
		if node.HasChildren() && isFigure(node.Parent()) {
			w.WriteString("</figcaption>\n")
		}
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	w.WriteString(`" alt="`)
	w.Write(n.Text(source))
	w.WriteByte('"')
	if n.Title != nil {
		w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		w.WriteString(" />")
	} else {
		w.WriteString(">")
	}
	// special case for figure
	if node.HasChildren() && isFigure(node.Parent()) {
		w.WriteString("\n<figcaption>")
		return ast.WalkContinue, nil
	}
	return ast.WalkSkipChildren, nil
}

func isFigure(node ast.Node) bool {
	var child = node.FirstChild()
	return node.Kind() == ast.KindParagraph &&
		child != nil &&
		child == node.LastChild() &&
		child.Kind() == ast.KindImage &&
		child.HasChildren()
}

// extension defines a goldmark.Extender for markdown image figures.
type extension struct{}

// Extend implement goldmark.Extender interface.
func (a *extension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(defaultRenderer, 100),
		),
	)
}

// Extension is a goldmark.Extender with markdown image figures support.
var Extension goldmark.Extender = new(extension)

// Enable is a goldmark.Option with image figures support.
var Enable = goldmark.WithExtensions(Extension)
