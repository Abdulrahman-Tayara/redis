package resp

const (
	CRLF = "\r\n"
	CR   = "\r"
	LF   = "\n"
)

const (
	stringSuffix      = "+"
	bulkStringsSuffix = "$"
	errorSuffix       = "-"
	intSuffix         = ":"
	arraySuffix       = "*"
	nullSuffix        = "_"
	boolSuffix        = "#"
	mapSuffix         = "%"
)
