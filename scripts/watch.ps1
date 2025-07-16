wgo `
    -file='.go' `
    -file='.templ' `
    -xfile='_templ.go' `
    templ generate -path ./webapp :: `
    go build -C ./webapp -o RecSIS.exe :: `
    ./webapp/RecSIS.exe --config ./webapp/config.dev.toml