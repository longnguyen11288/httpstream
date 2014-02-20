package main

// twitter oauth

import (
	"flag"
	"github.com/attilaolah/httpstream"
	"github.com/mrjones/oauth"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	maxCt    = flag.Int("maxct", 10, "Max # of messages")
	user     = flag.String("user", "", "twitter username")
	ck       = flag.String("ck", "", "Consumer Key")
	cs       = flag.String("cs", "", "Consumer Secret")
	ot       = flag.String("ot", "", "Oauth Token")
	osec     = flag.String("os", "", "OAuthTokenSecret")
	logLevel = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
	search   = flag.String("search", "android,golang,zeromq,javascript", "keywords to search for, comma delimtted")
	users    = flag.String("users", "", "list of twitter userids to filter for, comma delimtted")
)

func main() {

	flag.Parse()
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	// make a go channel for sending from listener to processor
	// we buffer it, to help ensure we aren't backing up twitter or else they cut us off
	stream := make(chan []byte, 1000)
	done := make(chan bool)

	c := oauth.NewConsumer(
		*ck,
		*cs,
		oauth.ServiceProvider{
			RequestTokenUrl:   "http://twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "http://twitter.com/oauth/authorize",
			AccessTokenUrl:    "http://twitter.com/oauth/access_token",
		},
	)

	//at := goauthcon.GetAccessToken(rt.Token, pin)
	at := oauth.AccessToken{ID: "",
		Token:    *ot,
		Secret:   *osec,
		UserRef:  *user,
		Verifier: "",
		Service:  "twitter",
	}
	// the stream listener effectively operates in one "thread"/goroutine
	// as the httpstream Client processes inside a go routine it opens
	// That includes the handler func we pass in here
	client := httpstream.NewOAuthClient(c, &at, httpstream.OnlyTweetsFilter(func(line []byte) {
		stream <- line
		// although you can do heavy lifting here, it means you are doing all
		// your work in the same thread as the http streaming/listener
		// by using a go channel, you can send the work to a
		// different thread/goroutine
	}))

	// find list of userids we are going to search for
	userIDs := make([]int64, 0)
	for _, userID := range strings.Split(*users, ",") {
		if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
			userIDs = append(userIDs, id)
		}
	}
	var keywords []string
	if search != nil && len(*search) > 0 {
		keywords = strings.Split(*search, ",")
	}
	err := client.Filter(userIDs, keywords, []string{"en"}, nil, false, done)
	if err != nil {
		httpstream.Log(httpstream.ERROR, err.Error())
	} else {

		go func() {
			// while this could be in a different "thread(s)"
			ct := 0
			for tw := range stream {
				println(string(tw))
				// heavy lifting
				ct++
				if ct > *maxCt {
					done <- true
				}
			}
		}()
		_ = <-done
	}

}
