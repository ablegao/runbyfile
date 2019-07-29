package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

var (
	run_count    = flag.Int("run_count", 0, "Run code , defualt CPU number -1")
	run_file     = flag.String("file", "", "run task by file")
	command_chan = make(chan string)
	once         = new(sync.WaitGroup)
)

func main() {
	flag.Parse()
	cpu := runtime.NumCPU()
	if *run_count == 0 {
		*run_count = cpu
	}

	for i := 0; i < *run_count; i++ {
		go runCommand()
	}

	file, err := os.Open(*run_file)

	if err != nil {
		log.Fatalln(err, "[", *run_file, "]")
		return
	}

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		command_chan <- string(line)
	}

	close(command_chan)
	once.Wait()
}

func runCommand() {
	once.Add(1)
	defer once.Done()
	for {
		line, ok := <-command_chan
		if !ok {
			log.Println(1)
			break
		}
		if line == "" {
			continue
		}
		cmd := exec.Command("/bin/sh", "-c", line)
		if cmd != nil {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			cmd.Wait()
		}
	}
	log.Println("runCommand:exit")
}
