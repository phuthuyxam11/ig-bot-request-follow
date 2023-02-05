package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-selenium.com/component/asyncjob"
	"go-selenium.com/component/db"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Follower struct {
	Id       string
	UserName string
}

type PublicPage struct {
	PageName string
}

type UserInfo struct {
	Data struct {
		User struct {
			ID                string `json:"id"`
			IsBusinessAccount bool   `json:"is_business_account"`
			IsPrivate         bool   `json:"is_private"`
			IsVerified        bool   `json:"is_verified"`
			Username          string `json:"username"`
		} `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}

type RequestGetFollower struct {
	Count string
	MaxId string
}

type FollowerDetail struct {
	Users []struct {
		HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture"`
		Pk                         string        `json:"pk"`
		PkID                       string        `json:"pk_id"`
		Username                   string        `json:"username"`
		FullName                   string        `json:"full_name"`
		IsPrivate                  bool          `json:"is_private"`
		IsVerified                 bool          `json:"is_verified"`
		ProfilePicID               string        `json:"profile_pic_id"`
		ProfilePicURL              string        `json:"profile_pic_url"`
		AccountBadges              []interface{} `json:"account_badges"`
		LatestReelMedia            int           `json:"latest_reel_media"`
		LinkedFbInfo               struct {
			LinkedFbUser struct {
				ID                    string      `json:"id"`
				Name                  string      `json:"name"`
				IsValid               bool        `json:"is_valid"`
				FbAccountCreationTime interface{} `json:"fb_account_creation_time"`
				LinkTime              interface{} `json:"link_time"`
			} `json:"linked_fb_user"`
		} `json:"linked_fb_info,omitempty"`
	} `json:"users"`
	BigList                    bool   `json:"big_list"`
	PageSize                   int    `json:"page_size"`
	NextMaxID                  string `json:"next_max_id"`
	HasMore                    bool   `json:"has_more"`
	ShouldLimitListOfFollowers bool   `json:"should_limit_list_of_followers"`
	Status                     string `json:"status"`
}

func Execute(listPage []PublicPage, info AuthInfo, store *sql.DB) error {
	// get page input info
	var pages []UserInfo
	for _, user := range listPage {
		err, userDetail := getPageIdByPageName(user.PageName, info)
		if err != nil {
			continue
		}
		if !userDetail.Data.User.IsPrivate {
			pages = append(pages, userDetail)
		}
	}

	var jobs []asyncjob.Job
	for _, page := range pages {
		jobInstance := asyncjob.NewJob(GetFollowers(page, info, store))
		jobs = append(jobs, jobInstance)
	}
	// get all user follow publish page
	group := asyncjob.NewGroup(true, jobs...)

	if err := group.Run(context.Background()); err != nil {
		return err
	}

	return nil
}

func getPageIdByPageName(pageName string, info AuthInfo) (error, UserInfo) {
	pageInfoEndpoint := fmt.Sprintf("https://i.instagram.com/api/v1/users/web_profile_info/?username=%s", pageName)

	err, resp := IgCurl(pageInfoEndpoint, info, http.MethodGet)
	if err != nil {
		return err, UserInfo{}
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr, UserInfo{}
	}

	userInfo := UserInfo{}
	jsonErr := json.Unmarshal(body, &userInfo)
	if jsonErr != nil {
		return jsonErr, UserInfo{}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return nil, userInfo
}

func GetFollowers(page UserInfo, authInfo AuthInfo, store *sql.DB) asyncjob.JobHandler {
	return func(ctx context.Context) error {
		fmt.Printf("Start get follower with page: %s\n", page.Data.User.Username)
		requestGetFollower := RequestGetFollower{
			Count: "200",
			MaxId: "",
		}

		err := ExeGetFollowers(requestGetFollower, authInfo, page.Data.User.ID, store)
		if err != nil {
			return err
		}
		return nil
	}
}

func ExeGetFollowers(request RequestGetFollower, authInfo AuthInfo, pageId string, store *sql.DB) error {
	endpoint := fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/followers/?count=%s&max_id=%s", pageId, request.Count, request.MaxId)
	err, resp := IgCurl(endpoint, authInfo, http.MethodGet)
	if err != nil {
		return err
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	follower := FollowerDetail{}
	jsonErr := json.Unmarshal(body, &follower)
	if jsonErr != nil {
		return jsonErr
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// exe follow action
	if len(follower.Users) > 0 {
		var userModelMap []db.UserModel
		for _, user := range follower.Users {
			userModelMap = append(userModelMap, db.UserModel{
				PkId:       user.PkID,
				UserName:   user.Username,
				IsPrivate:  checkBitSetVar(user.IsPrivate),
				FollowFlag: "no",
				CreatedAt:  time.Now().String(),
				ModifiedAt: time.Now().String(),
			})
		}
		err = db.SaveUser(store, userModelMap)
		if err != nil {
			return err
		}
	}

	if follower.NextMaxID == "" {
		fmt.Println("end next id")
		return nil
	} else {
		fmt.Printf("Next id is: %s\n", request.MaxId)
		newRequest := RequestGetFollower{
			Count: "200",
			MaxId: follower.NextMaxID,
		}
		err = ExeGetFollowers(newRequest, authInfo, pageId, store)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkBitSetVar(mybool bool) int {
	if mybool {
		return 1
	}
	return 0
}
