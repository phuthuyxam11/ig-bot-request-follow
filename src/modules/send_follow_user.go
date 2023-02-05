package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"go-selenium.com/component/asyncjob"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type sendFollowResponse struct {
	Result              string `json:"result"`
	Message             string `json:"message"`
	Spam                bool   `json:"spam"`
	FeedbackTitle       string `json:"feedback_title"`
	FeedbackMessage     string `json:"feedback_message"`
	FeedbackURL         string `json:"feedback_url"`
	FeedbackAppealLabel string `json:"feedback_appeal_label"`
	FeedbackIgnoreLabel string `json:"feedback_ignore_label"`
	FeedbackAction      string `json:"feedback_action"`
	Status              string `json:"status"`
}

func SendFollowRequest(followerDetail FollowerDetail, authInfo AuthInfo) error {
	var jobs []asyncjob.Job
	for _, user := range followerDetail.Users {
		endPoint := fmt.Sprintf("https://www.instagram.com/api/v1/web/friendships/%s/follow/", user.Pk)
		jobInstance := asyncjob.NewJob(ExeSendFollow(endPoint, authInfo, user.Username))
		jobs = append(jobs, jobInstance)
	}
	group := asyncjob.NewGroup(true, jobs...)
	if err := group.Run(context.Background()); err != nil {
		return err
	}

	return nil
}

func ExeSendFollow(endPoint string, authInfo AuthInfo, userName string) asyncjob.JobHandler {
	return func(ctx context.Context) error {
		time.Sleep(time.Second * 30)
		err, resp := IgCurl(endPoint, authInfo, http.MethodPost)
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}

		sendFollowResponse := sendFollowResponse{}
		jsonErr := json.Unmarshal(body, &sendFollowResponse)
		if jsonErr != nil {
			return jsonErr
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)
		if err != nil {
			return err
		}
		if sendFollowResponse.Status == "ok" && sendFollowResponse.Result == "following" {
			fmt.Printf("🚀 Done: send request follow to: %s\n", userName)
		} else {
			fmt.Printf("⚠️ Fail: request follow to: %s something wrong message response is: %s\n", userName, fmt.Sprint(sendFollowResponse))
		}
		return nil
	}
}
