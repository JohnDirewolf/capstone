package types

import (
	"database/sql"
	"html/template"
)

type UrlAction string

// These are all the legit actions that can be sent in the URL.

const (
	Start  UrlAction = "start"
	End    UrlAction = "end"
	Get    UrlAction = "get"    //items, key or sword
	Use    UrlAction = "use"    //key
	Attack UrlAction = "attack" //goblin
	North  UrlAction = "north"
	South  UrlAction = "south"
	West   UrlAction = "west"
	East   UrlAction = "east"
)

type DoorData struct {
	Id        int
	RoomId    int
	Direction string
	Locked    bool
	Guarded   bool
	KeyId     sql.NullInt16
}

type RoomData struct {
	Id          int
	Title       string
	Description string
	Discovered  bool
	Doors       map[string]DoorData
}

type ItemData struct {
	Id          int
	Name        string
	Article     string
	Description string
	Type        string
	CurLocation int
}

type CreatureData struct {
	Id           int
	Name         string
	Type         string
	Description  string
	IsAlive      bool
	VanquishedBy sql.NullInt16 //Item ID or Null
	CurLocation  int           //roomID
	Guards       sql.NullInt16 //Door ID or Null
}

type PageData struct {
	Title        string
	Rooms        template.HTML
	Compass      template.HTML
	Description  template.HTML
	Inventory    template.HTML
	Action       template.HTML
	Instructions template.HTML
}

type SpecialStatus struct {
	IsStart    bool
	IsLocked   bool
	Unlocked   bool
	IsGuarded  bool
	Vanquished bool
}
