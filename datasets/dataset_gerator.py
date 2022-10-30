from pathlib import Path
import os
from sklearn.datasets import make_blobs
import numpy as np
from pandas.plotting._matplotlib import scatter_matrix
from matplotlib import pyplot
from pandas import DataFrame
import argparse
import warnings

# TODO inserire riferimento a documentazione make_blobs sulla documentazione

argsParser = argparse.ArgumentParser(
    description='d-Dimensional Dataset Generator formed by K cluster.')

argsParser.add_argument('dimension', metavar='dim', type=int, nargs=1,
                        help='Number of features')
argsParser.add_argument('instances', metavar='N', type=int, nargs=1,
                        help='Number of instances')
argsParser.add_argument('centers', metavar='K', type=int, nargs=1,
                        help='Number of clusters')
argsParser.add_argument('std', metavar='std', type=float, nargs=1,
                        help='Standard deviation of the clusters')

args = argsParser.parse_args()
dimension = args.dimension[0]
instances = args.instances[0]
centers = args.centers[0]
std = args.std[0]

dir_path = os.path.dirname(os.path.realpath(__file__))
center_box = (-20.0, 20.0)

points, y = make_blobs(
    n_samples=instances, centers=centers, n_features=dimension, shuffle=True, center_box=center_box, cluster_std=std)

filename = f"dataset_{dimension}d_{centers}cluster_{instances}samples"
output_file = Path(f"{dir_path}/{filename}/{filename}.csv")
output_fig = Path(f"{dir_path}/{filename}/{filename}.png")
output_file.parent.mkdir(exist_ok=True, parents=True)

with open(output_file, "w") as file:
    for point in points:
        for value in range(dimension):
            if value == (dimension - 1):
                file.write(str(round(point[value], 4)))
            else:
                file.write(str(round(point[value], 4)) + ",")
        file.write("\n")

data = np.array(points)

# clustering plot
warnings.filterwarnings("ignore")
df = DataFrame(dict(x=points[:, 0], y=points[:, 1], label=y))
colors = {0: 'red', 1: 'blue', 2: 'green',
          3: 'black', 4: 'purple', 5: 'pink', 6: 'orange'}
fig, ax = pyplot.subplots()
grouped = df.groupby('label')
for key, group in grouped:
    group.plot(ax=ax, kind='scatter', x='x',
               y='y', label=key, color=colors[key])
pyplot.savefig(output_fig)
pyplot.show()
