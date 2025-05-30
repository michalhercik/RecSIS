try {
    $jobs = @(
        Start-ThreadJob {
            templ generate -path ./mock_cas
            go build -C ./mock_cas -o mockcas.exe
            ./mock_cas/mockcas.exe --cert server.crt --key server.key
        } -Name 'SomeJob'

        Start-ThreadJob {
            wgo `
                -file='.go' `
                -file='.templ' `
                -xfile='_templ.go' `
                templ generate -path ./src :: `
                go build -C ./src -o RecSIS.exe :: `
                ./src/RecSIS.exe --config ./src/config.dev.toml
        } -Name 'OtherJob'
    )

    $jobs | Receive-Job -Wait -AutoRemoveJob
}
finally {
 $jobs | Remove-Job -Force
}