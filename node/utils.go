package node

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Deb sid: "sancov-15"
const SANCOV = "sancov"

func ParseAsan(asan string) string {
	asan_lines := strings.Split(asan, "\n")
	for _, line := range asan_lines {
		sline := strings.Split(line, " ")
		if sline[0] == "SUMMARY:" {
			return sline[2]
		}
	}
	return ""
}

func PersistCrash(seed []byte, asan bytes.Buffer, crashN uint64, outDir string) {
	out := path.Join(outDir, fmt.Sprintf("crash%d", crashN))
	//Write input
	input := bytes.NewBuffer(seed)
	err := os.WriteFile(fmt.Sprintf("%s.in", out), input.Bytes(), 0444)
	if err != nil {
		log.Fatal(err)
	}
	//Write ASAN data
	report := bytes.NewBufferString("ASAN:\n\n")
	report.Write(asan.Bytes())
	err = os.WriteFile(fmt.Sprintf("%s.report", out), report.Bytes(), 0444)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCoverage(sancov_file string) ([]string, bool) {
	cov_cmd := exec.Command(SANCOV,
		"--print",
		sancov_file,
	)
	var out bytes.Buffer
	cov_cmd.Stdout = &out
	if err := cov_cmd.Run(); err != nil {
		return nil, false
	}
	go os.Remove(sancov_file)
	// Coverage tree parsing
	covered := strings.Split(strings.Trim(out.String(), "\n"), "\n")
	for i, v := range covered {
		covered[i] = v[2:]
	}
	return covered, true
}
