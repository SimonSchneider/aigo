package srv

import "net/http"

func WithHeader(key, value string) ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set(key, value)
			return next.RoundTrip(r)
		})
	}
}

type RoundTripperFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type ClientMiddleware func(next http.RoundTripper) http.RoundTripper

func ClientWith(clnt *http.Client, rts ...ClientMiddleware) *http.Client {
	nc := &http.Client{
		Transport:     clnt.Transport,
		CheckRedirect: clnt.CheckRedirect,
		Jar:           clnt.Jar,
		Timeout:       clnt.Timeout,
	}
	if nc.Transport == nil {
		nc.Transport = http.DefaultTransport
	}
	for _, rt := range rts {
		nc.Transport = rt(nc.Transport)
	}
	return nc
}
