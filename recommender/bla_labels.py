from data import TrainData

data = TrainData(1234)
data.fit()

print(data.povinn)

# B. Subdomain (more specific knowledge area)
# Algebra
# Mathematical Analysis
# Geometry & Topology
# Probability & Statistics
# Numerical Methods / Scientific Computing
# Optimization / Operations Research
# Discrete Mathematics & Graph Theory
# Logic & Foundations
# Differential Equations
# Dynamical Systems
# Econometrics
# Financial Mathematics
# Insurance Mathematics
# Software Systems
# Theoretical Computer Science
# AI / NLP (→ for linguistic / ML-related)
# Computational Modeling
# Physics – General
# Physics – Theoretical
# Physics – Applied
# Geophysics / Climate

# D. Course Type
# Core / Mandatory
# Elective
# Advanced / Graduate
# Doctoral
# Introductory / First-year
# General Education


# select p.povinn, p.panazev, t.memo from searchable_povinn sp
# left join povinn p on sp.povinn = p.povinn
# left join pamela t on p.povinn = t.povinn and t.typ = 'S' and t.jazyk='ENG'
# where p.pgarant != '32-STUD'


#             WITH istudium AS (
#                 SELECT
#                     soident, sident, sdruh, srokp, sobor, o.nazev sobor_nazev, splan
#                 FROM studium s
#                 LEFT JOIN obor o ON s.sobor = o.kod
#                 WHERE s.sobor like 'I%'
#                 AND sstav NOT IN ('Z', 'U')
#             ),
#             tmp_interactions AS (
#                 SELECT
#                     soident, sident, zpovinn, zskr::INT, zroc, zsem
#                 FROM istudium
#                 LEFT JOIN zkous z ON istudium.sident = z.zident
#                 WHERE z.zsplcelk = 'S'
#             ),
#             interactions AS (
#                 SELECT DISTINCT
#                     i1.soident, i1.sident, i2.zpovinn povinn, i2.zskr, i2.zroc, i2.zsem
#                 FROM tmp_interactions i1
#                 LEFT JOIN tmp_interactions i2 ON i1.soident = i2.soident
#             ),
#             povinn AS (
#                 SELECT DISTINCT
#                     p.povinn, p.pnazev, panazev, p.pgarant, vucit1, vucit2, vucit3
#                 FROM interactions i
#                 LEFT JOIN povinn p ON i.povinn = p.povinn
#                 WHERE p.pgarant != '32-STUD'
#             )
# select rp.povinn, rp.pnazev, count(*) from povinn p
# inner join preq r on p.povinn = r.povinn and r.reqtyp in ('P', 'K')
# inner join povinn rp on r.reqpovinn = rp.povinn
# group by rp.povinn, rp.pnazev
# order by count(*) desc
