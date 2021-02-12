// +build server

package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"net/http"
)

const page="<!doctype html>\n" +
	"<html lang=\"en\">\n<head>\n" +
	"    <meta charset=\"utf-8\">\n" +
	"    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\">\n" +
	"    <link rel=\"stylesheet\" href=\"external/bootstrap.min.css\">\n" +
	"    <link href=\"external/dashboard.css\" rel=\"stylesheet\">\n" +
	"    <title>Chillmode.info login page</title>\n" +
	"    <meta name=\"google-signin-client_id\" content=\""+sitedef.GoogleSigninClientId+"\">\n" +
	"</head>\n<body id=\"body\" style=\"overflow: hidden; background-color: rgba(0,0,0,0);\">\n<script src=\"external/pack.js\"></script>\n" +
	"<script src=\"https://apis.google.com/js/platform.js\" async defer></script>\n" +
	"<h1>Please sign in below:</h1>\n" +
	"<div id=\"signin-button\" class=\"g-signin2\" data-onsuccess=\"onSignIn\"></div>\n" +
	"<script>\n" +
	"    function onSignIn(googleUser) {\n" +
	"        var xhr=new XMLHttpRequest()\n" +
	"        xhr.open(\"PUT\", \"/rest/v0/associate\")\n" +
	"        xhr.setRequestHeader(\"Content-Type\", \"application/json\")\n" +
	"        xhr.addEventListener(\"load\", () => {\n" +
	"            var auth2=gapi.auth2.getAuthInstance();\n" +
	"            auth2.signOut().then(() => {\n" +
	"                document.getElementById(\"signin-button\").remove()\n" +
	"                document.getElementsByTagName(\"body\")[0].appendChild(document.createTextNode(\"Thanks, you can close this tab now\"))\n" +
	"            })\n" +
	"        })\n" +
	"        var query=\"\"\n" +
	"        var i = location.href.indexOf('?')\n" +
	"        if (i>=0) {\n" +
	"            query=location.href.substr(i+1)\n" +
	"        }\n" +
	"        xhr.send(JSON.stringify({\n" +
	"            \"idToken\": \"google_token:\"+googleUser.getAuthResponse().id_token,\n" +
	"            \"ReqMsg\": {\n" +
	"                \"ClientId\": query\n" +
	"            },\n" +
	"        }))\n" +
	"    }\n" +
	"</script>\n" +
	"</body>\n" +
	"</html>"

func handleAssociatePage(mux *http.ServeMux) {
	mux.HandleFunc("/associate.html", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(page))
	})
}