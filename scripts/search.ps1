$dict = @{"q"="NPFL129";"attributesToRetrieve"=@("code","cs.NAME");"facets"=@("credits","start_semester","semester_count","faculty_guarantor")}
# $dict = @{"q"="";"attributesToRetrieve"=@("code","cs.NAME", "credits");"filter"="credits IN [10,bla]"}
$response = Invoke-WebRequest -Uri "http://localhost:7700/indexes/courses/search" `
    -Method Post `
    -Headers @{ "Authorization" = "Bearer MASTER_KEY" } `
    -ContentType "application/json" `
    -Body ($dict | ConvertTo-Json)
echo "  $($response.Content)"