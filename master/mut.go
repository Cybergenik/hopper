package master

import (
    "os"
    "flag"
    "fmt"
    "math/rand"
    "math/bits"
    "time"
)

const (
    N = 7
    MUT = 0
    DEL = 1
    ADD = 2
    SWP = 3
    FLP = 4
    REV = 5
    ROT = 6
)

func main(){
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run mut.go -N <number of mutations> <input file> ")
        os.Exit(1)
    }
    N := flag.Int("N", 1, "Number of mutations to produce")
    H := flag.Int("H", 1, "Level of Havoc, number of individual mutation steps per mutation")
    flag.Parse()

    bytes, err := os.ReadFile(flag.Arg(0))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for i:=0; i<*N; i++ {
        fmt.Printf("%s \n", Mutator(bytes, *H))
    }
}

func Mutator(b []byte, havoc int) []byte {
    rand.Seed(time.Now().UnixNano())
    if len(b) == 0 {
        panic("received empty seed")
    }
    bytes := append([]byte{}, b...)
    for i:=0; i<havoc;i++{
        switch rand.Intn(N){
        case MUT:
            i := rand.Intn(len(bytes))
            rval := make([]byte, 1)
            rand.Read(rval)
            bytes[i] = rval[0]
        case DEL:
            i := rand.Intn(len(bytes))
            bytes = append(bytes[:i], bytes[i+1:]...)
        case ADD:
            i := rand.Intn(len(bytes))
            rval := make([]byte, 1)
            rand.Read(rval)
            bytes = append(bytes[:i+1], bytes[i:]...)
            bytes[i] = rval[0]
        case SWP:
            i := rand.Intn(len(bytes))
            j := rand.Intn(len(bytes))
            for j == i {
                j = rand.Intn(len(bytes))
            }
            bytes[i], bytes[j] = bytes[j], bytes[i]
        case FLP:
            i := rand.Intn(len(bytes))
            bytes[i] = bytes[i]^255
        case REV:
            i := rand.Intn(len(bytes))
            bytes[i] = bits.Reverse8(bytes[i])
        }
    }
    return bytes
}

