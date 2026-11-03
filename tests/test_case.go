package tests

import (
	"github.com/goravel/framework/testing"

	"elcomware/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
