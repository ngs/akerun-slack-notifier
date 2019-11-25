package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ngs/akerun-slack-notifier/akerun"
	"golang.org/x/oauth2"
)

// BlockSection .
type BlockSection struct {
	Type string    `json:"type"`
	Text BlockText `json:"text"`
}

// BlockText .
type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// BlocksInput .
type BlocksInput struct {
	Blocks []BlockSection `json:"blocks"`
}

// AccessConfiguration .
type AccessConfiguration struct {
	LastID int `json:"lastID"`
}

func checkOrganizationAccesses(client *http.Client, id string) error {
	bucket := os.Getenv("S3_BUCKET")
	key := "organization-" + id + ".json"
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	lastID := 0
	res, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil {
		var cfg AccessConfiguration
		err = json.NewDecoder(res.Body).Decode(&cfg)
		if err == nil {
			lastID = cfg.LastID
		}
	}
	resp, err := client.Get("https://api.akerun.com/v3/organizations/" + id + "/accesses?sort_by=id&sort_order=desc&id_after=" + strconv.Itoa(lastID))
	if err != nil {
		return err
	}
	var accessesResp *akerun.AccessesResponse
	err = json.NewDecoder(resp.Body).Decode(&accessesResp)
	if err != nil {
		return err
	}
	if len(accessesResp.Accesses) == 0 {
		return nil
	}
	for _, acc := range accessesResp.Accesses {
		text := accessText(acc)
		if text != "" {
			payload, _ := json.Marshal(&map[string]string{"text": text})
			slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
			resp, err = http.DefaultClient.Post(slackWebhookURL, "application/json", bytes.NewReader(payload))
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				return fmt.Errorf("%v", buf.String())
			}
		}
	}
	payload, _ := json.Marshal(&AccessConfiguration{LastID: accessesResp.Accesses[0].ID})
	_, err = svc.PutObject(&s3.PutObjectInput{
		Key:    &key,
		Bucket: &bucket,
		Body:   bytes.NewReader(payload),
	})
	return err
}

func accessText(access akerun.Access) string {
	deviceName := access.DeviceName
	akerunName := access.Akerun.Name
	icon := ""
	verb := ""
	if access.Action == "lock" {
		icon = ":lock:"
		verb = "locked"
	} else if access.Action == "unlock" {
		icon = ":unlock:"
		verb = "unlocked"
	} else {
		return ""
	}
	if access.DeviceType == "nfc_outside" {
		deviceName = "Outside NFC"
	} else if access.DeviceType == "nfc_inside" {
		deviceName = "Inside NFC"
	} else if access.DeviceType == "autolock" {
		return icon + " *" + akerunName + "* was locked automatically"
	}
	userName := "(unknown) :thinking_face:"
	if access.User != nil {
		userName = access.User.Name
	}
	return icon + " *" + userName + "* " + verb + " *" + akerunName + "* using " + deviceName
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) error {
	bucket := os.Getenv("S3_BUCKET")
	key := "token.json"
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	res, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var token *oauth2.Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return err
	}
	var orgResp *akerun.OrganizationsResponse
	client := akerun.Config.Client(ctx, token)
	resp, err := client.Get("https://api.akerun.com/v3/organizations")
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&orgResp)
	if err != nil {
		return err
	}
	for _, org := range orgResp.Organizations {
		id := org.ID
		err = checkOrganizationAccesses(client, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
