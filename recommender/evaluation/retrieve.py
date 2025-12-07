import csv
import os
from datetime import datetime

import pandas as pd
import psycopg2

OUT = "dataset.csv"
OUT_PATH = "dataset"


def main():
    expr = create_expr()
    password = os.getenv("RECSIS_RECOMMENDER_DB_PASS")
    conn = psycopg2.connect(
        dbname="recsis",
        user="recommender",
        host="localhost",
        password=password,
        port=5432,
    )
    cursor = conn.cursor()
    cursor.execute(expr)
    rows = cursor.fetchall()
    if cursor.description is None:
        raise ValueError("Query is empty")
    colnames = [desc[0] for desc in cursor.description]
    with open(OUT, "w", newline="", encoding="utf-8") as f:
        writer = csv.writer(f)
        writer.writerow(colnames)
        writer.writerows(rows)
    cursor.close()
    conn.close()

    df = pd.read_csv(OUT)
    print(df.shape)
    print(df.head())


def create_expr():
    expr = WITH_EXPR
    expr += "SELECT * FROM ("
    for type in range(1, 3):
        for year in range(1, 7):
            for semester in range(1, 3):
                ptr = type * 100 + year * 10 + semester
                next = 0
                if semester == 1:
                    next = ptr + 1
                elif year < 6:
                    next = type * 100 + (year + 1) * 10 + 1
                elif type < 2:
                    next = (type + 1) * 100 + 10 + 1
                else:
                    break
                expr += SELECT_TEMPL.format(ptr, next)
                expr += UNION_STMT
    expr = expr[: -len(UNION_STMT)]
    expr += ")"
    expr += "WHERE array_length(relevant, 1) > 0"
    expr += "\n"
    expr += "ORDER BY soident, sdruh, zroc, zsem"
    return expr


def out_file_path(name):
    datetime_tag = datetime.now().strftime("%y%m%d-%H%M%S")
    return f"{OUT_PATH}/{name}-{datetime_tag}.csv"


class Dataset1:
    """
    - relevant: all courses that student will enroll into
    - next: courses in next semester
    - replace non-searchable courses for searchable equivalent (interchangable) courses
    """

    def __init__(self):
        self.name = "Dataset1"
        self.with_expr = """
        WITH eqpovinn AS (
            SELECT DISTINCT
                p1.povinn p1, p2.povinn p2
            FROM preq preq1
            INNER JOIN preq preq2 ON preq1.povinn = preq2.reqpovinn AND preq1.reqpovinn = preq2.povinn AND preq1.reqtyp = preq2.reqtyp
            LEFT JOIN povinn p1 ON preq1.povinn = p1.povinn
            LEFT JOIN povinn p2 ON preq2.povinn = p2.povinn
            LEFT JOIN searchable_povinn ps1 ON p1.povinn = ps1.povinn
            LEFT JOIN searchable_povinn ps2 ON p2.povinn = ps2.povinn
            WHERE preq1.reqtyp = 'Z'
            AND p2.pvyucovan = 'V'
            AND ps1.povinn IS NULL
            AND ps2.povinn IS NOT NULL
        ),
        data AS (
            SELECT
                soident, zident, sdruh, zroc, zsem, zpovinn,
                (zroc::INT * 10 + zsem::INT) + (case when sdruh='B' then 100 when sdruh='N' then 200 end) ord
            FROM studium s
            LEFT JOIN zkous z ON s.sident = z.zident
            LEFT JOIN povinn p ON z.zpovinn = p.povinn
            WHERE s.sobor like 'I%'
            AND sstav = 'A'
            AND z.zsplcelk = 'S'
            AND p.pgarant != '32-STUD'
        )
        """
        self.select_templ = """
        SELECT x.soident, sdruh, x.zroc, x.zsem,
        x.finished, next_semester.next_semester, relevant.relevant,
        x.interchange_for, x.incompatible_with
        FROM (
            SELECT
                soident, max(sdruh) sdruh, substr('{0}', 2, 1)::INT zroc, substr('{0}', 3, 1)::INT zsem,
                array_agg(zpovinn) finished,
                array_agg(zpreq.povinn) FILTER (WHERE zpreq.povinn IS NOT NULL) interchange_for,
                array_agg(npreq.povinn) FILTER (WHERE npreq.povinn IS NOT NULL) incompatible_with
            FROM data
            LEFT JOIN preq zpreq ON data.zpovinn = zpreq.reqpovinn AND zpreq.reqtyp = 'Z'
            LEFT JOIN preq npreq ON data.zpovinn = npreq.reqpovinn AND npreq.reqtyp = 'N'
            WHERE ord <= {0}
            GROUP BY soident
        ) x
        LEFT JOIN (
            SELECT soident, array_agg(CASE WHEN eq.p1 IS NULL THEN zpovinn ELSE eq.p2 END) next_semester
            FROM data d
            LEFT JOIN eqpovinn eq ON d.zpovinn = eq.p1
            LEFT JOIN searchable_povinn ps ON d.zpovinn = ps.povinn
            WHERE ord = {1}
            AND (ps.povinn IS NOT NULL OR eq.p1 IS NOT NULL)
            GROUP BY soident
        ) next_semester ON next_semester.soident = x.soident
        LEFT JOIN (
            SELECT soident, array_agg(CASE WHEN eq.p1 IS NULL THEN zpovinn ELSE eq.p2 END) relevant
            FROM data d
            LEFT JOIN eqpovinn eq ON d.zpovinn = eq.p1
            LEFT JOIN searchable_povinn ps ON d.zpovinn = ps.povinn
            WHERE ord > {0}
            AND (ps.povinn IS NOT NULL OR eq.p1 IS NOT NULL)
            GROUP BY soident
        ) relevant ON relevant.soident = x.soident
        """
        self.misc_expr = """
        WHERE array_length(relevant, 1) > 0
        ORDER BY soident, sdruh, zroc, zsem
        """


