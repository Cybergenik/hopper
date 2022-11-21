package main

import (
    "fmt"
    "github.com/tidwall/gjson"
)

const json = `
{
"covered-points": [
  "4c6ad2",
  "4c6b10",
  "4c6b64"
],
"binary-hash": "67EF8FFC8C47FD5CD70061E06EB8393A0092252F",
"point-symbol-info": {
  "/home/luciano/hopper-test/cov.cc": {
    "foo()": {
      "4c6ad2": "4:0"
    },
    "main": {
      "4c6b10": "6:0",
      "4c6b47": "7:9",
      "4c6b64": "8:9"
    }
  }
}
`

func main() {

    fmt.Println("Hello there!")
    covered := gjson.Get(json, "covered-points").Array()
    for _, v := range covered {
        fmt.Println(v.Value())
        val := gjson.Get(json, fmt.Sprintf("point-symbol-info.*.*.%v", v.Value()))
        fmt.Println(val.Value())
    }

}
