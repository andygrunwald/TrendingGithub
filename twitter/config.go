package twitter

import (
	"log"
	"net/url"
	"time"
)

// LoadConfiguration requests the latest configuration
// settings from twitter and stores those in memory.
// Config settings like numbers of chars needed for the short url
// and so on are part of this. Thos settings
// are importent to know how much text we can tweet later
// See https://dev.twitter.com/rest/reference/get/help/configuration
func (client *Client) LoadConfiguration() error {
	v := url.Values{}
	conf, err := client.API.GetConfiguration(v)
	if err != nil {
		return err
	}

	client.Mutex.Lock()
	client.Configuration = &conf
	client.Mutex.Unlock()

	return nil
}

// SetupConfigurationRefresh sets up a scheduler and will refresh
// the configuration from twitter every duration d.
// The reason for this is that these values can change over time.
// See https://dev.twitter.com/rest/reference/get/help/configuration
func (client *Client) SetupConfigurationRefresh(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			client.LoadConfiguration()
		}
	}()
	log.Printf("Twitter Configuration refresh: Enabled ✅  (every %s)\n", d.String())
}
