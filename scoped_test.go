package gdit_test

import (
	"testing"

	"github.com/saweima12/gdit"
)

func TestScope(t *testing.T) {
	app := getTestApp()

	scope := app.GetScope("Hello").(*gdit.Scope)
	t.Run("scopeName should be `Hello`", func(t *testing.T) {
		if scope.Name != "Hello" {
			t.Fail()
		}
	})

	t.Run("scopeName should be `Hello`", func(t *testing.T) {
		if scope.Name != "Hello" {
			t.Fail()
		}
	})
}
