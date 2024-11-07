param(
    [switch]$Quiet=$false
)

Function IIf($If, $Right, $Wrong) {If ($If) {$Right} Else {$Wrong}}

docker compose up --build -d (IIf($Quiet, "--quiet-pull", ""))
if ($LASTEXITCODE -eq 0) {
    $error = docker cp recsis-webapp:/app . (IIf($Quiet, "--quiet", "")) 2>$1

    if ($LASTEXITCODE -eq 0) {
        mv ./app/*_templ.go . -Force
        rm ./app -Recurse 
        echo "Copied *_templ.go from recsis-webapp:/app"
    } else {
        echo $error
    }
}

