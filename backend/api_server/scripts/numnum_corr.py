import pandas as pd
import argparse
from pathlib import Path

def create_corr_matrix(df: pd.DataFrame)-> pd.DataFrame:
    # 수치형 데이터만 사용하여 상관계수 계산
    correlation_matrix = df.select_dtypes(include='number')
    correlation_matrix = correlation_matrix.corr().round(2)

    return correlation_matrix


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--data_path', type=str, help='data path')

    args = parser.parse_args()
    
    df_train = pd.read_csv(Path(args.data_path)/ 'train/train.csv')
    df_valid = pd.read_csv(Path(args.data_path)/ 'valid/valid.csv')
    df_test = pd.read_csv(Path(args.data_path)/ 'test/test.csv')

    heatmap_train: pd.DataFrame = create_corr_matrix(df_train)
    heatmap_valid: pd.DataFrame = create_corr_matrix(df_valid)
    heatmap_test: pd.DataFrame = create_corr_matrix(df_test)


    # DataFrame의 컬럼명을 feature로 설정
    features = heatmap_train.columns.tolist()  # 'Defect Type' 제외한 컬럼명

    # JSON 구조 생성
    numerical_heatmap = {
        "feature": features,
        "correlation": {
            "test": heatmap_test.to_dict(),
            "train": heatmap_train.to_dict(),  # 예시로 동일한 상관계수를 사용, 실제 데이터가 분리되어 있다면 그에 맞게 처리
            "valid": heatmap_valid.to_dict()   # 예시로 동일한 상관계수를 사용
        }
    }

    print(numerical_heatmap)
    # JSON 형식으로 출력
    #json_data = json.dumps(numerical_heatmap, ensure_ascii=False, indent=4)


if __name__=="__main__":
    main()
