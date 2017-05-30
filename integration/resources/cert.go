// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// +build integration

package resources

import (
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc/credentials"
)

const (
	caCrtResource      = "CertAuth.crt"
	clientCrtResources = "m3em_client.uberinternal.com.crt"
	clientKeyResources = "m3em_client.uberinternal.com.key"
	serverName         = "m3em_server.uberinternal.com"
	serverCrtResources = "m3em_server.uberinternal.com.crt"
	serverKeyResources = "m3em_server.uberinternal.com.key"
)

func caCertificate() (*x509.CertPool, error) {
	caCrt, err := Asset(caCrtResource)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(caCrt)
	if !ok {
		return nil, err
	}
	return certPool, nil
}

// ClientTransportCredentials return a DialOption for TLS Client communication
func ClientTransportCredentials() (credentials.TransportCredentials, error) {
	clientCrt, err := Asset(clientCrtResources)
	if err != nil {
		return nil, err
	}

	clientKey, err := Asset(clientKeyResources)
	if err != nil {
		return nil, err
	}

	certificate, err := tls.X509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, err
	}

	caCert, err := caCertificate()
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      caCert,
	}), nil

}

// ServerTransportCredentials return a DialOption for TLS Server communication
func ServerTransportCredentials() (credentials.TransportCredentials, error) {
	serverCrt, err := Asset(serverCrtResources)
	if err != nil {
		return nil, err
	}
	serverKey, err := Asset(serverKeyResources)
	if err != nil {
		return nil, err
	}
	certificate, err := tls.X509KeyPair(serverCrt, serverKey)
	if err != nil {
		return nil, err
	}

	caCert, err := caCertificate()
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    caCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}
