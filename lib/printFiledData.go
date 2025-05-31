package lib

import "fmt"

// Note: this function is only used for debugging purposes
func PrintFieldData(field *Field) {
	fmt.Println("Ants:")
	for _, ant := range field.Ants {
		fmt.Println(ant)
	}
	fmt.Println("Rooms:")
	for _, room := range field.Rooms {
		fmt.Println(room)
	}
	fmt.Println("Start room:", field.StartRoomName)
	fmt.Println("End room:", field.EndRoomName)
}
