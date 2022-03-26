# http-proxy

## Description:

http proxy deployed on port 8080, proxying http and https requests

### Examples

* HTTP `curl -x --insecure http://127.0.0.1:8080 http://mail.ru`
* HTTPS `curl -k https://mail.ru/ -x http://127.0.0.1:8080/ -vvv`