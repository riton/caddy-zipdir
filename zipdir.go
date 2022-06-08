package zipdir

import "github.com/caddyserver/caddy/v2"

func init() {
	caddy.RegisterModule(ZipDir{})
}

type ZipDir struct {
}

func (z ZipDir) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.zip_dir",
		New: func() caddy.Module {
			return ZipDir{}
		},
	}
}
