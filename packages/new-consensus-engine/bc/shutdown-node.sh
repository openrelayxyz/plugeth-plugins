[ -f "passwordfile" ] && rm passwordfile && [ -d "02/" ] && rm -r 02/ && [ -d "test02/" ] && rm -r test02/
mkdir -p test02/ 02/keystore 02/geth 02/plugins
cp ../engine.go ../shutdown.go test02/
cd test02/ 
go build -buildmode=plugin -o ../02/plugins
cd ../
cp UTC--2021-03-02T16-47-59.816632526Z--2cb2e3bdb066a83a7f1191eef1697da51793f631 02/keystore
cp nodekey02 02/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./02 genesis.json

$GETH --config config02.toml --authrpc.port 8556 --port 64484 --verbosity=4 --syncmode=full --nodiscover --networkid=6448 --datadir=./02/ --unlock 2cb2e3bdb066a83a7f1191eef1697da51793f631 --miner.etherbase 2cb2e3bdb066a83a7f1191eef1697da51793f631 --password passwordfile --ws --ws.port 8548 --ws.api eth,admin --http --http.api eth,debug,net --http.port 9547 --allow-insecure-unlock &

pid=$!

sleep 5

if ps -p $pid > /dev/null; then
  kill $pid
fi