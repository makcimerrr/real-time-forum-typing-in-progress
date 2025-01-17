package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var chat ChatData
	var isFalse bool

	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if chat.SendUser == "" || chat.ToUser == "" || chat.Message == "" {
		jsonResponse := map[string]interface{}{
			"success": false,
			"message": "Missing fields",
		}
		err := json.NewEncoder(w).Encode(jsonResponse)
		if err != nil {
			return
		}
		isFalse = true
	} else if chat.SendUser == chat.ToUser {
		jsonResponse := map[string]interface{}{
			"success": false,
			"message": "Can't send message to yourself",
		}
		err := json.NewEncoder(w).Encode(jsonResponse)
		if err != nil {
			return
		}
		isFalse = true
	}

	if !isFalse {

		insertMessage := "INSERT INTO private_message (sendUser, toUser, message, time) VALUES (?, ?, ?, ?)"

		_, err := Db.Exec(insertMessage, chat.SendUser, chat.ToUser, chat.Message, chat.Time)
		if err != nil {
			return
		}

		var chatData DisplayChat

		chatData.ReceiverUser = chat.ToUser
		chatData.SenderUser = chat.SendUser

		postDataChat := WebsocketMessage{Type: "chat", Data: chatData}
		broadcast <- postDataChat

		jsonResponse := map[string]interface{}{
			"success": true,
			"message": "Added Message",
		}
		err = json.NewEncoder(w).Encode(jsonResponse)
		if err != nil {
			return
		}
	}
}

func getChatMessages(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get sender and receiver
	var requestBody struct {
		Sender   string `json:"sender"`
		Receiver string `json:"receiver"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query database for messages between sender and receiver
	rows, err := Db.Query("SELECT * FROM private_message WHERE (sendUser = ? AND toUser = ?) OR (sendUser = ? AND toUser = ?) ORDER BY time", requestBody.Sender, requestBody.Receiver, requestBody.Receiver, requestBody.Sender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through rows and populate PrivateMessage structs
	var messages []ChatData
	for rows.Next() {
		var message ChatData
		err := rows.Scan(&message.ID, &message.SendUser, &message.ToUser, &message.Message, &message.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Vérifie les erreurs de la boucle rows.Next()
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Création de la réponse JSON
	jsonResponse := map[string]interface{}{
		"success": true,
		"message": messages,
	}

	// Encodage de la réponse JSON
	err = json.NewEncoder(w).Encode(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getConversations(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get sender and receiver
	var requestBody struct {
		Username string `json:"username"`
	}
	var ListUserTalk []string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := Db.Query("SELECT CASE WHEN sendUser = ? THEN toUser ELSE sendUser END AS user_name, MAX(time) AS last_message_time FROM private_message WHERE sendUser = ? OR toUser = ? GROUP BY user_name ORDER BY last_message_time DESC", requestBody.Username, requestBody.Username, requestBody.Username)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var lastMessageTime string
		err := rows.Scan(&username, &lastMessageTime)
		if err != nil {
			fmt.Println("Erreur lors de la lecture des données de la ligne :", err)
			return
		}

		ListUserTalk = append(ListUserTalk, username)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erreur lors du parcours des résultats :", err)
		return
	}

	jsonResponse := map[string]interface{}{
		"success": true,
		"message": ListUserTalk,
	}
	err = json.NewEncoder(w).Encode(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
