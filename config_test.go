package main

import (
	"reflect"
	"testing"
)

func TestConfig_NewConfiguration_WithWrongFilenames(t *testing.T) {
	mockData := []struct {
		Testcase string
		FileName string
	}{
		{"Empty filename", ""},
		{"Not existing file", "./not-existing-file"},
		{"Existing file, but wrong format", "./tests/config.xml"},
	}

	for _, data := range mockData {
		c, err := NewConfiguration(data.FileName)
		if c != nil {
			t.Errorf("Testcase %s: Configuration is not nil -> %+v", data.Testcase, c)
		}
		if err == nil {
			t.Errorf("Testcase %s: Err is nil, but expected an error", data.Testcase)
		}
	}

}

func TestConfig_NewConfiguration_WithExistingFile(t *testing.T) {
	configMock := &Configuration{
		Twitter: TwitterConfiguration{
			ConsumerKey:       "Consumer-Key-ABC",
			ConsumerSecret:    "Consumer-Secret-DEF",
			AccessToken:       "Access-Token-Foo",
			AccessTokenSecret: "Access-Token-Bar",
		},
		Redis: RedisConfiguration{
			URL:  "127.0.0.1:6379",
			Auth: "My-Secret-Password",
		},
	}

	c, err := NewConfiguration("./tests/config.json")
	if err != nil {
		t.Errorf("Err is not nil")
	}

	if reflect.DeepEqual(c, configMock) == false {
		t.Errorf("Configuration file is not the same as mock: %+v && %+v", c, configMock)
	}
}
