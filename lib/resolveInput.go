package lib

import (
	"fmt"
	"strconv"
	"strings"
)

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func ResolveInput(lines []string, field *Field) error {
	for index, line := range lines {
		if line[0] != '#' {
			// Assume that the line is the number of ants
			if isNumber(line) {
				if field.Ants != nil {
					return fmt.Errorf("ants already defined")
				}

				numberOfAnts, err := strconv.Atoi(line)
				if err != nil {
					return err
				}

				err = numberToAnts(numberOfAnts, field)
				if err != nil {
					return err
				}
				continue
			}

			// Assume that the line is a link between two rooms
			if strings.Contains(line, "-") {
				params := strings.Split(line, "-")
				if len(params) == 2 {
					if err := linkRooms(params[0], params[1], field); err != nil {
						return err
					}

					continue
				} else {
					return fmt.Errorf("invalid line %d", index)
				}
			}

			// Assume that the line is a room with coordinates
			if strings.Contains(line, " ") {
				params := strings.Split(line, " ")
				if len(params) == 3 && index > 0 {
					isStart := false
					isEnd := false

					if lines[index-1] == "##start" {
						isStart = true
					} else if lines[index-1] == "##end" {
						isEnd = true
					}

					if err := addRoom(params[0], isStart, isEnd, field); err != nil {
						return err
					}

					continue
				}
			} else {
				return fmt.Errorf("invalid line %d", index)
			}
		}
	}

	return nil
}
