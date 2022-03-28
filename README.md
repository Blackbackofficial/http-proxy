# http-proxy

## Description:

http proxy deployed on port 8080, proxying http and https requests
Web-api deployed on port 8000

### Examples proxy requests

* HTTP `curl -x http://127.0.0.1:8080 https://a06367.ru/`
* HTTPS `curl -k https://a06367.ru/ -x http://127.0.0.1:8080/ -vvv`

### Build

`./build.sh`

### API 
Parsing requests and response:
* HTTP method (GET/POST/PUT/HEAD)
* Path and GET parameters
* Headers, while separately parsing Cookies
* Request body, in case of application/x-www-form-urlencoded separate POST parameters
* gzip and other compression methods


### API Description
1. `GET /api/v1/requests` – List of requests;
2. `GET /api/v1/requests/{id}` – Output 1 request;
3. `GET /api/v1/repeat/{id}` – Resubmit request;
4. `GET /api/v1/scan/{id}` – Scan request;
