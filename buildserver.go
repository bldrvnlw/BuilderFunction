package builderfunction

// Could extend this to support a map of GitHub webhooks to build urls on Appveyor, Travis Azure etc.
import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var authStr = os.Getenv("BEARER_AUTH")
var webhookSecret = os.Getenv("WEBHOOK_SECRET")

type ForwardTarget struct {
	label string
	url string
	owner string
	repo string
	path string
}

var nullTarget = ForwardTarget{
	"", "", "", "", "",
}

// Webhook strings are mapped to a list of one or more http actions that
// will be executed when the webhok is recieved
var webhookActionMap = map[string]ForwardTarget{
	webhookSecret + "_new_hdsp_core": ForwardTarget{
	"Run new build",
	"https://api.github.com/",
	"bldrvnlw",
	"conan-hdps-core",
	"build_trigger.json",
	},
	webhookSecret + "_ImageLoaderPlugin": ForwardTarget{
		"Run new build",
		"https://api.github.com/",
		"bldrvnlw",
		"conan-ImageLoaderPlugin",
		"build_trigger.json",
	},
	webhookSecret + "_ImageViewerPlugin": ForwardTarget{
		"Run new build",
		"https://api.github.com/",
		"bldrvnlw",
		"conan-ImageViewerPlugin",
		"build_trigger.json",
	},
	webhookSecret + "_test": ForwardTarget{
		"***TEST***",
		"***TEST***",
		"bldrvnlw",
		"conan-hdps-core",
		"build_trigger.json",
	},
}

func ParseCommitJson(newCommitJson string) (map[string]interface{}, map[string]interface{}) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(newCommitJson), &jsonMap); err != nil {
		log.Fatal(err)
	}
	headCommit := jsonMap["head_commit"].(map[string]interface{})
	author :=  headCommit["author"].(map[string]interface{})
	fmt.Println(headCommit["message"])
	fmt.Println(author["name"])
	return headCommit, author
}

func ForwardWebhookJsonContext(newCommitJson string, target ForwardTarget) {
	head, author := ParseCommitJson(newCommitJson)

	commitMessage := fmt.Sprintf("%v", head["message"])
	authorName := fmt.Sprintf("%v", author["name"])
	authorEmail := fmt.Sprintf("%v", author["email"])

	fmt.Printf("%s %s %s", commitMessage, authorName, authorEmail)
	if target.url == "***TEST***" {
		return
	}
	ctx := context.Background()
	token := os.Getenv("BLDRVNLW_ACCESS_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	url, _ := url.Parse("https://api.github.com/")
	client.BaseURL = url
	client.UploadURL = url
	// Get the current file blob SHA - needed to update it
	fileContents, _, _, err := client.Repositories.GetContents(context.Background(), target.owner, target.repo, target.path, &github.RepositoryContentGetOptions{})
	_, _, err = client.Repositories.UpdateFile(context.Background(), target.owner, target.repo, target.path, &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: []byte(newCommitJson),
		SHA: fileContents.SHA,
		Branch: github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String(authorName), Email: github.String(authorEmail)},
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("New commit json submitted")
}

// Thank you Eli Bendersky! https://eli.thegreenplace.net/2019/github-webhook-payload-as-a-cloud-function/

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

// validateSignature validates the MAC signature in the request header using
// the secret key and the message body. Returns true iff valid.
func validateSignature(body []byte, r *http.Request) ForwardTarget {
	sigHeader := r.Header["X-Hub-Signature"]
	if len(sigHeader) < 1 {
		log.Println("signature header too short")
		return nullTarget
	}
	parts := strings.Split(sigHeader[0], "=")
	if len(parts) != 2 || parts[0] != "sha1" {
		log.Println("Expected signature header 'sha1=XXXXX'")
		return nullTarget
	}
	decoded, err := hex.DecodeString(parts[1])
	if err != nil {
		log.Println(err)
		return nullTarget
	}

	for k := range webhookActionMap {
		log.Printf("CheckMAC on key: %s\n", k)
		if checkMAC(body, decoded, []byte(k)) {
			return webhookActionMap[k]
		}
	}
	return nullTarget
}

// based on the incoming
func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	s := string(body)
	log.Printf("Length webhook body string %d \n", len(s))
	// log.Println(s)

	var target = validateSignature(body, r)
	ForwardWebhookJsonContext(s, target)
}
