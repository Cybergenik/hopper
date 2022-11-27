package main

import (
    "fmt"
    "os"
    "plugin"
    "log"
    "github.com/tidwall/gjson"
    "strings"
    "os/exec"
    "bytes"
    "path/filepath"
    //"container/ring"
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

func loadMutEngine(filename string) (func([]byte, int) []byte) {
    p, err := plugin.Open(filename)
    if err != nil {
        log.Fatalf("cannot load plugin %v", filename)
    }

    xmutf, err := p.Lookup("Mutator")
    if err != nil {
        log.Fatalf("cannot find Mutator in %v", filename)
    }
    
    return xmutf.(func([]byte, int) []byte)
}


func main() {

    fmt.Println("Hello there!")


    //loadMutEngine(os.Args[1])
    // Testing ring buffers
    //rbuf := ring.New(10)

    //for i:=0;i<6;i++{
    //    rbuf.Value = i
    //    rbuf = rbuf.Next()
    //}
    //
    //for i:=0;i<100;i++{
    //    fmt.Println(rbuf.Value, rbuf.Len())
    //    rbuf = rbuf.Next()
    //}

    //Replace string seed inputs
    s := "Hello World @@"
	s1 := strings.Replace(s, "@@", "vape juice", 1)
	fmt.Println(s, s1)


    // exec fuzz command
    cmd := exec.Command(os.Args[1], "te7^t@hellP0j2saB{.comN]")
    cmd.Env = append(os.Environ(),
        "ASAN_OPTIONS=coverage=1",
    )
    var errOut bytes.Buffer
    cmd.Stderr = &errOut
    err := cmd.Run()
    sancov_file := fmt.Sprintf("%s.%v.sancov", filepath.Base(os.Args[1]), cmd.Process.Pid)
    if err != nil {
        fmt.Println(parseAsan(errOut.String()))
    }
    cov_cmd := exec.Command("sancov-15",
        "--symbolize",
        sancov_file,
        os.Args[1],
    )
    var out bytes.Buffer
    cov_cmd.Stdout = &out
    if err := cov_cmd.Run(); err != nil {
        fmt.Println(err)
    }
    // Coverage tree parsing
    covered := gjson.Get(out.String(), "covered-points").Array()
    fmt.Printf("Edges: %v\n", len(covered))
    for _, v := range covered {
        fmt.Println(v.Value())
        val := gjson.Get(out.String(), fmt.Sprintf("point-symbol-info.*.*.%v", v.Value()))
        fmt.Println(val.Value())
    }
    err = os.MkdirAll("Node1", 0750)
    if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	err = os.WriteFile("Node1/crash1", out.Bytes(), 0660)
	if err != nil {
		log.Fatal(err)
	}
}
