package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
)

func generateUser() (string, time.Time, string){
    name := fake.Name()
    birthday := fake.DateRange(time.Unix(0, 0), time.Unix(1118793600, 0))
    avatar := fake.ImageURL(250, 250)
    return name, birthday, avatar
}

func generatePost() (string, string, time.Time){
    words := rand.Intn(5)
    sentences := rand.Intn(10)
    sentenceWords := rand.Intn(10)
    paragraphs := rand.Intn(5)
    title := fake.Sentence(words + 1)
    content := fake.Paragraph(paragraphs + 1, sentences + 5, sentenceWords + 5, "\n")
    postDate := fake.DateRange(time.Unix(1672531200, 0), time.Now())
    return title, content, postDate
}

func Seed(db *sql.DB) {
    schemaFile, err := os.ReadFile("sql/schema.sql")
    if err != nil {
        log.Fatalf("error reading schema file: %e", err)
    }
    schema := string(schemaFile)
    log.Println("setting up DB schema")
    _, err = db.Exec(schema)
    if err != nil {
       panic(err) 
    }
    log.Println("seeding users")
    userStmt := "INSERT INTO users (name, birthday, avatar) values ($1, $2, $3);"
    err = insert(db, userStmt, 50, func(stmt *sql.Stmt) error {
        name, birthday, avatar := generateUser()
        _, err := stmt.Exec(name, birthday.Unix(), avatar)
        return err
    })
    if err != nil {
        log.Fatalf("error seeding users: %e", err)
    }
    rows, err := db.Query("SELECT id FROM users")
    userIds := make([]int, 0)
    if err != nil {
        log.Fatalf("error fetching seeded user ids: %e", err)
    }
    defer rows.Close()
    for rows.Next() {
       var id int 
       err = rows.Scan(&id)
       if err != nil {
           log.Fatalf("error scanning user id: %e", err)
       }
       userIds = append(userIds, id)
    }
    log.Println("seeding posts")
    postStmt := "INSERT INTO posts (title, content, post_date, author_id) values ($1, $2, $3, $4);"
    err = insert(db, postStmt, 200, func(stmt *sql.Stmt) error {
        userId := userIds[rand.Intn(len(userIds))]
        title, content, postDate := generatePost()
        _, err := stmt.Exec(title, content, postDate.Unix(), userId)
        return err
    });
    if err != nil {
        log.Fatalf("error seeding posts: %e", err)
    }
    rows, err = db.Query("SELECT id, post_date FROM posts")
    postDates := make(map[int]int64)
    postIds := make([]int, 0)
    if err != nil {
        log.Fatalf("error fetching seeded post ids: %e", err)
    }
    defer rows.Close()
    for rows.Next() {
       var id int 
       var postDate int64
       err = rows.Scan(&id, &postDate)
       if err != nil {
           log.Fatalf("error scanning post id: %e", err)
       }
       postIds = append(postIds, id)
       postDates[id] = postDate
    }
    log.Println("seeding likes")
    likeStmt := "INSERT INTO likes (like_date, user_id, post_id) values ($1, $2, $3)"
    likeMap := make(map[string]bool)
    getIds := func() (int, int) {
        userId := userIds[rand.Intn(len(userIds))]
        postId := postIds[rand.Intn(len(postIds))]
        return userId, postId
    }
    err = insert(db, likeStmt, 1000, func(stmt *sql.Stmt) error {
        isValid := false
        var userId, postId int
        var key string
        for !isValid {
            userId, postId = getIds()
            key = fmt.Sprintf("%d-%d", userId, postId)
            isValid = !likeMap[key]
        }
        _, err := stmt.Exec(fake.DateRange(time.Unix(postDates[postId], 0), time.Now()).Unix(), userId, postId)
        likeMap[key] = true
        return err
    });
    if err != nil {
        log.Fatalf("error seeding likes: %e", err)
    }
}

func insert(db *sql.DB, str string, num int, cb func(stmt *sql.Stmt) error) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    stmt, err := tx.Prepare(str)
    if err != nil {
        return err
    }
    defer stmt.Close()
    for i := 1; i <= num; i++ {
        err := cb(stmt)
        if err != nil {
            return err
        }
    }
    err = tx.Commit()
    if err != nil {
        return err
    }
    return nil
}
