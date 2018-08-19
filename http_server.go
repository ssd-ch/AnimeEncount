package main

import (
  "fmt"
  "html/template"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
  "math/rand"
  "time"
  "strconv"
)

const startYear = 2014

/** JSONデコード用に構造体定義 */
type Anime struct {
  Id       int    `json:"id"`
  Title    string `json:"title"`
}

func currentSeason() int {
  return (int(time.Now().Month()) - 1) / 3 + 1
}

func htmlHandlerTop(w http.ResponseWriter, r *http.Request) {

  season := []string{"冬", "春", "夏", "秋"}

  // テンプレートをパース
  t := template.Must(template.ParseFiles("top.html.tpl"))

  str := strconv.Itoa(time.Now().Year()) + "年" + season[currentSeason()-1]

  // テンプレートを描画
  if err := t.ExecuteTemplate(w, "top.html.tpl", str); err != nil {
    log.Fatal(err)
  }
}

func htmlHandlerAnime(w http.ResponseWriter, r *http.Request) {

  now := time.Now()

  //乱数表初期化
  rand.Seed(now.UnixNano())

  //url
  year := startYear + rand.Intn(now.Year() - startYear)
  qr := rand.Intn(4) + 1
  if year == now.Year() {
    if currentSeason() != 1 {
      qr = rand.Intn(currentSeason()) + 1
    } else {
      qr = 1
    }
  }
  url := fmt.Sprintf("http://api.moemoe.tokyo/anime/v1/master/%d/%d", year , qr)
  fmt.Printf(url + "\n")

  //json取得
  resp, _ := http.Get(url)
  defer resp.Body.Close()
  byteArray, _ := ioutil.ReadAll(resp.Body)
  jsonBytes := ([]byte)(byteArray)

  // JSONデコード
  var animes []Anime
  if err := json.Unmarshal(jsonBytes, &animes); err != nil {
    log.Fatal(err)
  }

  // fmt.Fprintf(w, "Anime List \n")
  // for _, p := range animes {
  //   fmt.Fprintf(w, "%d : %s \n", p.Id, p.Title)
  // }

  // テンプレートをパース
  t := template.Must(template.ParseFiles("anime.html.tpl"))

  str := "読み込みに失敗しました"
  if len(animes) > 0 {
    str = animes[rand.Intn(len(animes))].Title
  }

  // テンプレートを描画
  if err := t.ExecuteTemplate(w, "anime.html.tpl", str); err != nil {
    log.Fatal(err)
  }
}

func main() {

  http.HandleFunc("/", htmlHandlerTop)

  http.HandleFunc("/anime", htmlHandlerAnime)

  log.Fatal(http.ListenAndServe(":8081", nil))
}
