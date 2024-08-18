package pg

import (
	"fmt"

	"github.com/SurkovIlya/chat-app/internal/models"
	"github.com/SurkovIlya/chat-app/pkg/postgres"
)

type PostgresStorage struct {
	storage *postgres.Database
}

func New(storage *postgres.Database) *PostgresStorage {
	return &PostgresStorage{
		storage: storage,
	}
}

func (ps *PostgresStorage) SaveUser(userName string) error {
	query := `INSERT INTO users (user_name) VALUES($1)`

	_, err := ps.storage.Conn.Exec(query, userName)
	if err != nil {
		return fmt.Errorf("error insert messages: %s", err)
	}

	return nil
}

func (ps *PostgresStorage) UserExist(userName string) (bool, error) {
	var exist bool
	query := `SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE user_name = $1);`
	row := ps.storage.Conn.QueryRow(query, userName)

	err := row.Scan(&exist)
	if err != nil {
		return exist, fmt.Errorf("error Scan:%s", err)
	}

	return exist, nil
}

func (ps *PostgresStorage) SaveRoom(roomName, userName string) error {
	query := `INSERT INTO rooms (room_name) VALUES($1)`

	_, err := ps.storage.Conn.Exec(query, roomName)
	if err != nil {
		return fmt.Errorf("error insert messages: %s", err)
	}

	err = ps.SaveMembersChat(roomName, userName)
	if err != nil {
		return fmt.Errorf("error SaveMembersChat: %s", err)
	}

	return nil
}

func (ps *PostgresStorage) getRoomID(roomName string) (int, error) {
	var roomID int

	queryRoom := `SELECT id FROM rooms WHERE room_name=$1`
	rowR := ps.storage.Conn.QueryRow(queryRoom, roomName)

	err := rowR.Scan(&roomID)
	if err != nil {
		return roomID, fmt.Errorf("error Scan user ID: %s", err)
	}

	return roomID, nil
}

func (ps *PostgresStorage) getUserID(userName string) (int, error) {
	var userID int

	queryUser := `SELECT id FROM users WHERE user_name=$1`
	rowU := ps.storage.Conn.QueryRow(queryUser, userName)

	err := rowU.Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("error Scan user ID: %s", err)
	}

	return userID, nil
}

func (ps *PostgresStorage) SaveMembersChat(roomName, userName string) error {
	userID, err := ps.getUserID(userName)

	if err != nil {
		return fmt.Errorf("error getsID: %s", err)
	}

	roomID, err := ps.getRoomID(roomName)
	if err != nil {
		return fmt.Errorf("error getsID: %s", err)
	}

	queryMembers := `INSERT INTO chat_members (room_id, user_id) VALUES ($1, $2)`
	_, err = ps.storage.Conn.Exec(queryMembers, roomID, userID)
	if err != nil {
		return fmt.Errorf("error Exec chat_members: %s", err)
	}

	return nil
}

func (ps *PostgresStorage) SaveMsg(roomName, userName, msg string) error {
	userID, err := ps.getUserID(userName)
	if err != nil {
		return fmt.Errorf("error getsUserID: %s", err)
	}

	roomID, err := ps.getRoomID(roomName)
	if err != nil {
		return fmt.Errorf("error getsRoomID: %s", err)
	}

	query := `INSERT INTO messages (room_id, user_id, content) VALUES ($1, $2, $3)`
	_, err = ps.storage.Conn.Exec(query, roomID, userID, msg)
	if err != nil {
		return fmt.Errorf("error Exec chat_members: %s", err)
	}

	return nil
}

func (ps *PostgresStorage) GetMsgs(roomName string) ([]models.RoomMsg, error) {
	oldMsgs := make([]models.RoomMsg, 0)

	roomID, err := ps.getRoomID(roomName)
	if err != nil {
		return nil, fmt.Errorf("error getsID: %s", err)
	}

	query := `SELECT u.user_name, m.content 
				FROM messages AS m 
				LEFT JOIN users AS u ON u.id = m.user_id
				WHERE m.room_id = $1`

	rows, err := ps.storage.Conn.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("error Query: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var msg models.RoomMsg

		err := rows.Scan(&msg.UserName, &msg.Content)
		if err != nil {
			return nil, fmt.Errorf("error Scan: %s", err)
		}

		oldMsgs = append(oldMsgs, msg)
	}

	return oldMsgs, nil
}

func (ps *PostgresStorage) GetAllRooms() ([]string, error) {
	rooms := make([]string, 0)

	query := `SELECT room_name FROM rooms`

	rows, err := ps.storage.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error Query: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var room string

		err := rows.Scan(&room)
		if err != nil {
			return nil, fmt.Errorf("error Scan: %s", err)
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}
