package controller

import (
	"database/sql"
	"encoding/json"
	"main/model"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users (username, email, password, profile_picture_url, bio, followers_count, followings_count) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.ProfilePictureURL, user.Bio, user.FollowersCount, user.FollowingsCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()
	user.UserID = int(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	UserID, err := strconv.Atoi(id)
	if err != nil || UserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user model.User
	stmt, err := db.Prepare("SELECT user_id, username, email, password, profile_picture_url, bio, followers_count, followings_count FROM users WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(UserID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ProfilePictureURL,
		&user.Bio,
		&user.FollowersCount,
		&user.FollowingsCount,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	UserID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("UPDATE users SET username = ?, email = ? WHERE user_id = ?", user.Username, user.Email, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("User Updated Successfully"))
}
