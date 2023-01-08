package parser

const (
	// operators
	equal    tokenType = "="
	question tokenType = "?"

	// delimiters
	comma tokenType = ","

	// keywords
	from      tokenType = "from"
	where     tokenType = "where"
	and       tokenType = "and"
	payload   tokenType = "payload"
	header    tokenType = "header"
	with      tokenType = "with"
	get       tokenType = "get"
	post      tokenType = "post"
	delete    tokenType = "delete"
	patch     tokenType = "patch"
	put       tokenType = "put"
	basicauth tokenType = "basicauth"
	timeout   tokenType = "timeout"

	// value
	value tokenType = "value"

	// start and EOF
	start tokenType = "start"
	eof   tokenType = "EOF"
)
