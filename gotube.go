package gotube

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const userAgent = `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:62.0) Gecko/20100101 Firefox/62.0`

type Youtube struct {
	videoID  string
	VideoURL string
	VideoInfo
}

type VideoInfo struct {
	URL     string
	Quality string
	Ext     string
	Name    string
}

func (y *Youtube) api(method string) ([]byte, error) {
	urlIndex := "http://www.youtube.com/" + method
	req, err := http.NewRequest(http.MethodGet, urlIndex, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)

	q := req.URL.Query()
	q.Add("video_id", y.videoID)
	req.URL.RawQuery = q.Encode()
	log.Printf("Request api url=%s", req.URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (y *Youtube) getVideoInfo() error {
	body, err := y.api("get_video_info")

	m, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}

	y.VideoInfo.Name = m.Get("title")

	qs := m.Get("url_encoded_fmt_stream_map")
	qs = strings.Replace(qs, ",", "&", -1)
	values, err := url.ParseQuery(qs)
	if err != nil {
		return err
	}

	y.VideoInfo.Quality = values.Get("quality")
	y.VideoInfo.URL = values.Get("url")
	ext := strings.Split(values.Get("type"), ";")
	ext = strings.Split(ext[0], "/")
	y.VideoInfo.Ext = ext[1]

	log.Printf("Video: %+v\n", y.VideoInfo)
	return nil
}

func (y *Youtube) GetVideo() error {
	parse := strings.Split(y.VideoURL, "?")
	paras, err := url.ParseQuery(parse[1])
	if err != nil {
		return err
	}

	y.videoID = paras.Get("v")
	err = y.getVideoInfo()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.%s", FileNameCorrect(y.Name), y.Ext)
	path := filepath.Join(".", fileName)
	if err := DownloadFile(path, y.URL); err != nil {
		return err
	}

	return nil
}

func FileNameCorrect(s string) string {
	var re = regexp.MustCompile(`[\/:*?"<>|]`)
	s = re.ReplaceAllString(s, " ")
	return s
}

func DownloadFile(filepath, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
