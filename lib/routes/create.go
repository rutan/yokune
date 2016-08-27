package routes

import (
	"../image"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"math/rand"
	"net/http"
	"unicode/utf8"
)

func Create(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	contents, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	blob, err := action(contents)
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	// TODO: ここでS3に投げるようにする
	w.Write(blob)
}

func action(contents []byte) ([]byte, error) {
	canvas := image.LoadFromBlob(contents)
	err := canvas.ResizeContain(480, 480)
	if err != nil {
		return nil, err
	}
	width, height := canvas.GetSize()
	text := selectText()
	setting := initTextSetting()
	setting.Size = width / float64(utf8.RuneCountInString(text))
	canvas.DrawText(text, width/2, height*0.92, &setting)
	return canvas.ExportBlob(), nil
}

func selectText() string {
	TEXT_LIST := []string{
		"LGTM",
		"いいね！",
		"よさそう",
		"みました",
	}
	return TEXT_LIST[rand.Intn(len(TEXT_LIST))]
}

func initTextSetting() image.TextSetting {
	setting := image.NewTextSetting()
	setting.Font = "./fonts/toroman.ttf"
	setting.FillColor = "#ffffff"
	setting.BorderColor = "#444444"
	setting.BorderWidth = 4
	return setting
}