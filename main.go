package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/russross/blackfriday/v2"
)

func main() {
	var inputFile string
	var outputFile string
	var launch bool

	flag.StringVar(&inputFile, "input", "", "Input markdown file")
	flag.StringVar(&outputFile, "output", "", "Output HTML file (optional)")
	flag.BoolVar(&launch, "launch", false, "Launch HTML file in default browser")
	flag.Parse()

	if inputFile == "" && flag.NArg() > 0 {
		inputFile = flag.Arg(0)
	}

	if inputFile == "" {
		fmt.Println("Usage: mdreader <input.md> [--output <output.html>] [--launch]")
		fmt.Println("       mdreader --input <input.md> [--output <output.html>] [--launch]")
		os.Exit(1)
	}

	if outputFile == "" {
		base := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile = base + ".html"
	}

	markdown, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	htmlContent := convertMarkdownToHTML(markdown)

	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("HTML file created: %s\n", outputFile)

	if launch {
		err = openBrowser(outputFile)
		if err != nil {
			log.Printf("Error opening browser: %v", err)
		}
	}
}

func convertMarkdownToHTML(markdown []byte) string {
	renderer := NewCustomHTMLRenderer()
	extensions := blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs | blackfriday.FencedCode | blackfriday.Tables
	body := blackfriday.Run(markdown, blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extensions))

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Preview</title>
    <style>
        %s
        %s
    </style>
</head>
<body>
    <div class="markdown-body">
        %s
    </div>
</body>
</html>`, getGithubCSS(), getChromaCSS(), string(body))

	return html
}

type CustomHTMLRenderer struct {
	*blackfriday.HTMLRenderer
}

func NewCustomHTMLRenderer() *CustomHTMLRenderer {
	return &CustomHTMLRenderer{
		HTMLRenderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CommonHTMLFlags,
		}),
	}
}

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if node.Type == blackfriday.CodeBlock {
		if entering {
			lang := string(node.CodeBlockData.Info)
			highlighted := highlightCode(string(node.Literal), lang)
			w.Write([]byte(highlighted))
		}
		return blackfriday.GoToNext
	}
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

func highlightCode(code, lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))
	
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true), html.PreventSurroundingPre(false))
	
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return fmt.Sprintf("<pre><code>%s</code></pre>", code)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return fmt.Sprintf("<pre><code>%s</code></pre>", code)
	}

	return buf.String()
}

func getChromaCSS() string {
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))
	var buf bytes.Buffer
	formatter.WriteCSS(&buf, style)
	return buf.String()
}

func getGithubCSS() string {
	return `
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #24292e;
            background-color: #ffffff;
            margin: 0;
            padding: 0;
        }

        .markdown-body {
            box-sizing: border-box;
            min-width: 200px;
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
        }

        .markdown-body h1 {
            padding-bottom: 0.3em;
            font-size: 2em;
            border-bottom: 1px solid #eaecef;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body h2 {
            padding-bottom: 0.3em;
            font-size: 1.5em;
            border-bottom: 1px solid #eaecef;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body h3 {
            font-size: 1.25em;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body h4 {
            font-size: 1em;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body h5 {
            font-size: 0.875em;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body h6 {
            font-size: 0.85em;
            color: #6a737d;
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }

        .markdown-body p {
            margin-top: 0;
            margin-bottom: 16px;
        }

        .markdown-body blockquote {
            margin: 0;
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin-top: 0;
            margin-bottom: 16px;
        }

        .markdown-body ul,
        .markdown-body ol {
            margin-top: 0;
            margin-bottom: 16px;
            padding-left: 2em;
        }

        .markdown-body ul ul,
        .markdown-body ul ol,
        .markdown-body ol ol,
        .markdown-body ol ul {
            margin-top: 0;
            margin-bottom: 0;
        }

        .markdown-body li {
            word-wrap: break-all;
        }

        .markdown-body li > p {
            margin-top: 16px;
        }

        .markdown-body li + li {
            margin-top: 0.25em;
        }

        .markdown-body code {
            padding: 0.2em 0.4em;
            margin: 0;
            font-size: 85%;
            background-color: rgba(27,31,35,0.05);
            border-radius: 3px;
            font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
        }

        .markdown-body pre {
            margin-top: 0;
            margin-bottom: 16px;
            font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
            font-size: 85%;
            line-height: 1.45;
            background-color: #f6f8fa;
            border-radius: 6px;
            overflow: auto;
            padding: 16px;
        }

        .markdown-body pre code {
            padding: 0;
            margin: 0;
            font-size: 100%;
            word-break: normal;
            white-space: pre;
            background: transparent;
            border: 0;
        }

        .markdown-body table {
            border-spacing: 0;
            border-collapse: collapse;
            margin-top: 0;
            margin-bottom: 16px;
            display: block;
            width: 100%;
            overflow: auto;
        }

        .markdown-body table th {
            font-weight: 600;
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
            background-color: #f6f8fa;
        }

        .markdown-body table td {
            padding: 6px 13px;
            border: 1px solid #dfe2e5;
        }

        .markdown-body table tr {
            background-color: #fff;
            border-top: 1px solid #c6cbd1;
        }

        .markdown-body table tr:nth-child(2n) {
            background-color: #f6f8fa;
        }

        .markdown-body hr {
            height: 0.25em;
            padding: 0;
            margin: 24px 0;
            background-color: #e1e4e8;
            border: 0;
        }

        .markdown-body a {
            color: #0366d6;
            text-decoration: none;
        }

        .markdown-body a:hover {
            text-decoration: underline;
        }

        .markdown-body img {
            max-width: 100%;
            box-sizing: content-box;
            background-color: #fff;
        }

        .markdown-body strong {
            font-weight: 600;
        }

        .markdown-body em {
            font-style: italic;
        }

        .markdown-body del {
            text-decoration: line-through;
        }

        @media (max-width: 767px) {
            .markdown-body {
                padding: 15px;
            }
        }
    `
}

func openBrowser(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	
	url := "file://" + absPath

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}