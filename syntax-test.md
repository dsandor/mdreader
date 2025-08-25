# Syntax Highlighting Test

This file tests syntax highlighting for various programming languages.

## Go Code

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

## JavaScript Code

```javascript
function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

const result = fibonacci(10);
console.log(`Fibonacci(10) = ${result}`);
```

## Python Code

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

## JSON

```json
{
  "name": "mdreader",
  "version": "1.0.0",
  "description": "Markdown to HTML converter",
  "features": ["syntax highlighting", "GitHub style", "live preview"]
}
```

## Plain text block

```
This is plain text without syntax highlighting
It should appear in a monospace font
But without any color coding
```