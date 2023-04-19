rm -r 01/ test01/ 
rm passwordfile
mkdir -p test01/ 01/keystore 01/geth 01/plugins
cp ../engine.go ../main.go ../hooks.go ../tracer.go ../live_tracer.go  test01/
cd test01/ 
go build -buildmode=plugin -o ../01/plugins
cd ../
cp UTC--2021-03-02T16-47-39.492920333Z--4204477bf7fce868e761caaba991ffc607717dbf 01/keystore
cp nodekey01 01/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./01 genesis.json

$GETH --config config01.toml --authrpc.port 8553 --port 64481 --verbosity=4 --syncmode=full --nodiscover --networkid=6448 --datadir=./01/ --unlock 4204477bf7fce868e761caaba991ffc607717dbf --miner.etherbase 4204477bf7fce868e761caaba991ffc607717dbf --password passwordfile --ws --ws.port 8546 --ws.api eth,admin --http --http.api eth,debug,net --http.port 9546 --allow-insecure-unlock

# wait
