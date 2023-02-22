package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func shell(cmd []string) string {
	stdout, _ := exec.Command(cmd[0], cmd[1:]...).Output()
	return strings.Trim(string(stdout), " \n")
}

func processLine(line string, output chan<- string) {
	pair := strings.SplitN(line, ": ", 2)
	switch n := len(pair); n {
	case 2:
		cmd := strings.Split(strings.Trim(pair[1], " "), " ")
		output <- pair[0] + ": " + shell(cmd)
	case 1:
		output <- pair[0]
	case 0:
		output <- ""
	}
}

func coRun(input <-chan string, output chan<- string, count int, process func(string, chan<- string)) {
	buf := make(chan string, count)
	n := 0

	for input != nil {
		if n == count {
			output <- <-buf
			n--
		}

		line, ok := <-input
		if ok {
			go process(line, buf)
			n++
		} else {
			input = nil
		}
	}

	for ; n > 0; n-- {
		output <- <-buf
	}

	close(buf)
	close(output)
}

func readFileToChannel(fileName string, lines chan<- string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	close(lines)
}

func writeChannelToFile(fileName string, lines <-chan string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for lines != nil {
		line, ok := <-lines
		if ok {
			_, err := fmt.Fprintln(writer, line)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			lines = nil
		}
	}
	writer.Flush()
}

func main() {
	inFile := flag.String("in", "unknown_file", "the input file")
	outFile := flag.String("out", "unknown_file", "the output file")
	count := flag.Int("np", 10, "the number of concurrent processors")
	flag.Parse()

	if *count <= 0 {
		log.Fatalf("Invalid number of concurrent processors: %d", *count)
	}

	input := make(chan string)
	output := make(chan string)

	go readFileToChannel(*inFile, input)
	go coRun(input, output, *count, processLine)

	writeChannelToFile(*outFile, output)
}
