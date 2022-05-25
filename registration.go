package main

import (
	"VRChat_Account_Generator/Shared"
	"encoding/json"
	"fmt"
	"github.com/thanhpk/randstr"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RegistrationPayload is the payload for the registration API
type RegistrationPayload struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	Year               string `json:"year"`
	Month              string `json:"month"`
	Day                string `json:"day"`
	CaptchaCode        string `json:"captchaCode"`
	Subscribe          bool   `json:"subscribe"`
	AcceptedTOSVersion int    `json:"acceptedTOSVersion"`
}

//RegisterVRC takes a 2Captcha Key and a Worker id, it registers an account on VRChat with randomly generated data
//Returns http.Response, username, password, email
func RegisterVRC(worker int) (*http.Response, string, string, string) {

	var Success bool
	for Success == false {

		//Generate Random Year in the range of 1980-2000
		var year = strconv.Itoa(randInt(1980, 2000))

		// Generate Random Month int between 1 and 12
		month := strconv.Itoa(randInt(1, 12))

		// Generate Random Day int between 1 and 31
		day := strconv.Itoa(randInt(1, 31))

		// Generate Random Username
		username := randstr.String(15)

		// Generate Random Password
		password := randstr.String(15)

		// Generate Random Email
		email := randstr.String(15) + "@" + "gmail" + ".com"

		code := SolveHCaptcha(worker)

		// Fill in Registration Payload
		Payload := RegistrationPayload{
			Username:           username,
			Password:           password,
			Email:              email,
			Year:               year,
			Month:              month,
			Day:                day,
			CaptchaCode:        code,
			Subscribe:          false,
			AcceptedTOSVersion: 7,
		}

		// Convert Payload to JSON
		jsonPayload, err := json.Marshal(Payload)

		//Setup our client and proxy transport
		transport, proxy, err := Shared.PManager.GetRandomProxyTransport()
		if err != nil {
			if proxy != nil {
				proxy.InUse = false
			}
			continue
		}
		client := &http.Client{Timeout: 10 * time.Second, Transport: transport}

		// Setup our request
		req, err := http.NewRequest("POST", "https://vrchat.com/api/1/auth/register?apiKey=JlE5Jldo5Jibnk5O5hTx6XVqsJu4WJ26", strings.NewReader(string(jsonPayload)))
		if err != nil {
			if proxy != nil {
				proxy.InUse = false
			}
			continue
		}

		if err != nil {
			fmt.Println(err)

		}

		// Set our headers
		req.Header.Add("Host", " vrchat.com")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
		req.Header.Add("Origin", "https://vrchat.com")
		req.Header.Add("Referer", "https://vrchat.com/home/register")

		//Register Account
		res, err := client.Do(req)
		if err != nil {
			if proxy != nil {
				proxy.InUse = false
			}
			continue
		}

		return res, username, password, email

	}
	return nil, "", "", ""
}

type RegistrationResponse struct {
	ID                             string        `json:"id"`
	Username                       string        `json:"username"`
	DisplayName                    string        `json:"displayName"`
	UserIcon                       string        `json:"userIcon"`
	Bio                            string        `json:"bio"`
	BioLinks                       []interface{} `json:"bioLinks"`
	ProfilePicOverride             string        `json:"profilePicOverride"`
	StatusDescription              string        `json:"statusDescription"`
	PastDisplayNames               []interface{} `json:"pastDisplayNames"`
	HasEmail                       bool          `json:"hasEmail"`
	HasPendingEmail                bool          `json:"hasPendingEmail"`
	ObfuscatedEmail                string        `json:"obfuscatedEmail"`
	ObfuscatedPendingEmail         string        `json:"obfuscatedPendingEmail"`
	EmailVerified                  bool          `json:"emailVerified"`
	HasBirthday                    bool          `json:"hasBirthday"`
	Unsubscribe                    bool          `json:"unsubscribe"`
	StatusHistory                  []string      `json:"statusHistory"`
	StatusFirstTime                bool          `json:"statusFirstTime"`
	Friends                        []interface{} `json:"friends"`
	FriendGroupNames               []interface{} `json:"friendGroupNames"`
	CurrentAvatarImageURL          string        `json:"currentAvatarImageUrl"`
	CurrentAvatarThumbnailImageURL string        `json:"currentAvatarThumbnailImageUrl"`
	CurrentAvatar                  string        `json:"currentAvatar"`
	CurrentAvatarAssetURL          string        `json:"currentAvatarAssetUrl"`
	FallbackAvatar                 string        `json:"fallbackAvatar"`
	AccountDeletionDate            interface{}   `json:"accountDeletionDate"`
	AcceptedTOSVersion             int           `json:"acceptedTOSVersion"`
	SteamID                        string        `json:"steamId"`
	SteamDetails                   struct {
	} `json:"steamDetails"`
	OculusID                 string        `json:"oculusId"`
	HasLoggedInFromClient    bool          `json:"hasLoggedInFromClient"`
	HomeLocation             string        `json:"homeLocation"`
	TwoFactorAuthEnabled     bool          `json:"twoFactorAuthEnabled"`
	TwoFactorAuthEnabledDate interface{}   `json:"twoFactorAuthEnabledDate"`
	State                    string        `json:"state"`
	Tags                     []interface{} `json:"tags"`
	DeveloperType            string        `json:"developerType"`
	LastLogin                time.Time     `json:"last_login"`
	LastPlatform             string        `json:"last_platform"`
	AllowAvatarCopying       bool          `json:"allowAvatarCopying"`
	Status                   string        `json:"status"`
	DateJoined               string        `json:"date_joined"`
	IsFriend                 bool          `json:"isFriend"`
	FriendKey                string        `json:"friendKey"`
	LastActivity             time.Time     `json:"last_activity"`
	AuthToken                string        `json:"authToken"`
}
