package application_test

import (
	"testing"

	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-logger"
)

func TestApp_demo1(t *testing.T) {
	t.SkipNow()
	application.New().
		Logger(logger.Default()).
		Modules(func() {
			t.Log("anonymous function")
		}).
		Run()
}
