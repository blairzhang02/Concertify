package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
)

type ConcertInfo struct {
	Artist string
	Venue  string
	Date   string // You might want to use a proper time.Time depending on your use case
	// Add other fields as necessary according to the Ticketmaster API response
}
type Config struct {
	SpotifyID       string `json:"SpotifyID"`
	SpotifySecret   string `json:"SpotifySecret"`
	TicketmasterKey string `json:"TicketmasterKey"`
}

func getConfig() Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Unable to open config file", err.Error())
	}
	defer configFile.Close()

	var config Config
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

var (
	redirectURL   = "http://localhost:8080/callback"
	authenticator = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserTopRead)
)

func main() {
	config := getConfig()

	if config.SpotifyID == "" || config.SpotifySecret == "" || config.TicketmasterKey == "" {
		log.Fatal("Configuration values cannot be empty. Please check your config file.")
	}

	// Set up the Spotify authenticator with the credentials from the configuration
	authenticator = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserTopRead)
	authenticator.SetAuthInfo(config.SpotifyID, config.SpotifySecret)

	state := "some-random-state"
	http.HandleFunc("/callback", completeAuth(state, config.TicketmasterKey))
	http.HandleFunc("/auth", startAuth(state))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startAuth(state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := authenticator.AuthURL(state)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func completeAuth(state string, ticketmasterKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := authenticator.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		// Create a client using the specified token
		client := authenticator.NewClient(token)

		// Save the client for future use or continue to make your requests
		handleTopArtistsConcerts(w, r, client, ticketmasterKey)
	}
}

func getTopArtists(client spotify.Client) ([]spotify.FullArtist, error) {
	// Retrieve results from Spotify
	results, err := client.CurrentUsersTopArtists()
	if err != nil {
		return nil, err
	}
	return results.Artists, nil
}

// Assuming you have a function to get artist names from Spotify data
func getConcertsForArtists(artists []string, ticketmasterApiKey string) ([]ConcertInfo, error) {
	var concerts []ConcertInfo

	for _, artist := range artists {
		url := fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json?keyword=%s&apikey=%s", artist, ticketmasterApiKey)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err // Handle error appropriately
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Non-OK status code, handle according to your policy
			continue
		}

		// Assuming the structure of the Ticketmaster response here. Adjust according to actual response.
		var response struct {
			Embedded struct {
				Events []struct {
					Name  string `json:"name"`
					Dates struct {
						Start struct {
							DateTime string `json:"dateTime"`
						} `json:"start"`
					} `json:"dates"`
					// ... other necessary fields
				} `json:"events"`
			} `json:"_embedded"`
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err // Handle error appropriately
		}

		for _, event := range response.Embedded.Events {
			concerts = append(concerts, ConcertInfo{
				Artist: artist,
				Venue:  event.Name, // Assuming 'Name' is the venue, adjust as necessary
				Date:   event.Dates.Start.DateTime,
			})
		}
	}

	return concerts, nil
}

func handleTopArtistsConcerts(w http.ResponseWriter, r *http.Request, client spotify.Client, ticketmasterApiKey string) {
	topArtists, err := getTopArtists(client)
	if err != nil {
		http.Error(w, "Failed to retrieve top artists", http.StatusInternalServerError)
		return
	}

	var artistNames []string
	for _, artist := range topArtists {
		artistNames = append(artistNames, artist.Name)
	}

	concerts, err := getConcertsForArtists(artistNames, ticketmasterApiKey)
	if err != nil {
		http.Error(w, "Failed to retrieve concerts", http.StatusInternalServerError)
		return
	}

	// Convert concerts to JSON and send it as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concerts)
}
