package api

import (
	"database/sql"
	"encoding/json"
	"github.com/cozy-software/interview-test/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
)

type PostTitle struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	PostDate int    `json:"post_date"`
}

type User struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Birthday int         `json:"birthday"`
	Avatar   string      `json:"avatar"`
	Posts    []PostTitle `json:"posts"`
}

type Post struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	PostDate       int    `json:"post_date"`
	AuthorID       int    `json:"author_id"`
	AuthorName     string `json:"author_name"`
	AuthorBirthday int    `json:"author_birthday"`
	AuthorAvatar   string `json:"author_avatar"`
	LikeCount      int    `json:"like_count"`
	LikedByUser    string `json:"liked_by_user"`
}

func Mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			return
		}
	})
	r.Get("/posts", GetPosts)
	r.Get("/posts/{id}", GetPostById)
	r.Get("/posts/{id}/likes", GetPostLikes)
	r.Get("/users/{id}", GetUserById)

	return r
}
func GetPosts(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("limit")
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "20"
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// Calculate offset and limit for pagination
	offset := (pageInt - 1) * pageSizeInt
	limit := pageSizeInt

	// Execute the SQL query to fetch posts
	rows, err := database.DB.Query("SELECT p.id, p.title, p.content, p.post_date, p.author_id, u.name AS author_name, u.birthday AS author_birthday, u.avatar AS author_avatar, COUNT(l.post_id) AS like_count "+
		"FROM posts p "+
		"JOIN users u ON p.author_id = u.id "+
		"LEFT JOIN likes l ON p.id = l.post_id "+
		"GROUP BY p.id, u.id "+
		"ORDER BY p.id "+
		"OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Prepare the response
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PostDate, &post.AuthorID, &post.AuthorName,
			&post.AuthorBirthday, &post.AuthorAvatar, &post.LikeCount)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//post.Author = author
		post.LikedByUser = "N/A"
		if userID != "" {
			userID, _ := strconv.Atoi(userID)
			post.LikedByUser = checkPostLiked(userID, post.ID)
		}

		posts = append(posts, post)
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		return
	}
}

func GetPostById(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := r.URL.Query().Get("user")

	row := database.DB.QueryRow("SELECT p.id, p.title, p.content, p.post_date, p.author_id, u.name AS author_name, u.birthday AS author_birthday, u.avatar AS author_avatar, COUNT(l.post_id) AS like_count "+
		"FROM posts p "+
		"JOIN users u ON p.author_id = u.id "+
		"LEFT JOIN likes l ON p.id = l.post_id "+
		"WHERE p.id = $1 "+
		"GROUP BY p.id, u.id", postID)

	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PostDate, &post.AuthorID, &post.AuthorName,
		&post.AuthorBirthday, &post.AuthorAvatar, &post.LikeCount)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//post.Author = author

	post.LikedByUser = "N/A"
	if userID != "" {
		userID, _ := strconv.Atoi(userID)
		post.LikedByUser = checkPostLiked(userID, post.ID)
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		return
	}
}

func GetPostLikes(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("limit")
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "20"
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	offset := (pageInt - 1) * pageSizeInt
	limit := pageSizeInt

	// Execute the SQL query to fetch the likes for the post
	rows, err := database.DB.Query(
		"SELECT user_id FROM likes WHERE post_id = $1 ORDER BY user_id OFFSET $2 LIMIT $3", postID, offset, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Prepare the response
	var likes []User
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Fetch the user details
		user := getUser(string(rune(userID)))
		if user != nil {
			likes = append(likes, *user)
		}
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(likes)
	if err != nil {
		return
	}
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	// Fetch the user details
	user := getUser(userID)
	if user == nil {
		http.NotFound(w, r)
		return
	}

	user.Posts = getUserPosts(userID)

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func checkPostLiked(userID int, postID int) string {
	rows, err := database.DB.Query("SELECT user_id FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID)
	if err != nil {
		log.Println(err)
		return "N/A"
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	//var likes []User
	for rows.Next() {
		var userIDtemp int
		err := rows.Scan(&userIDtemp)
		log.Println("UserID", userIDtemp)
		log.Println("PostID", postID)

		if err != nil {
			log.Println(err)
			return "N/A"
		}
		return "true"
	}
	return "false"
}

func getUser(userID string) *User {
	row := database.DB.QueryRow("SELECT id, name, birthday, avatar FROM users WHERE id = $1", userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Birthday, &user.Avatar)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &user
}

func getUserPosts(userID string) []PostTitle {
	rows, err := database.DB.Query("SELECT id, title, content, post_date FROM posts WHERE author_id = $1 ORDER BY post_date DESC LIMIT 5", userID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var posts []PostTitle
	for rows.Next() {
		var post PostTitle
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PostDate)
		if err != nil {
			log.Println(err)
			return nil
		}
		posts = append(posts, post)

	}
	return posts
}
