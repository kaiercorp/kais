import argparse
import os
import scipy.stats as ss
import numpy as np
import pandas as pd


parser = argparse.ArgumentParser()
parser.add_argument('--data_path', type=str, help='data path')
parser.add_argument('--selected_cat_features', nargs='+', help='list of selected categorical features')

args = parser.parse_args()

def cramers_v(x, y):
    contingency_table = pd.crosstab(x, y)
    chi2 = ss.chi2_contingency(contingency_table)[0]
    n = contingency_table.sum().sum()
    r, k = contingency_table.shape
    return np.sqrt(chi2 / (n * (min(r-1, k-1))))

def cramers_v_matrix(df, columns):
    n = len(columns)
    matrix = pd.DataFrame(np.full((n,n), np.nan), index=columns, columns=columns)
    
    for i in range(n):
        for j in range(i+1):
            if i == j:
                v = 1.0
            else:    
                v = cramers_v(df[columns[i]], df[columns[j]])
            matrix.iloc[i, j] = v
    
    return matrix

df_train = pd.read_csv(os.path.join(args.data_path, 'train/train.csv'))
df_valid = pd.read_csv(os.path.join(args.data_path, 'valid/valid.csv'))
df_test = pd.read_csv(os.path.join(args.data_path, 'test/test.csv'))

heatmap_train = cramers_v_matrix(df_train, args.selected_cat_features).to_dict(orient='index')
heatmap_valid = cramers_v_matrix(df_valid, args.selected_cat_features).to_dict(orient='index')
heatmap_test = cramers_v_matrix(df_test, args.selected_cat_features).to_dict(orient='index')

info_dict = {
    'feature' : args.selected_cat_features,
    'heatmap' : {
        'train' : {},
        'valid' : {},
        'test' : {}
    }
}

for k, v in heatmap_train.items():
    info_dict['heatmap']['train'][k] = v

for k, v in heatmap_valid.items():
    info_dict['heatmap']['valid'][k] = v

for k, v in heatmap_test.items():
    info_dict['heatmap']['test'][k] = v


print(info_dict)