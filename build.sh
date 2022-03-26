cd ./certs
chmod +x ./gen_ca.sh
./gen_ca.sh
cd ..
docker build -t blackbackofficial/http-proxy .
docker run -p 8080:8080 blackbackofficial/http-proxy
