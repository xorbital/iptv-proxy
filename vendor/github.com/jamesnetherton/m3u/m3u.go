package m3u

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Playlist is a type that represents an m3u playlist containing 0 or more tracks
type Playlist struct {
	Tracks []Track
}

// A Tag is a simple key/value pair
type Tag struct {
	Name  string
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
		 if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else {
			 track := &Track{fileName, 0, "", nil}
			playlist.Tracks = append(playlist.Tracks, *track)
			playlist.Tracks[0].URI = strings.Trim(line, " ")
		}
	}
	if len(playlist.Tracks) == 0 {
		err = errors.New("URI provided for playlist with no tracks")
	}
	return playlist, nil
}

// Marshall Playlist to an m3u file.
func Marshall(p Playlist) (io.Reader, error) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	if err := MarshallInto(p, w); err != nil {
		return nil, err
	}

	return buf, nil
}

// MarshallInto a *bufio.Writer a Playlist.
func MarshallInto(p Playlist, into *bufio.Writer) error {
	into.WriteString("#EXTM3U\n")
	for _, track := range p.Tracks {
		into.WriteString("#EXTINF:")
		into.WriteString(fmt.Sprintf("%d ", track.Length))
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				into.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value))
				continue
			}
			into.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value))
		}
		into.WriteString(", ")

		into.WriteString(fmt.Sprintf("%s\n%s\n", track.Name, track.URI))
	}

	return into.Flush()
}
