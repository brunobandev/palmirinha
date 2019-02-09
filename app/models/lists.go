package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// List of items
type List struct {
	ID          uint
	Description string
	channelID   string
	Items       []Item
}

// OpenList ...
func OpenList(description, channelID string) string {
	query := `
		INSERT cart 
		SET description = ?, status = ?, channel_id = ?
	`
	stmt, err := Connection.Mysql.Prepare(query)
	if err != nil {
		fmt.Println("Model OpenList [prepare]: ", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(description, 1, channelID)
	if err != nil {
		fmt.Println("Model OpenList [exec]: ", err)
	}

	id, _ := res.LastInsertId()
	idToString := strconv.FormatInt(int64(id), 10)

	return idToString
}

// CloseList ...
func CloseList(channelID string) bool {
	query := `
		UPDATE cart 
		SET status = ? 
		WHERE status = ? 
		AND channel_id = ?
	`
	stmt, err := Connection.Mysql.Prepare(query)
	if err != nil {
		fmt.Println("Model CloseList [prepare]: ", err)
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(0, 1, channelID)
	if err != nil {
		fmt.Println("Model CloseList [exec]: ", err)
		return false
	}

	return true
}

// CountOpenList ...
func CountOpenList(channelID string) uint {
	var count uint
	query := `
		SELECT COUNT(*) 
		FROM cart 
		WHERE status = 1 
		AND channel_id = ?
	`
	rows, err := Connection.Mysql.Query(query, channelID)
	if err != nil {
		fmt.Println(err)
	}

	if rows.Next() {
		rows.Scan(&count)
	}

	return count
}

// GetOpenListByChannelID ...
func GetOpenListByChannelID(channelID string) (list List, userErr error) {
	query := `
		SELECT id, description, channel_id 
		FROM cart 
		WHERE status = 1 
		AND channel_id = ?
	`
	row := Connection.Mysql.QueryRow(query, channelID)
	err := row.Scan(&list.ID, &list.Description, &list.channelID)
	if err == sql.ErrNoRows {
		userErr = errors.New("Não há um carrinho aberto, crie com o comando `!criar [nome]`")
	} else if err != nil {
		userErr = errors.New("Erro ao listar pedidos")
		fmt.Println("Model GetOpenListByChannelID [scan]: ", err, sql.ErrNoRows)
	}

	return list, userErr
}

// RaffleList ...
func RaffleList(channelID string) string {
	var Chosen string
	query := `
		SELECT item.discord_user_id 
		FROM cart
		JOIN item ON item.cart_id = cart.id 
		WHERE cart.status = 1 
		AND cart.channel_id = ? 
		ORDER BY RAND() 
		LIMIT 1
	`
	row := Connection.Mysql.QueryRow(query, channelID)
	err := row.Scan(&Chosen)
	if err != nil {
		fmt.Println("Model RaffleList [scan]: ", err)
	}

	return Chosen
}
