package types

import "database/sql"

type UrlAction string

// These are all the legit actions that can be sent in the URL.
const (
	Start UrlAction = "start"
	End   UrlAction = "end"
	North UrlAction = "north"
	South UrlAction = "south"
	West  UrlAction = "west"
	East  UrlAction = "east"
)

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

type ItemData struct {
	Id          int
	Name        string
	Description string
	CurLocation int
}
