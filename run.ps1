docker compose up --build -d
if ($LASTEXITCODE -eq 0) {
    docker cp recsis-webapp:/app . --quiet
}
if ($LASTEXITCODE -eq 0) {
    mv ./app/*_templ.go . -Force
    rm ./app -Recurse 
    echo "Copied *_templ.go from recsis-webapp:/app"
}

