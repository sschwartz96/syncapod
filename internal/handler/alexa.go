package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Intents
const (
	PlayPodcast       = "PlayPodcast"
	PlayLatestPodcast = "PlayLatestPodcast"
	PlayNthFromLatest = "PlayNthFromLatest"
	FastForward       = "FastForward"
	Rewind            = "Rewind"
	Pause             = "AMAZON.PauseIntent"
	Resume            = "AMAZON.ResumeIntent"

	// Directives
	DirPlay       = "AudioPlayer.Play"
	DirStop       = "AudioPlayer.Stop"
	DirClearQueue = "AudioPlayer.ClearQueue"
)

// Alexa handles all requests through /api/alexa endpoint
func (h *APIHandler) Alexa(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("couldn't read the body of the request")
		// TODO: proper response here
		return
	}

	var aData AlexaData
	err = json.Unmarshal(body, &aData)
	if err != nil {
		fmt.Println("couldn't unmarshal json to object: ", err)
		// TODO: proper response here
	}

	// get the person or user accessToken
	token, err := getAccessToken(&aData)
	if err != nil {
		fmt.Println("no accessToken: ", err)
		// TODO: proper response here
	}

	// validate the token and return user
	user, err := auth.ValidateAccessToken(h.dbClient, token)
	if err != nil {
		fmt.Println("error validating token: ", err)
	}

	name := aData.Request.Intent.AlexaSlots.Podcast.Value
	fmt.Println("request name of podcast: ", name)

	var response *AlexaResponseData
	var resText, directive string
	var pod *models.Podcast
	var epi *models.Episode
	var offset int64

	fmt.Println("the requested intent: ", aData.Request.Intent.Name)
	switch aData.Request.Intent.Name {
	case PlayPodcast:
		var podcasts []models.Podcast
		err = h.dbClient.Search(database.ColPodcast, name, &podcasts)
		if err != nil {
			resText = "Error occurred searching for podcast"
			break
		}
		if len(podcasts) > 0 {
			pod = &podcasts[0]
			epi = &pod.Episodes[0]
			directive = DirPlay
		} else {
			resText = "Podcast of the name: " + name + ", not found"
		}

	case PlayLatestPodcast:
		fmt.Println("playing latest")

	case PlayNthFromLatest:

	case FastForward:
		directive = DirPlay
		pod, epi, resText, offset = h.moveAudio(&aData, true)

	case Rewind:
		directive = DirPlay
		pod, epi, resText, offset = h.moveAudio(&aData, false)

	case Pause:
		audioTokens := strings.Split(aData.Context.AudioPlayer.Token, "-")
		if len(audioTokens) > 1 {
			podID, _ := primitive.ObjectIDFromHex(audioTokens[1])
			epiID, _ := primitive.ObjectIDFromHex(audioTokens[2])
			directive = DirStop
			defer podcast.UpdateOffset(h.dbClient, user.ID, podID,
				epiID, aData.Context.AudioPlayer.OffsetInMilliseconds)
		} else {
			resText = "Please play a podcast first"
		}

	case Resume:
		splitID := strings.Split(aData.Context.AudioPlayer.Token, "-")
		if len(splitID) > 1 {
			podID, _ := primitive.ObjectIDFromHex(splitID[1])
			epiID := splitID[2]
			err := h.dbClient.FindByID(database.ColPodcast, podID, &pod)
			if err != nil {
				fmt.Println("couldn't find podcast from ID: ", err)
				resText = "Please try playing new podcast"
				break
			}

			for i, _ := range pod.Episodes {
				if pod.Episodes[i].ID.Hex() == epiID {
					epi = &pod.Episodes[i]
					break
				}
			}
		} else {
			pod, epi, offset, err = podcast.FindUserLastPlayed(h.dbClient, user.ID)
			if err != nil {
				fmt.Println("couldn't find user last played: ", err)
				resText = "Couldn't find any currently played podcast, please play new one"
				break
			}
		}

		if epi != nil {
			directive = DirPlay
			resText = "Resuming"
			if offset == 0 {
				offset = aData.Context.AudioPlayer.OffsetInMilliseconds
			}
		} else {
			resText = "Episode not found, please try playing new podcast"
		}

	default:
		resText = "This command is currently not supported, please request"
	}

	// If we are creating an alexa audio repsonse
	if directive != "" {
		// get details from non-nil episode
		if epi != nil {
			if resText == "" {
				resText = "Playing " + pod.Title + ", " + epi.Title
			}
			if offset == 0 {
				offset = podcast.FindOffset(h.dbClient, user, epi)
			}
			fmt.Println("offset: ", offset)

			response = createAudioResponse(directive, user.ID.Hex(),
				resText, pod, epi, offset)
		} else {
			response = createPauseResponse(directive)
		}
	} else {
		response = createEmptyResponse(resText)
	}

	jsonRes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("couldn't marshal alexa response: ", err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonRes)
}

