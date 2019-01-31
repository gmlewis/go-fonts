# inconsolata

![inconsolata](inconsolata.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/inconsolata"
)

func main() {
	// ...
	render, err := Text(x, y, xs, ys, message, "inconsolata"),
	// ...
}
```
