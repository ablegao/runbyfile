package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	run_count    = flag.Int("run_count", 0, "Run code , defualt CPU number -1")
	run_file     = flag.String("file", "", "run task by file")
	command_chan chan string
	once         = new(sync.WaitGroup)
	startTask    = new(sync.WaitGroup)
)

func runTaskRuntinue(count int) {
	log.Println("count:", count)
	startTask.Add(1)
	defer startTask.Done()
	command_chan = make(chan string)
	for i := 0; i < count; i++ {
		go runCommand(i)
	}
}

func writeChan(s string) {
	command_chan <- s
}
func closeChan() {
	if command_chan == nil {
		return
	}
	select {
	case _, ok := <-command_chan:
		if !ok {

		}
	default:
		close(command_chan)
	}
}
func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	cpu := runtime.NumCPU() - 1
	if *run_count == 0 {
		*run_count = cpu
	}
	// runTaskRuntinue()

	file, err := os.Open(*run_file)

	if err != nil {
		log.Fatalln(err, "[", *run_file, "]")
		return
	}

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			closeChan()
			break
		}

		s_line := string(line)
		line_len := len(s_line)
		if line_len == 0 {
			continue
		}
		if line_len >= 4 {
			switch s_line[:4] {
			case "run ": // run 5
				out := strings.Split(s_line, " ")
				if len(out) == 2 {
					l, err := strconv.Atoi(out[1])
					if err != nil {
						log.Println(err)
						break
					}
					if l == 0 {
						l = *run_count
					}
					runTaskRuntinue(l)
					startTask.Wait() // 等待runTask 启动完成
				}
			case "exit": // 执行到退出指令时
				closeChan()
				once.Wait()
			default:
				writeChan(s_line)
			}

		} else {
			writeChan(s_line)
		}
	}

	once.Wait()
}

func runCommand(i int) {
	once.Add(1)
	log.Println("runCommand:start", i)
	defer once.Done()
	for {
		line, ok := <-command_chan
		if !ok {
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
	log.Println("runCommand:exit", i)
}
