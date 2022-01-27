package seasons

const (
	queryCreate = `INSERT INTO seasons (id) VALUES (?)`
	queryFind   = `SELECT id FROM seasons WHERE id = ? LIMIT 1`
	queryList   = `SELECT id FROM seasons %s ORDER BY %s`
	queryDelete = `DELETE FROM seasons WHERE id = ?`
)
