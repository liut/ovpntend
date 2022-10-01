package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/render"

	"github.com/liut/ovpntend/pkg/ipip"
	"github.com/liut/ovpntend/pkg/settings"
	"github.com/liut/ovpntend/ui"
)

var (
	avatarReplacer = strings.NewReplacer("/0", "/60")

	base  string
	inDev bool

	FindCity = ipip.FindCity
)

func init() {
	base = "/api/vpn/"
	if settings.Current.ServerPlace != "" {
		base += settings.Current.ServerPlace + "/"
	}
}

func renderHTML(w http.ResponseWriter, r *http.Request, name string, data interface{}) (err error) {
	instance := tpl(name)
	if m, ok := data.(render.M); ok {
		user, ok := UserFromContext(r.Context())
		if ok {
			m["user"] = user
		}
		err = instance.Execute(w, m)
	} else {
		err = instance.Execute(w, data)
	}
	return
}

func tpl(name string) *template.Template {

	// var tpl *template.Template
	// tpl, err = template.New("default").Parse(string(blob))
	// if err != nil {
	// 	logger().Warnf("parse template err %s", err)
	// 	return
	// }
	// if t, ok := cachedTemplates[name]; ok {
	// 	return t
	// }

	t := template.New("_base.html").Funcs(template.FuncMap{
		"duration":    time.Since,
		"formatBytes": formatBytes,
		"findPlace":   FindPlace,
		"isOffice":    IsOfficeIP,
		"urlFor":      URLFor,
		"avatarHTML":  AvatarHTML,
	})
	if inDev {
		t = template.Must(t.ParseFiles(
			filepath.Join("ui/templates/_base.html"),
			filepath.Join("ui/templates", name),
		))
	} else {
		tmp, err := t.ParseFS(ui.FS(), "templates/_base.html")
		if err != nil {
			panic(err)
		}

		tmp, err = tmp.ParseFS(ui.FS(), "templates/"+name)
		if err != nil {
			panic(err)
		}
		t = tmp
	}

	return t
}

// AvatarHTML 生成头像的HTML标签，目前仅支持微信头像
func AvatarHTML(s string) template.HTML {
	if len(s) == 0 {
		return ""
	}
	if strings.HasSuffix(s, "/") {
		s = s + "0"
	}
	if strings.HasPrefix(s, "/bizmail") || strings.HasPrefix(s, "/wwhead") { // wechat avatar
		s = "//p.qlogo.cn" + avatarReplacer.Replace(s)
	}
	return template.HTML("<img class=\"avatar img-thumbnail\" src=\"" + s + "\">")
}

// URLFor ...
func URLFor(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(path, "/"))
}

// FindPlace ...
func FindPlace(ip string) string {
	city, pro, _ := ipip.FindCity(ip)
	if city != "" {
		if pro != "" && pro != city {
			return pro + city
		}
		return city
	}
	return "[未知地区]"
}

// IsOfficeIP ...
func IsOfficeIP(ip string) bool {
	// TODO:
	return false
}

// formatBytes ...
func formatBytes(num int) string {
	return FormatBytes(float64(num), "")
}
