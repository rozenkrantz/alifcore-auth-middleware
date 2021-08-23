package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dequinox/alifcore-auth-middleware/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewPublicKey),
	fx.Invoke(Worker),
)

type Params struct {
	fx.In
	Config config.Config
}

type PubKeyResponse struct {
	Meta     Meta     `json:"meta"`
	Response Response `json:"response"`
}

type Meta struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	StatusCode uint64 `json:"statusCode"`
}

type Response struct {
	PublicKey string `json:"public_key"`
}

func NewPublicKey(p Params) (*rsa.PublicKey, error) {

	url := p.Config.GetString("PUB_KEY_URI")
	data := p.Config.GetString("PUB_KEY_DATA")

	return GetPublicKey(url, data)
}

func GetPublicKey(url, data string) (*rsa.PublicKey, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response PubKeyResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if response.Meta.Error {
		err = errors.New("error retrieving public key")
		return nil, err
	}

	pubPem, _ := pem.Decode([]byte(response.Response.PublicKey))
	if pubPem == nil {
		err = errors.New("error decoding public key")
		return nil, err
	}

	if pubPem.Type != "PUBLIC KEY" {
		err = errors.New("error type of public key")
		return nil, err
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		err = errors.New("error getting public key")
		return nil, err
	}

	return pubKey, nil
}

type WorkerParams struct {
	fx.In
	Config    config.Config
	PublicKey *rsa.PublicKey
}

func Worker(w WorkerParams) {
	var err error
	url := w.Config.GetString("PUB_KEY_URI")
	data := w.Config.GetString("PUB_KEY_DATA")

	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				w.PublicKey, err = GetPublicKey(url, data)
				if err != nil {
					log.Fatalf("unable to get public keys")
				} else {
					log.Println("public keys updated")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
