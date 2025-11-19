package rpeat

import (
	"fmt"
	"os"
	"os/exec"
)

func Edit(editor string, file string, update bool, host string, port string) {
	// lock file - possibly protect from reading?
	// alternate approach is to send message to server API signaling file is locked for editing
	// and send a completed request after finished.
	// if update == true, jobs will update jobs after edit/validation
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	// unlock file
	fmt.Println(err)
}

func editJob() {
	// in-browser editor
}
