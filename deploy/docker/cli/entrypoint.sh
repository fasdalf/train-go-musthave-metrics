!/bin/sh

#RUN pwd;
envsubst < pgen.yaml.dist > pgen.yaml;
#RUN pgen;
pgen;
rm -f pgen.yaml;
chmod u+x ./build.sh;
./build.sh;
go mod tidy;
./build.sh;

./proxy serve -v
