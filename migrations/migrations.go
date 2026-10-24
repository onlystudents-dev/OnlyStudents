package migrations

import "embed"

//go:embed *.sql
var FS embed.FS

//go:embed seeds/*.sql
var SeedsFS embed.FS
