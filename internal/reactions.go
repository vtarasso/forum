package internal

import (
	"database/sql"
	"errors"
)

type ReactionsModel struct {
	DB *sql.DB
}

func (r *ReactionsModel) LikePost(user_id, post_id int) error {
	stmt := `SELECT reaction FROM LikeReact WHERE user_id = ? AND post_id = ?`
	var reaction int

	err := r.DB.QueryRow(stmt, user_id, post_id).Scan(&reaction)
	if errors.Is(err, sql.ErrNoRows) {
		_ = r.InsertReactions(user_id, post_id, 1)
		return nil
	} else if err != nil {
		return err
	}
	if reaction == -1 {
		err = r.ReplaceReactions(user_id, post_id, 1)
		if err != nil {
			return err
		}
	}
	if reaction == 1 {
		err = r.DropReaction(user_id, post_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionsModel) DislikePost(user_id, post_id int) error {
	stmt := `SELECT reaction FROM LikeReact WHERE user_id = ? AND post_id = ?`
	var reaction int
	err := r.DB.QueryRow(stmt, user_id, post_id).Scan(&reaction)
	if errors.Is(err, sql.ErrNoRows) {
		_ = r.InsertReactions(user_id, post_id, -1)
		return nil
	} else if err != nil {
		return err
	}

	if reaction == 1 {
		err = r.ReplaceReactions(user_id, post_id, -1)
		if err != nil {
			return err
		}
	}

	if reaction == -1 {

		err = r.DropReaction(user_id, post_id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReactionsModel) InsertReactions(user_id, post_id, reaction int) error {
	stmt := `INSERT INTO LikeReact (post_id, user_id, reaction) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(stmt, post_id, user_id, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionsModel) ReplaceReactions(user_id, post_id, reaction int) error {
	_, err := r.DB.Exec(`UPDATE LikeReact SET reaction = ? WHERE post_id = ? AND user_id = ?`, reaction, post_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionsModel) DropReaction(user_id, post_id int) error {
	stmt := `DELETE FROM LikeReact WHERE LikeReact.user_id = ? AND LikeReact.post_id = ?`
	_, err := r.DB.Exec(stmt, user_id, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionsModel) LikeComment(userID, commentID int) error {
	stmt := `SELECT (type) FROM ComentReact WHERE comment_id = ? AND user_id = ?`
	var reaction int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = r.InsertComReactions(userID, commentID, 1)
			return nil
		} else {
			return err
		}
	}
	if reaction == -1 {
		err = r.ReplaceComReactions(userID, commentID, 1)

		if err != nil {
			return err
		}
	}
	if reaction == 1 {
		err = r.DropComReaction(userID, commentID)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionsModel) DislikeComment(userID, commentID int) error {
	stmt := `SELECT type FROM ComentReact WHERE user_id = ? AND comment_id = ?`

	var reaction int
	err := r.DB.QueryRow(stmt, userID, commentID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = r.InsertComReactions(userID, commentID, -1)
		} else {
			return err
		}
	}

	if reaction == 1 {
		err = r.ReplaceComReactions(userID, commentID, -1)

		if err != nil {
			return err
		}
	}
	if reaction == -1 {
		err = r.DropComReaction(userID, commentID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionsModel) InsertComReactions(user_id, commentID, reaction int) error {
	stmt := `INSERT INTO ComentReact (comment_id, user_id, type) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(stmt, commentID, user_id, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionsModel) ReplaceComReactions(user_id, commentID, reaction int) error {
	_, err := r.DB.Exec(`UPDATE ComentReact SET type = ? WHERE comment_id = ? AND user_id = ?`, reaction, commentID, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReactionsModel) DropComReaction(user_id, commentID int) error {
	stmt := `DELETE FROM ComentReact WHERE ComentReact.user_id = ? AND ComentReact.comment_id = ?`
	_, err := r.DB.Exec(stmt, user_id, commentID)
	if err != nil {
		return err
	}
	return nil
}
