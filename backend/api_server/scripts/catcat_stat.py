from typing import Dict, List
import pandas as pd
import json
from collections import Counter
from pathlib import Path


def load_data(train_path, valid_path, test_path):
    train_df = pd.read_csv(Path(train_path))
    valid_df = pd.read_csv(Path(valid_path))
    test_df = pd.read_csv(Path(test_path))
    return train_df, valid_df, test_df


def get_feature_statistics(df: pd.DataFrame, column):
    return (Counter(df[column]))


def get_max_min_mode(counts: Counter):
    max_mode = counts.most_common(1)[0][0]
    min_count = min(counts.values())
    min_modes = [k for k, v in counts.items() if v == min_count]
    return max_mode, min_modes


def generate_statistics(train_df: pd.DataFrame, valid_df: pd.DataFrame,
                        test_df: pd.DataFrame, columns):
    df = pd.concat([train_df, valid_df, test_df])
    stats = {
        "feature_stat": [],
        "feature_detail": {

        }
    }

    for column in columns:
        train_counts = get_feature_statistics(train_df, column)
        valid_counts = get_feature_statistics(valid_df, column)
        test_counts = get_feature_statistics(test_df, column)
        total_counts = get_feature_statistics(df, column)

        test_max_mode, test_min_modes = get_max_min_mode(test_counts)
        train_max_mode, train_min_modes = get_max_min_mode(train_counts)
        valid_max_mode, valid_min_modes = get_max_min_mode(valid_counts)
        feature_stat: List[Dict] = stats["feature_stat"]
        feature_stat.append({
            "feature": column,
            "test_maxMode_size": test_max_mode,
            "test_minMode_size": ",".join(test_min_modes),
            "test_size": len(test_df),
            "total": len(df),
            "train_maxMode_size": train_max_mode,
            "train_minMode_size": ",".join(train_min_modes),
            "train_size": len(train_df),
            "valid_maxMode_size": valid_max_mode,
            "valid_minMode_size": ",".join(valid_min_modes),
            "valid_size": len(valid_df)
        })
        stats["feature_detail"][column] = [
            {
                "feature": key,
                "test": test_counts.get(key, 0),
                "total": total_counts.get(key, 0),
                "train": train_counts.get(key, 0),
                "valid": valid_counts.get(key, 0)
            }
            for key in total_counts.keys()
        ]

    return stats


def main():
    train_path = "./train/train`.csv"  # 실제 파일명을 사용하세요
    valid_path = "./valid/valid`.csv"
    test_path = "./test/test`.csv"
    train_df, valid_df, test_df = load_data(train_path, valid_path, test_path)
    cleaned_cat_features = ['Defect Type']
    stats = generate_statistics(
        train_df, valid_df, test_df, columns=cleaned_cat_features)
    with open("catcat_stat_output.json", "w", encoding="utf-8") as json_file:
        json.dump(stats, json_file, indent=4, ensure_ascii=False)


if __name__ == "__main__":
    main()
