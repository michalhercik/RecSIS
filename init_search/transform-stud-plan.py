import pandas as pd

plans = pd.read_csv('./init_search/stud-plans.csv')
# plans = plans.apply(lambda x: x.to_dict(), axis=1)
plans = plans.reset_index().rename(columns={"index": "id"})

# plans.to_csv('./init_search/stud-plans-transformed.csv', index=False)
plans.to_json('./init_search/degree-plans-transformed.json', orient='records', lines=True)