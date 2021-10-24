package main

import (
	"knative.dev/pkg/injection/sharedmain"

	loudvents "github.com/odacremolbap/loudvents/pkg/reconciler/loudvents/controller"
)

func main() {
	sharedmain.Main("loudvents-controller", loudvents.NewController)
}
