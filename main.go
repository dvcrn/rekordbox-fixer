package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// MUSIC_FILEPATH is the path to search for music
const MUSIC_FILEPATH = "/Users/david/Music/Music/Media.localized/Music"
const OUT_FILE = "/Users/david/Downloads/rekordbox_new.xml"
const IN_FILE = "/Users/david/Downloads/rekordbox.xml"

type Tempo struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	XMLName xml.Name   `xml:"TEMPO"`
}
type WildcardElement struct {
	Attrs []xml.Attr `xml:",any,attr"`
}

type Track struct {
	Attrs        []xml.Attr         `xml:",any,attr"`
	XMLName      xml.Name           `xml:"TRACK"`
	TrackID      string             `xml:"TrackID,attr"`
	Name         string             `xml:"Name,attr"`
	Artist       string             `xml:"Artist,attr"`
	Composer     string             `xml:"Composer,attr"`
	Album        string             `xml:"Album,attr"`
	Grouping     string             `xml:"Grouping,attr"`
	Genre        string             `xml:"Genre,attr"`
	Kind         string             `xml:"Kind,attr"`
	Size         string             `xml:"Size,attr"`
	DiscNumber   string             `xml:"DiscNumber,attr"`
	TrackNumber  string             `xml:"TrackNumber,attr"`
	Year         string             `xml:"Year,attr"`
	AverageBpm   string             `xml:"AverageBpm,attr"`
	DateAdded    string             `xml:"DateAdded,attr"`
	BitRate      string             `xml:"BitRate,attr"`
	SampleRate   string             `xml:"SampleRate,attr"`
	Comments     string             `xml:"Comments,attr"`
	PlayCount    string             `xml:"PlayCount,attr"`
	Rating       string             `xml:"Rating,attr"`
	Location     string             `xml:"Location,attr"`
	Remixer      string             `xml:"Remixer,attr"`
	Tonality     string             `xml:"Tonality,attr"`
	Label        string             `xml:"Label,attr"`
	Mix          string             `xml:"Mix,attr"`
	Tempo        []*WildcardElement `xml:"TEMPO"`
	PositionMark []*WildcardElement `xml:"POSITION_MARK"`
}

type Collection struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	XMLName xml.Name   `xml:"COLLECTION"`
	Entries string     `xml:"Entries,attr"`
	// Attrs  []xml.Attr
	Tracks []*Track `xml:"TRACK"`
}

type Node struct {
	Attrs []xml.Attr         `xml:",any,attr"`
	Track []*WildcardElement `xml:"TRACK"`
	Node  []*Node            `xml:"NODE"`
}

type Playlists struct {
	Attrs []xml.Attr `xml:",any,attr"`
	Nodes []*Node    `xml:"NODE"`
}

type DJPlaylists struct {
	Attrs      []xml.Attr       `xml:",any,attr"`
	XMLName    xml.Name         `xml:"DJ_PLAYLISTS"`
	Collection *Collection      `xml:"COLLECTION"`
	Product    *WildcardElement `xml:"PRODUCT"`
	Playlists  *Playlists       `xml:"PLAYLISTS"`
}

func unescapeFilePath(path string) (string, error) {
	result, err := url.PathUnescape(path)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(result, "file://localhost", ""), nil
}

func escapeFilePath(path string) string {
	return "file://localhost" + url.PathEscape(path)
}

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func searchForFile(name string) string {
	foundPath := ""
	err := filepath.Walk(MUSIC_FILEPATH, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		base := filepath.Base(path)
		ext := filepath.Ext(base)
		baseWithoutExit := strings.ReplaceAll(base, ext, "")

		searchForBase := filepath.Base(name)
		searchForExt := filepath.Ext(name)
		searchForBaseWithoutExit := strings.ReplaceAll(searchForBase, searchForExt, "")

		if baseWithoutExit == searchForBaseWithoutExit {
			// fmt.Printf("===> Found matching file name: %v\n", name)
			// fmt.Printf("===> %v == %v?\n", baseWithoutExit, searchForBaseWithoutExit)

			foundPath = path
			return nil
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", err)
		return ""
	}

	return foundPath
}

func main() {
	xmlFile, err := os.Open(IN_FILE)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var djPlaylists DJPlaylists
	xml.Unmarshal(byteValue, &djPlaylists)

	newTracks := []*Track{}

	for _, track := range djPlaylists.Collection.Tracks {
		// skip catalog stuff, don't know what that is
		if strings.Contains(track.Location, "/v4/catalog") {
			continue
		}

		path, err := unescapeFilePath(track.Location)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// check if the file already exists, if yes we don't need to do anything
		if checkFileExists(path) {
			continue
		}

		base := filepath.Base(path)
		fmt.Printf("[>] Searching for: %v (%s)\n", base, path)
		newPath := searchForFile(base)
		if newPath == "" {
			fmt.Printf("[x] Couldn't find a matching file")
			continue
		}

		fmt.Printf("[O] Found matching file\n")
		fmt.Printf("[O] New: %v\n", newPath)
		fmt.Printf("[O] Old: %v\n", path)

		fmt.Printf("[O] New: %v\n", escapeFilePath(newPath))
		fmt.Printf("[O] Old: %v\n", track.Location)

		fmt.Printf("\n")

		track.Location = escapeFilePath(newPath)
		newTracks = append(newTracks, track)
	}

	// swap out tracks with new list
	djPlaylists.Collection.Tracks = newTracks

	b, err := xml.MarshalIndent(djPlaylists, "", " ")
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = os.WriteFile(OUT_FILE, b, 0666)
	if err != nil {
		log.Fatalf("error writing new file: %v\n", err)
	}
	fmt.Printf("Done writing file!\n")
}
