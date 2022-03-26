cd ./certs
chmod +x ./gen_ca.sh
./gen_ca.sh
cd ..
docker build -t blackbackofficial .
docker run -p 8080:8080 -p 9432:5432 blackbackofficial