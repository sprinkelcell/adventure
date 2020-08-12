package adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template
var defaultHandlerTmpl = `<!DOCTYPE html>
<html>

<head>
    <meta charset='utf-8'>

    <title>Adventure story</title>

</head>

<body>
	<section class="page">
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
    <ul>
        {{range .Options}}
        <li>
            <a href="/{{.Arc}}">{{.Text}}</a>
        </li>
        {{end}}
	</ul>
	</section>
	<style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
</body>

</html>`

type HandlerOption func(h *handler)

func WithTemplate(temp *template.Template) HandlerOption {
	return func(h *handler) {
		h.temp = temp
	}
}

func NewHandler(story Story, opts ...HandlerOption) http.Handler {
	h := handler{story, tpl, defaultPathfn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	story    Story
	temp     *template.Template
	pathFunc func(r *http.Request) string
}

func defaultPathfn(r *http.Request) string {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)
	if arc, ok := h.story[path]; ok {
		err := h.temp.Execute(w, arc)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func JsonParser(file io.Reader) (Story, error) {
	decodedFile := json.NewDecoder(file)
	var story Story
	if err := decodedFile.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string    `json:"title"`
	Paragraphs []string  `json:"story"`
	Options    []Options `json:"options"`
}
type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
