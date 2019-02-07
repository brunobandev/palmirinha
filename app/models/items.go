package models

import "fmt"

// Item of the list
type Item struct {
	ID            uint
	CartID        uint
	DiscordUserID string
	Description   string
}

// AddItem ...
func AddItem(item *Item) bool {
	query := `INSERT item SET description = ?, cart_id = ?, discord_user_id = ?`
	stmt, err := Connection.Mysql.Prepare(query)
	if err != nil {
		fmt.Println("Model AddItem [prepare]: ", err)
		return false
	}

	fmt.Println(item)

	_, err = stmt.Exec(item.Description, item.CartID, item.DiscordUserID)
	if err != nil {
		fmt.Println("Model AddItem [exec]: ", err)
		return false
	}

	return true
}

// RemoveItem ...
func RemoveItem(item *Item) bool {
	query := `DELETE FROM item WHERE id = ?`
	stmt, err := Connection.Mysql.Prepare(query)
	if err != nil {
		fmt.Println("Model RemoveItem [prepare]: ", err)
		return false
	}

	_, err = stmt.Exec(&item.ID)
	if err != nil {
		fmt.Println("Model RemoveItem [exec]: ", err)
		return false
	}

	return true
}

// GetItem ...
func GetItem(item *Item, channelID string) Item {
	var getItem Item
	query := `SELECT 
				item.id, item.description, item.discord_user_id 
			FROM 
				item
			INNER JOIN cart ON item.cart_id = cart.id
			WHERE cart.status = TRUE
			AND cart.channel_id = ?
			AND item.discord_user_id = ?`

	row := Connection.Mysql.QueryRow(query, item.DiscordUserID, channelID)
	err := row.Scan(&getItem.ID, &getItem.Description, &getItem.DiscordUserID)
	if err != nil {
		fmt.Println("Model GetItem [scan]: ", err)
	}

	return getItem
}
