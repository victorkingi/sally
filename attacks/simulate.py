import subprocess
import json
import time
import _thread
import time

def run_simulation(from_, prKey, to, name):
    process = subprocess.Popen(['node', 'tx_simulate.js', from_, prKey, to],
                                   stdout=subprocess.PIPE,
                                   universal_newlines=True)

    while True:
        start_time = time.time()
        output = process.stdout.readline()
        val = output.strip()
        if val:
            print(val)
            total = time.time() - start_time
            a_lock = _thread.allocate_lock()
            with a_lock:
                with open('block_time.txt', 'a') as file_block_time:
                    file_block_time.write(name+':  '+str(total)+'\n')

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
    i = 0
    for key in data:
        name = "Thread-"+str(i)
        try:
            print("%s: %s" % ( name, time.ctime(time.time())))
            _thread.start_new_thread( run_simulation, (key, data[key], '0x9acd24d607cb3d561e37f47837b651b9eba7d729', name))
            i += 1
        except BaseException as err:
            print(f"Unexpected {err=}, {type(err)=}")
            raise

    while 1:
        pass
