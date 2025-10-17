package controller

import (
	"database/sql"
	"encoding/json"
	"main/model"
	"net/http"
	"strconv"
)

var _ = (*sql.DB)(nil)

func CreateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var v model.Video
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare(`INSERT INTO videos (user_id, video_caption, duration, video_url, thumbnail_url, is_public) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(v.VideoCaption, v.Duration, v.VideoURL, v.ThumbnailURL, v.IsPublic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	v.VideoID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func GetVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var v model.Video
	err = db.QueryRow(`SELECT video_id, video_caption, upload_date, duration, video_url, thumbnail_url, likes_count, comments_count, views_count, is_public FROM videos WHERE video_id = ?`, id).
		Scan(&v.VideoID, &v.VideoCaption, &v.UploadDate, &v.Duration, &v.VideoURL, &v.ThumbnailURL, &v.LikesCount, &v.CommentsCount, &v.ViewsCount, &v.IsPublic)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
