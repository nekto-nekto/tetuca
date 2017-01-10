package db

import "database/sql"

// Stores generated prepared statements
var prepared = make(map[string]*sql.Stmt, len(protoPrepared))

// Queries to be converted into prepared statements
var protoPrepared = map[string]string{
	"isLoggedIn": `
		SELECT EXISTS(
			SELECT true FROM sessions
				WHERE account = $1 AND token = $2
		);`,

	"writePost": `
		INSERT INTO posts (
			editing, spoiler, id, board, op, time, body, name, trip, auth,
			SHA1, imageName, links, backlinks, commands
		) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,

	"writeImage": `
		INSERT INTO images (
			APNG, audio, video, fileType, thumbType, dims, length, size, MD5,
			SHA1
		) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,

	"writeOP": `
		INSERT INTO threads (
			board, log, id, postCtr, imageCtr, replyTime, bumpTime, subject
		) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`,

	"newPostID": `SELECT nextval('post_id')`,

	"getAllBoard": `
		SELECT t.board, t.id, t.postCtr, t.imageCtr, t.replyTime, t.bumpTime,
				t.subject, p.spoiler, p.time, p.name, p.trip, p.auth,
				p.imageName,
				(SELECT array_length(t.log, 1)) AS logCtr,
				i.*
			FROM threads AS t
			INNER JOIN posts AS p
				ON t.id = p.id AND (p.deleted is NULL OR p.deleted = 'false')
			LEFT OUTER JOIN images AS i
				ON p.SHA1 = i.SHA1
			ORDER BY bumpTime DESC`,

	"getBoard": `
		SELECT t.board, t.id, t.postCtr, t.imageCtr, t.replyTime, t.bumpTime,
				t.subject, p.spoiler, p.time, p.name, p.trip, p.auth,
				p.imageName,
				(SELECT array_length(t.log, 1)) AS logCtr,
				i.*
			FROM threads AS t
			INNER JOIN posts AS p
				ON t.id = p.id AND (deleted IS NULL OR deleted = 'false')
			LEFT OUTER JOIN images AS i
				ON p.SHA1 = i.SHA1
			WHERE t.board = $1
			ORDER BY bumpTime DESC`,

	"getPost": `
		SELECT op, board, editing, banned, id, time, body, name, trip, auth,
				links, backlinks, commands, images.*
			FROM posts
			LEFT OUTER JOIN images
				ON posts.SHA1 = images.SHA1
			WHERE id = $1 AND (deleted IS NULL OR deleted = 'false')`,

	"getThread": `
		SELECT board, postCtr, imageCtr, replyTime, bumpTime, subject,
				(SELECT array_length(log, 1)) AS logCtr
			FROM threads
			WHERE id = $1`,

	"getOP": `
		SELECT op, board, editing, banned, id, time, body, name, trip, auth,
				commands, images.*
			FROM posts
			LEFT OUTER JOIN images
				ON posts.SHA1 = images.SHA1
			WHERE id = $1 AND (deleted IS NULL OR deleted = 'false')`,
}
