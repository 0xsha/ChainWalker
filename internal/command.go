package internal

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

func ExecuteCommand(command string, timeOut time.Duration, args ...string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, command, args...)
	out, _ := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("command timed out")
	}

	// will do
	//if err != nil {
	//	return "", err
	//}

	return string(out), nil

}
