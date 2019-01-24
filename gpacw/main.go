package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/darren/gpac"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s test.pac curl example.com\n", os.Args[0])
	os.Exit(1)
}

func run(cmds string, oargs []string, proxy string) error {
	var args = oargs
	if proxy != "" {
		switch cmds {
		case "wget":
			args = append([]string{
				"-e", "http_proxy=" + proxy,
				"-e", "https_proxy=" + proxy,
			}, oargs...)
		case "curl":
			args = append([]string{"-x", proxy}, oargs...)
		}
	}

	log.Printf("Invoke %s %v", cmds, args)

	cmd := exec.Command(cmds, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	go func() {
		io.Copy(os.Stdout, stdout)
	}()

	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		usage()
	}

	pacf := os.Args[1]
	cmds := os.Args[2]
	args := os.Args[3:]
	dst := os.Args[len(os.Args)-1]

	if strings.HasPrefix(dst, "http://") &&
		strings.HasPrefix(dst, "https://") {
		dst = "http://" + dst
	}

	switch cmds {
	case "curl", "wget":
	default:
		log.Fatalf("Only curl or wget is supported, %s is given", cmds)
	}

	pac, err := gpac.From(pacf)
	if err != nil {
		log.Fatalf("Fail to load pac file: %v", err)
	}

	p, err := pac.FindProxy(dst)
	if err != nil {
		log.Fatal(err)
	}

	for _, x := range p {
		err = run(cmds, args, x.URL())
		if err == nil {
			break
		}
	}
}
