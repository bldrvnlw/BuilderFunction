package builderfunction

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

var testJson = `{"commits":[{"id":"375dd8e2c924110e621c95f65138fd9065e2d428","tree_id":"5e03f4b5167a113689c0a57ac08a7c818a9d2076","distinct":true,"message":"Made log view sortable","timestamp":"2019-09-03T20:13:02+02:00","url":"https://github.com/hdps/core/commit/375dd8e2c924110e621c95f65138fd9065e2d428","author":{"name":"Niels Dekker","email":"N.Dekker@lumc.nl","username":"N-Dekker"},"committer":{"name":"Niels Dekker","email":"N.Dekker@lumc.nl","username":"N-Dekker"},"added":[],"removed":[],"modified":["HDPS/src/LogItemModel.cpp","HDPS/src/LogItemModel.h","HDPS/src/Logger.cpp","HDPS/src/Logger.h"]}],"head_commit":{"id":"375dd8e2c924110e621c95f65138fd9065e2d428","tree_id":"5e03f4b5167a113689c0a57ac08a7c818a9d2076","distinct":true,"message":"Made log view sortable","timestamp":"2019-09-03T20:13:02+02:00","url":"https://github.com/hdps/core/commit/375dd8e2c924110e621c95f65138fd9065e2d428","author":{"name":"Niels Dekker","email":"N.Dekker@lumc.nl","username":"N-Dekker"},"committer":{"name":"Niels Dekker","email":"N.Dekker@lumc.nl","username":"N-Dekker"},"added":[],"removed":[],"modified":["HDPS/src/LogItemModel.cpp","HDPS/src/LogItemModel.h","HDPS/src/Logger.cpp","HDPS/src/Logger.h"]},"ref":"refs/heads/MakeLogViewSortable","before":"0000000000000000000000000000000000000000","after":"375dd8e2c924110e621c95f65138fd9065e2d428","created":true,"deleted":false,"forced":false,"base_ref":null,"compare":"https://github.com/hdps/core/commit/375dd8e2c924","repository":{"id":197026558,"node_id":"MDEwOlJlcG9zaXRvcnkxOTcwMjY1NTg=","name":"core","full_name":"hdps/core","private":true,"owner":{"name":"hdps","email":null,"login":"hdps","id":52745266,"node_id":"MDEyOk9yZ2FuaXphdGlvbjUyNzQ1MjY2","avatar_url":"https://avatars0.githubusercontent.com/u/52745266?v=4","gravatar_id":"","url":"https://api.github.com/users/hdps","html_url":"https://github.com/hdps","followers_url":"https://api.github.com/users/hdps/followers","following_url":"https://api.github.com/users/hdps/following{/other_user}","gists_url":"https://api.github.com/users/hdps/gists{/gist_id}","starred_url":"https://api.github.com/users/hdps/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/hdps/subscriptions","organizations_url":"https://api.github.com/users/hdps/orgs","repos_url":"https://api.github.com/users/hdps/repos","events_url":"https://api.github.com/users/hdps/events{/privacy}","received_events_url":"https://api.github.com/users/hdps/received_events","type":"Organization","site_admin":false},"html_url":"https://github.com/hdps/core","description":null,"fork":false,"url":"https://github.com/hdps/core","forks_url":"https://api.github.com/repos/hdps/core/forks","keys_url":"https://api.github.com/repos/hdps/core/keys{/key_id}","collaborators_url":"https://api.github.com/repos/hdps/core/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/hdps/core/teams","hooks_url":"https://api.github.com/repos/hdps/core/hooks","issue_events_url":"https://api.github.com/repos/hdps/core/issues/events{/number}","events_url":"https://api.github.com/repos/hdps/core/events","assignees_url":"https://api.github.com/repos/hdps/core/assignees{/user}","branches_url":"https://api.github.com/repos/hdps/core/branches{/branch}","tags_url":"https://api.github.com/repos/hdps/core/tags","blobs_url":"https://api.github.com/repos/hdps/core/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/hdps/core/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/hdps/core/git/refs{/sha}","trees_url":"https://api.github.com/repos/hdps/core/git/trees{/sha}","statuses_url":"https://api.github.com/repos/hdps/core/statuses/{sha}","languages_url":"https://api.github.com/repos/hdps/core/languages","stargazers_url":"https://api.github.com/repos/hdps/core/stargazers","contributors_url":"https://api.github.com/repos/hdps/core/contributors","subscribers_url":"https://api.github.com/repos/hdps/core/subscribers","subscription_url":"https://api.github.com/repos/hdps/core/subscription","commits_url":"https://api.github.com/repos/hdps/core/commits{/sha}","git_commits_url":"https://api.github.com/repos/hdps/core/git/commits{/sha}","comments_url":"https://api.github.com/repos/hdps/core/comments{/number}","issue_comment_url":"https://api.github.com/repos/hdps/core/issues/comments{/number}","contents_url":"https://api.github.com/repos/hdps/core/contents/{+path}","compare_url":"https://api.github.com/repos/hdps/core/compare/{base}...{head}","merges_url":"https://api.github.com/repos/hdps/core/merges","archive_url":"https://api.github.com/repos/hdps/core/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/hdps/core/downloads","issues_url":"https://api.github.com/repos/hdps/core/issues{/number}","pulls_url":"https://api.github.com/repos/hdps/core/pulls{/number}","milestones_url":"https://api.github.com/repos/hdps/core/milestones{/number}","notifications_url":"https://api.github.com/repos/hdps/core/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/hdps/core/labels{/name}","releases_url":"https://api.github.com/repos/hdps/core/releases{/id}","deployments_url":"https://api.github.com/repos/hdps/core/deployments","created_at":1563204934,"updated_at":"2019-09-03T14:05:26Z","pushed_at":1567534392,"git_url":"git://github.com/hdps/core.git","ssh_url":"git@github.com:hdps/core.git","clone_url":"https://github.com/hdps/core.git","svn_url":"https://github.com/hdps/core","homepage":null,"size":27286,"stargazers_count":0,"watchers_count":0,"language":"C++","has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":1,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":3,"license":null,"forks":1,"open_issues":3,"watchers":0,"default_branch":"develop","stargazers":0,"master_branch":"develop","organization":"hdps"},"pusher":{"name":"N-Dekker","email":"N.Dekker@lumc.nl"},"organization":{"login":"hdps","id":52745266,"node_id":"MDEyOk9yZ2FuaXphdGlvbjUyNzQ1MjY2","url":"https://api.github.com/orgs/hdps","repos_url":"https://api.github.com/orgs/hdps/repos","events_url":"https://api.github.com/orgs/hdps/events","hooks_url":"https://api.github.com/orgs/hdps/hooks","issues_url":"https://api.github.com/orgs/hdps/issues","members_url":"https://api.github.com/orgs/hdps/members{/member}","public_members_url":"https://api.github.com/orgs/hdps/public_members{/member}","avatar_url":"https://avatars0.githubusercontent.com/u/52745266?v=4","description":""},"sender":{"login":"N-Dekker","id":27005366,"node_id":"MDQ6VXNlcjI3MDA1MzY2","avatar_url":"https://avatars0.githubusercontent.com/u/27005366?v=4","gravatar_id":"","url":"https://api.github.com/users/N-Dekker","html_url":"https://github.com/N-Dekker","followers_url":"https://api.github.com/users/N-Dekker/followers","following_url":"https://api.github.com/users/N-Dekker/following{/other_user}","gists_url":"https://api.github.com/users/N-Dekker/gists{/gist_id}","starred_url":"https://api.github.com/users/N-Dekker/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/N-Dekker/subscriptions","organizations_url":"https://api.github.com/users/N-Dekker/orgs","repos_url":"https://api.github.com/users/N-Dekker/repos","events_url":"https://api.github.com/users/N-Dekker/events{/privacy}","received_events_url":"https://api.github.com/users/N-Dekker/received_events","type":"User","site_admin":false}}`

func generateSignature(secret string, data []byte) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(data)
	result := mac.Sum(nil)
	return "sha1=" + hex.EncodeToString(result)
}

func TestWebHookHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(testJson))
	req.Header.Add("Content-Type", "application/json")
	signature := generateSignature(webhookSecret + "_test", []byte(testJson))
	req.Header.Add("X-Hub-Signature", signature)
	rr := httptest.NewRecorder()
	WebHookHandler(rr, req)

	_, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	//if !strings.Contains(string(out), "Signature validation failed") {
	//	t.Errorf("Got %s, expected to see 'reply num'", string(out))
	//}
}

func TestServerJsonParse(t *testing.T) {
	head, committer := ParseCommitJson(testJson)
	fmt.Println(head["message"])
	fmt.Println(committer["name"])
}