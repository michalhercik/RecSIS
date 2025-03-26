# $req = @{"q"="NPFL129";"attributesToRetrieve"=@("code","cs.NAME");"facets"=@("credits","start_semester","semester_count","faculty_guarantor")}
# $response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/multi-search" `
#     -Method Post `
#     -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
#     -ContentType "application/json" `
#     -Body ($req | ConvertTo-Json)
# echo "  $($response.Content)"

# $dict = @{"q"="";"attributesToRetrieve"=@("code","cs.NAME", "credits");"filter"="credits IN [10,bla]"}
$req1 = @{"indexUid"="courses";"q"="NPFL129";"filter"="credits IN [48]";"attributesToRetrieve"=@("code","cs.NAME");"facets"=@("credits","start_semester","semester_count","faculty_guarantor")}
$req2 = @{"indexUid"="courses";"q"="NPFL129";"limit"=0;"facets"=@("credits","start_semester","semester_count","faculty_guarantor")}
$req = @{"queries"=@($req1, $req2)}
$response = Invoke-WebRequest -Uri "http://localhost:7700/multi-search" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($req | ConvertTo-Json -Depth 5)
echo "$($response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 5)"