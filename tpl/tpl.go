package tpl

import (
	_ "embed"
)

//go:embed root.gotpl
var RootTemplate string
