import pandas as pd
import matplotlib.pyplot as plt
import pandas as pd
import csv
import argparse
import warnings

warnings.filterwarnings("ignore")

argsParser = argparse.ArgumentParser(
    description='Checking tool for K-Means algorithm results.')

argsParser.add_argument('dataset', metavar='dataset', type=str, nargs=1,
                        help='Path of original CSV dataset')
argsParser.add_argument('centroids', metavar='centroids', type=str, nargs=1,
                        help='Path of CSV file with K-Means results')

args = argsParser.parse_args()
datasetFilename = args.dataset[0]
resultsFilename = args.centroids[0]

with open(datasetFilename, 'r') as csvfile:
    dialectDataset = csv.Sniffer().sniff(csvfile.readline())

with open(resultsFilename, 'r') as csvfile:
    dialactCentroids = csv.Sniffer().sniff(csvfile.readline())

centroids = pd.read_csv(
    resultsFilename, delimiter=dialactCentroids.delimiter, header=None)
dataset = pd.read_csv(datasetFilename,
                      delimiter=dialectDataset.delimiter, header=None)


plt.scatter(dataset.iloc[:, 0], dataset.iloc[:, 1], c='#3f81ba',
            label='Dataset')
plt.scatter(x=centroids.iloc[:, 0], y=centroids.iloc[:, 1], c='#d94638',
            label='K-Means centroids')
plt.xlabel("x")
plt.ylabel("y")
plt.legend()

output_fig_path = str(datasetFilename).split(sep='/')

# output_fig_path will now contain only the parents dirs
filename = output_fig_path.pop().replace('.csv', '')

plt.savefig(
    f"{'/'.join(output_fig_path)}/{filename.replace('dataset', 'kmeans_results')}.png")
plt.show()
