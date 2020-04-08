package recaptcha

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"time"
)

// Endpoint to reCAPTCHA API.
const Endpoint = "https://google.com/recaptcha/api/siteverify"

// Request holds the payload to be sent to reCAPTCHA's API.
type Request struct {
	// Required. The shared key between your site and reCAPTCHA.
	Secret string

	// Required. The user response token provided by the reCAPTCHA client-side integration on your site.
	Response string

	// Optional. The user's IP address.
	RemoteIP string
}

// Validate validates whether or not this request contains valid parameters. It returns an error should they be
// invalid. Validation right now has to do with checking whether or not required parameters are filled.
func (r Request) Validate() error {
	if r.Secret == "" {
		return errors.New("secret must not be empty")
	}

	if r.Response == "" {
		return errors.New("response must not be empty")
	}

	return nil
}

// MarshalTo fills a *fasthttp.Request with its parameters. Make sure to Validate() to validate this Request's parameters
// beforehand.
func (r Request) MarshalTo(req *fasthttp.Request) {
	args := req.PostArgs()
	args.Set("secret", r.Secret)
	args.Set("response", r.Response)

	if r.RemoteIP != "" {
		args.Set("remoteip", r.RemoteIP)
	}
}

// Response represents a response from reCAPTCHA's API.
type Response struct {
	// Whether this request was a valid reCAPTCHA token for your site.
	Success bool `json:"success"`

	// (v3) The score for this request (0.0 - 1.0). 0.0 signifies high likelihood of a bot, and 1.0 signifies
	// otherwise.
	Score float64 `json:"score"`

	// (v3) The action name for this request (important to verify).
	Action string `json:"action"`

	// Timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ).
	ChallengeTS time.Time `json:"challenge_ts"`

	// The hostname of the site, or package name of the Android APK where the reCAPTCHA was solved.
	Hostname string `json:"hostname"`

	// Optional error codes.
	ErrorCodes []string `json:"error-codes"`
}

// Parse parses either a reCAPTCHA v2 or v3 response from Google reCAPTCHA's API.
func Parse(body []byte) (Response, error) {
	var r Response

	val, err := fastjson.ParseBytes(body)
	if err != nil {
		return r, fmt.Errorf("failed to parse api json response: %w", err)
	}

	var v *fastjson.Value

	if v = val.Get("success"); v == nil {
		return r, errors.New("json field 'success' not found")
	}
	if r.Success, err = v.Bool(); err != nil {
		return r, fmt.Errorf("failed to parse json field 'success': %w", err)
	}

	if !r.Success {
		if v = val.Get("error-codes"); v == nil {
			return r, errors.New("json field 'error-codes' not found")
		}
		codes, err := v.Array()
		if err != nil {
			return r, fmt.Errorf("failed to parse json field 'error-codes': %w", err)
		}
		r.ErrorCodes = make([]string, len(codes))
		for i := range r.ErrorCodes {
			r.ErrorCodes[i] = string(codes[i].GetStringBytes())
		}

		return r, nil
	}

	/** START v3 FIELDS **/

	if v = val.Get("score"); v != nil {
		r.Score, err = v.Float64()
		if err != nil {
			return r, fmt.Errorf("failed to parse json field 'score': %w", err)
		}
	}

	if v = val.Get("action"); v != nil {
		r.Action = string(v.GetStringBytes())
	}

	/** END v3 FIELDS **/

	if v = val.Get("challenge_ts"); v == nil {
		return r, errors.New("json field 'challenge_ts' not found")
	}
	if r.ChallengeTS, err = time.Parse(time.RFC3339, string(v.GetStringBytes())); err != nil {
		return r, fmt.Errorf("failed to parse json field 'challenge_ts': %w", err)
	}

	if v = val.Get("hostname"); v != nil {
		r.Hostname = string(v.GetStringBytes())
	} else if v = val.Get("apk_package_name"); v != nil {
		r.Hostname = string(v.GetStringBytes())
	} else {
		return r, fmt.Errorf("both json fields 'hostname' and 'apk_package_name' not found")
	}

	return r, nil
}

// Do verifies a reCAPTCHA v2 or v3 attempt against reCAPTCHA's API. It errors if an error occurred pinging
// Google's servers, or on validating the response from Google's servers.
func Do(params Request) (Response, error) {
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

	if err := fasthttp.Do(req, res); err != nil {
		return Response{}, fmt.Errorf("failed to reach google recaptcha api: %w", err)
	}

	return Parse(res.Body())
}

// DoTimeout verifies a reCAPTCHA v2 or v3 attempt against reCAPTCHA's API. It errors if either the request times
// out after timeout, or if an error occurred pinging Google's servers, or on validating the response from Google's
// servers.
func DoTimeout(params Request, timeout time.Duration) (Response, error) {
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

	if err := fasthttp.DoTimeout(req, res, timeout); err != nil {
		return Response{}, fmt.Errorf("failed to reach google recaptcha api: %w", err)
	}

	return Parse(res.Body())
}
