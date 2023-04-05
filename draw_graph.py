import matplotlib.pyplot as plt
import numpy as np
import seaborn

# Use seaborn to style the plots (optional)


topo=1
mu=0.1
lambda_=0.05
genType=0
procType=0
NetType=1
#netSpeed=0.1
duration=100000
bufferSize=1
Cores=4
netSpeeds=[0.05,0.08,0.1,0.15,0.2,0.3,0.4]
#lambda_s=[0.005,0.008,0.01,0.02,0.03,0.04,0.05]#
#mus=[0.01,0.02,0.04,0.05,0.08,0.1,0.12,0.15,0.2]
def read_file(file_name):
    with open(file_name, 'r') as file:
        lines = file.readlines()
        
        # Read the first line, split by space, and convert each element to an integer
        numbers_line = list(map(float, lines[0].split()))
        
        # Read the remaining lines and convert each line to an integer
        single_numbers = [float(line.strip()) for line in lines[1:]]
        
    return numbers_line, single_numbers

# Read the file and get the numbers
attributes=[]
serviceTimes=[]
#for lambda_ in lambda_s:
for netSpeed in netSpeeds:
#for mu in mus:
    filename='./save/{}_{:.4f}_{:.4f}_{}_{}_{}_{:.4f}_{:.2f}_{}_{}.txt'.format(topo, mu,lambda_,genType,procType,NetType,netSpeed,duration,bufferSize,Cores)
    attribute,serviceTime=read_file(filename)

    indices_to_delete = [1, 3, 8]
    for index in sorted(indices_to_delete, reverse=True):
        del attribute[index]
    attributes.append(attribute)
    serviceTimes.append(serviceTime)

seaborn.set()
#################### draw pictures ##################################

lists =serviceTimes
item_names=["Count","AVG","50th", "90th","95th", "99th"]
# Function to calculate the CDF for a given list of integers
def calculate_cdf(data):
    sorted_data = np.sort(data)
    cumulative = np.cumsum(sorted_data)
    cdf = cumulative / cumulative[-1]
    return sorted_data, cdf

# Set up subplots
num_plots = len(lists)

fig, axs = plt.subplots(num_plots,1 , figsize=(10,num_plots * 6 ), sharey=True, sharex=True)

# Plot the CDF for each list
for i, data in enumerate(lists):
    x, y = calculate_cdf(data)
    axs[i].plot(x, y)
    axs[i].set_title("Changeable param:{}".format(netSpeeds[i]))
    #axs[i].set_xlabel("Value")
    if i == 0:
        axs[i].set_ylabel("CDF")
plt.savefig('./figures/cdf.png')
plt.clf()


########################### grouped bar graph #########################################

fig = plt.figure(figsize=(30, 5))
# Sample data: a list of lists, each containing 4 integers
data = attributes

# Number of groups and bars per group
num_groups = len(data)
num_bars = len(data[0])

# Set the bar width and positions
bar_width = 1/(num_bars+1)
group_positions = np.arange(num_groups)

# Create the figure and axes
fig, ax1 = plt.subplots()
ax2 = ax1.twinx()

# Plot the first bar for each group using the left y-axis
ax1.bar(group_positions, [group[0] for group in data], width=bar_width, label="Count")

# Plot the other 3 bars for each group using the right y-axis
for i in range(1, num_bars):
    ax2.bar(group_positions + i * bar_width, [group[i] for group in data], width=bar_width, label=item_names[i])

# Configure the x-axis
plt.xticks(group_positions + bar_width, netSpeeds)

# Add labels and legends
ax1.set_ylabel("Count")
ax2.set_ylabel("latency")
ax1.legend(loc="upper left")
ax2.legend(loc="upper right")

plt.savefig('./figures/grouped_bar.png')
plt.clf()