# ostrichsansblack

![ostrichsansblack](ostrichsansblack.png)

To use this font in your code, simply import it:

```go
import (
  . "github.com/gmlewis/go-fonts/fonts"
  _ "github.com/gmlewis/go-fonts/fonts/ostrichsansblack"
)

func main() {
  // ...
  render, err := fonts.Text(xPos, yPos, xScale, yScale, message, "ostrichsansblack", Center)
  if err != nil {
    return err
  }
  log.Printf("MBB: %v", render.MBB)
  for _, poly := range render.Polygons {
    // ...
  }
  // ...
}
```
