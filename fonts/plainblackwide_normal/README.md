# plainblackwide_normal

![plainblackwide_normal](plainblackwide_normal.png)

To use this font in your code, simply import it:

```go
import (
  . "github.com/gmlewis/go-fonts/fonts"
  _ "github.com/gmlewis/go-fonts/fonts/plainblackwide_normal"
)

func main() {
  // ...
  render, err := fonts.Text(xPos, yPos, xScale, yScale, message, "plainblackwide_normal")
  if err != nil {
    return err
  }
  log.Printf("MBB: (%.2f,%.2f)-(%.2f,%.2f)", render.Xmin, render.Ymin,render.Xmax, render.Ymax)
  for _, poly := range render.Polygons {
    // ...
  }
  // ...
}
```