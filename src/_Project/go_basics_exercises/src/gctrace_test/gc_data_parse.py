from typing import Tuple

import sys
sys.path.append("/home/urarao/.local/lib/python3.10/site-packages")
sys.path.append("/home/urarao/.local/lib/python3.10/site-packages/regex")
sys.path.append("/home/urarao/.local/lib/python3.10/site-packages/matplotlib")

import regex as reg
import matplotlib.pyplot as plt


# gc 3 @0.015s 10%: 1.6+0.57+0 ms clock, 13+0/1.0/0.50+0 ms cpu, 0->16->0 MB, 16 MB goal, 0 MB stacks, 0 MB globals, 8 P (forced)
# gc 6 @0.022s 38%: 11+0.63+0 ms clock, 89+1.2/1.2/0+0 ms cpu, 16->80->64 MB, 80 MB goal, 0 MB stacks, 0 MB globals, 8 P
gc_log_pattern = reg.compile(r'gc[^\d]+(?P<count>[\d]+)[^\+]+\+(?P<ph_time>[^\+]+)[^/]+/(?P<cpu_time>[^/]+)[^,]+,\s+(?P<bfgc_mem>\d+)\->[^>]+>(?P<aftgc_mem>\d+)[^\(]+(?P<force_gc_flag>\(forced\))?')


def parseline(line:str)->Tuple[bool, int, float, float, int, int, bool]:
    badlineRes =( False, 0, 0, 0, 0, 0, 0)
    if not line.startswith('gc'):
        return badlineRes
    match = gc_log_pattern.search(line)
    if not match:
        return badlineRes
    
    return True, int(match.group("count"))-1, float(match.group("ph_time")), float(match.group("cpu_time")), int(match.group("bfgc_mem")), int(match.group("aftgc_mem")), match.group("force_gc_flag") == "(forced)"


idx = []
ph_times = []
cpu_times = []
bfgc_mems = []
aftgc_mems = []
force_gc_flags = []

def load_file(gc_log_f:str):
    global idx, ph_times, cpu_times, bfgc_mems, aftgc_mems, force_gc_flags
    idx = []
    ph_times = []
    cpu_times = []
    bfgc_mems = []
    aftgc_mems = []
    force_gc_flags = []
    i = 0
    with open(gc_log_f) as f:
        while True:
            
            l = f.readline()
            if l == '':
                break
            if i % 10 == 0:
                print(f"parse line:{i}\n")
            ok, count, ph_time, cpu_time, bfgc_mem, aftgc_mem, force_gc_flag = parseline(l)
            if not ok:
                continue
            idx.append(count)
            ph_times.append(ph_time)
            cpu_times.append(cpu_time)
            bfgc_mems.append(bfgc_mem)
            aftgc_mems.append(aftgc_mem)
            force_gc_flags.append(force_gc_flag)

            i+=1

gc1 = "/home/urarao/git_playground/chaidaxuan/go_basics_exercises/src/gctrace_test/gc1.log"
gc2 = "/home/urarao/git_playground/chaidaxuan/go_basics_exercises/src/gctrace_test/gc2.log"

# CPU耗时
# load_file(gc1)
# plt.plot(idx, ph_times, 'r-', idx, cpu_times, 'b-')
# # plt.show()
# plt.savefig('/home/urarao/git_playground/chaidaxuan/go_basics_exercises/src/gctrace_test/CPU_time2.png')

# GC内存占用
# load_file(gc1)
# plt.plot(idx, ph_times, 'r-', idx, cpu_times, 'b-')
# bar_width = 0.35
# plt.bar(idx, bfgc_mems, bar_width, label="gc前内存占用(M)", color="b")
# plt.bar(idx, aftgc_mems, bar_width, label="gc后内存占用(M)", color="r")
# plt.legend()
# plt.savefig('/home/urarao/git_playground/chaidaxuan/go_basics_exercises/src/gctrace_test/GC_memory.png')


