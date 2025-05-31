package lib

import "fmt"

func linkRooms(firstRoom string, secondRoom string, field *Field) error {
	// Check if the rooms exist
	var firstRoomObj interface{}
	var secondRoomObj interface{}

	for _, room := range field.Rooms {
		if room.Name == firstRoom {
			firstRoomObj = room
		}

		if room.Name == secondRoom {
			secondRoomObj = room
		}
	}

	if firstRoomObj == nil || secondRoomObj == nil {
		return fmt.Errorf("Room %s or %s does not exist", firstRoom, secondRoom)
	}

	// Check if the rooms are already linked
	for _, link := range firstRoomObj.(*Room).ConnectedRooms {
		if link == secondRoom {
			return fmt.Errorf("rooms %s and %s are already linked", firstRoom, secondRoom)
		}
	}

	for _, link := range secondRoomObj.(*Room).ConnectedRooms {
		if link == firstRoom {
			return fmt.Errorf("rooms %s and %s are already linked", firstRoom, secondRoom)
		}
	}

	// Link the rooms
	firstRoomObj.(*Room).ConnectedRooms = append(firstRoomObj.(*Room).ConnectedRooms, secondRoom)
	secondRoomObj.(*Room).ConnectedRooms = append(secondRoomObj.(*Room).ConnectedRooms, firstRoom)

	return nil
}
