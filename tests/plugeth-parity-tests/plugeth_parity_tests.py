import json
import requests
import time
from datetime import datetime

geth = 'http://127.0.0.1:8555'
parity = 'http://127.0.0.1:8545'

def get_transactions():
	data = {"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", False],"id":1}
	r = requests.post(geth, json=data).json()
	return r['result']['transactions']

def trace(tx, target):
	data = {"jsonrpc":"2.0","method":"trace_replayTransaction","params":[f"{tx}", ["trace"]],"id":1}
	g = requests.post(f'{geth}', json=data).json()
	p = requests.post(f'{parity}', json=data).json()
	try:
		g_result = (g['result']['output'], g['result']['trace'])
		p_result = (p['result']['output'], p['result']['trace'])
	except (KeyError, UnboundLocalError) as e:
		target.write(f"{e}, {tx}\n")
		return
	try:
		if g_result[0] != p_result[0]:
			target.write(f"outputs dont match, {tx}") 
	except (IndexError, UnboundLocalError) as e:
		target.write(f"{e}, {tx}\n")
		return
	try:
		if len(g_result[1]) != len(p_result[1]):
			target.write(f"traces of different lengths, {tx}\n")
	except (IndexError, UnboundLocalError) as e:
		target.write(f"{e}, {tx}")
		return
	try:
		for i, (g_val, p_val) in enumerate(zip(g_result[1], p_result[1])):
			if g_val != p_val:
				target.write(f"trace item does not match, {i}, {tx}\n")
	except (IndexError, UnboundLocalError) as e:
		target.write(f"{e}, {tx}")
		return
	target.write("test complete\n")
	return target

def load(target):
	l = get_transactions()
	if l != None:
		for item in l:
			trace(item, target)

def test1():
	t= str(datetime.now())[12:-7]
	f = open("auto/" + t + "_test.txt", "a")
	for i in range(6):
		time.sleep(30)
		load(f)
	f.close()


if __name__ == "__main__":
	test1()
   
