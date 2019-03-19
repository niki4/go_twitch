package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type StreamInfo struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	GameID      string `json:"game_id"`
	Title       string `json:"title"`
	ViewerCount int    `json:"viewer_count"`
	ThumbURL    string `json:"thumbnail_url"`
}

type StreamsList struct {
	Streams []StreamInfo `json:"data"`
}

type Video struct {
	PlayerWidth  int
	PlayerHeight int
	ChannelName  string
}

type Chat struct {
	FrameBorderWidth int
	Scrolling        string
	ChannelID        string
	SrcURL           string
	Height           int
	Width            int
}

type EventsFrame struct {
	Total  int
	Events []string
}

type Frames struct {
	Video
	Chat
	EventsFrame
}

// ListStreams dispatch request to list all Streams
func (r *Router) ListStreams(ctx *routing.Context) error {
	authToken := ctx.Request.Header.Cookie("Authorization")
	r.logger.Info("ListStreams: Cookie", zap.ByteString("Authorization", authToken))

	resp, err := getStreamList(authToken)

	r.logger.Info("ListStreams: Twitch API response", zap.String("status", string(resp)))

	streamLst := new(StreamsList)
	if err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(streamLst); err != nil {
		r.logger.Error("ListStreams: JSON decode error", zap.Error(err))
		return err
	}

	r.logger.Info("ListStreams: trying to iterate over received list of streams",
		zap.Int("total", len(streamLst.Streams)))

	// set size for thumbnails
	for k, v := range streamLst.Streams {
		streamLst.Streams[k].ThumbURL = strings.Replace(v.ThumbURL, `{width}x{height}`, `50x50`, -1)
	}

	ctx.SetContentType("text/html")
	tmpl := template.Must(template.ParseFiles("templates/stream_list.html"))
	if err := tmpl.Execute(ctx, streamLst); err != nil {
		r.logger.Error("ListStreams: Cannot render template", zap.Error(err))
	}

	return nil
}

// ShowStreamPage dispatch request for specified Stream
func (r *Router) ShowStreamPage(ctx *routing.Context) error {
	channelName := ctx.Param("name")
	r.logger.Info("ShowStreamPage:", zap.String("ChannelName", channelName))

	channelID := ctx.QueryArgs().Peek("id")
	clientID := []byte("zhhxr55p8a8ft88s88mp0nng3ssqhd")
	chEvents, err := getChannelEvents(channelID, clientID)
	if err != nil {
		r.logger.Error("ShowStreamPage: Failure on get Channel Events", zap.Error(err))
		return err
	}

	events := new(EventsFrame)
	if err = json.NewDecoder(bytes.NewBuffer(chEvents)).Decode(events); err != nil {
		r.logger.Error("ShowStreamPage: JSON decode error", zap.Error(err))
		return err
	}

	embFrames := Frames{
		Video{
			PlayerWidth:  1280,
			PlayerHeight: 720,
			ChannelName:  channelName,
		},
		Chat{
			FrameBorderWidth: 0,
			Scrolling:        "yes",
			ChannelID:        channelName,
			SrcURL:           fmt.Sprintf("https://www.twitch.tv/embed/%s/chat", channelName),
			Height:           720,
			Width:            450,
		},
		EventsFrame{
			Total:  events.Total,
			Events: events.Events,
		},
	}

	ctx.SetContentType("text/html")
	tmpl := template.Must(template.ParseFiles("templates/stream_embed.html"))
	if err := tmpl.Execute(ctx, embFrames); err != nil {
		r.logger.Error("ShowStreamPage: Cannot render template", zap.Error(err))
	}

	return nil
}

func requestWithCookies(method, url string, cookies map[string][]byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range cookies {
		req.Header.Set(k, string(v))
	}

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

func getStreamList(authToken []byte) ([]byte, error) {
	return requestWithCookies(
		"GET",
		"https://api.twitch.tv/helix/streams?first=20",
		map[string][]byte{"Authorization": authToken})
}

func getChannelEvents(streamID, clientID []byte) ([]byte, error) {
	url := fmt.Sprintf("https://api.twitch.tv/v5/channels/%s/events", streamID)
	return requestWithCookies(
		"GET",
		url,
		map[string][]byte{"Client-ID": clientID})
}
