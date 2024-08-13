package validation

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

// REQUIREMENTS:
// yarn, so we can prepare https://codeload.github.com/Saw-mon-and-Natalie/clones-with-immutable-args/tar.gz/105efee1b9127ed7f6fedf139e1fc796ce8791f2
func TestGenesisPredeploys(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	// Get the directory of the current file
	dir := filepath.Dir(filename)

	monorepoDir := path.Join(dir, "../../optimism")
	contractsDir := path.Join(monorepoDir, "packages/contracts-bedrock")

	// chainId := 34443 // Mode mainnet

	monorepoCommit := "d80c145e0acf23a49c6a6588524f57e32e33b91c"

	executeCommandInDir(t, monorepoDir, exec.Command("git", "checkout", monorepoCommit))
	executeCommandInDir(t, monorepoDir, exec.Command("git", "fetch", "--recurse-submodules"))

	// TODO unskip these, I am skipping to save time in development
	// executeCommandInDir(t, monorepoDir, exec.Command("rm", "-rf", "node_modules"))
	// executeCommandInDir(t, contractsDir, exec.Command("rm", "-rf", "node_modules"))

	// possible optimization
	// executeCommandInDir(t, monorepoDir, exec.Command("echo", "'recursive-install=true'", ">>", ".npmrc"))
	// executeCommandInDir(t, contractsDir, exec.Command("pnpm", "install"))

	executeCommandInDir(t, contractsDir, exec.Command("pnpm", "install"))
	executeCommandInDir(t, contractsDir, exec.Command("forge", "build"))
}

func executeCommandInDir(t *testing.T, dir string, cmd *exec.Cmd) {
	t.Logf("executing %s", cmd.String())
	cmd.Dir = dir
	var outErr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		// error case : status code of command is different from 0
		fmt.Println(outErr.String())
		t.Fatal(err)
	}
}
