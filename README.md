### Gcloud function to link a Webhook from a source repo to an external conan build repo

<sub><sup>(This is derived from [Eli Bendersky's blog entry: GitHub webhook payload as a cloud function ](https://eli.thegreenplace.net/2019/github-webhook-payload-as-a-cloud-function/) )</sup></sub>

To upload the function do the following


```shell script
gcloud functions deploy builderfunction --entry-point WebHookHandler --runtime go111 --trigger-http --set-env-vars WEBHOOK_SECRET=<Secret for the webhook (is appended with the operation)>,BEARER_AUTH="<Bearer xxx string for Appveyor>",BLDRVNLW_ACCESS_TOKEN="<Access token for repo>"
```

The incoming github Webhook with json body is captured and the body forwarded to a conan-* 
build repository to trigger the conan based build.

The basic Webhook text is suffixed on the source repo with a string to indicate which conan build is used

Webhook operation prefixes:

Webhook suffix | Conan build repo
--- | --- 
_new_hdsp_core  | github.com/bldrvnlw/conan-hdps-core
_ImageLoaderPlugin | github.com/bldrvnlw/conan-ImageLoaderPlugin
_ImageViewerPlugin | github.com/bldrvnlw/conan-ImageViewerPlugin

For example if the secret Webhook  is "My$ecret#ebhook" to trigger the conan-hdps-core build enter 
"My$ecret#ebhook_new_hdsp_core" to trigger a conan-hdps-core build

Get logs for the cloud function by

```shell script
gcloud functions logs read builderfunction
```