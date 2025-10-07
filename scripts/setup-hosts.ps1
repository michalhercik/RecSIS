Start-Process powershell -Verb RunAs -ArgumentList "
    if (-not (Select-String -Path $env:SystemRoot\System32\drivers\etc\hosts -Pattern '127.0.0.1 mockcas')) {
        Add-Content -Path $env:SystemRoot\System32\drivers\etc\hosts -Value '`n127.0.0.1 mockcas'
    }
"
