package migrations

import "embed"

// FS carries the .sql files into the binary so `potash migrate` works in a
// container that has nothing but the executable.
//
//go:embed *.sql
var FS embed.FS
