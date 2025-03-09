package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/hajimehoshi/oto"
	"github.com/dhowden/tag"
)

type MusicPlayer struct {
	songs         []string
	currentSong   string
	audioPlayer   *oto.Player
	shuffleState  bool
	currentIndex  int
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		shuffleState: false,
		currentIndex: -1,
	}
}

// Load songs from a specified folder
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

// Play a song
func (mp *MusicPlayer) PlaySong(songPath string) error {
	// Open the file and start playing (using oto)
	mp.currentSong = songPath
	file, err := os.Open(songPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tag, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}

	// Print song info (title and artist)
	fmt.Printf("Playing: %s by %s\n", tag.Title(), tag.Artist())
	mp.currentIndex++
	return nil
}

// Play next song
func (mp *MusicPlayer) PlayNextSong() {
	if len(mp.songs) == 0 {
		fmt.Println("No songs loaded!")
		return
	}

	// If shuffle is enabled, shuffle the list
	if mp.shuffleState {
		mp.ShuffleSongs()
	}

	// If currentIndex is the last song, restart from the beginning
	if mp.currentIndex >= len(mp.songs)-1 {
		mp.currentIndex = -1
	}

	// Play the next song
	nextSong := mp.songs[mp.currentIndex+1]
	err := mp.PlaySong(nextSong)
	if err != nil {
		log.Println("Error playing next song:", err)
	}
}

// Shuffle the playlist
func (mp *MusicPlayer) ShuffleSongs() {
	// Implement the shuffle logic here
	// You can use the rand package to shuffle the song list
}

// Toggle shuffle mode
func (mp *MusicPlayer) ToggleShuffle() {
	mp.shuffleState = !mp.shuffleState
}

// GTK UI Initialization
func initializeUI(mp *MusicPlayer) {
	// Initialize GTK
	gtk.Init(nil)

	// Create a new window
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Simple Music Player")
	window.SetDefaultSize(400, 200)

	// Create a Box to hold widgets (buttons, labels, etc.)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	window.Add(box)

	// Load songs button
	loadButton, _ := gtk.ButtonNewWithLabel("Load Songs")
	loadButton.Connect("clicked", func() {
		folderDialog := gtk.FileChooserDialogNewWith1Button("Select Music Folder", window, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER, "Select", gtk.RESPONSE_ACCEPT)
		folderDialog.Connect("response", func(dlg *gtk.FileChooserDialog, response int) {
			if response == gtk.RESPONSE_ACCEPT {
				folderPath := dlg.GetFilename()
				err := mp.LoadSongsFromFolder(folderPath)
				if err != nil {
					fmt.Println("Error loading songs:", err)
				}
			}
		})
		folderDialog.Run()
		folderDialog.Destroy()
	})

	// Play next song button
	playNextButton, _ := gtk.ButtonNewWithLabel("Play Next")
	playNextButton.Connect("clicked", func() {
		mp.PlayNextSong()
	})

	// Shuffle toggle button
	shuffleButton, _ := gtk.ButtonNewWithLabel("Toggle Shuffle")
	shuffleButton.Connect("clicked", func() {
		mp.ToggleShuffle()
		if mp.shuffleState {
			shuffleButton.SetLabel("Shuffle ON")
		} else {
			shuffleButton.SetLabel("Shuffle OFF")
		}
	})

	// Add buttons to the box
	box.PackStart(loadButton, false, false, 0)
	box.PackStart(playNextButton, false, false, 0)
	box.PackStart(shuffleButton, false, false, 0)

	// Show all UI elements
	window.ShowAll()

	// Start the GTK main loop
	gtk.Main()
}

func main() {
	// Initialize music player
	player := NewMusicPlayer()

	// Initialize and run the UI
	initializeUI(player)
}
