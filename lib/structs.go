package lib

type Ant struct {
	ID          int
	CurrentRoom string
	IsFinished  bool
}

type Room struct {
	Name           string
	ConnectedRooms []string
}

type Field struct {
	Ants          []*Ant
	Rooms         []*Room
	StartRoomName string
	EndRoomName   string
}
