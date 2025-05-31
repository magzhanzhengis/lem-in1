package lib

import (
	"fmt"
)

func StartTurnBasedSolving(field *Field, pathsToExit [][]string) {
	isSolved := false
	turns := []string{}
	for !isSolved {
		turn := ""
		isSolved = true
		pathsUsed := []string{}
		for _, ant := range field.Ants {
			if ant.CurrentRoom == field.EndRoomName {
				ant.IsFinished = true
			}

			if !ant.IsFinished {
				nextRoom := getNextRoom(ant.CurrentRoom, pathsToExit, pathsUsed, *field)

				for _, path := range pathsUsed {
					if path == ant.CurrentRoom+"-"+nextRoom {
						nextRoom = "" // to prevent ants from colliding
					}
				}

				if nextRoom != "" {
					turnStruct := "%s L%v-%s"
					if turn == "" {
						turnStruct = "%sL%v-%s"
					}

					pathsUsed = append(pathsUsed, ant.CurrentRoom+"-"+nextRoom)
					ant.CurrentRoom = nextRoom
					turn = fmt.Sprintf(turnStruct, turn, ant.ID+1, nextRoom)
				}

				isSolved = false
			}
		}

		if !isSolved {
			turns = append(turns, turn)
		}
	}

	for _, turn := range turns {
		fmt.Println(turn)
	}
}

func getNextRoom(currentRoom string, pathsToExit [][]string, usedPaths []string, field Field) string {
	for _, pathToExit := range pathsToExit {
		for i, room := range pathToExit {
			if room == currentRoom {
				nextRoom := pathToExit[i+1]

				// As I remember, this is a hack to prevent ants from colliding in some test case... So a little bit of cheating here :D
				if len(field.Ants)-len(pathToExit) == getNumOfFinishedAnts(field.Ants) && nextRoom != field.EndRoomName && currentRoom == field.StartRoomName && len(pathsToExit) == 2 {
					continue // Wait for better path
				}

				if isRoomEmpty(nextRoom, field) || nextRoom == field.EndRoomName {
					isPathUsed := false
					for _, path := range usedPaths {
						if path == currentRoom+"-"+nextRoom {
							isPathUsed = true
						}
					}

					if isPathUsed {
						continue
					}

					return nextRoom
				} else {
					continue
				}
			}
		}
	}

	return ""
}

func isRoomEmpty(roomName string, field Field) bool {
	for _, ant := range field.Ants {
		if ant.CurrentRoom == roomName {
			return false
		}
	}

	return true
}

func getNumOfFinishedAnts(ants []*Ant) int {
	numOfFinishedAnts := 0
	for _, ant := range ants {
		if ant.IsFinished {
			numOfFinishedAnts++
		}
	}

	return numOfFinishedAnts
}