// moveAudio takes pointer to aData and bool for direction
// returns pointers to podcast and episode, response text and offset in millis
func (h *APIHandler) moveAudio(aData *AlexaData, forward bool) (*models.Podcast, *models.Episode, string, int64) {
	var pod *models.Podcast
	var epi *models.Episode
	var resText string
	var offset int64
	var err error

	audioTokens := strings.Split(aData.Context.AudioPlayer.Token, "-")
	if len(audioTokens) > 1 {
		pID, _ := primitive.ObjectIDFromHex(audioTokens[1])
		eID, _ := primitive.ObjectIDFromHex(audioTokens[2])

		// find episode
		pod, epi, err = podcast.FindPodcastEpisode(h.dbClient, pID, eID)
		if err != nil {
			fmt.Println("error finding podcast episode", err)
			resText = "Error occurred, please try again"
			return nil, nil, resText, 0
		}

		// get the current time and duration to move
		curTime := aData.Context.AudioPlayer.OffsetInMilliseconds
		dura := convertISO8601ToMillis(aData.Request.Intent.AlexaSlots.Duration.Value)
		durString := durationToText(time.Millisecond * time.Duration(dura))

		fmt.Printf("cur time: %v, aData: %v, duration calculated: %v\n", curTime, aData.Request.Intent.AlexaSlots.Duration.Value, dura)

		fmt.Println("durString: ", durString)

		if forward {
			offset = curTime + dura
			resText = "Fast-forwarded " + durString
		} else {
			offset = curTime - dura
			resText = "Rewound " + durString
		}

		if offset < 0 {
			offset = 1
		} else {
			// check if duration does not exist
			if epi.Duration == 0 {
				epi.Duration = podcast.FindLength(epi.Enclosure.MP3)
				go func() {
					err := podcast.UpdateEpisode(h.dbClient, pod, epi)
					if err != nil {
						fmt.Println("error updating episoe: ", err)
					}
				}()
			}

			// check if we are trying to fast forward past end of episode
			if epi.Duration < offset {
				tilEnd := time.Duration(epi.Duration-curTime) * time.Millisecond
				resText = "Cannot fast forward further than: " + durationToText(tilEnd)
				offset = 1
			}
		}
	} else {
		resText = "Please play a podcast first"
	}

	return pod, epi, resText, offset
}

func durationToText(dur time.Duration) string {
	bldr := strings.Builder{}
	if int(dur.Hours()) == 1 {
		bldr.WriteString("1 hour, ")
	} else if dur.Hours() > 1 {
		bldr.WriteString(strconv.Itoa(int(dur.Hours())))
		bldr.WriteString(" hours, ")
	}
	dur = dur - dur.Truncate(time.Hour)

	if int(dur.Minutes()) == 1 {
		bldr.WriteString("1 minute, ")
	} else if dur.Minutes() > 1 {
		bldr.WriteString(strconv.Itoa(int(dur.Minutes())))
		bldr.WriteString(" minutes, ")
	}
	dur = dur - dur.Truncate(time.Minute)

	if int(dur.Seconds()) == 1 {
		bldr.WriteString("1 second, ")
	} else if dur.Seconds() > 1 {
		bldr.WriteString(strconv.Itoa(int(dur.Seconds())))
		bldr.WriteString(" seconds, ")
	}

	return bldr.String()
}

