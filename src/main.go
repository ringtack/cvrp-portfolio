package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ALNS_CVRP = filepath.Join("cvrp", "run.sh")
	HGS_CVRP  = filepath.Join("cvrp-gen", "run.sh")
	CMD       = "./run.sh"
)

var input = flag.String("input", "input.txt", "Input file for the CVRP problem")

// Define output struct format
type CVRPOutput struct {
	Instance string  `json:"Instance"`
	Time     float64 `json:"Time"`
	Result   float64 `json:"Result"`
	Solution string  `json:"Solution"`
}

func main() {
	flag.Parse()
	if _, err := os.Stat(*input); err != nil {
		log.Fatalf("Input %v does not exist: %v", *input, err)
	}

	if _, err := os.Stat(ALNS_CVRP); err != nil {
		log.Fatalf("%v does not exist: %v", ALNS_CVRP, err)
	}
	if _, err := os.Stat(HGS_CVRP); err != nil {
		log.Fatalf("%v does not exist: %v", HGS_CVRP, err)
	}

	log.Printf("Running CVRP on %v", *input)

	// Run two processes, cvrp/run.sh and cvrp-gen/run.sh, and collect their outputs;
	// start both in go-routines, and collect output
	res := make(chan string, 2)
	go func() {
		// Change directory to the cvrp-gen directory
		hgsCmd := exec.Command(CMD, *input)
		hgsCmd.Dir = filepath.Dir(HGS_CVRP)
		log.Printf("Dir: %v", hgsCmd.Dir)
		hgsOut, err := hgsCmd.CombinedOutput()
		if err != nil {
			log.Printf("Output: %v", string(hgsOut))
			log.Fatalf("Error running %v: %v", HGS_CVRP, err)
		}
		res <- string(hgsOut)
	}()
	go func() {
		alnsCmd := exec.Command(CMD, *input)
		alnsCmd.Dir = filepath.Dir(ALNS_CVRP)
		log.Printf("Dir: %v", alnsCmd.Dir)
		alnsOut, err := alnsCmd.CombinedOutput()
		if err != nil {
			log.Printf("Output: %v", string(alnsOut))
			log.Fatalf("Error running %v: %v", ALNS_CVRP, err)
		}
		res <- string(alnsOut)
	}()

	// Wait for both to finish
	hgsOut := <-res
	alnsOut := <-res

	// Split into lines, and get last line
	hgsLines := strings.Split(hgsOut, "\n")
	hgsLine := hgsLines[len(hgsLines)-2]
	fmt.Printf("HGS: %v\n", hgsLine)
	alnsLines := strings.Split(alnsOut, "\n")
	alnsLine := alnsLines[len(alnsLines)-2]
	fmt.Printf("ALNS: %v\n", alnsLine)

	// Parse into JSONs
	var hgsOutput, alnsOutput CVRPOutput
	if err := json.Unmarshal([]byte(hgsLine), &hgsOutput); err != nil {
		log.Fatalf("Error parsing HGS output %v: %v", hgsLine, err)
	}
	if err := json.Unmarshal([]byte(alnsLine), &alnsOutput); err != nil {
		log.Fatalf("Error parsing ALNS output %v: %v", alnsLine, err)
	}

	// If HGS is better, print HGS, else print ALNS
	if hgsOutput.Result < alnsOutput.Result {
		fmt.Println(hgsLine)
	} else {
		fmt.Println(alnsLine)
	}
}
