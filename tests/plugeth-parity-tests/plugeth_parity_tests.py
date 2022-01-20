import json
import requests

geth = 'http://127.0.0.1:8555'
parity = 'http://127.0.0.1:8545'

def get_transactions():
	data = {"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", False],"id":1}
	r = requests.post(geth, json=data).json()
	return r['result']['transactions']

def trace(tx):
	data = {"jsonrpc":"2.0","method":"trace_replayTransaction","params":[f"{tx}", ["trace"]],"id":1}
	g = requests.post(f'{geth}', json=data).json()
	p = requests.post(f'{parity}', json=data).json()
	g_result = (g['result']['output'], g['result']['trace'])
	p_result = (p['result']['output'], p['result']['trace'])
	if g_result[0] != p_result[0]:
		return ("outputs dont match", tx) 
	if len(g_result[1]) != len(p_result[1]):
		return ("traces of different lengths", tx)
	for i in range(len(g_result[1])):
		if g_result[1][i] != p_result[1][i]:
			print("trace item does not match", i, tx)

def test1():
	l = get_transactions()
	for item in l:
		trace(item)
	print("test complete")


if __name__ == "__main__":
	test1()
   
