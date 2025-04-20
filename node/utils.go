package node

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Deb sid: "sancov-15"
var SANCOV = "sancov"

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

func GetCoverage(pid int) ([]string, bool) {
	sancovGlob := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("*.%v.sancov", pid),
	)
	sancovMatches, err := filepath.Glob(sancovGlob)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	allCoveredEdges := []string{}
	for _, sancovFile := range sancovMatches {
		defer func(f string) {
			if err := os.Remove(f); err != nil {
				log.Println(err)
			}
		}(sancovFile)
		cov_cmd := exec.Command(SANCOV,
			"--print",
			sancovFile,
		)
		var out bytes.Buffer
		cov_cmd.Stdout = &out
		if err := cov_cmd.Run(); err != nil {
			log.Println(err)
			continue
		}
		// Coverage tree parsing
		covered := strings.Split(strings.Trim(out.String(), "\n"), "\n")
		for _, v := range covered {
			if len(v) > 2 {
				allCoveredEdges = append(allCoveredEdges, v[2:])
			}
		}
	}
	return allCoveredEdges, len(allCoveredEdges) > 0
}
