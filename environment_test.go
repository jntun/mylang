package main

import (
	"fmt"
	"testing"
)

func TestEnvPushPop(t *testing.T) {
	env := NewEnvironment()
	env.push("1")
	env.push("2")
	env.push("3")
	env.push("4")
	env.pop()
	fmt.Println(env)
}
