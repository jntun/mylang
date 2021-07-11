package lang

import (
	"fmt"
	"testing"
)

func TestEnvPushPop(t *testing.T) {
	env := NewEnvironment("env-test")
	env.push("1")
	env.push("2")
	env.push("3")
	env.push("4")
	if len(env.blocks) != 5 {
		t.Fail()
	}
	env.pop()
	if len(env.blocks) != 4 {
		t.Fail()
	}
	fmt.Println(env)
}
