[ -f "passwordfile" ] && rm passwordfile && [ -d "00/" ] && rm -r 00/ && [ -d "test00/" ] && rm -r test00/ && [ -d "01/" ] && rm -r 01/ && [ -d "test01/" ] && rm -r test01/

mkdir -p test00/ test01 00/keystore 01/keystore  00/geth 01/geth  00/plugins 01/plugins 


cp ../engine.go test00/ 
cp ../engine.go ../main.go test01/
cd test00/ 
go build -buildmode=plugin -o ../00/plugins
cd ../
cd test01/ 
go build -buildmode=plugin -o ../01/plugins
cd ../

cp UTC--2021-03-02T16-47-49.510918858Z--f2c207111cb6ef761e439e56b25c7c99ac026a01 00/keystore
cp UTC--2021-03-02T16-47-39.492920333Z--4204477bf7fce868e761caaba991ffc607717dbf 01/keystore

cp nodekey00 00/geth/nodekey
cp nodekey01 01/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./00 bor-genesis.json
$GETH init --datadir=./01 bor-genesis.json

$GETH --cache.preimages --config=./config00.toml --authrpc.port 8552 --port 64480 --verbosity=0  --nodiscover --networkid=6448 --datadir=./00/ --mine --miner.noverify --miner.etherbase f2c207111cb6ef761e439e56b25c7c99ac026a01 --unlock f2c207111cb6ef761e439e56b25c7c99ac026a01 --http --http.api eth,debug,net,admin --http.port 9545 --password passwordfile --allow-insecure-unlock --rpc.allow-unprotected-txs &
pid0=$!
$GETH --cache.preimages --config=./config01.toml --authrpc.port 8553 --port 64481 --verbosity=3 --syncmode=full --nodiscover --networkid=6448 --datadir=./01/ --unlock 4204477bf7fce868e761caaba991ffc607717dbf --miner.etherbase 4204477bf7fce868e761caaba991ffc607717dbf --password passwordfile --ws --ws.port 8546 --ws.api eth,admin --http --http.api eth,debug,net,admin --http.port 9547 --allow-insecure-unlock --rpc.allow-unprotected-txs &

sleep 2

python3 connect_nodes.py # This will need to be revisited in the future. The node keys from the config.toml *should be configuarable to make this unecessary. 

sleep 5

if ps -p $pid0 > /dev/null; then
  kill $pid0
fi

sleep 1

rm -rf 00 01 test00 test01
rm -r passwordfile
