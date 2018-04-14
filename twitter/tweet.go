package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

// Tweet will ... tweet the text :D ... Badum ts
func (client *Client) Tweet(text string) (*anaconda.Tweet, error) {
	v := url.Values{}
	tweet, err := client.API.PostTweet(text, v)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}
