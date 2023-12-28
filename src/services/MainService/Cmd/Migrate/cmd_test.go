package Migrate

import (
	"testing"
)

func TestModelsMigrate(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		ModelsMigrate()
	})
}

func TestESSetup(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		ESSetup()
	})
}

func TestInsertAllPostersInES(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		InsertAllPostersInES()
	})
}

func TestDropModels(t *testing.T) {
	t.Run("Good Case", func(t *testing.T) {
		DropModels()
	})
}
