package sqlitestorage

const (
	createTablequery = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date VARCHAR(8) NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128)
	                                     )
`

	createIndexQuery = `
	CREATE INDEX IF NOT EXISTS idx_scheduler ON scheduler (date)
`

	insertTaskQuery = `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	getTasksQuery = `
	SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE date >= ?
	ORDER BY date ASC
	LIMIT ?
`

	getTasksByDateQuery = `
	SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE date = ?
	ORDER BY date ASC
	LIMIT ?
`

	getTasksByStringQuery = `
	SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE title LIKE ? OR comment LIKE ? 
	ORDER BY date ASC
	LIMIT ?
`

	getTaskByIdQuery = `
	SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE id = ?	
`

	updateTaskQuery = `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
`

	deleteTaskQuery = `DELETE FROM scheduler WHERE id = ?`
)
