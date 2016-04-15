package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/gorilla/mux"
	"github.com/schmatz/hn-messenger-bot/messenger"
)

var (
	pageAccessToken   = flag.String("page-access-token", "", "The page access token")
	verificationToken = flag.String("token", "", "The challenge verification token")
	port              = flag.Uint("port", 3000, "The port to listen on")
	bot               *messenger.Bot
)

func main() {
	flag.Parse()

	if *pageAccessToken == "" || *verificationToken == "" {
		flag.Usage()
		log.Fatal("Page access token and verification token are required parameters")
	}

	bot = messenger.New(*pageAccessToken, *verificationToken, handleMessaging)

	r := mux.NewRouter()
	r.HandleFunc("/webhook/", bot.HandleVerificationChallenge).Methods("GET")
	r.HandleFunc("/webhook/", bot.HandleWebhookPost).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func handleMessaging(m messenger.Messaging) (err error) {
	topStories, err := getHNTopStoryIDs()
	if err != nil {
		return
	}

	var templateItems []messenger.GenericTemplateElement

	numStories := 5
	for i := 0; i < numStories; i++ {
		topStory := topStories[i]
		title, url, description, err := getHNStoryDetails(topStory)
		if err != nil {
			return err
		}

		item := messenger.GenericTemplateElement{
			Title:    title,
			Subtitle: description,
			ItemURL:  url,
		}
		templateItems = append(templateItems, item)
	}

	err = bot.SendGenericTemplateReply(m.Sender.ID, templateItems)

	return
}

func getHNTopStoryIDs() (ids []int64, err error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ids)

	return
}

func getHNStoryDetails(storyID int64) (title string, url string, description string, err error) {
	var story struct {
		Title  string `json:"title"`
		URL    string `json:"url"`
		Author string `json:"by"`
		Points int64  `json:"score"`
		Time   int64  `json:"time"`
	}

	resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyID))
	if err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&story)
	if err != nil {
		return
	}

	title = story.Title
	url = story.URL
	description = fmt.Sprintf("%d points by %s %s", story.Points, story.Author, humanize.Time(time.Unix(story.Time, 0)))

	return
}
