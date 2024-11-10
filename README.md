# Boot.dev Capstone Project

## Description
This is an example of Go programming a server and database driven program as a very simple fantasy web based game. 
In the game you explore the "Hatrock Dungeon" finding a magic key and a magic sword to defeat a goblin and escape the dungeon's maze. 
This program creates a local server using ServeMux as a file server and handler for game actions.
A PostgreSQL database is then used to store game data.
All art assets were created by myself.
This is my Capstone Project for the Boot.dev back-end developers course.

## Installation & Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/JohnDirewolf/capstone
   cd capstone
   ```

2. Install PostgreSQL if needed (see below.)

3. Install dependencies:
        ```bash
        go mod download
        ```

4. Create a .env file based on .env.example and configure your environment variables. (See below.)

5. Running the Application
   To run the application:
        ```bash
        go run .
        ```

## Required:
* Postgres (PostgreSQL) 14.13

### Install PostgreSQL 14.13 on Ubuntu:
To install PostgreSQL on Ubuntu, run the following commands in your terminal:
    ```bash
    sudo apt update
    sudo apt install postgresql
    ```

### Set up your PostgreSQL database connection string:
1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Open the .env file in your text editor:
    ```bash
    nano .env
    ```

3. Update the DB_URL variable with your PostgreSQL connection string:

postgres://<username>:<password>@<host>:<port>/<database_name>?sslmode=disable

Example format explanation:
- username: your PostgreSQL username
- password: your PostgreSQL password
- host: localhost (if running locally) or your database host
- port: 5432 (default PostgreSQL port)
- database_name: name of your database (game_database)

4. Save the file (in nano: Ctrl + O, then Ctrl + X to exit)

