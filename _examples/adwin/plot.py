import os
import pandas as pd
import argparse
import matplotlib.pyplot as plt

parser = argparse.ArgumentParser()
parser.add_argument("--base", default="example_log", help="Directory of log base")
args = parser.parse_args()

def newest(path):
    return os.path.join(path, sorted(os.listdir(path))[-1])

def load(path):
    log = newest(path)
    return pd.read_csv(log, names=['x', 'mu', 'sum', 'size'])

path = newest(args.base)
df = load(path)

fig, ax1 = plt.subplots()
ax2 = ax1.twinx()

ax1.plot(df.index, df['mu'], c='C0', label=r'$\mu$')
ax1.scatter(df.index, df['x'], s=1, alpha=0.3, c='C0')
ax1.set_ylim(0,1)
ax2.plot(df.index, df['size'], c='C1', label='Size of W')
ax1.legend(loc='upper left')
ax2.legend(loc='upper right')
plt.title('ADWIN')
plt.savefig('example_adwin.png')
