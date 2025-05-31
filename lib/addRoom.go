package lib

import (
	"fmt"
	"strings"
)

func addRoom(roomName string, isStart bool, isEnd bool, field *Field) error {
	// Check if roomName is valid
	if strings.Contains(roomName, " ") {
		return fmt.Errorf("room name cannot contain spaces")
	} else if roomName[0] == 'L' {
		return fmt.Errorf("room name cannot start with 'L'")
	}

	// Check if the room already exists
	for _, room := range field.Rooms {
		if room.Name == roomName {
			return fmt.Errorf("Room %s already exists", roomName)
		}
	}

	// Create the room
	room := Room{
		Name:           roomName,
		ConnectedRooms: []string{},
	}

	// Add the room to the field
	field.Rooms = append(field.Rooms, &room)

	// Set the start and end rooms
	if isStart {
		// Check if the start room already exists
		if field.StartRoomName != "" {
			return fmt.Errorf("start room already exists")
		}

		field.StartRoomName = roomName
		for _, ant := range field.Ants {
			ant.CurrentRoom = roomName
			// fmt.Println(ant)
		}

	} else if isEnd {
		// Check if the end room already exists
		if field.EndRoomName != "" {
			return fmt.Errorf("end room already exists")
		}

		field.EndRoomName = roomName
	}

	return nil
}
