// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("go", "test")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Fatalf("error $v\n", err)
	}

	version := getGitVersion()
	args := []string{}
	args = append(args, "install")
	args = append(args, "-ldflags", "-X github.com/timothyham/dicewords.VersionString="+version)
	args = append(args, "github.com/timothyham/dicewords/cmd/dicewords")
	cmd = exec.Command("go", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("error %v: %v\n", err, string(out))
	}

	cgi := []string{}
	cgi = append(cgi, "build", "-o", "dicewords.cgi")
	cgi = append(cgi, "-ldflags", "-X github.com/timothyham/dicewords.VersionString="+version)
	cgi = append(cgi, "github.com/timothyham/dicewords/cmd/dicewords-cgi")
	cmd = exec.Command("go", cgi...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("error %v: %v\n", err, string(out))
	}

}

func getGitVersion() string {
	_, err := exec.LookPath("git")
	if err != nil {
		return "git binary not found during build"
	}

	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Couldn't get git version: %v", err)
	}
	return strings.TrimSpace(string(out))
}
