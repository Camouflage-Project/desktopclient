package proxy

import (
	"context"
	"crypto/tls"
	"desktopClient/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

func InitializeForwardProxy(c *config.Configuration, stdLogger *log.Logger, logger *zap.Logger) {
	logger.Info("initializing forward proxy")
	var (
		flagCertPath = ""
		flagKeyPath  = ""
		flagAddr     = "localhost:" + strconv.Itoa(c.ProxyPort)
		flagAuthUser = ""
		flagAuthPass = ""
		flagAvoid    = ""

		flagDestDialTimeout         = 10*time.Second
		flagDestReadTimeout         = 5*time.Second
		flagDestWriteTimeout        = 5*time.Second
		flagClientReadTimeout       = 5*time.Second
		flagClientWriteTimeout      = 5*time.Second
		flagServerReadTimeout       = 30*time.Second
		flagServerReadHeaderTimeout = 30*time.Second
		flagServerWriteTimeout      = 30*time.Second
		flagServerIdleTimeout       = 30*time.Second

		flagLetsEncrypt = false
		flagLEWhitelist = ""
		flagLECacheDir  = "/tmp"
	)

	p := &Proxy{
		ForwardingHTTPProxy: NewForwardingHTTPProxy(stdLogger),
		Logger:              logger,
		AuthUser:            flagAuthUser,
		AuthPass:            flagAuthPass,
		DestDialTimeout:     flagDestDialTimeout,
		DestReadTimeout:     flagDestReadTimeout,
		DestWriteTimeout:    flagDestWriteTimeout,
		ClientReadTimeout:   flagClientReadTimeout,
		ClientWriteTimeout:  flagClientWriteTimeout,
		Avoid:               flagAvoid,
	}

	s := &http.Server{
		Addr:              flagAddr,
		Handler:           p,
		ErrorLog:          stdLogger,
		ReadTimeout:       flagServerReadTimeout,
		ReadHeaderTimeout: flagServerReadHeaderTimeout,
		WriteTimeout:      flagServerWriteTimeout,
		IdleTimeout:       flagServerIdleTimeout,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	if flagLetsEncrypt {
		if flagLEWhitelist == "" {
			p.Logger.Fatal("error: no -lewhitelist flag set")
		}
		if flagLECacheDir == "/tmp" {
			p.Logger.Info("-lecachedir should be set, using '/tmp' for now...")
		}

		m := &autocert.Manager{
			Cache:      autocert.DirCache(flagLECacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(flagLEWhitelist),
		}

		s.Addr = ":https"
		s.TLSConfig = m.TLSConfig()
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		p.Logger.Info("Server shutting down")
		if err := s.Shutdown(context.Background()); err != nil {
			p.Logger.Error("Server shutdown failed", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	var svrErr error
	if flagCertPath != "" && flagKeyPath != "" || flagLetsEncrypt {
		svrErr = s.ListenAndServeTLS(flagCertPath, flagKeyPath)
	} else {
		svrErr = s.ListenAndServe()
	}

	if svrErr != http.ErrServerClosed {
		p.Logger.Error("Listening for incoming connections failed", zap.Error(svrErr))
	}

	<-idleConnsClosed
	p.Logger.Info("Server stopped")
}

