rm -r 00/ 01/ test00/ test01/ 
rm passwordfile
mkdir -p test00/ test01/ 00/keystore 01/keystore 00/geth 01/geth 00/plugins 01/plugins
cp ../engine.go test00/ 
cp ../engine.go ../main.go ../plugins.go  test01/
cd test00/ 
go build -buildmode=plugin -o ../00/plugins
cd ../test01/ 
go build -buildmode=plugin -o ../01/plugins
cd ../
cp UTC--2021-03-02T16-47-49.510918858Z--f2c207111cb6ef761e439e56b25c7c99ac026a01 00/keystore
cp UTC--2021-03-02T16-47-39.492920333Z--4204477bf7fce868e761caaba991ffc607717dbf 01/keystore
# cp UTC--2021-03-02T16-47-59.816632526Z--2cb2e3bdb066a83a7f1191eef1697da51793f631 02/keystore
cp nodekey00 00/geth/nodekey
cp nodekey01 01/geth/nodekey
# cp nodekey02 02/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./00 genesis.json
$GETH init --datadir=./01 genesis.json
# $GETH init --datadir=./02 bn2.json

$GETH --config config00.toml --authrpc.port 8552 --port 64480 --verbosity 5 --nodiscover --networkid=6448 --datadir=./00/ --mine --miner.etherbase f2c207111cb6ef761e439e56b25c7c99ac026a01 --unlock f2c207111cb6ef761e439e56b25c7c99ac026a01 --http --http.api eth,plugeth,debug,net --http.port 9545 --password passwordfile --allow-insecure-unlock & 
$GETH --config config01.toml --authrpc.port 8553 --port 64481 --nodiscover --networkid=6448 --datadir=./01/ --unlock 4204477bf7fce868e761caaba991ffc607717dbf --miner.etherbase 4204477bf7fce868e761caaba991ffc607717dbf --password passwordfile --ws --ws.port 8546 --ws.api eth,admin --http --http.api eth,plugeth,debug,net --allow-insecure-unlock &
# $GETH --config config02.toml --verbosity 6 --authrpc.port 8554 --port 64482 --nodiscover --networkid=6448 --datadir=./02/ --unlock 2cb2e3bdb066a83a7f1191eef1697da51793f631 --password passwordfile --ws --ws.port 8547 --ws.api eth,admin --allow-insecure-unlock &

wait
