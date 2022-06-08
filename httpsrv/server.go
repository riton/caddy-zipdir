package httpsrv

import (
	"archive/zip"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/riton/dirzipper/fileslist"
)

type ServerOptions struct {
	FilesListProcessor fileslist.Processor
	ZIPFileUrl         string
	ZIPFilename        string
	Debug              bool
}

type httpServer struct {
	filesListProcessor fileslist.Processor
	zipFilename        string
	e                  *echo.Echo
}

func NewHTTPServerWithOptions(opts ServerOptions) *httpServer {
	h := &httpServer{
		filesListProcessor: opts.FilesListProcessor,
		zipFilename:        opts.ZIPFilename,
	}

	h.setupHTTPEngine(opts)

	return h
}

func (h *httpServer) setupHTTPEngine(opts ServerOptions) {
	e := echo.New()
	if opts.Debug {
		e.Debug = true
		e.Logger.SetLevel(log.DEBUG)
	}

	e.Use(middleware.Logger())

	e.GET("/"+opts.ZIPFileUrl, h.serveHTTP)

	h.e = e
}

func (h *httpServer) ListenAndServe(listen string) error {
	return h.e.Start(listen)
}

func (h *httpServer) serveHTTP(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", h.zipFilename))

	zipWriter := zip.NewWriter(c.Response().Writer)

	files, err := h.filesListProcessor.GetFiles()
	if err != nil {
		c.Logger().Errorf("fail to list files: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	for _, file := range files {
		c.Logger().Debugf("adding file %s to zip", file)
		if err := addFileToZip(zipWriter, file); err != nil {
			c.Logger().Errorf("fail to add file %s to zip: %w", file, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	if err := zipWriter.Close(); err != nil {
		c.Logger().Errorf("fail to close zip writer: %w", err)
		return errors.Wrap(err, "closing zip writer")
	}

	return nil
}
