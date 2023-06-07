[ -f "passwordfile" ] && rm passwordfile && [ -d "00/" ] && rm -r 00/ && [ -d "test00/" ] && rm -r test00/
mkdir -p test00/ 00/keystore  00/geth  00/plugins 
cp ../engine.go test00/ 
cd test00/ 
go build -buildmode=plugin -o ../00/plugins
cd ../
cp UTC--2021-03-02T16-47-49.510918858Z--f2c207111cb6ef761e439e56b25c7c99ac026a01 00/keystore
cp nodekey00 00/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./00 genesis.json

$GETH --config config00.toml --authrpc.port 8552 --port 64480 --verbosity 5 --nodiscover --networkid=6448 --datadir=./00/ --mine --miner.etherbase f2c207111cb6ef761e439e56b25c7c99ac026a01 --unlock f2c207111cb6ef761e439e56b25c7c99ac026a01 --http --http.api eth,debug,net --http.port 9545 --password passwordfile --allow-insecure-unlock

# pid=$!

# sleep 8

# if ps -p $pid > /dev/null; then
#   kill $pid
# fi
