package controller

import "github.com/tuananhlai/brevity-go/internal/otelsdk"

const (
	packageName = "github.com/tuananhlai/brevity-go/internal/controller"
)

var (
	appTracer = otelsdk.Tracer(packageName)
	appLogger = otelsdk.Logger(packageName)
)
