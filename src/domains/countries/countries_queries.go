package countries

const (
	queryCreate = `INSERT INTO countries(
		id,
		code,
		name,
		flag,
		active)
	VALUES (
		:id
		:code
		:name
		:flag
		:active)`
)

