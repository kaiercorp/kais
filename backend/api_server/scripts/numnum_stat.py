import time
import pandas as pd
import numpy as np
import json
from scipy.stats import gaussian_kde
from pathlib import Path


def main():
    df_train = pd.read_csv(Path('train/train.csv'))
    df_test = pd.read_csv(Path('test/test.csv'))
    df_valid = pd.read_csv(Path('valid/valid.csv'))

    sample_size = 1000
    random_state = 42
    sample_train = df_train.sample(n=sample_size, random_state=random_state)
    sample_test = df_test.sample(n=sample_size, random_state=random_state)
    sample_valid = df_valid.sample(n=sample_size, random_state=random_state)


    # 숫자형 열 추출
    numerical_columns = df_train.select_dtypes(include="number").columns

    # JSON 구조 생성
    numerical_feature_statistics = {
        "feature_stat": [],
        "feature_detail": {},
        "box_plot": {},
        "pdf": {}
    }

    # 각 numerical column에 대해 계산
    for col in numerical_columns:
        # 각 컬럼의 통계값 계산
        col_values_train = df_train[col].dropna()
        col_values_test = df_test[col].dropna()
        col_values_valid = df_valid[col].dropna()
        
        # feature_stat
        len_train = len(col_values_train)
        len_test = len(col_values_test)
        len_valid = len(col_values_valid)
        feature_stat = {
            "feature": col,
            "test": len_test,  # 테스트 데이터 예시
            "total": len_test+len_train+len_valid,
            "train": len_train,  # 임의로 분리한 예시
            "valid": len_valid  # 임의로 설정한 예시
        }
        numerical_feature_statistics["feature_stat"].append(feature_stat)

        # feature_detail
        feature_detail = [
            {"category": "mean", "test":float(round(col_values_test.mean(),2)), "train": float(round(col_values_train.mean(),2)), "valid": float(round(col_values_valid.mean(),2))},
            {"category": "median", "test": float(round(col_values_test.median(),2)), "train": float(round(col_values_train.median(),2)), "valid": float(round(col_values_valid.median(),2))},
            {"category": "min", "test": float(round(col_values_test.min(),2)), "train": float(round(col_values_train.min(),2)), "valid": float(round(col_values_valid.min(),2))},
            {"category": "max", "test": float(round(col_values_test.max(),2)), "train": float(round(col_values_train.max(),2)), "valid": float(round(col_values_valid.max(),2))},
            {"category": "stdev", "test": float(round(col_values_test.std(),2)), "train": float(round(col_values_train.std(),2)), "valid": float(round(col_values_valid.std(),2))}
        ]
        numerical_feature_statistics["feature_detail"][col] = feature_detail

        # box_plot
        numerical_feature_statistics["box_plot"][col] = {
            "test": col_values_test.tolist(),  # 예시로 그대로 넣음, 실제 데이터를 기반으로 조정
            "train": col_values_train.tolist(),
            "valid": col_values_valid.tolist()
        }

        col_values_train = sample_train[col].dropna()
        col_values_test = sample_test[col].dropna()
        col_values_valid = sample_valid[col].dropna()

        # pdf 계산: Kernel Density Estimation (KDE) 사용
        try:
            kde = gaussian_kde(col_values_train)
            x_grid_train = col_values_train.sort_values().unique()
            x_range = np.linspace(min(x_grid_train), max(x_grid_train), 1000)
            y_grid_train = kde.pdf(x_range)
        except Exception as e:
            print(e)

        # pdf 계산: Kernel Density Estimation (KDE) 사용
        try:
            kde = gaussian_kde(col_values_test)
            x_grid_test = col_values_test.sort_values().unique()
            x_range = np.linspace(min(x_grid_test), max(x_grid_test), 1000)
            y_grid_test = kde.pdf(x_range)
        except Exception as e:
             print(e)



        # pdf 계산: Kernel Density Estimation (KDE) 사용
        try:
            kde = gaussian_kde(col_values_valid)
            x_grid_valid = col_values_valid.sort_values().unique()
            x_range = np.linspace(min(x_grid_valid), max(x_grid_valid), 1000)
            y_grid_valid = kde.pdf(x_range)
        except Exception as e:
             print(e)





        
        # pdf 데이터 추가
        numerical_feature_statistics["pdf"][col] = {
            "test": {"xData": x_grid_test.tolist(), "yData": y_grid_test.tolist()},
            "valid": {"xData": x_grid_valid.tolist(), "yData": y_grid_valid.tolist()},
            "train": {"xData": x_grid_train.tolist(), "yData": y_grid_train.tolist()}
        }

    # JSON 형식으로 변환
    with open("numnum_stat_output.json", "w", encoding="utf-8") as json_file:
        json.dump(numerical_feature_statistics, json_file, ensure_ascii=False, indent=4)

    # 출력
    #print(numerical_feature_statistics)


if __name__=="__main__":
    start_time = time.time()
    main()
    end_time = time.time()
    elapsed_time = end_time - start_time


    print(f"Time taken to compute the PDFs: {elapsed_time:.2f} seconds")

