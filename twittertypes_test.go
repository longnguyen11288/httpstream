package httpstream

import (
	"encoding/json"
	//"github.com/bsdf/twitter"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	tweets    = make([]string, 0)
	tweetData []interface{}
)

func init() {
	SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), "debug")
	loadJSONData()
}

func loadJSONData() {
	// load the tweet data
	if jsonb, err := ioutil.ReadFile("data/testdata.json"); err == nil {
		parts := bytes.Split(jsonb, []byte("\n\n"))
		for _, part := range parts {
			tweets = append(tweets, string(bytes.Trim(part, "\n \t\r")))
		}
	}
}
func prettyJSON(js string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(js), &m); err == nil {
		if b, er := json.MarshalIndent(m, "", "  "); er == nil {
			log.Println(string(b))
		} else {
			log.Println(er)
		}
	} else {
		log.Println(err)
	}
}

func TestDecodeTweet1Test(t *testing.T) {
	twlist := make([]Tweet, 0)
	for i := 0; i < len(tweets); i++ {
		//log.Println(tweets[i])
		//for i := 3; i < 4; i++ {
		tw := Tweet{}
		err := json.Unmarshal([]byte(tweets[i]), &tw)
		if err != nil {
			t.Error(err)
			log.Println(tweets[i][0:100])
		}
		log.Println(i, " ", err, tw.Text)
		twlist = append(twlist, tw)
	}
	/*
		twx := twlist[1]
		for _, url := range twx.URLs() {
			Debug(url)
		}
		twx = twlist[1]
		u := twx.Entities.URLs[0]
		log.Println(twx.URLs())
		log.Println(u.ExpandedURL)
	*/
	//prettyJSON(tweet3)
	//tw2 := twitter.Tweet{}
	//err = json.Unmarshal([]byte(tweet2), &tw2)
	//log.Println(err)
}
