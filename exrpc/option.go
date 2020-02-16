package exrpc

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

type option func(receiver OptionReceiver) error

type OptionReceiver interface {
	setCredentials(credentials credentials.TransportCredentials)
	setTls(tls bool)
}

func WithCredentialsFromFile(certFile, keyFile string) option {
	cred, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	return func(receiver OptionReceiver) error {
		if err != nil {
			return err
		}
		receiver.setCredentials(cred)
		return nil
	}
}

func WithCredentials(certPEMBlock, keyPEMBlock []byte) option {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	cred := credentials.NewServerTLSFromCert(&cert)
	return func(receiver OptionReceiver) error {
		if err != nil {
			return err
		}
		receiver.setCredentials(cred)
		return nil
	}
}

func WithTls(tls bool) option {
	return func(receiver OptionReceiver) error {
		receiver.setTls(tls)
		return nil
	}
}
