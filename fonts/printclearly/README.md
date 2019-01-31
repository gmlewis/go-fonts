# printclearly

![printclearly](printclearly.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/printclearly"
)

func main() {
	// ...
	render, err := Text(x, y, xs, ys, message, "printclearly"),
	// ...
}
```
