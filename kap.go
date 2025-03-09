package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/bogem/id3v2"
	"github.com/hajimehoshi/oto"
	"github.com/dhowden/tag"
	"github.com/gookit/color"
)

type MusicPlayer struct {
	songs         []string
	currentSong   string
	audioPlayer   *oto.Player
	shuffleState  bool
	lyrics        string
	equalizer     map[string]int
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		shuffleState: false,
		equalizer:    make(map[string]int),
	}
}

// Load songs from the specified directory
func (mp *MusicPlayer) LoadSongsFromFolder(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp3" {
			mp.songs = append(mp.songs, filepath.Join(folderPath, file.Name()))
		}
	}
	return nil
}

// Shuffle songs
func (mp *MusicPlayer) ShuffleSongs() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(mp.songs), func(i, j int) {
		mp.songs[i], mp.songs[j] = mp.songs[j], mp.songs[i]
	})
}

// Play a song
func (mp *MusicPlayer) PlaySong(songPath string) error {
	// Here we should add functionality to play audio, currently just outputting song name
	mp.currentSong = songPath

	// Open MP3 file and extract metadata (song name)
	file, err := os.Open(songPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read metadata from MP3
	tag, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}

	// Print the song name and artist (if available)
	color.Green("Now playing: %s - %s", tag.Title(), tag.Artist())
	mp.fetchLyrics(tag.Title(), tag.Artist())
	return nil
}

// Fetch lyrics (this is a simplified version)
func (mp *MusicPlayer) fetchLyrics(songTitle, artist string) {
	// In a real case, you'd fetch lyrics from an API like Genius.
	// For now, we'll just simulate this.
	mp.lyrics = fmt.Sprintf("Fetching lyrics for %s by %s...", songTitle, artist)
	color.Yellow(mp.lyrics)
}

// Play next song in folder
func (mp *MusicPlayer) PlayNextSong() {
	if len(mp.songs) == 0 {
		color.Red("No songs loaded!")
		return
	}

	if mp.shuffleState {
		mp.ShuffleSongs()
	}

	nextSong := mp.songs[0]
	err := mp.PlaySong(nextSong)
	if err != nil {
		log.Println("Error playing next song:", err)
	}
}

// Toggle shuffle mode
func (mp *MusicPlayer) ToggleShuffle() {
	mp.shuffleState = !mp.shuffleState
	state := "off"
	if mp.shuffleState {
		state = "on"
	}
	color.Cyan("Shuffle is now %s", state)
}

func main() {
	// Initialize the music player
	player := NewMusicPlayer()

	// Specify your music folder path
	musicFolder := "./music" // Adjust to your folder

	// Load songs from folder
	err := player.LoadSongsFromFolder(musicFolder)
	if err != nil {
		log.Fatal(err)
	}

	// Simulating user actions
	player.PlayNextSong()    // Play first song
	player.ToggleShuffle()   // Enable shuffle
	player.PlayNextSo
