// This tool executes suricata-update on windows.  It embeds knowledge of the locations of the suricata
// executable and suricata paths path in the expanded 'zdeps/suricata'
// directory inside a Brim installation.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// These paths are relative to the zdeps/suricata directory.
var (
	execRelPath = "bin/suricata-update.exe"
)

// zdepsSuricataDirectory returns the absolute path of the zdeps/suricata directory,
// based on the assumption that this executable is located directly in it.
func zdepsSuricataDirectory() (string, error) {
	execFile, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execFile), nil
}

func makeConfig(baseDir, dest string) error {
	ruleConfig := fmt.Sprintf(`
data-directory: %s\var\lib\suricata
dist-rule-directory: %s\share\suricata\rules
`, baseDir, baseDir)

	return ioutil.WriteFile(filepath.Join(baseDir, dest), []byte(ruleConfig), 0644)
}

func runSuricataUpdate(baseDir, execPath string) error {
	cmd := exec.Command(execPath,
		"--suricata", filepath.Join(baseDir, "bin/suricata.exe"),
		"--config", filepath.Join(baseDir, "update.yaml"),
		"--suricata-conf", filepath.Join(baseDir, "brim-conf.yaml"),
		"--no-test",
		"--no-reload")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	path := fmt.Sprintf("PATH=%s;%s", filepath.Join(baseDir, "dlls"), os.Getenv("PATH"))
	cmd.Env = append(os.Environ(), path)

	return cmd.Run()
}

func main() {
	baseDir, err := zdepsSuricataDirectory()
	if err != nil {
		log.Fatalln("zdepsSuricataDirectory failed:", err)
	}

	if err := makeConfig(baseDir, "update.yaml"); err != nil {
		log.Fatalln("makeConfig failed:", err)
	}

	execPath := filepath.Join(baseDir, filepath.FromSlash(execRelPath))
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalln("suricata-update executable not found at", execPath)
	}

	err = runSuricataUpdate(baseDir, execPath)
	if err != nil {
		log.Fatalln("launchSuricata failed", err)
	}
}