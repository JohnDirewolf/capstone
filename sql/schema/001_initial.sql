-- +goose Up

CREATE TABLE rooms (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL DEFAULT 'Maze Runner!',
    description TEXT NOT NULL DEFAULT 'No Description Provided',
    discovered BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL DEFAULT 'Item Name',
    article TEXT NOT NULL DEFAULT 'an ',
    description TEXT NOT NULL DEFAULT 'A Mysterious Object',
    type TEXT NOT NULL DEFAULT 'item',
    cur_location INTEGER,
    FOREIGN KEY (cur_location) REFERENCES rooms(id)
);

CREATE TABLE doors (
    id INTEGER PRIMARY KEY,
    room_id INTEGER NOT NULL, 
    direction TEXT CHECK (direction IN ('north', 'south', 'west', 'east')),
    locked BOOLEAN NOT NULL DEFAULT false,
    key_id INTEGER NULL,  
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (key_id) REFERENCES items(id)
);

CREATE TABLE creatures (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL DEFAULT 'Tim',
    type TEXT NOT NULL DEFAULT 'Enchanter',
    description TEXT NOT NULL DEFAULT 'Ni',
    is_alive BOOLEAN NOT NULL DEFAULT true,
    vanquished_by INTEGER NULL,
    cur_location INTEGER,
    FOREIGN KEY (vanquished_by) REFERENCES items(id),
    FOREIGN KEY (cur_location) REFERENCES rooms(id)
);

-- +goose Down
DROP TABLE creatures;
DROP TABLE doors;
DROP TABLE items;
DROP TABLE rooms;