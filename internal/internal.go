package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"path/filepath"
)

type Config struct {
	TokenEndpoint     string `json:"token_endpoint"`
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	Scope             string `json:"token_scope"`
	MFAEndpoint       string `json:"mfa_list_endpoint"`
	MFAAcceptEndpoint string `json:"mfa_accept_endpoint"`
	LoginPath         string `json:"login_path"`
	ConfigPath        string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type MFAKey struct {
	User          string `json:"user"`
	AuthorizedKey string `json:"authorized_key"`
}

type Token string

func LoadCfg(path string) (cfg Config, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return cfg, err
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	cfg.ConfigPath = path
	return
}

func GetToken(cfg Config) (token Token, err error) {
	resp, e := http.PostForm(cfg.TokenEndpoint,
		url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {cfg.ClientID},
			"client_secret": {cfg.ClientSecret},
			"scope":         {cfg.Scope}})
	if e != nil {
		err = e
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Failed to get token")
		return
	}
	defer resp.Body.Close()
	var tokenData tokenResponse
	if e := json.NewDecoder(resp.Body).Decode(&tokenData); e != nil {
		err = e
		return
	}
	token = Token(tokenData.AccessToken)
	return
}

func GetMFAKeys(cfg Config, token Token) (keys []MFAKey, err error) {
	client := http.Client{}
	req, e := http.NewRequest(http.MethodGet, cfg.MFAEndpoint, nil)
	if e != nil {
		err = e
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, e := client.Do(req)
	if e != nil {
		err = e
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		err = errors.New("Failed to get MFA keys")
		return
	}
	defer resp.Body.Close()
	if e := json.NewDecoder(resp.Body).Decode(&keys); e != nil {
		err = e
		return
	}
	return
}

func GetKeyLines(cfg Config, keys []MFAKey) (result []string) {
	for _, key := range keys {
		result = append(result, fmt.Sprintf("command=\"%s --config %s --user %s\" %s\n", cfg.LoginPath, cfg.ConfigPath, key.User, key.AuthorizedKey))
	}
	return
}

func MFAAccept(cfg Config, token Token, user string) error {
	body := fmt.Sprintf("{\"user\": \"%s\"}", user)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, cfg.MFAAcceptEndpoint, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		err = errors.New("Failed to accept challenge")
		return err
	}
	return nil
}
