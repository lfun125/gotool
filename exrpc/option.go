package exrpc

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

type Option func(receiver OptionReceiver) error

type OptionReceiver interface {
	setCredentials(credentials credentials.TransportCredentials)
	setTls(tls bool)
}

func WithServerCredentialsFromFile(certFile, keyFile string) Option {
	cred, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	return func(receiver OptionReceiver) error {
		if err != nil {
			return err
		}
		receiver.setCredentials(cred)
		receiver.setTls(true)
		return nil
	}
}

func WithServerCredentials(certPEMBlock, keyPEMBlock []byte) Option {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	cred := credentials.NewServerTLSFromCert(&cert)
	return func(receiver OptionReceiver) error {
		if err != nil {
			return err
		}
		receiver.setCredentials(cred)
		receiver.setTls(true)
		return nil
	}
}

func WithClientCredentialsFromFile(certFile, serverNameOverride string) Option {
	cred, err := credentials.NewClientTLSFromFile(certFile, serverNameOverride)
	return func(receiver OptionReceiver) error {
		if err != nil {
			return err
		}
		receiver.setCredentials(cred)
		return nil
	}
}

func WithTls(tls bool) Option {
	return func(receiver OptionReceiver) error {
		receiver.setTls(tls)
		return nil
	}
}
