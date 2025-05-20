package pkg

import "github.com/jaevor/go-nanoid"

// `A-Za-z0-9-`.
var (
	alphabet      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	GenID, _      = nanoid.CustomUnicode(string(alphabet), 21)
	GenShortID, _ = nanoid.CustomUnicode(string(alphabet), 16)
)
