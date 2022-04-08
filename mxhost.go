package sms

import (
	"fmt"
	"net"
	"strings"
)

func ResolveMXHost(address string) (string, error) {
	domain, found := domainPart(address)
	if !found {
		return "", fmt.Errorf("invalid address %s", address)
	}

	mxs, err := net.LookupMX(domain)
	if err != nil {
		return "", fmt.Errorf("lookup mx records of %s: %w", domain, err)
	}

	mx, ok := First(mxs)
	if !ok {
		return "", fmt.Errorf("mx records are empty: domain: %s", domain)
	}

	return mx.Host, nil
}

func domainPart(address string) (part string, found bool) {
	_, after, found := strings.Cut(address, "@")
	return after, found
}

func First[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[0], true
}
