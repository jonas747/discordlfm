package main

import (
	"errors"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/shkh/lastfm-go/lastfm"
	"log"
	"sync"
	"time"
)

var (
	flagDiscordToken   string
	flagLFMAPIKey      string
	flagLFMUsername    string
	flagNoSong         string
	flagNoSongDuration int
)

func init() {
	flag.StringVar(&flagDiscordToken, "t", "", "Discord token")
	flag.StringVar(&flagLFMAPIKey, "l", "", "Last.fm api key")
	flag.StringVar(&flagLFMUsername, "u", "", "Last.fm username")
	flag.StringVar(&flagNoSong, "g", "", "Game to set to if there hasn't been a new song for a while") //"" clears playing status
	flag.IntVar(&flagNoSongDuration, "n", 60*10, "Number of seconds without a new song for it to be considered nothing.")
	flag.Parse()
}

func main() {
	if flagDiscordToken == "" {
		log.Fatal("No discord token specified")
	}

	if flagLFMAPIKey == "" {
		log.Fatal("No lastfm api key specified")
	}

	if flagLFMUsername == "" {
		log.Fatal("No last.fm username specified")
	}

	log.Println("Starting up...")

	// Setup lastfm
	lfm := lastfm.New(flagLFMAPIKey, "")

	// Setup discord
	dsession, err := discordgo.New(flagDiscordToken)
	if err != nil {
		log.Fatal("Error creating discord session:", err)
	}

	var wg sync.WaitGroup
	dsession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		wg.Done()
	})
	wg.Add(1)

	err = dsession.Open()
	if err != nil {
		log.Fatal("Error opening discord ws conn:", err)
	}

	wg.Wait()

	log.Println("Ready received! Ctrl-c to stop.")
	run(dsession, lfm)
}

func run(s *discordgo.Session, lfm *lastfm.Api) {
	// Run continously untill somethign catches fire or an std
	ticker := time.NewTicker(time.Second * 10)

	lastPlaying := ""
	var lastPlayingTime time.Time
	setFallback := false

	for {
		<-ticker.C

		playing, err := check(lfm)
		if err != nil {
			log.Println("Error checking:", err)
			continue
		}

		if playing == lastPlaying {

			if !setFallback && time.Since(lastPlayingTime).Seconds() > float64(flagNoSongDuration) {

				err = s.UpdateStatus(0, flagNoSong)
				if err != nil {
					log.Println("Error updating status:", err)
				} else {
					if flagNoSong == "" {
						log.Println("Cleared playing status")
					} else { 
						log.Println("Updated status to:", flagNoSong)
					}
					setFallback = true
				}

			}
		} else {

			err = s.UpdateStatus(0, playing)
			if err != nil {
				log.Println("Error updating status:", err)
			} else {
				log.Println("Updated status to:", playing)
				lastPlayingTime = time.Now()
				setFallback = false
				lastPlaying = playing
			}
		}
	}
}

func check(lfm *lastfm.Api) (string, error) {
	recent, err := lfm.User.GetRecentTracks(map[string]interface{}{"user": flagLFMUsername})
	if err != nil {
		return "", err
	}

	if len(recent.Tracks) < 1 {
		return "", errors.New("No tracks")
	}

	track := recent.Tracks[0]

	return "♩" + track.Name + " by " + track.Artist.Name + "♩", nil
}
