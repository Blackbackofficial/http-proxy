# http-proxy

## Description:
http-proxy deployed on port 8080, proxying http and https requests
<br/><br/>
Web-api deployed on port 8000 && implemented param-miner vulnerability scanner


### Examples proxy requests
* HTTP `curl -x http://127.0.0.1:8080 http://a06367.ru/`
* HTTPS `curl -k https://a06367.ru/ -x http://127.0.0.1:8080/ -vvv`

### Build
Certificates are generated and docker is started

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
4. `GET /api/v1/scan/{id}` – Request vulnerability scanner (param-miner);

### Param-miner
Param-miner - add to the request in turn each GET parameter from the [params from PortSwigger](https://github.com/PortSwigger/param-miner/blob/master/resources/params) dictionary with a random value (?param=shefuisehfuishe) 
<br/><br/>
We are looking for the specified random value in the response, if found, we display its name of the hidden parameter.