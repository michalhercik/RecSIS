from data_repository import DataRepository

class hercikmi():
    def recommend(self, user_id: str, user_study_info, blueprint) -> list[str]:
        df = DataRepository().get_degree_plan("NIPVS19B", 2020)
        print(df)
        return ["NPFL129", "NPRG031", "NDBI021", "NSWI177"]

class whatever:
    def recommend(self, user_id: str, user_study_info, blueprint) -> list[str]:
        return ["NPFL129", "NPRG031", "NDBI021", "NSWI177", "NDMI037", "NPFL129", "NPRG031", "NDBI021", "NSWI177", "NDMI037"]
    
class bla:
    def recommend(self, user_id: str, user_study_info, blueprint) -> list[str]:
        return ["NPFL129", "NPRG031", "NDBI021", "NSWI177", "NDMI037", "NPFL129", "NPRG031", "NDBI021", "NSWI177", "NDMI037"]
    

