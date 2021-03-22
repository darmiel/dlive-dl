package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var pattern *regexp.Regexp

func init() {
	var err error
	if pattern, err = regexp.Compile(`"playbackUrl":"((.*)\.m3u8)"`); err != nil {
		log.Fatalln("Error compiling pattern:", err)
	}
}

func cmdDownload(context *cli.Context) error {
	url := context.String("url")
	log.Println("☁️ Parsing:", url)

	// make request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	str := string(body)
	var u string

	allString := pattern.FindAllStringSubmatch(str, 2)
	for _, s := range allString {
		va := s[1]
		u = va
		break
	}

	if u == "" {
		return errors.New("playback not found")
	}

	return dlM3u8AndAskForPlaylist(u)
}

func dlM3u8AndAskForPlaylist(url string) (err error) {
	url, err = strconv.Unquote("\"" + url + "\"")
	if err != nil {
		return
	}
	log.Println("☁️ Parsing segments from:", url)

	// read url
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if listType != m3u8.MASTER {
		return errors.New("playback is not a master")
	}

	master := p.(*m3u8.MasterPlaylist)

	var opts []string
	optVal := make(map[string]*m3u8.Variant)
	for _, v := range master.Variants {
		elems := v.Resolution + " (" + v.Video + "), Codec: " + v.Codecs
		optVal[elems] = v
		opts = append(opts, elems)
	}
	q := &survey.Select{
		Message: "Select Variant:",
		Options: opts,
	}
	var variantName string
	if err = survey.AskOne(q, &variantName); err != nil {
		return
	}

	variant, ok := optVal[variantName]
	if !ok || variant == nil {
		return errors.New("unknown variant selected")
	}
	if variant.URI == "" {
		return errors.New("variant has no url")
	}

	plURI := variant.URI

	return dlM3u8Playlist(plURI)
}

func dlM3u8Playlist(url string) (err error) {
	log.Println("☁️ Parsing playlist from:", url)

	// read url
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if listType != m3u8.MEDIA {
		return errors.New("playback is not media")
	}

	media := p.(*m3u8.MediaPlaylist)

	var segments []string
	for _, seg := range media.Segments {
		if seg == nil {
			continue
		}

		uri := seg.URI
		segments = append(segments, uri)
	}

	log.Println("  🚀 Parsed", len(segments), "to download.")
	return dlM3u8Segments(segments)
}

func dlM3u8Segments(segments []string) (err error) {
	// ask for download location
	var file string
	var f *os.File

	for {
		prompt := &survey.Input{
			Message: "Output File:",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}
		if err = survey.AskOne(prompt, &file); err != nil {
			return
		}
		if !strings.HasSuffix(file, ".ts") {
			file = file + ".ts"
		}

		log.Println("☁️ Downloading to:", file)

		if _, err := os.Stat(file); os.IsExist(err) {
			log.Println("🤬 File already exists")
			continue
		}

		if f, err = os.Create(file); err != nil {
			log.Println("🤬 Error creating file:", err)
			continue
		}

		break
	}

	log.Println("☁️ Downloading", len(segments), "segments ...")
	fmt.Println()

	bar := progressbar.Default(int64(len(segments)))
	errno := 0

	for _, seg := range segments {
		if err := bar.Add(1); err != nil {
			log.Println("error updating progressbar:", err)
		}

		res, err := http.Get(seg)
		if err != nil {
			errno++
			fmt.Println(" ERROR: ", err)
			continue
		}

		var buff bytes.Buffer
		if _, err := buff.ReadFrom(res.Body); err != nil {
			errno++
			fmt.Println(" ERROR: ", err)
			res.Body.Close()
			continue
		}

		if _, err := f.Write(buff.Bytes()); err != nil {
			errno++
			fmt.Println(" ERROR: ", err)
			res.Body.Close()
			continue
		}

		res.Body.Close()
	}

	if err := f.Close(); err != nil {
		log.Println("ERROR closing file:", err)
	}

	fmt.Println()
	log.Println("☁️ All segments downloaded!", errno, "errors encountered.")

	return nil
}
