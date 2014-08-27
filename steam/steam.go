package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ID_CS_GO          = 710
	ID_CS_SOURCE      = 260
	ID_DOTA           = 570
	ID_DOTA_BETA_TEST = 205790
	ID_DOTA_TEST      = 816
	ID_PORTAL_2       = 620
	ID_PORTAL_2_BETA  = 841
	ID_TF2            = 440
	ID_TF2_BETA       = 520
	app_cache         []App
)

type Client struct {
	key string
}

type App struct {
	AppId int
	Name  string
}

type appListResponse struct {
	Applist struct {
		Apps []App `json:"apps"`
	} `json:"applist"`
}

func NewClient(key string) *Client {
	return &Client{key: key}
}

func (c *Client) GetAppList(skipCache bool) ([]App, error) {
    if !skipCache && app_cache != nil {
        return app_cache, nil
    }
	url := "http://api.steampowered.com/ISteamApps/GetAppList/v2"
	res, err := http.Get(url + "?key=" + c.key)
	if err != nil {
		return nil, fmt.Errorf("error fetching in GetAppList: %v", err)
	}
	defer res.Body.Close()

	var parsed appListResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error parsing GetAppList response: %v", err)
	}
    app_cache = parsed.Applist.Apps
	return parsed.Applist.Apps, nil
}

type appNewsResponse struct {
	AppNews struct {
		AppId     int        `json:"appid"`
		NewsItems []NewsItem `json:"newsitems"`
	} `json:"appnews"`
}

type NewsItem struct {
	Gid        string `json:"gid"`
	Title      string `json:"title"`
	Url        string `json:"url"`
	IsExternal bool   `json:"is_external_url"`
	Author     string `json:"author"`
	Contents   string `json:"contents"`
	FeedLabel  string `json:"feedlabel"`
	Date       uint64 `json:"date"`
	Feedname   string `json:"feedname"`
}

func (c *Client) GetNewsForApp(id int) ([]NewsItem, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamNews/GetNewsForApp/v0002?key=%s&appid=%d",
		c.key, id)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching in GetNewsForApp: %v", err)
	}
	defer res.Body.Close()

	var parsed appNewsResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error parsing response in GetNewsForApp: %v", err)
	}
	return parsed.AppNews.NewsItems, nil
}

type friendsListResponse struct {
	FriendsList struct {
		Friends []Friend `json:"friends"`
	} `json:"friendslist"`
}

type Friend struct {
	SteamId      string `json:"steamid"`
	Relationship string `json:"relationship"`
	Since        uint64 `json:"friend_since"`
}

func (c *Client) GetFriendList(steamid string, roles string) ([]Friend, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetFriendList/v1?key=%s&steamid=%s&relationship=%s",
		c.key, steamid, roles)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching GetFriendList: %v", err)
	}
	defer res.Body.Close()

	var parsed friendsListResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error parsing GetFriendList response: %v", err)
	}
	return parsed.FriendsList.Friends, nil
}

func (c *Client) GetGlobalStatsForGame(appid int) (*string, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetGlobalStatsForGame/v0001?key=%s&appid=%d&count=1", c.key, appid)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching GetGlobalStatsForGame: %v", err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	return nil, nil
}

func (c *Client) GetSchemaForGame(appid int) (*string, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v2?key=%s&appid=%d", c.key, appid)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching GetSchemaForGame: %v", err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	return nil, nil
}

func (c *Client) GetUserStatsForGame(userid uint64, appid int) (*string, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v2?key=%s&steamid=%d&appid=%d", c.key, userid, appid)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching GetSchemaForGame: %v", err)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	return nil, nil
}
