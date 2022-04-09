import js2py
import _thread
import random
import time
import json
import subprocess

# Define a function for the thread
def run_js_script( threadName, addr, prKey):
  print("%s: %s" % ( threadName, time.ctime(time.time()) ))
  print(" ".join(['node', 'tx_simulate.js', addr, prKey, '0x9acd24d607cb3d561e37f47837b651b9eba7d729']))
  process = subprocess.Popen(['node', 'tx_simulate.js', addr, prKey, '0x9acd24d607cb3d561e37f47837b651b9eba7d729'],
                     stdout=subprocess.PIPE, 
                     stderr=subprocess.PIPE,
                     universal_newlines=True)
  
  while True:
    output = process.stdout.readline()
    print(output.strip())
    # Do something else
    return_code = process.poll()
    if return_code is not None:
        print('RETURN CODE', return_code)
        # Process has finished, read rest of the output 
        for output in process.stdout.readlines():
            print(output.strip())
        break


if __name__ == '__main__':
  f = open('accounts.json', 'r')
  data = json.load(f)
  f.close()
  run_js_script('hello', '0x1e90BBca26fd37DD5433a03bCcB2beCF69B2b05a', '0x8c838d4243ffaf7f59a2726c7e9097762b6270033bb9411aec1aa05e8dbd231b')
  '''for key in data:
    try:
      thread_name = "Thread-"+key[:6]
      _thread.start_new_thread( run_js_script, (thread_name, key, data[key]) )
    except:
      print("Error starting thread") '''
  
