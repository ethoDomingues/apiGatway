package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/ethoDomingues/braza"
	"github.com/ethoDomingues/gateway/models"
	"gorm.io/gorm"
)

func beforeRequest(ctx *braza.Ctx) {
	ctx.Global["db"] = models.DBSession()
}

func main() {
	app := braza.NewApp(&braza.Config{
		Servername:           "localhost",
		DisableParseFormBody: true,
	})
	app.BeforeRequest = beforeRequest

	app.AddRoute(&braza.Route{
		Name:    "newApp",
		Func:    newApp,
		Methods: []string{"POST"},
		Schema:  &SchemaNewApp{},
	})

	app.Mount(&braza.Router{
		Name:      "gate",
		Subdomain: "{sub}",
		Routes:    routes,
	})

	db := models.DBSession()
	db.AutoMigrate(&models.GateApp{})

	app.Listen()
}

var routes = []*braza.Route{
	{
		Name:    "redirect",
		Url:     "/{url:path}",
		Func:    prox,
		Methods: []string{"GET", "POST", "PUT", "DELETE", "CONNECT", "TRACE", "PATCH"},
		Schema:  &SchemaProx{},
	},
}

func newApp(ctx *braza.Ctx) {
	ga := &models.GateApp{}
	db := models.DBSession()
	sch := ctx.Schema.(*SchemaNewApp)

	db.Where("host = ? OR name = ?", sch.Host, sch.Name).Find(ga)
	if ga.UUID != "" {
		ctx.JSON(ga, 400)
	}

	ga.Host = sch.Host
	ga.Name = sch.Name

	db.Save(ga)
	ctx.JSON(ga, 201)
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func prox(ctx *braza.Ctx) {
	sch := ctx.Schema.(*SchemaProx)
	ga := &models.GateApp{}
	db := ctx.Global["db"].(*gorm.DB)
	db.Where("name = ?", sch.Sub).Find(ga)
	if ga.UUID == "" {
		ctx.NotFound()
	}

	rq := ctx.Request

	delHopHeaders(rq.Header)
	if clientIP, _, err := net.SplitHostPort(rq.RemoteAddr); err == nil {
		appendHostToXForwardHeader(rq.Header, clientIP)
	}
	client := &http.Client{}

	uri := path.Join(ga.Host, rq.URL.Path)
	nrq, _ := http.NewRequest(rq.Method, "http://"+uri, rq.Body)
	copyHeader(nrq.Header, rq.Header)
	rsp, err := client.Do(nrq)

	if err != nil {
		fmt.Println(err)
		ctx.InternalServerError()
	}
	defer rsp.Body.Close()

	delHopHeaders(rsp.Header)
	copyHeader(ctx.Response.Header(), rsp.Header)
	ctx.Response.StatusCode = rsp.StatusCode
	io.Copy(ctx, rsp.Body)
	ctx.ReadFrom(rsp.Body)
}

type SchemaNewApp struct {
	Host string `braza:"required"`
	Name string `braza:"required"`
}

type SchemaProx struct {
	Sub string `braza:"in=subdomain,required"`
}