func createAudioResponse(directive, userID, text string,
	pod *models.Podcast, epi *models.Episode, offset int64) *AlexaResponseData {

	imgURL := epi.Image.URL
	if imgURL == "" {
		imgURL = pod.Image.URL
		if imgURL == "" {
			// TODO: add custom generic defualt image
			imgURL = "https://emby.media/community/uploads/inline/355992/5c1cc71abf1ee_genericcoverart.jpg"
		}
	}

	return &AlexaResponseData{
		Version: "1.0",
		Response: AlexaResponse{
			Directives: []AlexaDirective{
				{
					Type:         directive,
					PlayBehavior: "REPLACE_ALL",
					AudioItem: AlexaAudioItem{
						Stream: AlexaStream{
							URL:                  epi.Enclosure.MP3,
							Token:                userID + "-" + pod.ID.Hex() + "-" + epi.ID.Hex(),
							OffsetInMilliseconds: offset,
						},
						Metadata: AlexaMetadata{
							Title:    epi.Title,
							Subtitle: epi.Subtitle,
							Art: AlexaArt{
								Sources: []AlexaURL{
									AlexaURL{
										URL:    imgURL,
										Height: 144,
										Width:  144,
									},
								},
							},
						},
					},
				},
			},
			OutputSpeech: AlexaOutputSpeech{
				Type: "PlainText",
				Text: text,
			},
			ShouldEndSession: true,
		},
	}
}

func createPauseResponse(directive string) *AlexaResponseData {
	return &AlexaResponseData{
		Version: "1.0",
		Response: AlexaResponse{
			Directives: []AlexaDirective{
				{
					Type: directive,
				},
			},
			OutputSpeech: AlexaOutputSpeech{
				Type: "PlainText",
				Text: "Paused",
			},
			ShouldEndSession: true,
		},
	}
}

func createEmptyResponse(text string) *AlexaResponseData {
	return &AlexaResponseData{
		Version: "1.0",
		Response: AlexaResponse{
			Directives: nil,
			OutputSpeech: AlexaOutputSpeech{
				Type:         "PlainText",
				Text:         text,
				PlayBehavior: "REPLACE_ENQUEUE",
			},
			ShouldEndSession: true,
		},
	}
}

func convertISO8601ToMillis(data string) int64 {
	data = data[2:]

	var durRegArr [3]*regexp.Regexp
	var durStrArr [3]string
	var durIntArr [3]int64

	durRegArr[0], _ = regexp.Compile("([0-9]+)H")
	durRegArr[1], _ = regexp.Compile("([0-9]+)M")
	durRegArr[2], _ = regexp.Compile("([0-9]+)S")

	for i, _ := range durStrArr {
		durStrArr[i] = durRegArr[i].FindString(data)
		if len(durStrArr[i]) > 1 {
			str := durStrArr[i]
			val, _ := strconv.Atoi(str[:len(str)-1])
			durIntArr[i] = int64(val)
		}
	}

	return (durIntArr[0])*int64(3600000) +
		(durIntArr[1])*int64(60000) +
		(durIntArr[2])*int64(1000)
}

func getAccessToken(data *AlexaData) (string, error) {
	if data.Context.System.Person.AccessToken != "" {
		return data.Context.System.Person.AccessToken, nil
	} else if data.Context.System.User.AccessToken != "" {
		return data.Context.System.User.AccessToken, nil
	}
	return "", errors.New("no accessToken")
}

// AlexaData contains all the informatino and data from request sent from alexa
type AlexaData struct {
	Version string       `json:"version,omitempty"`
	Context AlexaContext `json:"context,omitempty"`
	Request AlexaRequest `json:"request,omitempty"`
}

