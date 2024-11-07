package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

type API struct {
	client *workloadapi.Client
	ctx    context.Context
}

func NewAPI(ctx context.Context, client *workloadapi.Client) (*API, error) {
	return &API{
		client: client,
		ctx:    ctx,
	}, nil
}

// HealthCheckHandler is a simple health check handler that returns a 200 OK status
// if server can get the X.509 SVID from the Workload API client.
func (a *API) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := a.client.FetchX509SVID(a.ctx)
	if err != nil {
		http.Error(w, "Error fetching X.509 SVID", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) GetJwtHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	svid, err := a.client.FetchJWTSVID(ctx, jwtsvid.Params{
		Audience: "spirl.com",
		ExtraAudiences: []string{
			"spiffe://example.org/foo",
			"spiffe://acme.com/bar",
		},
	})
	if err != nil {
		str := "Error fetching JWT SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)
	log.Printf("JWT SVID fetched in %s", elapsed)

	// Convert the JWT-SVID to a JSON response
	response := struct {
		Token string `json:"token"`
	}{
		Token: svid.Marshal(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response as JSON: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("JWT: %s", svid.Marshal())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetX509Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	svid, err := a.client.FetchX509SVID(ctx)
	if err != nil {
		str := "Error fetching X509 SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)
	log.Printf("X509 SVID fetched in %s", elapsed)

	cert, key, err := svid.MarshalRaw()
	if err != nil {
		str := "Error marshalling X509 SVID: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	// Convert the X509-SVID to a JSON response
	response := struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
	}{
		Cert: base64.StdEncoding.EncodeToString(cert),
		Key:  base64.StdEncoding.EncodeToString(key),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response as JSON: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("X.509-SVID: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetX509TrustBundleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	bundles, err := a.client.FetchX509Bundles(ctx)
	if err != nil {
		str := "Error fetching bundles: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	bundleMap := make(map[string][]string)
	for _, bundle := range bundles.Bundles() {
		trustDomain := bundle.TrustDomain().IDString()
		for _, authority := range bundle.X509Authorities() {
			encoded := base64.StdEncoding.EncodeToString(authority.Raw)
			bundleMap[trustDomain] = append(bundleMap[trustDomain], encoded)
		}
	}

	jsonResponse, err := json.MarshalIndent(map[string]interface{}{"bundles": bundleMap}, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Trust bundle: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (a *API) GetJwtTrustBundleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	bundles, err := a.client.FetchJWTBundles(ctx)
	if err != nil {
		str := "Error fetching bundles: " + err.Error()
		log.Printf("%v", str)
		http.Error(w, str, http.StatusInternalServerError)
		return
	}

	bundleMap := make(map[string]json.RawMessage)
	for _, bundle := range bundles.Bundles() {
		trustDomain := bundle.TrustDomain().IDString()
		jwks, err := bundle.Marshal()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling bundle: %v", err), http.StatusInternalServerError)
			return
		}
		bundleMap[trustDomain] = jwks
	}

	jsonResponse, err := json.MarshalIndent(map[string]interface{}{"bundles": bundleMap}, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Trust bundle: %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
