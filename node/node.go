package main

import (
    "os"
    "os/exec"
    "fmt"
    "log"
    "flag"
    "strings"
    "path/filepath"
    "net/rpc"
    "bytes"
    "github.com/tidwall/gjson"
    c "github.com/Cybergenik/hopper/common"
)

//Deb sid: "sancov-15"
const SANCOV = "sancov"

type HopperNode struct {
    name        string
    id          int
    target      string
    args        string
    env         []string
    stdin       bool
    master      string
    crashN      int
    conn        *rpc.Client
}

func (n *HopperNode) getFTask() (c.FTask, bool) {
    args := c.FTaskArgs{}
    t := c.FTask{}

    if ok := n.call("Hopper.GetFTask", &args, &t); !ok {
        log.Println("Error Getting FTask!")
        return t, ok
    }
    return t, true
}

func (n *HopperNode) updateFTask(ut c.UpdateFTask) {
    reply := c.UpdateReply{}
    if ok := n.call("Hopper.UpdateFTask", &ut, &reply); !ok {
        log.Println("Error Updating FTask!")
    }
}

func parseAsan(asan string) string{
    asan_lines := strings.Split(asan, "\n")
    for _, line := range asan_lines {
        sline := strings.Split(line, " ")
        if sline[0] == "SUMMARY:"{
            return sline[2]
        }
    }
    return ""
}

func (n *HopperNode) persistCrash(seed []byte, asan bytes.Buffer, crashN int) {
    report := bytes.NewBufferString("Seed:\n\n")
    report.Write(seed)
    report.WriteString("\n\nASAN:\n\n")
    report.Write(asan.Bytes())
    err := os.WriteFile(fmt.Sprintf("%s/crash%d", n.name, crashN), report.Bytes(), 0660)
    if err != nil {
        log.Fatal(err)
    }
}

func (n *HopperNode) getCov(cov_cmd *exec.Cmd, sancov_file string) ([]string, bool){
    var out bytes.Buffer
    cov_cmd.Stdout = &out
    if err := cov_cmd.Run(); err != nil {
        return nil, false
    }
    go os.Remove(sancov_file)
    // Coverage tree parsing
    covered := gjson.Get(out.String(), "covered-points").Array()
    cov_s := []string{}
    for _, v := range covered {
        edge := gjson.Get(out.String(), fmt.Sprintf("point-symbol-info.*.*.%v", v.Value()))
        cov_s = append(cov_s, fmt.Sprintf("%s", edge.Value()))
    }
    return cov_s, true
}

func (n *HopperNode) fuzz(t c.FTask) {
    //Run seed
    var fuzzCommand []string
    if n.stdin {
        fuzzCommand = strings.Split(n.args, " ")
    } else {
        fuzzCommand = strings.Split(
            strings.Replace(
                n.args,
                "@@",
                string(t.Seed),
                1,
            ),
            " ",
        )
    }
    update := c.UpdateFTask{
        NodeId: n.id,
        Id:     t.Id,
    }
    cmd := exec.Command(n.target, fuzzCommand...)
    cmd.Env = append(os.Environ(), n.env...)
    // Gather err output
    var errOut bytes.Buffer
    var stdin  bytes.Buffer
    cmd.Stderr = &errOut
    if n.stdin {
        cmd.Stdin = &stdin
    }
    if err := cmd.Start(); err != nil{
        log.Println(err)
        update.Ok = false
        go n.updateFTask(update)
        cmd.Wait()
        return
    }
    if n.stdin {
        stdin.Write(t.Seed)
    }
    err := cmd.Wait()
    sancov_file := fmt.Sprintf("%s.%v.sancov",
        filepath.Base(n.target),
        cmd.Process.Pid,
    )
    //Crash Detected
    if err != nil {
        update.Crash = parseAsan(errOut.String())
        go n.persistCrash(t.Seed, errOut, n.crashN)
        n.crashN++
    }
    cov_cmd := exec.Command(SANCOV,
        "--symbolize",
        sancov_file,
        n.target,
    )
    //Generate Coverage data
    cov_s, ok := n.getCov(cov_cmd, sancov_file)
    update.Ok = ok
    if !ok {
        go n.updateFTask(update)
    } else {
        update.CovEdges = len(cov_s)
        update.CovHash = c.Hash([]byte(strings.Join(cov_s, "-")))
        go n.updateFTask(update)
    }
}

func Node(id int, target string, args string, env string, stdin bool, master string) {
    n := HopperNode{
        name:    fmt.Sprintf("Node%d", id),
        id:      id,
        target:  target,
        args:    args,
        env:     strings.Split(env, ";"),
        stdin:   stdin,
        master:  master,
        crashN:  0,
    }
    
    // Check target executable exists
    if _, err := os.Stat(target); err != nil {
        log.Fatal(err)
    }
    // Env vars
    n.env = append(n.env, "ASAN_OPTIONS=coverage=1")
    // Init TCP/IP connection to master
    c, err := rpc.DialHTTP("tcp", n.master)
    if err != nil {                 
        log.Fatal("dialing:", err)
    }
    n.conn = c
    defer n.conn.Close()
    // Create node out dir
    if err := os.MkdirAll(n.name, 0750); err != nil && !os.IsExist(err) {
        log.Fatal(err)
    }
    //Infinite loop, request Task -> do Task
    log.Printf("Started Node: %v\n", id)
    for {
        ftask, ok := n.getFTask()
        if !ok || ftask.Die {
            return
        }
        n.fuzz(ftask)
    }
}

func main() {
    id      := flag.Int("I", 0, "Node ID, usually just a unique int")
    target  := flag.String("T", "", "instrumented target binary")
    args    := flag.String("args", "", "args to use against target, ex: --depth=1 @@")
    env     := flag.String("env", "", "env variables for target seperated by a `;`, ex: ARG1=foo;ARG2=bar;")
    stdin   := flag.Bool("stdin", false, "seed should be fed as stdin or as an argument")
    master  := flag.String("M", "localhost", "IP/address of Master")
    port    := flag.Int("P", 6969, "Port of Master")
        
    flag.Parse()
    err := ""
    if *id == 0 {
        err += "Hopper Node: Please provide a unique Node Id greater than 0: -I \n"
    }
    if *target == "" {
        err += "Hopper Node: Please provide a command to run the target: -C\n"
    }
    if !strings.Contains(*args, "@@") && !*stdin {
        err += "Hopper Node: must provide an @@ input in args if not using stdin: ex --stdin or --args @@\n"
    }
    if err != "" {
        fmt.Println(err)
        os.Exit(1)
    }
    _, Err := exec.LookPath(SANCOV)
	if Err != nil {
        log.Fatalf("Hopper Node: Node requires clang-tools utils: sanvoc")
	}
    masterNode := fmt.Sprintf("%s:%d", *master, *port)
    Node(*id, *target, *args, *env, *stdin, masterNode)
}

func (n *HopperNode) call(rpcname string, args interface{}, reply interface{}) bool {
    err := n.conn.Call(rpcname, args, reply)
    if err != nil {                    
        log.Print(err)
        return false
    }                                                          

    return true           
}
    
