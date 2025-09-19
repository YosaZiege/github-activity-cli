package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// CreateEvent : Creation of a branch or Tag
// DeleteEvent : Deletion of a branch or a Tag
// ForkEvent : fork a repository
// GollumEvent : A wiki page is Created or Update
// IssueCommentEvent : Creation, Edtion , Deletion of an Issue
// IssuesEvent : Activity to an Issue Either opened , edited , closed , reopened , assigned , unassigned, labeled , unlabeled
// MemberEvent : Activity related to to repository Collaborators
// Public Event : When a repo is turned public
// ......
type Activity struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		DisplayLogin string `json:"display_login"`
	} `json:"actor"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		RefType string `json:"ref_type"`
		Size    int    `json:"size"`
	}
	CreatedAt time.Time `json:"created_at"`
}

// 1 : No action checking , 2 : action checking
const (
	Issue                    = "IssuesEvent"                   // 2
	IssueComment             = "IssueCommentEvent"             // 2
	Create                   = "CreateEvent"                   // 1
	Delete                   = "DeleteEvent"                   // 1
	CommitComment            = "CommitCommentEvent"            // 2
	Fork                     = "ForkEvent"                     // 1
	Wiki                     = "GollumEvent"                   // 1
	Member                   = "MemberEvent"                   // 2
	Public                   = "PublicEvent"                   // 1
	Push                     = "PushEvent"                     // 1
	PullRequest              = "PullRequestEvent"              // 2
	PullRequestReview        = "PullRequestReviewEvent"        // 2
	PullRequestReviewComment = "PullRequestReviewCommentEvent" // 2
	PullRequestReviewThread  = "PullRequestReviewThreadEvent"  // 2
	Star                     = "WatchEvent"                    // 1
)

var IssueActions = []string{"opened", "edited", "closed", "reopened", "assigned", "unassigned", "labeled", "unlabeled"}
var IssueCommentActions = []string{"created", "edited", "deleted"}
var CommitCommentActions = []string{"created"}
var MemberActions = []string{"added"}
var PullRequestActions = []string{"opened", "edited", "closed", "reopened", "assigned", "unassigned", "review_requested", "review_request_removed", "labeled", "unlabeled", "synchronize"}
var PullRequestReviewActions = []string{"created"}
var PullRequestReviewCommentActions = []string{"created"}
var PullRequestReviewThreadActions = []string{"resolved", "unresolved"}

func contains(slice []string, item string) bool {

	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false

}
func FetchActivity(user string) []byte {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", user)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error Calling github api %v \n", err)
		return nil
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error Calling github api %v\n", err)
		return nil
	}
	return body
}

func ActionPerformed(event string, user string, action string, repo string) {

	// EventTypeActions := map[string][]string{
	// 	Issue:                    IssueActions,
	// 	IssueComment:             IssueCommentActions,
	// 	CommitComment:            CommitCommentActions,
	// 	Member:                   MemberActions,
	// 	PullRequest:              PullRequestActions,
	// 	PullRequestReview:        PullRequestReviewActions,
	// 	PullRequestReviewComment: PullRequestReviewCommentActions,
	// 	PullRequestReviewThread:  PullRequestReviewThreadActions,
	// }
	DisplayType := map[string]string{
		Issue:                    "an issue",
		IssueComment:             "an issue comment",
		CommitComment:            "a commit comment",
		PullRequest:              "a pull request",
		PullRequestReview:        "a pull request review",
		PullRequestReviewComment: "a pull request review comment",
		PullRequestReviewThread:  "a pull request review thread",
	}

	fmt.Printf("%s has %s %s in %s\n", user, action, DisplayType[event], repo)
}
func HandleGithubData(body []byte, user string) {

	eventTypeCounts := map[string]int{
		Issue:                    0,
		IssueComment:             0,
		Create:                   0,
		Delete:                   0,
		CommitComment:            0,
		Fork:                     0,
		Wiki:                     0,
		Member:                   0,
		Public:                   0,
		Push:                     0,
		PullRequest:              0,
		PullRequestReview:        0,
		PullRequestReviewComment: 0,
		PullRequestReviewThread:  0,
		Star:                     0,
	}
	if body == nil {
		fmt.Println("No data received")
		return
	}
	var activ []Activity
	var relatedActivity []Activity
	now := time.Now()
	json.Unmarshal(body, &activ)

	// Filter Activites in the Time span We want
	for _, a := range activ {
		if now.Sub(a.CreatedAt) < 2*24*time.Hour {
			relatedActivity = append(relatedActivity, a)
		}
	}
	// Collect the Counts of the Different Activites
	for _, a := range relatedActivity {
		switch a.Type {
		case Issue:
			eventTypeCounts[Issue]++
		case Create:
			eventTypeCounts[Create]++
		case Delete:
			eventTypeCounts[Delete]++
		case CommitComment:
			eventTypeCounts[CommitComment]++
		case Fork:
			eventTypeCounts[Fork]++
		case Wiki:
			eventTypeCounts[Wiki]++
		case IssueComment:
			eventTypeCounts[IssueComment]++
		case Member:
			eventTypeCounts[Member]++
		case Public:
			eventTypeCounts[Public]++
		case Push:
			eventTypeCounts[Push]++
		case PullRequest:
			eventTypeCounts[PullRequest]++
		case PullRequestReview:
			eventTypeCounts[PullRequestReview]++
		case PullRequestReviewComment:
			eventTypeCounts[PullRequestReviewComment]++
		case PullRequestReviewThread:
			eventTypeCounts[PullRequestReviewThread]++
		default:
			fmt.Println("Unknown event type:", a.Type)
		}
	}
	fmt.Printf("Activity in the Last 48 Hours of the User : %s \n", user)
	for _, a := range relatedActivity {

		switch a.Type {
		case Issue, IssueComment, CommitComment, PullRequest,
			PullRequestReview, PullRequestReviewComment, PullRequestReviewThread:
			// Events with actions
			ActionPerformed(a.Type, user, a.Payload.Action, a.Repo.Name)
		case Member:
			fmt.Printf("%s has joined %s \n", user, a.Repo.Name)
		case Star:
			fmt.Printf("%s has Starred a %s \n", user, a.Repo.Name)
		case Create:
			fmt.Printf("%s has Created a new %s \n", user, a.Payload.RefType)
		case Delete:
			fmt.Printf("%s has Deleted a %s \n", user, a.Payload.RefType)
		case Fork:
			fmt.Printf("%s has Forked %s \n", user, a.Repo.Name)
		case Wiki:
			fmt.Printf("%s has Updated/Created a wiki \n", user)
		case Public:
			fmt.Printf("%s has made the Repo : %s Public \n", user, a.Repo.Name)
		case Push:
			// Events without specific actions
			fmt.Printf("%s Pushed %d commits to %s\n", user, a.Payload.Size, a.Repo.Name)

		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	user := os.Args[1]
	body := FetchActivity(user)
	HandleGithubData(body, user)
}
