package recaptcha

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

// Client wraps around a fasthttp.Client for verifying reCAPTCHA requests.
type Client fasthttp.Client

// Do verifies a reCAPTCHA v2 or v3 attempt against reCAPTCHA's API. It errors if an error occurred pinging
// Google's servers, or on validating the response from Google's servers.
func (c *Client) Do(params Request) (Response, error) {
	if err := params.Validate(); err != nil {
		return Response{}, fmt.Errorf("failed to validate recaptcha request: %w", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	req.SetRequestURI(Endpoint)

	params.MarshalTo(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := (*fasthttp.Client)(c).Do(req, res); err != nil {
		return Response{}, fmt.Errorf("failed to reach google recaptcha api: %w", err)
	}

	return Parse(res.Body())
}

// DoTimeout verifies a reCAPTCHA v2 or v3 attempt against reCAPTCHA's API. It errors if either the request times
// out after timeout, or if an error occurred pinging Google's servers, or on validating the response from Google's
// servers.
func (c *Client) DoTimeout(params Request, timeout time.Duration) (Response, error) {
	if err := params.Validate(); err != nil {
		return Response{}, fmt.Errorf("failed to validate recaptcha request: %w", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	req.SetRequestURI(Endpoint)

	params.MarshalTo(req)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	if err := (*fasthttp.Client)(c).DoTimeout(req, res, timeout); err != nil {
		return Response{}, fmt.Errorf("failed to reach google recaptcha api: %w", err)
	}

	return Parse(res.Body())
}
