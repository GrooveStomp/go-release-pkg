package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func notEnoughArgs() {
	fmt.Fprintln(os.Stderr, "Not enough arguments")
	fmt.Fprintln(os.Stderr, "Try --help")
	os.Exit(0)
}

func usage() {
	text := `
Usage: go run build.go <os>
       go run build.go release <os>

Builds confluence_tool according to <os> parameter.
if 'release' is specified, then package the executable as a gzipped tar file.

os:

     osx: OSX 64-bit x86
     linux: Linux 64-bit x86

Options:

     help: Show this help text
`
	fmt.Println(text)
	os.Exit(0)
}

func buildExecutable(osParam, outputName string) {
	cmd := exec.Command("go", "build", "-o", outputName, "main.go")
	cmd.Env = append(os.Environ(), "GOARCH=amd64")

	if osParam == "osx" {
		cmd.Env = append(cmd.Env, "GOOS=darwin")
	} else if osParam == "linux" {
		cmd.Env = append(cmd.Env, "GOOS=linux")
	} else {
		fmt.Println("Unsupported OS specified.")
		fmt.Fprintln(os.Stderr, "Try --help")
		os.Exit(0)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func tarGzip(src string) *bytes.Buffer {
	fi, err := os.Stat(src)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	gzw := gzip.NewWriter(buffer)
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	header := &tar.Header{
		Name: fi.Name(),
		Mode: 0774,
		Size: int64(fi.Size()),
	}

	if err := tw.WriteHeader(header); err != nil {
		panic(err)
	}

	f, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(tw, f); err != nil {
		panic(err)
	}

	return buffer
}

func main() {
	flagHelp := flag.Bool("help", false, "Show help")
	flag.Parse()

	args := flag.Args()

	if *flagHelp {
		usage()
	}

	if len(args) < 1 {
		notEnoughArgs()
	}

	doRelease := false
	osArg := os.Args[1]
	if osArg == "help" {
		usage()
	} else if osArg == "release" {
		doRelease = true
		if len(args) < 2 {
			notEnoughArgs()
		}
		osArg = os.Args[2]
	}

	executableName := "confluence_tool"
	buildExecutable(osArg, executableName)

	if doRelease {
		buf := tarGzip(executableName)

		outfile := fmt.Sprintf("%v.%v.tar.gz", executableName, osArg)
		err := ioutil.WriteFile(outfile, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}

		err = os.Remove(executableName)
		if err != nil {
			panic(err)
		}
	}
}
