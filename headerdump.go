package headerdump

import (
	"context"
	"net/http"
	"log"
	"io"
	"os"
)

type Config struct {
	Prefix string `json:"string,omitempty"`
}

// user's don't need to declare a prefix, so default to HDlog
func CreateConfig() *Config {
	return &Config{
		Prefix: "HDlog",
	}
}

type Headerdump struct {
	next     http.Handler
	name     string
	config   *Config
	hdlog    *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	// create logger with user's prefix, write to stdout, then pass it along
	hdlog := log.New(io.Discard, config.Prefix + ": ", log.Lmsgprefix|log.Ldate|log.Ltime)

	hdlog.SetOutput(os.Stdout);

	return &Headerdump{
		next:     next,
		name:     name,
		config:   config,
		hdlog:    hdlog,
	}, nil
}

func (a *Headerdump) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// dump all request headers to stdout
	for key, values := range req.Header {
		// gotta do a nest, otherwise it will flood our stdout for every single value repeatedly
		for _, value := range values {
			a.hdlog.Println(req.Host, key, value);
		}
	}
	// no modifications, continue with the next mw
	a.next.ServeHTTP(rw, req)
}
