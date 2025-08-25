package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

func runUI(initialFile string) {
	initialContent := ""
	if initialFile != "" {
		content, err := os.ReadFile(initialFile)
		if err == nil {
			initialContent = string(content)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, generateUIHTML(initialContent))
	})

	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Filename string `json:"filename"`
			Content  string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := os.WriteFile(data.Filename, []byte(data.Content), 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	http.HandleFunc("/api/load", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Filename string `json:"filename"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content, err := os.ReadFile(data.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"content": string(content),
		})
	})

	port := "8080"
	url := fmt.Sprintf("http://localhost:%s", port)
	
	fmt.Printf("Starting MD Reader UI on %s\n", url)
	fmt.Println("Press Ctrl+C to stop")

	// Open browser after a short delay
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowserUI(url)
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}

		switch msg.Type {
		case "convert":
			html := convertMarkdownToHTMLForUI([]byte(msg.Content))
			if len(html) > 100 {
				log.Printf("Generated HTML preview (%d chars): %s...", len(html), html[:100])
			}
			response := Message{
				Type:    "preview",
				Content: html,
			}
			conn.WriteJSON(response)
		}
	}
}

func openBrowserUI(url string) error {
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

func generateUIHTML(initialContent string) string {
	// Use JSON encoding to properly escape the content for JavaScript
	contentJSON, _ := json.Marshal(initialContent)
	
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>MD Reader - Editor</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            height: 100vh;
            display: flex;
            flex-direction: column;
            background: #1e1e1e;
        }

        .toolbar {
            background: #2d2d30;
            padding: 10px;
            display: flex;
            gap: 10px;
            align-items: center;
            border-bottom: 1px solid #3e3e42;
        }

        .toolbar button {
            padding: 6px 12px;
            background: #0e639c;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }

        .toolbar button:hover {
            background: #1177bb;
        }

        .toolbar input {
            flex: 1;
            padding: 6px 10px;
            background: #3c3c3c;
            color: #cccccc;
            border: 1px solid #3e3e42;
            border-radius: 4px;
            font-size: 14px;
        }

        .toolbar .separator {
            width: 1px;
            height: 24px;
            background: #3e3e42;
            margin: 0 5px;
        }

        .container {
            display: flex;
            flex: 1;
            overflow: hidden;
        }

        .pane {
            flex: 1;
            display: flex;
            flex-direction: column;
            position: relative;
        }

        .pane-header {
            background: #2d2d30;
            padding: 8px 15px;
            font-size: 13px;
            color: #cccccc;
            border-bottom: 1px solid #3e3e42;
            font-weight: 500;
        }

        #editor {
            flex: 1;
            padding: 20px;
            font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
            font-size: 14px;
            line-height: 1.6;
            border: none;
            outline: none;
            resize: none;
            background: #1e1e1e;
            color: #d4d4d4;
            tab-size: 4;
        }

        #editor::selection {
            background: #264f78;
        }

        #preview-frame {
            flex: 1;
            border: none;
            background: white;
        }

        .divider {
            width: 4px;
            background: #2d2d30;
            cursor: col-resize;
            position: relative;
        }

        .divider:hover {
            background: #007acc;
        }

        .status-bar {
            background: #007acc;
            color: white;
            padding: 4px 15px;
            font-size: 12px;
            display: flex;
            justify-content: space-between;
        }

        .file-dialog {
            display: none;
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: #2d2d30;
            border: 1px solid #3e3e42;
            border-radius: 6px;
            padding: 20px;
            z-index: 1000;
            min-width: 400px;
        }

        .file-dialog h3 {
            color: #cccccc;
            margin-bottom: 15px;
        }

        .file-dialog input {
            width: 100%;
            padding: 8px;
            background: #1e1e1e;
            color: #cccccc;
            border: 1px solid #3e3e42;
            border-radius: 4px;
            margin-bottom: 15px;
        }

        .file-dialog-buttons {
            display: flex;
            gap: 10px;
            justify-content: flex-end;
        }

        .file-dialog button {
            padding: 6px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }

        .file-dialog .primary {
            background: #0e639c;
            color: white;
        }

        .file-dialog .secondary {
            background: #3c3c3c;
            color: #cccccc;
        }

        .overlay {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.5);
            z-index: 999;
        }
    </style>
</head>
<body>
    <div class="toolbar">
        <button onclick="newFile()">New</button>
        <button onclick="openFileDialog()">Open</button>
        <button onclick="saveFileDialog()">Save</button>
        <button onclick="saveAsDialog()">Save As</button>
        <div class="separator"></div>
        <button onclick="exportHTML()">Export HTML</button>
        <div class="separator"></div>
        <input type="text" id="current-file" placeholder="Untitled.md" value="Untitled.md">
    </div>

    <div class="container">
        <div class="pane">
            <div class="pane-header">MARKDOWN EDITOR</div>
            <textarea id="editor" placeholder="Start typing markdown..." spellcheck="false"></textarea>
        </div>
        
        <div class="divider" id="divider"></div>
        
        <div class="pane">
            <div class="pane-header">PREVIEW</div>
            <iframe id="preview-frame"></iframe>
        </div>
    </div>

    <div class="status-bar">
        <span id="status-text">Ready</span>
        <span id="cursor-pos">Line 1, Column 1</span>
    </div>

    <div class="overlay" id="overlay"></div>
    
    <div class="file-dialog" id="open-dialog">
        <h3>Open File</h3>
        <input type="text" id="open-filename" placeholder="Enter filename (e.g., document.md)">
        <div class="file-dialog-buttons">
            <button class="secondary" onclick="closeDialogs()">Cancel</button>
            <button class="primary" onclick="openFile()">Open</button>
        </div>
    </div>

    <div class="file-dialog" id="save-dialog">
        <h3>Save As</h3>
        <input type="text" id="save-filename" placeholder="Enter filename (e.g., document.md)">
        <div class="file-dialog-buttons">
            <button class="secondary" onclick="closeDialogs()">Cancel</button>
            <button class="primary" onclick="saveFile()">Save</button>
        </div>
    </div>

    <script>
        let currentFilename = 'Untitled.md';
        let isDirty = false;
        let lastSavedContent = '';
        let ws = null;

        const editor = document.getElementById('editor');
        const preview = document.getElementById('preview-frame');
        const statusText = document.getElementById('status-text');
        const cursorPos = document.getElementById('cursor-pos');
        const currentFileInput = document.getElementById('current-file');

        // Initialize WebSocket connection
        function connectWebSocket() {
            ws = new WebSocket('ws://localhost:8080/ws');
            
            ws.onopen = () => {
                console.log('WebSocket connected');
                statusText.textContent = 'Connected';
                updatePreview();
            };

            ws.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                if (msg.type === 'preview') {
                    preview.srcdoc = msg.content;
                    statusText.textContent = 'Preview updated';
                }
            };

            ws.onclose = () => {
                console.log('WebSocket disconnected, reconnecting...');
                statusText.textContent = 'Reconnecting...';
                setTimeout(connectWebSocket, 1000);
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                statusText.textContent = 'Connection error';
            };
        }

        // Set initial content
        const initialContent = ` + string(contentJSON) + `;
        editor.value = initialContent;

        // Initialize
        connectWebSocket();
        lastSavedContent = editor.value;
        updateCursorPosition();

        // Update preview on input
        let updateTimer;
        editor.addEventListener('input', () => {
            clearTimeout(updateTimer);
            updateTimer = setTimeout(updatePreview, 300);
            checkDirty();
        });

        // Update cursor position
        editor.addEventListener('click', updateCursorPosition);
        editor.addEventListener('keyup', updateCursorPosition);

        // Handle tab key
        editor.addEventListener('keydown', (e) => {
            if (e.key === 'Tab') {
                e.preventDefault();
                const start = editor.selectionStart;
                const end = editor.selectionEnd;
                const value = editor.value;
                editor.value = value.substring(0, start) + '    ' + value.substring(end);
                editor.selectionStart = editor.selectionEnd = start + 4;
            }
        });

        function updatePreview() {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: 'convert',
                    content: editor.value
                }));
            }
        }

        function updateCursorPosition() {
            const text = editor.value.substring(0, editor.selectionStart);
            const lines = text.split('\n');
            const line = lines.length;
            const column = lines[lines.length - 1].length + 1;
            cursorPos.textContent = 'Line ' + line + ', Column ' + column;
        }

        function checkDirty() {
            isDirty = editor.value !== lastSavedContent;
            updateTitle();
        }

        function updateTitle() {
            const indicator = isDirty ? '• ' : '';
            currentFileInput.value = indicator + currentFilename;
        }

        function newFile() {
            if (isDirty) {
                if (!confirm('You have unsaved changes. Continue without saving?')) {
                    return;
                }
            }
            editor.value = '';
            currentFilename = 'Untitled.md';
            lastSavedContent = '';
            isDirty = false;
            updateTitle();
            updatePreview();
            statusText.textContent = 'New file created';
        }

        function openFileDialog() {
            document.getElementById('overlay').style.display = 'block';
            document.getElementById('open-dialog').style.display = 'block';
            document.getElementById('open-filename').focus();
        }

        function saveFileDialog() {
            if (currentFilename === 'Untitled.md') {
                saveAsDialog();
            } else {
                saveCurrentFile();
            }
        }

        function saveAsDialog() {
            document.getElementById('overlay').style.display = 'block';
            document.getElementById('save-dialog').style.display = 'block';
            document.getElementById('save-filename').value = currentFilename;
            document.getElementById('save-filename').focus();
        }

        function closeDialogs() {
            document.getElementById('overlay').style.display = 'none';
            document.getElementById('open-dialog').style.display = 'none';
            document.getElementById('save-dialog').style.display = 'none';
        }

        async function openFile() {
            const filename = document.getElementById('open-filename').value;
            if (!filename) return;

            try {
                const response = await fetch('/api/load', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({filename})
                });
                
                const data = await response.json();
                if (data.status === 'success') {
                    editor.value = data.content;
                    currentFilename = filename;
                    lastSavedContent = data.content;
                    isDirty = false;
                    updateTitle();
                    updatePreview();
                    closeDialogs();
                    statusText.textContent = 'Opened: ' + filename;
                } else {
                    alert('Error opening file');
                }
            } catch (error) {
                alert('Error opening file: ' + error.message);
            }
        }

        async function saveFile() {
            const filename = document.getElementById('save-filename').value;
            if (!filename) return;

            await saveToFile(filename);
            closeDialogs();
        }

        async function saveCurrentFile() {
            await saveToFile(currentFilename);
        }

        async function saveToFile(filename) {
            try {
                const response = await fetch('/api/save', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        filename: filename,
                        content: editor.value
                    })
                });
                
                const data = await response.json();
                if (data.status === 'success') {
                    currentFilename = filename;
                    lastSavedContent = editor.value;
                    isDirty = false;
                    updateTitle();
                    statusText.textContent = 'Saved: ' + filename;
                } else {
                    alert('Error saving file');
                }
            } catch (error) {
                alert('Error saving file: ' + error.message);
            }
        }

        async function exportHTML() {
            const htmlFilename = currentFilename.replace(/\.md$/, '.html');
            const suggestedName = prompt('Export HTML as:', htmlFilename);
            if (!suggestedName) return;

            try {
                const response = await fetch('/api/save', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        filename: suggestedName,
                        content: preview.srcdoc
                    })
                });
                
                const data = await response.json();
                if (data.status === 'success') {
                    statusText.textContent = 'Exported HTML: ' + suggestedName;
                } else {
                    alert('Error exporting HTML');
                }
            } catch (error) {
                alert('Error exporting HTML: ' + error.message);
            }
        }

        // Handle window resize
        const divider = document.getElementById('divider');
        let isResizing = false;

        divider.addEventListener('mousedown', (e) => {
            isResizing = true;
            document.body.style.cursor = 'col-resize';
        });

        document.addEventListener('mousemove', (e) => {
            if (!isResizing) return;
            
            const container = document.querySelector('.container');
            const containerRect = container.getBoundingClientRect();
            const percentage = ((e.clientX - containerRect.left) / containerRect.width) * 100;
            
            if (percentage > 20 && percentage < 80) {
                const panes = document.querySelectorAll('.pane');
                panes[0].style.flex = percentage + '%';
                panes[1].style.flex = (100 - percentage) + '%';
            }
        });

        document.addEventListener('mouseup', () => {
            isResizing = false;
            document.body.style.cursor = '';
        });

        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey || e.metaKey) {
                switch(e.key) {
                    case 's':
                        e.preventDefault();
                        if (e.shiftKey) {
                            saveAsDialog();
                        } else {
                            saveFileDialog();
                        }
                        break;
                    case 'o':
                        e.preventDefault();
                        openFileDialog();
                        break;
                    case 'n':
                        e.preventDefault();
                        newFile();
                        break;
                }
            }
        });
    </script>
</body>
</html>`
}