// AlexaContext contains system
type AlexaContext struct {
	System      AlexaSystem      `json:"System,omitempty"`
	AudioPlayer AlexaAudioPlayer `json:"AudioPlayer,omitempty"`
}

// AlexaSystem is the container for person and user
type AlexaSystem struct {
	Person AlexaPerson `json:"person,omitempty"`
	User   AlexaUser   `json:"user,omitempty"`
}

// AlexaAudioPlayer contains info of the currently played track if available
type AlexaAudioPlayer struct {
	OffsetInMilliseconds int64  `json:"offsetInMilliseconds,omitempty"`
	Token                string `json:"token,omitempty"`
	PlayActivity         string `json:"playActivity,omitempty"`
}

// AlexaPerson holds the info about the person who explicitly called the skill
type AlexaPerson struct {
	PersonID    string `json:"personId,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

// AlexaUser contains info about the user that holds the skill
type AlexaUser struct {
	UserID      string `json:"userId,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
}

// AlexaRequest holds all the information and data
type AlexaRequest struct {
	Type                 string      `json:"type,omitempty"`
	RequestID            string      `json:"requestId,omitempty"`
	Timestamp            time.Time   `json:"timestamp,omitempty"`
	Token                string      `json:"token,omitempty"`
	OffsetInMilliseconds int64       `json:"offsetInMilliseconds,omitempty"`
	Intent               AlexaIntent `json:"intent,omitempty"`
}

// AlexaIntent holds information and data of intent sent from alexa
type AlexaIntent struct {
	Name       string     `json:"name,omitempty"`
	AlexaSlots AlexaSlots `json:"slots,omitempty"`
}

// AlexaSlots are the container for the slots
type AlexaSlots struct {
	Nth      AlexaSlot `json:"nth,omitempty"`
	Episode  AlexaSlot `json:"episode,omitempty"`
	Podcast  AlexaSlot `json:"podcast,omitempty"`
	Duration AlexaSlot `json:"duration,omitempty"`
}

// AlexaSlot holds information of the slot for the intent
type AlexaSlot struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// AlexaResponseData contains the version and response
type AlexaResponseData struct {
	Version  string        `json:"version,omitempty"`
	Response AlexaResponse `json:"response,omitempty"`
}

// AlexaResponse contains the actual response
type AlexaResponse struct {
	Directives       []AlexaDirective  `json:"directives,omitempty"`
	OutputSpeech     AlexaOutputSpeech `json:"outputSpeech,omitempty"`
	ShouldEndSession bool              `json:"shouldEndSession,omitempty"`
}

// AlexaDirective tells alexa what to do
type AlexaDirective struct {
	Type         string         `json:"type,omitempty"`
	PlayBehavior string         `json:"playBehavior,omitempty"`
	AudioItem    AlexaAudioItem `json:"audioItem,omitempty"`
}

// AlexaAudioItem holds information of audio track
type AlexaAudioItem struct {
	Stream   AlexaStream   `json:"stream,omitempty"`
	Metadata AlexaMetadata `json:"metadata,omitempty"`
}

type AlexaStream struct {
	Token                string `json:"token,omitempty"`
	URL                  string `json:"url,omitempty"`
	OffsetInMilliseconds int64  `json:"offsetInMilliseconds,omitempty"`
}

type AlexaMetadata struct {
	Title    string   `json:"title,omitempty"`
	Subtitle string   `json:"subtitle,omitempty"`
	Art      AlexaArt `json:"art,omitempty"`
}

type AlexaArt struct {
	Sources []AlexaURL `json:"sources,omitempty"`
}

type AlexaURL struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

// AlexaOutputSpeech takes type: "PlainText", text, and playBehavior: REPLACE_ENQUEUE
type AlexaOutputSpeech struct {
	Type         string `json:"type,omitempty"`
	Text         string `json:"text,omitempty"`
	PlayBehavior string `json:"playBehavior,omitempty"`
}
