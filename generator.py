import random

semester1 = [1,3,4,5,6,7,12,13,28,18,19,20,21,22,29,23,24,25,26,27,30,31,32,33,34,35,38,41,57,42,44,45,48,49,50,51,58,52,53,54,55,59,60,61,62,63,73,74,76,79,80,82,84,86,88,90,92,95,96]
semester2 = [1,2,4,5,8,9,10,11,14,15,16,17,23,36,37,39,40,43,44,45,46,47,56,64,65,66,67,68,69,70,71,72,75,77,78,81,83,85,86,87,89,91,93,94,97,98,99,100]

for _ in range(13):
    choice = random.choice(semester1)
    semester1.remove(choice)
    print(f"{choice},1")

for _ in range(14):
    choice = random.choice(semester1)
    semester1.remove(choice)
    print(f"{choice},2")

