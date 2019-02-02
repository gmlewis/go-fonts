# Render open source fonts to polygons in Go

This is an experimental package used to render open source fonts to
polygons using Go.

## Example usage

To use one or more fonts within a Go program, import the main
package and the font(s) you want, like this:

```go
import (
  "github.com/gmlewis/go-fonts/fonts"
  _ "github.com/gmlewis/go-fonts/fonts/ubuntumonoregular"
)
```

Then render the text to polygons and use them however you want:

```go
  render, err := fonts.Text(xPos, yPos, xScale, yScale, message, "ubuntumonoregular")
  if err != nil {
    return err
  }
  log.Printf("MBB: (%.2f,%.2f)-(%.2f,%.2f)", render.Xmin, render.Ymin,render.Xmax, render.Ymax)
  for _, poly := range render.Polygons {
    // ...
  }
```

See https://github.com/gmlewis/go-gerber for an example application
that uses this package.

## Status
[![GoDoc](https://godoc.org/github.com/gmlewis/go-fonts/fonts?status.svg)](https://godoc.org/github.com/gmlewis/go-fonts/fonts)
[![Build Status](https://travis-ci.org/gmlewis/go-fonts.png)](https://travis-ci.org/gmlewis/go-fonts)

----------------------------------------------------------------------

Enjoy!

----------------------------------------------------------------------

# License

Copyright 2019 Glenn M. Lewis. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
