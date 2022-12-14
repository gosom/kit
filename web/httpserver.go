package web

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gosom/kit/logging"
	"golang.org/x/crypto/acme/autocert"
)

type HttpHandler interface {
	http.Handler
	MethodFunc(method string, pattern string, h http.HandlerFunc)
}

/* Generate self signed certificate note
	openssl req -x509 -out localhost.crt -keyout localhost.key   -newkey rsa:2048 -nodes -sha256   -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth") --keyout certs/localhost.key --out certs/localhost.crt
*/

// ServerConfig contains the configuration of the web server
// Essentially, only the Router is mandatory in order to
// have a working http server with some sane defaults
type ServerConfig struct {
	// Host defaults to :8080
	Host string
	// Router is the http router
	Router HttpHandler
	// ReadTimeout defaults to 1m
	ReadTimeout time.Duration
	// IdleTimeout if it's zero Go uses by default the ReadTimeout
	IdleTimeout time.Duration
	// ReadHeaderTimeout  defaults to 20s
	ReadHeaderTimeout time.Duration
	// WriteTimeout defaults to 2m
	WriteTimeout time.Duration
	// MaxHeaderBytes defaults to 1MB
	MaxHeaderBytes int
	// ExitSignals defaults to os.Interrupt
	// Here define the OS signals for which the http server should perform
	// a graceful exit
	ExitSignals []os.Signal
	// Domain by default is localhost. It is used when UseTLS = true
	// If you have a valid domain then it fetches a certificate from
	// let's encrypt. It firsts looks for a certificate in a certs folder.
	// See the getSelfSignedOrLetsEncryptCert function
	// thanks to :  https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt
	Domain string
	// UseTLS when true it uses TLS
	UseTLS bool
	// LogLevel defaults to logger.InfoLevel
	LogLevel logging.Level
}

// setDefaults sets the default values for the web server
func setDefaults(cfg ServerConfig) ServerConfig {
	if len(cfg.Host) == 0 {
		cfg.Host = ":8080"
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = time.Minute
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = time.Second * 20
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = time.Minute * 2
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = 1 << 20
	}
	if len(cfg.ExitSignals) == 0 {
		cfg.ExitSignals = append(cfg.ExitSignals, os.Interrupt)
	}
	if len(cfg.Domain) == 0 {
		cfg.Domain = "localhost"
	}
	return cfg
}

type HttpServer struct {
	srv  *http.Server
	cfg  ServerConfig
	sigs chan os.Signal
}

// NewHttpServer creates a new http server
func NewHttpServer(cfg ServerConfig) *HttpServer {
	cfg = setDefaults(cfg)
	ans := HttpServer{
		srv: &http.Server{
			Addr:              cfg.Host,
			Handler:           cfg.Router,
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
		},
		cfg:  cfg,
		sigs: make(chan os.Signal, 1),
	}
	if cfg.UseTLS {
		// thanks to https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
			Cache:      autocert.DirCache("certs"),
		}

		ans.srv.TLSConfig = certManager.TLSConfig()
		ans.srv.TLSConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(&certManager)
	}
	signal.Notify(ans.sigs, cfg.ExitSignals...)
	return &ans
}

// ListenAndServe starts the http server
func (o *HttpServer) ListenAndServe(ctx context.Context) error {
	if o.srv.Handler == nil {
		return errors.New("no router defined")
	}
	var err error
	defer func() {
		if err != nil {
			logging.Log(logging.ERROR, "http server exited with error",
				"component", "http", "error", err)
		} else {
			logging.Log(logging.INFO, "http server exited gracefully",
				"component", "http")
		}
	}()
	serverShutdown := func() error {
		const timeout = 5 * time.Second
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		if err := o.srv.Shutdown(ctx2); err != nil {
			if err := o.srv.Close(); err != nil {
				return err
			}
		}
		return nil
	}
	errs := make(chan error, 1)
	go func() {
		switch o.srv.TLSConfig {
		case nil:
			errs <- o.srv.ListenAndServe()
		default:
			errs <- o.srv.ListenAndServeTLS("", "")
		}
	}()
	select {
	case <-ctx.Done():
		if err = serverShutdown(); err != nil {
			return err
		}
		err = <-errs
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return err
	case <-o.sigs:
		if err = serverShutdown(); err != nil {
			return err
		}
		err := <-errs
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return err
	case err := <-errs:
		return err
	}
}

// getSelfSignedOrLetsEncryptCert returns a function that returns a certificate
func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			return certManager.GetCertificate(hello)
		}
		return &certificate, err
	}
}
