package json2go_test

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
	"testing"

	"github.com/stretchr/testify/require"
)

func compile(t *testing.T, src []byte) any {
	t.Helper()
	tmpdir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(tmpdir)
	data := fmt.Sprintf("package main\nvar A = *new(%s)", src)
	os.WriteFile("a.go", []byte(data), 0400)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "a.so", "a.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	require.NoError(t, err)
	plug, err := plugin.Open("a.so")
	require.NoError(t, err)
	a, err := plug.Lookup("A")
	require.NoError(t, err)
	return a
}
