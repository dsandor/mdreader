# MD Reader

A command-line utility that converts Markdown files to beautifully formatted HTML with GitHub-style rendering and syntax highlighting.

## Features

- GitHub-style markdown rendering
- Syntax highlighting for code blocks
- Support for tables, blockquotes, and all standard markdown elements
- Optional browser launching after conversion
- Clean, responsive HTML output
- **Interactive UI mode** with split-pane editor and live preview

## Installation

### From Source

Requires Go 1.20 or higher.

```bash
git clone https://github.com/dsandor/mdreader.git
cd mdreader
go build -o mdreader
```

### Quick Start

```bash
# Build the tool
go build -o mdreader

# Convert a markdown file
./mdreader README.md
```

## Usage

### Basic Usage

Convert a markdown file to HTML:

```bash
mdreader <input.md>
```

This creates an HTML file with the same name as your input file (e.g., `input.md` → `input.html`).

### Command-Line Options

#### Input File

Specify the markdown file to convert. Can be provided as a positional argument or with the `--input` flag:

```bash
mdreader document.md
# or
mdreader --input document.md
```

#### Output File (`--output`)

Specify a custom name for the output HTML file:

```bash
mdreader document.md --output custom-name.html
```

If not specified, the output file will have the same base name as the input file with an `.html` extension.

#### Launch Browser (`--launch`)

Automatically open the generated HTML file in your default web browser:

```bash
mdreader document.md --launch
```

#### Interactive UI Mode (`--ui`)

Launch an interactive markdown editor with live preview in your web browser:

```bash
# Start UI with empty editor
mdreader --ui

# Start UI with a file loaded
mdreader --ui document.md
```

The UI mode provides:
- **Split-pane interface**: Markdown editor on the left, live preview on the right
- **Real-time preview**: See changes as you type
- **File operations**: New, Open, Save, Save As
- **Export to HTML**: Export the rendered HTML to a file
- **Keyboard shortcuts**:
  - `Ctrl/Cmd + S`: Save file
  - `Ctrl/Cmd + Shift + S`: Save As
  - `Ctrl/Cmd + O`: Open file
  - `Ctrl/Cmd + N`: New file
- **Resizable panes**: Drag the divider to adjust editor/preview sizes
- **Syntax highlighting** in the preview pane
- **Line and column position** tracking

### Examples

#### Simple Conversion
```bash
# Convert README.md to README.html
mdreader README.md
```

#### Custom Output Name
```bash
# Convert notes.md to documentation.html
mdreader notes.md --output documentation.html
```

#### Convert and View
```bash
# Convert and immediately open in browser
mdreader README.md --launch
```

#### Interactive Editing
```bash
# Launch the UI editor
mdreader --ui

# Edit an existing file in the UI
mdreader --ui README.md
```

#### Full Options
```bash
# Use all options together (for CLI mode)
mdreader --input notes.md --output final-doc.html --launch
```

## Supported Markdown Features

- **Headers** (H1-H6)
- **Text formatting**: bold, italic, strikethrough
- **Lists**: ordered and unordered, nested lists
- **Links and images**
- **Code blocks** with syntax highlighting for popular languages:
  - Go, Python, JavaScript, TypeScript, Java, C/C++, Rust
  - HTML, CSS, JSON, YAML, SQL
  - Shell/Bash scripts
  - And many more...
- **Tables** with GitHub-style formatting
- **Blockquotes**
- **Horizontal rules**
- **Inline code**

## Output

The generated HTML includes:

- Embedded GitHub-style CSS
- Syntax highlighting styles
- Responsive design that works on desktop and mobile
- Self-contained file (no external dependencies)

## Examples of Supported Languages

The syntax highlighter automatically detects and highlights code blocks:

````markdown
```go
func main() {
    fmt.Println("Hello, World!")
}
```

```python
def greet(name):
    return f"Hello, {name}!"
```

```javascript
const greet = (name) => {
    return `Hello, ${name}!`;
};
```
````

## Building from Source

```bash
# Clone the repository
git clone https://github.com/dsandor/mdreader.git
cd mdreader

# Install dependencies
go mod download

# Build the binary
go build -o mdreader

# Optional: Install to PATH
go install
```

## Requirements

- Go 1.20 or higher (for building from source)
- A modern web browser (for viewing the HTML output)

## License

[Add your license information here]

## Contributing

[Add contribution guidelines if applicable]
