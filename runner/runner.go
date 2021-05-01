package main

import (
	"log"
	"os/exec"
)

func main() {
	result := make(chan []byte)
	commands := make([]*exec.Cmd, 12)
	for i := 0; i < 12; i += 1 {
		cmd := exec.Command("./seed")
		commands[i] = cmd
		go func(cmd *exec.Cmd) {
			out, err := cmd.Output()
			if err != nil {
				log.Fatalf("Seed failed to execute. %v\n", err)
			}
			result <- out
		}(cmd)
	}
	log.Println(string(<-result))
	for _, cmd := range commands {
		cmd.Process.Kill()
	}
}
