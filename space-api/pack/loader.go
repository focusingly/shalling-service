package pack

import "embed"

//go:embed static/*
var SpaResource embed.FS
