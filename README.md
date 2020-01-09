## Gcloud function to link a Webhook from a source repo to an external conan build repo

<sub><sup>(This is derived from [Eli Bendersky's blog entry: GitHub webhook payload as a cloud function ](https://eli.thegreenplace.net/2019/github-webhook-payload-as-a-cloud-function/) )</sup></sub>

### Uploading this function to the Google Cloud

(You need to install the Google Cloud SDK.)
To upload the function do the following in the root of this BuilderFunction project.

```shell script
gcloud functions deploy builderserver --entry-point WebHookHandler --runtime go111 --trigger-http --set-env-vars=WEBHOOK_SECRET=<Secret for the webhook (is appended with the operation)>,BEARER_AUTH="<Bearer xxx string for Appveyor>",BLDRVNLW_ACCESS_TOKEN="<Access token for repo>"
```

Note - nospaces between the environment vars and all vars enclosed in quotes
The incoming github Webhook with json body is captured and the body forwarded to a conan-* 
build repository to trigger the conan based build.

### Routing webhooks from multiple repos
The basic Webhook text is suffixed on the source repo with a string to indicate which conan build is used

Webhook operation prefixes:

Webhook suffix | Conan build repo
--- | --- 
_new_hdsp_core  | github.com/bldrvnlw/conan-hdps-core
_ImageLoaderPlugin | github.com/bldrvnlw/conan-ImageLoaderPlugin
_ImageViewerPlugin | github.com/bldrvnlw/conan-ImageViewerPlugin

### Access keys
**WEBHOOK_SECRET**

For example if the secret Webhook  is "My$ecret#ebhook" to trigger the conan-hdps-core build enter 
**_My$ecret#ebhook_new_hdsp_core_** to trigger a conan-hdps-core build

**BLDRVNLW_ACCESS_TOKEN**

The function needs to push the build hook json to a repo using an access token. Currently only repos accessible 
with a token from bldrvnlw are supported bu this could be extended.

**BEARER_AUTH**

Unused - can be removed in the next version.

Get logs for the cloud function by
```shell script
gcloud functions logs read builderserver
```