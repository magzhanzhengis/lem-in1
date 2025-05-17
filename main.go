package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name    string
	X, Y    int
	Links   []*Room
	Visited bool
	Parent  *Room
}

var (
	rooms     = make(map[string]*Room)
	startRoom *Room
	endRoom   *Room
	nAnts     int
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	filename := os.Args[1]
	lines, err := readLines(filename)
	if err != nil {
		fmt.Println("ERROR: could not read file")
		return
	}

	if !parseInput(lines) {
		fmt.Println("ERROR: invalid data format")
		return
	}

	path := bfs(startRoom, endRoom)
	if path == nil {
		fmt.Println("ERROR: no path found")
		return
	}

	simulateAnts(path)
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parseInput(lines []string) bool {
	var err error
	state := "ants"

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##")) {
			continue
		}

		switch line {
		case "##start":
			if i+1 >= len(lines) {
				return false
			}
			i++
			room := parseRoom(lines[i])
			if room == nil {
				return false
			}
			startRoom = room
			rooms[room.Name] = room

		case "##end":
			if i+1 >= len(lines) {
				return false
			}
			i++
			room := parseRoom(lines[i])
			if room == nil {
				return false
			}
			endRoom = room
			rooms[room.Name] = room

		default:
			if state == "ants" {
				nAnts, err = strconv.Atoi(line)
				if err != nil || nAnts <= 0 {
					return false
				}
				state = "rooms"
			} else if strings.Contains(line, " ") {
				room := parseRoom(line)
				if room == nil || rooms[room.Name] != nil {
					return false
				}
				rooms[room.Name] = room
			} else if strings.Contains(line, "-") {
				link := strings.Split(line, "-")
				if len(link) != 2 {
					return false
				}
				a, b := rooms[link[0]], rooms[link[1]]
				if a == nil || b == nil || a == b {
					return false
				}
				a.Links = append(a.Links, b)
				b.Links = append(b.Links, a)
			}
		}
	}
	return startRoom != nil && endRoom != nil
}

func parseRoom(line string) *Room {
	if strings.HasPrefix(line, "L") || strings.HasPrefix(line, "#") {
		return nil
	}
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil
	}

	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return nil
	}
	return &Room{Name: parts[0], X: x, Y: y}
}
