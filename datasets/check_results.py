import pandas as pd
import matplotlib.pyplot as plt
import pandas as pd
import csv
import argparse
import warnings
import os

warnings.filterwarnings("ignore")

argsParser = argparse.ArgumentParser(
    description='Checking tool for K-Means algorithm results.')

argsParser.add_argument('dataset', metavar='dataset', type=str, nargs=1,
                        help='Path of original CSV dataset')

args = argsParser.parse_args()
datasetFilename = args.dataset[0]

dirPath = os.path.abspath(__file__).replace("check_results.py", str(
    datasetFilename).replace(".csv", "/"))
datasetPath = dirPath + datasetFilename

resultsPath = dirPath + str(datasetFilename).replace("dataset", "centroids")

with open(datasetPath, 'r') as csvfile:
    dialectDataset = csv.Sniffer().sniff(csvfile.readline())

with open(resultsPath, 'r') as csvfile:
    dialactCentroids = csv.Sniffer().sniff(csvfile.readline())

centroids = pd.read_csv(
    resultsPath, delimiter=dialactCentroids.delimiter, header=None)
dataset = pd.read_csv(datasetPath,
                      delimiter=dialectDataset.delimiter, header=None)


plt.scatter(dataset.iloc[:, 0], dataset.iloc[:, 1], c='#3f81ba',
            label='Dataset')
plt.scatter(x=centroids.iloc[:, 0], y=centroids.iloc[:, 1], c='#d94638', edgecolors='black',
            label='K-Means centroids')
plt.xlabel("x")
plt.ylabel("y")
plt.legend()

output_fig_path = resultsPath.replace(
    "centroids", "kmeans_reuslts").replace(".csv", ".png")
plt.savefig(output_fig_path)
plt.show()
