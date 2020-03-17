package theater

import (
	"errors"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetMovieDuration(uri string) (duration int64, err error) {
	command := exec.Command("ffprobe", "-show_streams", uri)
	output, err := command.Output()
	if err != nil {
		return 0, err
	}
	strOutput := string(output)
	dRegex := regexp.MustCompile(`(?m)duration=([0-9]+).([0-9]+)`)
	if len(dRegex.FindStringIndex(strOutput)) > 0 {
		strDuration := strings.ReplaceAll(dRegex.FindString(strOutput), "duration=", "")
		splits := strings.Split(strDuration, ".")
		durationInt, err := strconv.Atoi(splits[0])
		if err != nil {
			return 0, err
		}
		return int64(durationInt), nil
	}
	err = errors.New("could not get the duration parameter")
	return
}

func GetMovieFileSize(uri string) (int64, error) {

	resp, err := http.Head(uri)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}

	sizeInt, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}

	return int64(sizeInt), nil
}