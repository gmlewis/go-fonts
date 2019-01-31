# acme_regular

![acme_regular](acme_regular.png)

To use this font in your code, simply import it:

```go
import (
	. "github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/acme_regular"
)

func main() {
	// ...
	render, err := Text(x, y, xs, ys, message, "acme_regular"),
	// ...
}
```
