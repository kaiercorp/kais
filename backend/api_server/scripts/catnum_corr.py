import scipy.stats as ss
import pandas as pd
import numpy as np
import argparse
import os
from sklearn import preprocessing

parser = argparse.ArgumentParser()
parser.add_argument('--data_path', type=str, help='data path')
parser.add_argument('--selected_cat_features', nargs='+', help='list of selected categorical features')
parser.add_argument('--selected_nume_features', nargs='+', help='list of selected numerical features')

args = parser.parse_args()


def pointbiserial_matrix(df, cat_columns, nume_columns):
    m = len(cat_columns)
    n = len(nume_columns)
    matrix = pd.DataFrame(np.full((m,n), np.nan), index=cat_columns, columns=nume_columns)
    le = preprocessing.LabelEncoder()

    for i in range(m):
        df[cat_columns[i]] = le.fit_transform(df[cat_columns[i]])
        for j in range(n):  
            matrix.iloc[i, j] = ss.pointbiserialr(df[cat_columns[i]], df[nume_columns[j]]).statistic
    
    return matrix


df_train = pd.read_csv(os.path.join(args.data_path, 'train/train.csv'))
df_valid = pd.read_csv(os.path.join(args.data_path, 'valid/valid.csv'))
df_test = pd.read_csv(os.path.join(args.data_path, 'test/test.csv'))

cleaned_nume_features = [feature for feature in args.selected_nume_features]
cleaned_cat_features = [feature for feature in args.selected_cat_features]

heatmap_train = pointbiserial_matrix(df_train, cleaned_cat_features, cleaned_nume_features).to_dict(orient='index')
heatmap_valid = pointbiserial_matrix(df_valid, cleaned_cat_features, cleaned_nume_features).to_dict(orient='index')
heatmap_test = pointbiserial_matrix(df_test, cleaned_cat_features, cleaned_nume_features).to_dict(orient='index')

info_dict = {
    'feature' : {
        'categorical': cleaned_cat_features,
        'numerical': cleaned_nume_features
    },
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