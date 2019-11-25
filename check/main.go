package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ngs/akerun-slack-notifier/akerun"
	"golang.org/x/oauth2"
)

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
	resp, err := client.Get("https://api.akerun.com/v3/organizations/" + id + "/accesses?sort_by=id&sort_order=asc&id_after=" + strconv.Itoa(lastID))
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
		err := notifyAccess(acc)
		if err != nil {
			return err
		}
	}
	payload, _ := json.Marshal(&AccessConfiguration{LastID: accessesResp.Accesses[len(accessesResp.Accesses)-1].ID})
	_, err = svc.PutObject(&s3.PutObjectInput{
		Key:    &key,
		Bucket: &bucket,
		Body:   bytes.NewReader(payload),
	})
	return err
}

func notifyAccess(access akerun.Access) error {
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	attachment := slack.Attachment{}
	deviceName := access.DeviceName
	akerunName := access.Akerun.Name
	var color string
	var text string
	var icon string
	var verb string
	if access.Action == "lock" {
		color = "good"
		icon = ":lock:"
		verb = "locked"
	} else if access.Action == "unlock" {
		color = "danger"
		icon = ":unlock:"
		verb = "unlocked"
	} else {
		return nil
	}
	attachment.Color = &color
	if access.DeviceType == "autolock" {
		text = icon + " *" + akerunName + "* was locked automatically"
	} else {
		if access.DeviceType == "nfc_outside" {
			deviceName = "Outside NFC"
		} else if access.DeviceType == "nfc_inside" {
			deviceName = "Inside NFC"
		}
		userName := "(unknown) :thinking_face:"
		if access.User != nil {
			userName = access.User.Name
		}
		text = icon + " *" + userName + "* " + verb + " *" + akerunName + "*"
		attachment.AddField(slack.Field{Title: "Device", Value: deviceName})
	}
	ts := access.AccessedAt.Unix()
	attachment.Timestamp = &ts
	payload := slack.Payload{
		Text:        text,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.Send(slackWebhookURL, "", payload)
	if len(err) > 0 {
		return err[0]
	}
	return nil
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
