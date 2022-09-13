package main

import (
	"github.com/mgutz/logxi/v1"
	"html/template"
	"net/http"
)

const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
            <meta charset="utf-8">
            <title>Список фильмов с www.afisha.ru </title>
        </head>
        <body>
            {{if .}}
				<ol>
                {{range .}}
					<li>
                    <a href="https://weather.rambler.ru/v-moskve/{{.Ref}}">{{.Title}}</a>
                    <br/>
					</li>
                {{end}}
				</ol>
            {{else}}
                Не удалось загрузить новости!
            {{end}}
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, downloadNews()); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
