package types

import "log"

const (
	LOG_FLAGS = log.Ldate | log.Ltime

	LOG_MODEL   = "  \x1b[34;1m[MODEL]\x1b[39;m "
	LOG_SCANNER = "\x1b[34;1m[SCANNER]\x1b[39;m "
	LOG_SENDER  = " \x1b[34;1m[SENDER]\x1b[39;m "
	LOG_MONITOR = "\x1b[34;1m[MONITOR]\x1b[39;m "
	LOG_ERRS    = "  \x1b[31;1m[ERROR]\x1b[39;m "
)
