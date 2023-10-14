package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Filter struct {
	Mask    *regexp.Regexp
	Letters []string
	// Includes []string
	// Excludes []string
	Includes map[string]interface{}
	Excludes map[string]interface{}
	// regexp *regexp.Regexp
}

func NewFilter() *Filter {
	filter := &Filter{
		Letters: []string{"[a-z]", "[a-z]", "[a-z]", "[a-z]", "[a-z]"},
		// Includes: make([]string, 0),
		// Excludes: make([]string, 0),
		Includes: make(map[string]interface{}),
		Excludes: make(map[string]interface{}),
	}
	filter.updateMask()

	return filter
}

func (this *Filter) updateMask() {
	this.Mask = regexp.MustCompile("^" + strings.Join(this.Letters, "") + "$")
}

func (this *Filter) Mark(pos int, letter string) {
	if len(letter) != 1 {
		log.Panicf("Invalid letter %s. Required length is 1", letter)
	}

	this.Letters[pos] = letter
	this.updateMask()
}

func (this *Filter) Include(word string) {
	for _, l := range word {
		// this.Includes = append(this.Includes, string(l))
		sl := string(l)
		_, ok := this.Includes[sl]
		if ok {
			delete(this.Includes, sl)
		} else {
			var i interface{}
			this.Includes[sl] = i
		}
	}
}

func (this *Filter) Exclude(word string) {
	for _, l := range word {
		// this.Excludes = append(this.Excludes, string(l))
		sl := string(l)
		_, ok := this.Excludes[sl]
		if ok {
			delete(this.Excludes, sl)
		} else {
			var i interface{}
			this.Excludes[sl] = i
		}
	}
}

func (this *Filter) filter(word string) bool {
	r := this.Mask.Match([]byte(word))

	for l, _ := range this.Includes {
		r = r && strings.Contains(word, l)
	}

	for l, _ := range this.Excludes {
		r = r && !strings.Contains(word, l)
	}

	return r
}

func List(filter *Filter) {
	file, err := os.Open("/usr/share/dict/words")

	if err != nil {
		log.Panic("Failed to read file ", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		word := scanner.Text()
		if filter.filter(word) {
			fmt.Println(word)
		}
	}
}

func main() {
	inReader := bufio.NewReader(os.Stdin)
	var cmd string
	var args []string
	f := NewFilter()

	for {
		fmt.Print("\n> ")
		input, err := inReader.ReadString('\n')

		if err != nil {
			log.Panic(err)
		}

		input = strings.Trim(input, "\n")

		cmd = strings.Split(input, " ")[0]
		args = strings.Split(input, " ")[1:]

		switch cmd {
		case "exit":
			os.Exit(0)
		case "list":
			List(f)
		case "include", "in":
			f.Include(args[0])
			fmt.Println("Included:", args[0])
		case "exclude", "ex":
			f.Exclude(args[0])
			fmt.Println("Excluded:", args[0])
		case "mark":
			i, e := strconv.Atoi(args[0])
			if e != nil {
				fmt.Println("Invalid position ", i)
			}
			f.Mark(i, args[1])
		case "mask":
			fmt.Print("Mask:", f.Letters)
			fmt.Print("\nIncludes: ")
			for k, _ := range f.Includes {
				fmt.Print(k, " ")
			}
			fmt.Print("\nExcludes: ")
			for k, _ := range f.Excludes {
				fmt.Print(k, " ")
			}
		case "clear":
			f = NewFilter()
		default:
			fmt.Println("Invalid command: ", cmd)
		}
	}
}
