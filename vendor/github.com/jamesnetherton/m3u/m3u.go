package m3u

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Playlist is a type that represents an m3u playlist containing 0 or more tracks
type Playlist struct {
	Tracks []Track
}

// A Tag is a simple key/value pair
type Tag struct {
	Name string
	Value string
}

// Track represents an m3u track with a Name, Lengh, URI and a set of tags
type Track struct {
	Name   string
	Length int
	URI    string
	Tags   []Tag
}

// Parse parses an m3u playlist with the given file name and returns a Playlist
func Parse(fileName string) (playlist Playlist, err error) {
	var f io.ReadCloser
	var data *http.Response
	if strings.HasPrefix(fileName, "http://") || strings.HasPrefix(fileName, "https://") {
		data, err = http.Get(fileName)
		f = data.Body
	} else {
		f, err = os.Open(fileName)
	}
	
	if err != nil {
		err = errors.New("Unable to open playlist file")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		onFirstLine = false

	if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else if len(playlist.Tracks) == 0 {
			err = errors.New("URI provided for playlist with no tracks")
			return
		} else {
			playlist.Tracks[len(playlist.Tracks)-1].URI = strings.Trim(line, " ")
		}
	}

	return playlist, nil
}
