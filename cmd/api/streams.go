package api

import (
	"bytes"
	"encoding/json"
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

type VideoFrame struct {
	PlayerWidth  int
	PlayerHeight int
	ChannelName  string
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
	channelName := ctx.Param("id")
	r.logger.Info("ShowStreamPage:", zap.String("ChannelName", channelName))

	vFrame := VideoFrame{
		PlayerWidth:  640,
		PlayerHeight: 480,
		ChannelName:  channelName,
	}

	ctx.SetContentType("text/html")
	tmpl := template.Must(template.ParseFiles("templates/stream_embed.html"))
	if err := tmpl.Execute(ctx, vFrame); err != nil {
		r.logger.Error("ShowStreamPage: Cannot render template", zap.Error(err))
	}

	return nil
}

func requestWithAuthorization(method, url string, authToken []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(authToken))

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
	return requestWithAuthorization("GET", "https://api.twitch.tv/helix/streams?first=20", authToken)
}
