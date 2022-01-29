import json
import requests
import time
import os
from os.path import exists

geth = 'http://127.0.0.1:8555'
parity = 'http://127.0.0.1:9545'

def adjust_values(g_trace, p_trace, tx): 
	for (g_item, p_item) in zip(g_trace['ops'], p_trace['ops']):
		p_item['Op'] = g_item['Op'] 
		g_item['tx'] = str(tx)
		p_item['tx'] = str(tx)
		if g_item['Op'] == "DELEGATECALL":
			del g_item['ex']['used'] 
			del p_item['ex']['used'] 
		if g_item['sub'] is not None: 
			adjust_values(g_item['sub'], p_item['sub'], tx) 
	return g_trace, p_trace 

def prepare_reserve(g_rsv, count): 
	count += 1
	for item in g_rsv['ops']: 
		item['depth'] = str(count)
		if item['Op'] == "DELEGATECALL":
			del item['ex']['used']  
		if item['sub'] is not None: 
			parepare_reserve(item['sub'], count) 
	return g_rsv

def index(dct, count):  
	count += 1  
	for item in dct['ops']:  
	    item['depth'] = str(count)  
	    if item['sub'] is not None:  
	    	index(item['sub'], count)

def get_transactions():
	data = {"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", False],"id":1}
	r = requests.post(geth, json=data).json()
	bk = r['result']['number']
	txs = r['result']['transactions']
	time.sleep(10)
	print(f"{bk}")
	return txs

def get_lists(tx):
	data = {"jsonrpc":"2.0","method":"trace_replayTransaction","params":[f"{tx}", ["vmTrace"]],"id":1}	
	try:
		time.sleep(5)
		g = requests.post(geth, json=data).json()
		# reserve = requests.post(geth, json=data).json()
		p = requests.post(parity, json=data, headers={"host": "localhost:8545"}).json()
	except (Exception) as e:
		print(tx, e, "initial parsing")
		return
	try:	
		g_tr, p_tr = g['result']['vmTrace'], p['result']['vmTrace'] 
		# g_reserve = reserve['result']['vmTrace']
	except Exception as e:
		print(tx, e, "getting geth, parity vmTrace", p)
		return
	try:
		if len(g_tr['ops']) != len(p_tr['ops']):
			print(e, "divergnet length")
			return
	except (Exception) as e:
		print(e, "checking length")
		return
	try:
		if g_tr is not None:
			gt, pt = adjust_values(g_tr, p_tr, tx)
		# grsv = parepare_reserve(g_reserve, 0)
	except Exception as e:
		print(tx, e, "recursive function")
		return
	try:
		package = (gt, pt)
	except Exception as e:
		print(tx, e, "making package")
	return package

def compare_and_report(pkg, g_target, p_target):
	for i, (g_val, p_val) in enumerate(zip(pkg[0]['ops'], pkg[1]['ops'])):
		if g_val != p_val:
			g_target.write(f'{json.dumps(g_val, sort_keys=True)}\n')
			p_target.write(f'{json.dumps(p_val, sort_keys=True)}\n')
			if g_val['sub'] is not None:
				compare_and_report((g_val['sub'], p_val['sub']), g_target, p_target)

def load(target_0, target_1):
	t = get_transactions()
	if t != None:
		print(len(t))
		print(t)
		for tx in t:
			pkg = get_lists(tx)
			compare_and_report(pkg, target_0, target_1)

def test1():
	f_0 = open("geth.txt", "a")
	f_1 = open("parity.txt", "a")
	load(f_0, f_1)
	f_0.close()
	f_1.close()


if __name__ == "__main__":
	test1()


# def trace(tx, target):
# 	data = {"jsonrpc":"2.0","method":"trace_replayTransaction","params":[f"{tx}", ["vmTrace"]],"id":1}	
# 	try:
# 		g = requests.post(geth, json=data).json()
# 		p = requests.post(parity, json=data, headers={"host": "localhost:8545"}).json()
# 	except (Exception) as e:
# 		print(tx, e, "initial parsing")
# 		return
# 	try:	
# 		g_tr, p_tr = g['result']['vmTrace'], p['result']['vmTrace'] 
# 	except Exception as e:
# 		print(tx, e, "getting geth, parity vmTrace", p)
# 		return
# 	try:
# 		gt, pt = delete_value(g_tr, p_tr)
# 	except Exception as e:
# 		print(tx, e, "recursive function")
# 		return
# 	try:	
# 		g_results = (g['result']['output'], g['result']['vmTrace']['code'], gt['ops'])
# 		p_results = (p['result']['output'], p['result']['vmTrace']['code'], pt['ops'])
# 	except (Exception) as e:
# 		print(tx, e, "creating tuples")
# 		return
# 	try:
# 		if g_results[0] != p_results[0]:
# 			target.write(f"{str(tx)}\n")
# 	except (Exception) as e:
# 		print(tx, e, "checking output")
# 		return
# 	try:
# 		if g_results[1] != p_results[1]:
# 			target.write(f"{str(tx)}\n")
# 	except (Exception) as e:
# 		print(tx, e, "checking code")
# 		return 
# 	try:
# 		if len(g_results[2]) != len(p_results[2]):
# 			target.write(f"{str(tx)}\n")
# 	except (Exception) as e:
# 		print("geth parity divergent length", e)
# 		return
# 	try:
# 		for (g_val, p_val) in zip(g_results[2], p_results[2]):
# 			if g_val != p_val:
# 				target.write(f'{str(tx)}\n')
# 				# target.write(f'{(json.dumps(g_val), json.dumps(p_val))}\n')
# 			# divergent_g_ops.append(g_reserve[i])
# 			# divergent_p_ops.append(p_val)
# 			# return (divergent_g_ops, divergent_p_ops)
# 	except (Exception) as e:
# 		print(tx, e, "checking values")
# 		return
# 	return target

# def load(target_0, target_1):
# 	t = get_transactions()
# 	if t != None:
# 		print(len(t))
# 		for tx in t:
# 			trace(tx, target_0, target_1)

# def test1():
# 	txns = 'transactions'
# 	ops = 'ops'
# 	files = (txns, ops)
# 	for file in files:
# 		if exists(f'./{file}.txt') == True:
# 			os.remove(f'./{file}.txt')
# 			f_0 = open(f"{file_name}.txt", "a")
# 			f_1 = open(f"{file_name}.txt", "a")
# 	# for i in range(12):
# 	# 	time.sleep(5)
# 	# 	load(f)
# 	load(f)
# 	f.close()


# if __name__ == "__main__":
# 	test1()

# def index(lst, count):  
# 	count += 1  
# 	for item in lst:  
# 	    item['depth'] = str(count)  
# 	    if len(item['depth']) > 0:  
# 	    	index(item['depth'], count) 




# 	def isolate_ops(g_trace, p_trace, g_trace2):
# 		for (g_item, p_item) in zip(g_trace['ops'], p_trace['ops']): 
# 		if g_item['Op'] == "DELEGATECALL":
# 			del g_item['ex']['used'] 
# 			del p_item['ex']['used'] 
# 		del g_item['Op'] 
# 		if g_item['sub'] is not None: 
# 			delete_value(g_item['sub'], p_item['sub']) 
# 	return g_trace, p_trace 


# def isolate_ops(g_trace, p_trace):
# 	for i, (g_item, p_item) in enumerate(zip(g_trace, p_trace)):

