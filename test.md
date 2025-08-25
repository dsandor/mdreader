# Test Markdown File

This is a test markdown file to demonstrate GitHub-style rendering with syntax highlighting.

## Features

- **Bold text** and *italic text*
- ~~Strikethrough text~~
- [Links](https://github.com)
- `Inline code`

### Code Blocks with Syntax Highlighting

#### Go Code
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    for i := 0; i < 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

#### JavaScript Code
```javascript
function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

const result = fibonacci(10);
console.log(`Fibonacci(10) = ${result}`);
```

#### Python Code
```python
def quicksort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quicksort(left) + middle + quicksort(right)

numbers = [3, 6, 8, 10, 1, 2, 1]
print(f"Sorted: {quicksort(numbers)}")
```

## Tables

| Language | Extension | Popularity |
|----------|-----------|------------|
| Go       | .go       | High       |
| Python   | .py       | Very High  |
| JavaScript | .js     | Very High  |
| Rust     | .rs       | Growing    |

## Blockquotes

> This is a blockquote
> with multiple lines
> 
> And another paragraph in the quote

## Lists

### Ordered List
1. First item
2. Second item
   1. Nested item
   2. Another nested item
3. Third item

### Unordered List
- Item one
- Item two
  - Nested item
  - Another nested
- Item three

---

That's all for the test!