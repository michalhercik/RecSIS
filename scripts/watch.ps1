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
                templ generate -path ./webapp :: `
                go build -C ./webapp -o RecSIS.exe :: `
                ./webapp/RecSIS.exe --config ./webapp/config.dev.toml
        } -Name 'OtherJob'
    )

    $jobs | Receive-Job -Wait -AutoRemoveJob
}
finally {
 $jobs | Remove-Job -Force
}