package main

import (
    "os"
    "os/exec"
    "fmt"
    "log"
    "flag"
    "strings"
    "strconv"
    "path/filepath"
    "net/rpc"
    "bytes"
    "github.com/tidwall/gjson"
    c "github.com/Cybergenik/hopper/common"
)

type HopperNode struct {
    name        string
    id          int
    target      string
    args        string
    env         []string
    stdin       bool
    master      string
    crashN      int
}

func (n *HopperNode) getFTask() c.FTask {
    args := c.FTaskArgs{}
	t := c.FTask{}
	if ok := n.call("Hopper.GetFTask", &args, &t); !ok {
		log.Println("Error Getting FTask!")
        t.Id = 0
	}
	return t
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
    err := os.WriteFile(n.name+"/"+"crash"+strconv.Itoa(crashN), report.Bytes(), 0660)
    if err != nil {
        log.Fatal(err)
    }
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
    //Generate Coverage data
    cov_cmd := exec.Command("sancov-15",
        "--symbolize",
        sancov_file,
        n.target,
    )
    var out bytes.Buffer
    cov_cmd.Stdout = &out
    if err := cov_cmd.Run(); err != nil {
        log.Println(err)
        update.Ok = false
        go n.updateFTask(update)
        return
    }
    go os.Remove(sancov_file)
    // Coverage tree parsing
    covered := gjson.Get(out.String(), "covered-points").Array()
    update.CovEdges = len(covered)
    cov_s := []string{}
    for _, v := range covered {
        edge := gjson.Get(out.String(), fmt.Sprintf("point-symbol-info.*.*.%v", v.Value()))
        cov_s = append(cov_s, fmt.Sprintf("%s", edge.Value()))
    }
    update.CovHash = c.Hash([]byte(strings.Join(cov_s, "-")))
    update.Ok = true
    go n.updateFTask(update)
}

func Node(id int, target string, args string, env string, stdin bool, master string) {
    n := HopperNode{
        name:    "Node"+strconv.Itoa(id),
        id:      id,
        target:  target,
        args:    args,
        env:     strings.Split(env, ";"),
        stdin:   stdin,
        master:  master,
        crashN:  0,
    }
    err := os.MkdirAll(n.name, 0750)
    if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
    n.env = append(n.env, "ASAN_OPTIONS=coverage=1")
    //Infinite loop, request Task -> do Task
	for {
		ftask := n.getFTask()
        for ftask.Id == 0 {
            ftask = n.getFTask()
        }
		if ftask.Die {
			os.Exit(0)
		}	
        //fmt.Printf("Fuzzing: %s\n", ftask.Seed)
        n.fuzz(ftask)
    }
}

func main() {
    id      := flag.Int("I", 0, "Node ID, usually just a unique int")
    target  := flag.String("T", "", "instrumented target binary")
    args    := flag.String("args", "", "args to use against target, ex: --depth=1 @@")
    env     := flag.String("env", "", "env variables for target seperated by a `;`, ex: ARG1=foo;ARG2=bar;")
    stdin   := flag.Bool("stdin", false, "seed should be fed as stdin or as an argument")
    master  := flag.String("M", "localhost", "instrumented target binary")
    port    := flag.Int("P", 6969, "instrumented target binary")
        
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
        log.Fatal(err)
    }
    fmt.Printf("Starting Node: %v\n", *id)
    //fmt.Printf("Args: %v\n", *args)
    masterNode := *master+":"+strconv.Itoa(*port)
    Node(*id, *target, *args, *env, *stdin, masterNode)
}

func (n *HopperNode) call(rpcname string, args interface{}, reply interface{}) bool {
    c, err := rpc.DialHTTP("tcp", n.master)
    if err != nil {                 
        log.Fatal("dialing:", err)
    }
    defer c.Close()       
                                                            
    err = c.Call(rpcname, args, reply)
    if err != nil {                    
        log.Print(err)
        return false
    }                                                          

    return true           
}
    
