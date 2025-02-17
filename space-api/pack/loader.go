package pack

import "embed"

//go:embed static/*
var SpaResource embed.FS

//go:embed templates/check-billing-fault.html
var CheckBillingFaultTemplate []byte

//go:embed templates/billing-detail.html
var BillingSubsCostTemplate []byte
