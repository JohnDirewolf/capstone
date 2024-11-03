package types

import "database/sql"

// Direction constants
const North string = "north"
const South string = "south"
const West string = "west"
const East string = "east"

type DoorData struct {
	Exists bool
	Locked bool
	KeyID  sql.NullInt16
}

type RoomData struct {
	Id          int
	Title       string
	Description string
	Discovered  bool
	North       DoorData
	South       DoorData
	West        DoorData
	East        DoorData
}