UNION_STMT = "UNION"

WITH_EXPR = """
WITH eqpovinn AS (
    SELECT DISTINCT
        p1.povinn p1, p2.povinn p2
    FROM preq preq1
    INNER JOIN preq preq2 ON preq1.povinn = preq2.reqpovinn AND preq1.reqpovinn = preq2.povinn AND preq1.reqtyp = preq2.reqtyp
    LEFT JOIN povinn p1 ON preq1.povinn = p1.povinn
    LEFT JOIN povinn p2 ON preq2.povinn = p2.povinn
    LEFT JOIN searchable_povinn ps1 ON p1.povinn = ps1.povinn
    LEFT JOIN searchable_povinn ps2 ON p2.povinn = ps2.povinn
    WHERE preq1.reqtyp = 'Z'
    AND p2.pvyucovan = 'V'
    AND ps1.povinn IS NULL
    AND ps2.povinn IS NOT NULL
),
data AS (
    SELECT
        soident, zident, sdruh, zroc, zsem, o.nazev, zpovinn,
        (zroc::INT * 10 + zsem::INT) + (case when sdruh='B' then 100 when sdruh='N' then 200 end) ord
    FROM studium s
    LEFT JOIN zkous z ON s.sident = z.zident
    LEFT JOIN povinn p ON z.zpovinn = p.povinn
    LEFT JOIN obor o ON o.kod = s.sobor
    WHERE s.sobor like 'I%'
    AND sstav = 'A'
    AND z.zsplcelk = 'S'
    AND p.pgarant != '32-STUD'
)
"""
SELECT_TEMPL = """
SELECT x.soident, x.sdruh, x.sobor, x.zroc, x.zsem,
x.finished, next_semester.next_semester, relevant.relevant,
x.interchange_for, x.incompatible_with
FROM (
    SELECT
        soident, max(sdruh) sdruh, max(nazev) sobor, substr('{0}', 2, 1)::INT zroc, substr('{0}', 3, 1)::INT zsem,
        array_agg(zpovinn) finished,
        array_agg(zpreq.povinn) FILTER (WHERE zpreq.povinn IS NOT NULL) interchange_for,
        array_agg(npreq.povinn) FILTER (WHERE npreq.povinn IS NOT NULL) incompatible_with
    FROM data
    LEFT JOIN preq zpreq ON data.zpovinn = zpreq.reqpovinn AND zpreq.reqtyp = 'Z'
    LEFT JOIN preq npreq ON data.zpovinn = npreq.reqpovinn AND npreq.reqtyp = 'N'
    WHERE ord <= {0}
    GROUP BY soident
) x
LEFT JOIN (
    SELECT soident, array_agg(CASE WHEN eq.p1 IS NULL THEN zpovinn ELSE eq.p2 END) next_semester
    FROM data d
    LEFT JOIN eqpovinn eq ON d.zpovinn = eq.p1
    LEFT JOIN searchable_povinn ps ON d.zpovinn = ps.povinn
    WHERE ord = {1}
    AND (ps.povinn IS NOT NULL OR eq.p1 IS NOT NULL)
    GROUP BY soident
) next_semester ON next_semester.soident = x.soident
LEFT JOIN (
    SELECT soident, array_agg(CASE WHEN eq.p1 IS NULL THEN zpovinn ELSE eq.p2 END) relevant
    FROM data d
    LEFT JOIN eqpovinn eq ON d.zpovinn = eq.p1
    LEFT JOIN searchable_povinn ps ON d.zpovinn = ps.povinn
    WHERE ord > {0}
    AND (ps.povinn IS NOT NULL OR eq.p1 IS NOT NULL)
    GROUP BY soident
) relevant ON relevant.soident = x.soident
"""

if __name__ == "__main__":
    main()
