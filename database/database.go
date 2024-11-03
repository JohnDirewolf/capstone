package database

import (
	"database/sql"
	"log"

	"github.com/JohnDirewolf/capstone/shared/types"

	_ "github.com/lib/pq"
)

const dbURL string = "postgres://postgres:postgres@localhost:5432/game_database?sslmode=disable"

var heart *sql.DB

func Init() error {
	var err error
	heart, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Database, Init: Error in connecting to database: %v", err)
		return err
	}

	// Verify we can connect
	err = heart.Ping()
	if err != nil {
		heart.Close()
		log.Printf("Database, Init: Failed to ping database: %v", err)
		return err
	}
	return nil
}

func Clear() error {
	_, err := heart.Exec("DELETE FROM rooms;")
	if err != nil {
		log.Printf("Database, Clear: Error inserting values for Room: %v", err)
	}
	return err
}

func GetRoom(roomID int) (types.RoomData, error) {
	var room types.RoomData
	// Temporary variables for scanning
	var northExists, northLocked, southExists, southLocked, westExists, westLocked, eastExists, eastLocked bool
	var northKeyID, southKeyID, westKeyID, eastKeyID sql.NullInt16

	query := `SELECT id, title, description, discovered,
			doorNorth, doorNorthLocked, doorNorthKey_id,
			doorSouth, doorSouthLocked, doorSouthKey_id,
			doorWest, doorWestLocked, doorWestKey_id,
			doorEast, doorEastLocked, doorEastKey_id
			FROM rooms WHERE id=$1`
	roomRecord := heart.QueryRow(query, roomID)

	// Use the retrieved variables to build the nested DoorData structures
	err := roomRecord.Scan(&room.Id, &room.Title, &room.Description, &room.Discovered,
		&northExists, &northLocked, &northKeyID,
		&southExists, &southLocked, &southKeyID,
		&westExists, &westLocked, &westKeyID,
		&eastExists, &eastLocked, &eastKeyID)
	if err != nil {
		log.Printf("Database, GetRoom: Error getting data for room: %v", err)
		return types.RoomData{}, err
	}

	// Construct the DoorData objects and assign to room
	room.North = types.DoorData{Exists: northExists, Locked: northLocked, KeyID: northKeyID}
	room.South = types.DoorData{Exists: southExists, Locked: southLocked, KeyID: southKeyID}
	room.West = types.DoorData{Exists: westExists, Locked: westLocked, KeyID: westKeyID}
	room.East = types.DoorData{Exists: eastExists, Locked: eastLocked, KeyID: eastKeyID}

	return room, nil
}

func GetDiscoveredRooms() ([]int, error) {
	var rooms []int
	var roomID int
	rows, err := heart.Query("SELECT id FROM rooms WHERE discovered = TRUE;")
	if err != nil {
		log.Printf("Database, GetDiscoveredRooms: Error getting list of rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&roomID)
		if err != nil {
			log.Printf("Database, GetDiscoveredRooms: Error building list of room IDs: %v", err)
			return nil, err
		}
		rooms = append(rooms, roomID)
	}

	return rooms, nil
}

func InsertRoom(roomInfo types.RoomData) error {
	_, err := heart.Exec(`
		INSERT INTO rooms (
			id, 
			title,
			description,
			discovered,
			doorNorth,
			doorNorthLocked,
			doorNorthKey_id,  
			doorSouth,
			doorSouthLocked,
			doorSouthKey_id, 
			doorWest,
			doorWestLocked,
			doorWestKey_id, 
			doorEast,
			doorEastLocked,
			doorEastKey_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
		roomInfo.Id,
		roomInfo.Title,
		roomInfo.Description,
		roomInfo.Discovered,
		roomInfo.North.Exists,
		roomInfo.North.Locked,
		roomInfo.North.KeyID,
		roomInfo.South.Exists,
		roomInfo.South.Locked,
		roomInfo.South.KeyID,
		roomInfo.West.Exists,
		roomInfo.West.Locked,
		roomInfo.West.KeyID,
		roomInfo.East.Exists,
		roomInfo.East.Locked,
		roomInfo.East.KeyID)
	if err != nil {
		log.Printf("Database, InsertRoom: Error inserting values for Room: %v", err)
	}
	return err
}

func DiscoverRoom(roomID int) error {
	_, err := heart.Exec("UPDATE rooms SET discovered=TRUE WHERE id=$1;", roomID)
	if err != nil {
		log.Printf("Database, DiscoverRoom: Error updating discovered for room: %v", err)
	}
	return err
}

func Close() error {
	err := heart.Close()
	if err != nil {
		log.Printf("Database, Close: Error closing database: %v", err)
	}
	return err
}
