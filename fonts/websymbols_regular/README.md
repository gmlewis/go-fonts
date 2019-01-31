# websymbols_regular

![websymbols_regular](websymbols_regular.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/websymbols_regular"
)

func main() {
	// ...
	render, err := Text(x, y, xs, ys, message, "websymbols_regular"),
	// ...
}
```
