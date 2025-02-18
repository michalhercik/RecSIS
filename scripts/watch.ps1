$idk=Start-Process powershell -PassThru -NoNewWindow -ArgumentList "-NoProfile -Command `"wgo -file='.go' -file='.templ' -xfile='_templ.go' templ generate :: go run .`""
try {
    Write-Host "Process started with ID: $($idk.Id)"
    while ($true) {
        Start-Sleep -Seconds 1
    }
} finally {
    Stop-Process -Id $idk.Id -Force
    Write-Host "Removing temporary files..."
    Remove-Item -Recurse -Force ~\AppData\Local\Temp\go-build*
    Write-Host "Done"
}