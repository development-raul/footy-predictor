package countries

const (
	queryCreate = `INSERT INTO countries(
		as_id,
		code,
		name,
		flag,
		active)
	VALUES (
		:as_id,
		:code,
		:name,
		:flag,
		:active)`

	queryUpdate = `UPDATE countries
	  SET
		code = :code,
		name = :name,
		flag = :flag,
		active = :active
	  WHERE
		id = :id`

	queryFindByID   = `SELECT * FROM countries WHERE id = ? LIMIT 1`
	queryFindByAsID = `SELECT * FROM countries WHERE as_id = ? LIMIT 1`

	queryList      = `SELECT * FROM countries %s ORDER BY %s %s`
	queryListTotal = `SELECT count(id) FROM countries %s`

	queryDelete = `DELETE FROM countries WHERE id = ?`
)
