package internal

import (
	"database/sql"
	"errors"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Create(f *Snippet) (int, error) {
	id, err := m.InsertPost(f)
	if err != nil {
		return 0, err
	}
	err = m.InsertCategory(id, f.CatID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) GetFilterPost(catID []int, UserID int) ([]*Snippet, error) {
	snippets, err := m.Filter(catID)
	if err != nil {
		return nil, err
	}
	for _, post := range snippets {
		post.Catigories, err = m.CatigoriesById(post.ID)
		if err != nil {
			return nil, err
		}
	}
	return snippets, err
}

func (m *SnippetModel) Filter(catID []int) ([]*Snippet, error) {
	newpost := []*Snippet{}
	query := `SELECT snippets.id, snippets.author_name, snippets.title, snippets.content, snippets.created FROM snippets 
	JOIN post_cat ON snippets.id=post_cat.post_id
	WHERE post_cat.cat_id=?;`
	for _, i := range catID {
		rows, err := m.DB.Query(query, i)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			p := &Snippet{}
			err := rows.Scan(&p.ID, &p.UserName, &p.Title, &p.Content, &p.Created)
			if err != nil {
				return nil, err
			}
			if newpost != nil && IsUnique(p.ID, newpost) {
				newpost = append(newpost, p)
			}
		}
	}
	return newpost, nil
}

func IsUnique(postID int, posts []*Snippet) bool {
	for _, post := range posts {
		if postID == post.ID {
			return false
		}
	}
	return true
}

func (m *SnippetModel) InsertPost(f *Snippet) (int, error) {
	stmt := `INSERT INTO snippets (author_name, title, content, created) VALUES (?, ?, ?, datetime('now', '+6 hours'))`

	result, err := m.DB.Exec(stmt, f.UserName, f.Title, f.Content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) GetById(post_id int) (*Snippet, error) {
	snippet, err := m.PostById(post_id)
	if err != nil {
		return nil, err
	}
	snippet.Catigories, err = m.CatigoriesById(post_id)
	if err != nil {
		return nil, err
	}
	snippet, err = m.ReactionsByPostId(post_id, snippet)
	if err != nil {
		return nil, err
	}
	return snippet, nil
}

func (m *SnippetModel) ReactionsByPostId(id int, s *Snippet) (*Snippet, error) {
	stmt := `SELECT reaction FROM LikeReact WHERE post_id = ?`
	rows, err := m.DB.Query(stmt, id)
	var a int
	if err != nil {
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&a)
		if err != nil {
			return nil, err
		}
		if a == 1 {
			s.Likes++
		} else if a == -1 {
			s.Dislikes++
		}
	}
	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	stmt := `SELECT id, author_name, title, content, created FROM snippets ORDER BY created`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var snippets []*Snippet
	defer rows.Close()
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.UserName, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		s.Catigories, _ = m.CatigoriesById(s.ID)
		// Добавляем структуру в срез.
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) ReturMyPosts(username string) ([]*Snippet, error) {
	stmt := `SELECT id, author_name, title, content, created FROM snippets WHERE author_name = ?`

	rows, err := m.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}
	var snippets []*Snippet
	defer rows.Close()
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.UserName, &s.Title, &s.Content, &s.Created)
		if err != nil {
			return nil, err
		}
		s.Catigories, _ = m.CatigoriesById(s.ID)

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) GetLikedPost(userID int) ([]*Snippet, error) {
	posts, err := m.GetAllLikedPost(userID)
	if err != nil {
		return nil, err
	}
	var snippets []*Snippet
	for _, post := range posts {
		s, err := m.PostById(post)
		if err != nil {
			return nil, err
		}
		s.Catigories, _ = m.CatigoriesById(s.ID)
		snippets = append(snippets, s)
	}
	return snippets, nil
}

func (m *SnippetModel) GetAllCatigories(data *templateData) ([]string, error) {
	categories := []string{}
	query := `SELECT category FROM categories`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c *string
		err := rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *c)
	}
	data.Catigories = categories
	return categories, nil
}

func (m *SnippetModel) InsertCategory(postID int, catID []int) error {
	query := `INSERT INTO post_cat (post_id, cat_id) VALUES(?, ?)`
	for _, i := range catID {
		_, err := m.DB.Exec(query, postID, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *SnippetModel) PostById(post_id int) (*Snippet, error) {
	query := `SELECT id, author_name, title, content, created FROM snippets WHERE id = ?`
	s := &Snippet{}
	err := m.DB.QueryRow(query, post_id).Scan(&s.ID, &s.UserName, &s.Title, &s.Content, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) CatigoriesById(post_id int) ([]string, error) {
	query := `SELECT cat_id, (
	SELECT category FROM categories WHERE categories.id = post_cat.cat_id
	)
FROM post_cat WHERE post_id=?`
	var catigories []string
	rows, err := m.DB.Query(query, post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var catid int64
		var cat string
		err = rows.Scan(&catid, &cat)
		catigories = append(catigories, cat)
		if err != nil {
			return nil, err
		}

	}
	return catigories, nil
}

func (m *SnippetModel) GetComments(post_id int) ([]*SnippetCom, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT * FROM comments WHERE post_id = ?`

	row, err := m.DB.Query(stmt, post_id)
	if err != nil {
		return nil, err
	}

	comments := []*SnippetCom{}

	for row.Next() {
		s := &SnippetCom{}
		err := row.Scan(&s.ID, &s.PostID, &s.UserId, &s.UserName, &s.Content, &s.Created)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
		s.LikesCom, _ = m.LikeByComID(s.ID)
		s.DislikesCom, _ = m.DislikeByComID(s.ID)
		comments = append(comments, s)
	}
	// Если все хорошо, возвращается объект Snippet.
	return comments, nil
}

func (m *SnippetModel) InsertComments(c *FormComment) error {
	stmt := `INSERT INTO comments (post_id, user_id, author_name, content, created) VALUES (?, ?, ?, ?, datetime('now', '+6 hours'))`

	_, err := m.DB.Exec(stmt, &c.PostID, &c.UserId, &c.UserName, &c.Content)
	if err != nil {
		return err
	}

	return nil
}

func (m *SnippetModel) GetAllLikedPost(userID int) ([]int, error) {
	stmt := `SELECT post_id, reaction FROM LikeReact WHERE user_id = ?`
	var posts []int
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post int
		var reaction int
		err := rows.Scan(&post, &reaction)
		if err != nil {
			return nil, err
		}
		if reaction == 1 {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (m *SnippetModel) LikeByComID(comment_id int) (int, error) {
	stmt := `SELECT type FROM ComentReact WHERE comment_id = ?`
	rows, err := m.DB.Query(stmt, comment_id)
	var a int
	var reaction int
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&a)
		if err != nil {
			return 0, err
		}
		if a == 1 {
			reaction++
		}
	}
	return reaction, nil
}

func (m *SnippetModel) DislikeByComID(comment_id int) (int, error) {
	stmt := `SELECT type FROM ComentReact WHERE comment_id = ?`
	rows, err := m.DB.Query(stmt, comment_id)
	var a int
	var reaction int
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&a)
		if err != nil {
			return 0, err
		}
		if a == -1 {
			reaction++
		}
	}
	return reaction, nil
}
