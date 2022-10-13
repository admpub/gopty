package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/admpub/gopty"
)

func main() {

	proc, err := gopty.New(120, 60)
	if err != nil {
		panic(err)
	}
	defer proc.Close()

	args := []string{gopty.GetBash(), gopty.GetFlagVar()}

	if runtime.GOOS == "windows" {
		args = append(args, "dir")
	} else {
		args = append(args, "ls -lah --color")
	}

	if err := proc.Start(args); err != nil {
		panic(err)
	}

	go func() {
		_, err = io.Copy(os.Stdout, proc)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()

	if _, err := proc.Wait(); err != nil {
		log.Printf("Wait err: %v\n", err)
	}
}
