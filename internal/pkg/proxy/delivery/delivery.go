package delivery

import (
	"bufio"
	"crypto/tls"
	"http-proxy/internal/models"
	"http-proxy/internal/pkg/proxy"
	"http-proxy/internal/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type StartServer struct {
	repo proxy.RepoProxy
	port string
}

func NewProxyServer(ProxyRepo proxy.RepoProxy, port string) *StartServer {
	return &StartServer{repo: ProxyRepo, port: port}
}

func (ps *StartServer) ListenAndServe() error {
	server := http.Server{
		Addr: ps.port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				ps.proxyHTTPS(w, r)
			} else {
				ps.proxyHTTP(w, r)
			}
		}),
	}

	return server.ListenAndServe()
}

func (ps *StartServer) proxyHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")

	reqId, err := ps.repo.SaveRequest(r)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	_, err = ps.repo.SaveResponse(reqId, resp)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ps *StartServer) proxyHTTPS(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	localConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Printf("hijacking error: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	_, err = localConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		log.Printf("handshaking failed: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		localConn.Close()
		return
	}
	defer localConn.Close()

	host := strings.Split(r.Host, ":")[0]

	tlsConfig, err := utils.GenTLSConf(host, r.URL.Scheme)
	if err != nil {
		log.Printf("error getting cert: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tlsLocalConn := tls.Server(localConn, &tlsConfig)
	err = tlsLocalConn.Handshake()
	if err != nil {
		tlsLocalConn.Close()
		log.Printf("handshaking failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tlsLocalConn.Close()

	remoteConn, err := tls.Dial("tcp", r.URL.Host, &tlsConfig)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer remoteConn.Close()

	reader := bufio.NewReader(tlsLocalConn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Printf("error getting request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestByte, err := httputil.DumpRequest(request, true)
	if err != nil {
		log.Printf("failed to dump request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = remoteConn.Write(requestByte)
	if err != nil {
		log.Printf("failed to write request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serverReader := bufio.NewReader(remoteConn)
	response, err := http.ReadResponse(serverReader, request)
	if err != nil {
		log.Printf("failed to read response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawResponse, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Printf("failed to dump response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tlsLocalConn.Write(rawResponse)
	if err != nil {
		log.Printf("fail to write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.URL.Scheme = "https"
	hostAndPort := strings.Split(r.URL.Host, ":")
	request.URL.Host = hostAndPort[0]

	reqId, err := ps.repo.SaveRequest(request)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}

	_, err = ps.repo.SaveResponse(reqId, response)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}
}

func (ps *StartServer) ProxyHTTP(r *http.Request) *models.Response {
	r.Header.Del("Proxy-Connection")

	reqId, err := ps.repo.SaveRequest(r)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	response, err := ps.repo.SaveResponse(reqId, resp)
	if err != nil {
		log.Printf("fail save to db: %v", err)
	}
	return response
}
