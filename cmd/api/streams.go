package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
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

// ListStreams dispatch request to list all Streams
func (r *Router) ListStreams(ctx *routing.Context) error {
	authToken := ctx.Request.Header.Cookie("Authorization")
	r.logger.Info("ListStreams: Cookie", zap.ByteString("Authorization", authToken))

	resp, err := requestStreamList(authToken)

	r.logger.Info("ListStreams: Twitch API response", zap.String("status", string(resp)))

	streamLst := new(StreamsList)
	if err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(streamLst); err != nil {
		r.logger.Error("ListStreams: JSON decode error", zap.Error(err))
		return err
	}

	r.logger.Info("ListStreams: trying to iterate over received list of streams", zap.Int("total", len(streamLst.Streams)))
	for k, v := range streamLst.Streams {
		fmt.Fprintf(ctx, "%v - %v\n", k, v)
	}

	return nil
}

// ShowStreamPage dispatch request for specified Stream
func (r *Router) ShowStreamPage(ctx *routing.Context) error {

	return nil
}

func requestStreamList(authToken []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?first=20", nil)
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
