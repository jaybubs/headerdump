package headerdump

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Prefix string
	TLS bool
}

// user's don't need to declare a prefix, so default to HDlog
func CreateConfig() *Config {
	return &Config{
		Prefix: "HDlog",
		TLS: false,
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
	a.hdlog.Println("--------------------Start of header dump--------------------")
	for key, values := range req.Header {
		// gotta do a nest, otherwise it will flood our stdout for every single value repeatedly
		for _, value := range values {
			a.hdlog.Println(req.Host, key, value);
		}
	}
	// tls dump

	// declare to allow guard clause
	var tlsconn []byte
	var err error
	if ( !a.config.TLS ) {goto ENDOFDUMP}

	a.hdlog.Println("---------- Pretty TLS Connection struct ----------")
	tlsconn, err = json.MarshalIndent(req.TLS, "", "\t")
	if err != nil {
		a.hdlog.Println("something went tits up")
		goto ENDOFDUMP
	}

	// spit out certs in der and pem formats
	a.hdlog.Printf("\n%+v\n", string(tlsconn))
	for _, cert := range req.TLS.PeerCertificates {
		// der
		a.hdlog.Printf("TLS DER CERT: \n%+v\n", cert.Raw);
		// pem
		certpem := string(pem.EncodeToMemory(&pem.Block{
			Type: "CERTIFICATE",
			Bytes: cert.Raw,
		}))
		a.hdlog.Printf("TLS PEM CERT: \n%+v\n", certpem)
	}

	ENDOFDUMP:
	a.hdlog.Println("--------------------End of header dump--------------------")

	// no modifications, continue with the next mw
	a.next.ServeHTTP(rw, req)
}
