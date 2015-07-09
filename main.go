package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Gopher struct {
	Id         int
	TotatlBW   int
	Files      []FileEntry
	Stdin      io.WriteCloser
	ListLoaded chan struct{}
	Done       chan struct{}
	Wg         sync.WaitGroup
}

type FileEntry struct {
	Name  string
	Size  int
	Value int
}

type FilesByValue []FileEntry
type FilesBySize []FileEntry

func hackMe(room string, g *Gopher, mainThread bool) {
	cmd := exec.Command("nc", "gophercon2015.coreos.com", "4001")
	var err error
	g.Stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer g.Stdin.Close()
	defer stdout.Close()
	//	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	io.WriteString(g.Stdin, room+"\n")
	io.WriteString(g.Stdin, "/msg Glenda help\n")
	io.WriteString(g.Stdin, "/list\n")
	scanner := bufio.NewScanner(stdout)
	//totalbw := 0
	//var files []FileEntry
	for scanner.Scan() {
		if true {
			rawline := scanner.Text()
			fmt.Println(rawline)
			if strings.Contains(rawline, "Secrecy") {
				continue
			}
			line := strings.Fields(scanner.Text())
			if strings.Contains(rawline, "Remaining Bandwidth") {
				//totalbw, _ = strconv.Atoi(line[5])
				g.TotatlBW, _ = strconv.Atoi(line[5])
				continue
			}
			if line[0] == "list" {
				size, _ := strconv.Atoi(strings.TrimSuffix(line[4], "KB"))
				value, _ := strconv.Atoi(line[5])
				fe := &FileEntry{
					Name:  line[3],
					Size:  size,
					Value: value,
				}
				g.Files = append(g.Files, *fe)
			}
			if len(g.Files) == 15 {
				sort.Sort(FilesByValue(g.Files))
				run := true
				i := 0
				for run == true {
					//if (totalbw - files[i].Size) < 0 {
					if g.TotatlBW < g.Files[i].Size {
						sort.Sort(FilesBySize(g.Files))
						for c := 0; c < len(g.Files); c++ {
							//Try to find an important file that fits
							fmt.Println(g.TotatlBW, "-", g.Files[c].Size)
							if g.TotatlBW > g.Files[c].Size {
								fmt.Println("Bonus: ", g.Files[c])
								ss := fmt.Sprintf("/send Glenda %s\n", g.Files[c].Name)
								io.WriteString(g.Stdin, ss)
								g.TotatlBW -= g.Files[c].Size
							}
						}
						run = false
						break
					}
					ss := fmt.Sprintf("/send Glenda %s\n", g.Files[i].Name)
					fmt.Println(g.Files[i])
					fmt.Println(len(g.Files))
					io.WriteString(g.Stdin, ss)
					g.TotatlBW -= g.Files[i].Size
					if len(g.Files) == 1 {
						//g.Files = append(g.Files[:0], g.Files[0:0]...)
						g.Files = g.Files[:0]
						run = false
						break
					}
					if len(g.Files) > 1 {
						g.Files = append(g.Files[:0], g.Files[1:]...)
					}
				}
				fmt.Println("Finished with ", g.TotatlBW, "KB to spare")
				io.WriteString(g.Stdin, "/msg Glenda done\n")
				g.Wg.Done()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard output:", err)
	}
	cmd.Wait()
}

func main() {
	room := "myroom"
	i := 0
	var wg sync.WaitGroup
	for i = 0; i < 2; i++ {
		wg.Add(1)
		g := &Gopher{
			Id: i + 1,
			Wg: wg,
		}
		go hackMe(room, g, false)
	}
	wg.Add(1)
	g := &Gopher{
		Id: i + 1,
		Wg: wg,
	}
	hackMe(room, g, true)
	//Calculate And move
}

func (s FilesByValue) Len() int {
	return len(s)
}
func (s FilesByValue) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s FilesByValue) Less(i, j int) bool {
	return (s[i].Value > s[j].Value)
}

func (s FilesBySize) Len() int {
	return len(s)
}
func (s FilesBySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s FilesBySize) Less(i, j int) bool {
	return (s[i].Size < s[j].Size)
}
