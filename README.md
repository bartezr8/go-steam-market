# CSGO Steam Market Pricing API Wrapper

This simply wraps the Steam Market Pricing API to enable doing json requests to the market from Go

# Usage

Start off by getting this package:

`go get github.com/Gacnt/go-steam-market`


To use this there is one simple function:

```
package main

import (
        "fmt"
        "github.com/Gacnt/go-steam-market"
)

func main() {
        jsonResp := gosm.GetSinglePrice(false, "M4A1-S", "Master Piece", gosm.FT, "G")
        fmt.Println(jsonResp)
}

// Response:
{true $88.69 39 $80.95}

```

The function takes 5 parameters:

```
GetSinglePrice(StatTrak Bool, "Weapon Type", "Skin Name", "Skin Wear", "G for Gun or K for Knife")
```

You can see the constants that are the skin types to pass to the function [here](http://godoc.org/github.com/Gacnt/go-steam-market#pkg-constants)

If item is a knife with no skin, e.g. it's JUST a Gut Knife just put an empty string `""` for the Skin Wear parameters as Golang does not support optional Params


You can see more about the API here: [GoDoc](http://godoc.org/github.com/Gacnt/go-steam-market)
