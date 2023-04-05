import threading
import os
import numpy as np
# gstreamer_pipeline returns a GStreamer pipeline for capturing from the CSI camera
# Defaults to 1280x720 @ 60fps
# Flip the image by setting the flip_method (most common values: 0 and 2)
# display_width and display_height determine the size of the window on the screen
import argparse

exitFlag = 0

topo=1
mu=0.1
lambda_=0.05
genType=0
procType=0
NetType=1
#netSpeed=0.1
duration=100000
#lambda_s=[0.005,0.008,0.01,0.02,0.03,0.04,0.05]#
netSpeeds=[0.05,0.08,0.1,0.15,0.2,0.3,0.4]
#mus=[0.01,0.02,0.04,0.05,0.08,0.1,0.12,0.15,0.2]
class ParaClass:
    def __init__(self, topo, mu,lambda_,genType,procType,NetType,netSpeed,duration):
        self.topo = topo
        self.mu = mu
        self.lambda_ = lambda_
        self.genType = genType
        self.procType = procType
        self.NetType = NetType
        self.netSpeed = netSpeed
        self.duration = duration



class Thread_shell (threading.Thread):
    def __init__(self, threadID, name,paras):
        threading.Thread.__init__(self)
        self.threadID = threadID
        self.name = name
        self.paras=paras
    def run(self):
        #print("--lambda",lambda_)
        os.system("schedsim --topo={} --mu={} --lambda={} --genType={} --procType={}  --NetType={} --netSpeed={} --duration {}"\
                  .format(self.paras.topo,self.paras.mu,self.paras.lambda_,self.paras.genType,self.paras.procType,self.paras.NetType,self.paras.netSpeed,self.paras.duration))   



thread_id=1
model_names=['avmnist','mmimdb','sarcasm','mosiclass']#
sub_modules={'avmnist': ["01","10"], 
        'mmimdb': ["01","10"], 
        'sarcasm': ["012","021","102","120","201","210"], 
        'mosiclass':["012","021","102","120","201","210"]  }


total_count=0
Thread_pool=[]
#for lambda_ in lambda_s:
for netSpeed in netSpeeds:  
#for mu in mus:
    paras=ParaClass(topo, mu,lambda_,genType,procType,NetType,netSpeed,duration)
    Thread_pool.append(Thread_shell(thread_id, "Thread_measure", paras))
    thread_id+=1
    total_count+=1
    if total_count>100:
        total_count=0
        for thread in Thread_pool:
            thread.start()
        for thread in Thread_pool:
            thread.join()
        Thread_pool=[]

for thread in Thread_pool:
    thread.start()
for thread in Thread_pool:
    thread.join()
Thread_pool=[] 
 