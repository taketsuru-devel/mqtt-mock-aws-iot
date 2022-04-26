package connection

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewTLSConfig(PemCa *string, PemCrt *string, PemKey *string) *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	//pemCerts, err := ioutil.ReadFile("../cert/awsRootCa.pem")
	pemCerts, err := ioutil.ReadFile(*PemCa)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	} else {
		panic(err)
	}

	// Import client certificate/key pair
	//cert, err := tls.LoadX509KeyPair("../cert/660c4e2a8c-certificate.pem.crt", "../cert/660c4e2a8c-private.pem.key")
	cert, err := tls.LoadX509KeyPair(*PemCrt, *PemKey)
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	//fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

/*
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
    fmt.Printf("TOPIC: %s\n", msg.Topic())
    fmt.Printf("MSG: %s\n", msg.Payload())
}
*/

func GetConnection(clientID string, PemCa *string, PemCrt *string, PemKey *string, EndPoint *string) *MQTT.Client {
	tlsconfig := NewTLSConfig(PemCa, PemCrt, PemKey)

	opts := MQTT.NewClientOptions()
	//コンソールの設定：カスタムエンドポイント
	opts.AddBroker(*EndPoint)
	//他のリソースと重複しなければいいらしい
	opts.SetClientID(clientID).SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})

	// Start the connection
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &c
}
