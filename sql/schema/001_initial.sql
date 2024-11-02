-- +goose Up
CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    description TEXT NOT NULL DEFAULT 'Mysterious Object',
    initialLocation INTEGER,
    currentLocation INTEGER
);

CREATE TABLE rooms (
    id INTEGER PRIMARY KEY,
    description TEXT NOT NULL DEFAULT 'No Description Provided',
    doorNorth BOOLEAN NOT NULL DEFAULT false,
    doorNorthLocked BOOLEAN NOT NULL DEFAULT false,
    doorNorthKey INTEGER,  
    doorSouth BOOLEAN NOT NULL DEFAULT false,
    doorSouthLocked BOOLEAN NOT NULL DEFAULT false,
    doorSouthKey INTEGER, 
    doorWest BOOLEAN NOT NULL DEFAULT false,
    doorWestLocked BOOLEAN NOT NULL DEFAULT false,
    doorWestKey INTEGER, 
    doorEast BOOLEAN NOT NULL DEFAULT false,
    doorEastLocked BOOLEAN NOT NULL DEFAULT false,
    doorEastKey INTEGER,
    FOREIGN KEY (doorNorthKey) REFERENCES items(id),
    FOREIGN KEY (doorSouthKey) REFERENCES items(id),
    FOREIGN KEY (doorWestKey) REFERENCES items(id),
    FOREIGN KEY (doorEastKey) REFERENCES items(id)
);

-- +goose Down
DROP TABLE rooms;
DROP TABLE items;