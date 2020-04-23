package minnow

import (
	"fmt"
	"io"
	"os"
	"time"
)

func makeRandomPath(baseDir Path) (Path, error) {
	name := fmt.Sprintf("%d", time.Now().UnixNano())
	path := baseDir.JoinPath(Path(name))

	// If the path somehow already exists, try again
	for path.Exists() {
		name = fmt.Sprintf("%d", time.Now().UnixNano())
		path = baseDir.JoinPath(Path(name))
	}

	err := path.Mkdir()

	if err != nil {
		return baseDir, err
	}

	return path, nil
}

func CopyFile(source, destination Path) error {
	if source.IsDir() {
		return fmt.Errorf("Cannot use Copy for a directory: %s", source)
	}

	if destination.IsDir() {
		destination = destination.JoinPath(Path(source.Name()))
	}

	if source == destination {
		return fmt.Errorf("Copy source and destination are the same.  Nothing to do.")
	}

	from, err := os.Open(string(source))

	if err != nil {
		return err
	}

	defer from.Close()

	perms, err := source.Permissions()

	if err != nil {
		return err
	}

	// Copy source permissions, making sure that the destination is
	// at least readable and writable.
	perms = perms | 0600
	to, err := os.OpenFile(string(destination), os.O_RDWR|os.O_CREATE, perms)

	if err != nil {
		return err
	}

	defer to.Close()

	_, err = io.Copy(to, from)

	if err != nil {
		return err
	}

	return nil
}
