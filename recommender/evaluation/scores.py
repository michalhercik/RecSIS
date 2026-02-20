import numpy as np
import pandas as pd


def eval_vector(y_true, y_pred) -> np.ndarray:
    y_true_set = set(y_true)
    return np.array([1 if code in y_true_set else 0 for code in y_pred], dtype=np.int64)


def precision_at_k(y_true, y_pred, k):
    y_pred = y_pred[:k]
    if len(y_pred) == 0:
        return 0.0
    eval_vec = eval_vector(y_true, y_pred)
    return eval_vec.sum() / k


def recall_at_k(y_true, y_pred, k):
    if len(y_true) == 0:
        return 0.0
    y_pred = y_pred[:k]
    eval_vec = eval_vector(y_true, y_pred)
    return eval_vec.sum() / len(y_true)


def avg_prec_at_k(y_true, y_pred, k, normalize=False):
    result = 0
    counter = 0
    iterate_over = min(len(y_pred), k)
    # if iterate_over == len(y_pred):
    #     print("WARNING: len(y_pred) <", k)
    for i in range(iterate_over):
        if y_pred[i] in y_true:
            counter += 1
            result += counter / (i + 1)
    divisor = counter
    if normalize:
        divisor = min(len(y_true), k)
        print("divisor", divisor)
    return result / divisor if divisor > 0 else 0
