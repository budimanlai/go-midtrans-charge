package main

import (
	"time"

	"encoding/base64"

	goargs "github.com/budimanlai/go-args"
	goconfig "github.com/budimanlai/go-config"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	cfg := &goconfig.Config{}
	e := cfg.Open("config/midtrans.conf")

	if e != nil {
		panic(e)
	}

	args := &goargs.Args{}
	args.Parse()

	skey := cfg.GetString("midtrans.server_key")
	server_key := base64.StdEncoding.EncodeToString([]byte(skey))
	mode := cfg.GetString("midtrans.mode")
	port := args.GetString("port")

	var url string
	if mode == "sandbox" {
		url = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	} else if mode == "production" {
		url = "https://app.midtrans.com/snap/v1/transactions"
	} else {
		panic("Invalid mode. Must be sandbox or production")
	}

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Page not found or wrong HTTP request method is used A")
	})
	app.Post("/", func(c *fiber.Ctx) error {
		resp, e := chargeAPI(url, server_key, c.Body())
		if e != nil {
			return c.SendString(e.Error())
		}

		return c.Send(resp.Body())
	})

	app.Listen(":" + port)
}

func chargeAPI(url string, server_key string, body []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.DisableNormalizing()

	req.Header.SetContentType(`application/json`)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+server_key)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(body)

	respClone := &fasthttp.Response{}
	e := fasthttp.DoTimeout(req, resp, 60*time.Second)
	resp.CopyTo(respClone)

	return respClone, e
}
