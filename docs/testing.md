# Testing

In this application, we use only integration tests, as it is a web application. We test almost all routes, with exceptions written in issues. The testing is straightforward and for more information, please refer `main_test.go` file, which contains all the tests and setup.

One important note is that you need to define environment variables for the tests to run successfully. In VS Code, this can be done by adding the following configuration to your `settings.json` file:

```json
"go.testEnvVars": {
  "RECSIS_WEBAPP_DB_PASS": "your_password_defined_in_env",
  "MEILI_MASTER_KEY": "your_meili_master_key_defined_in_env"
},
```

If you run tests from command line, you need to set the environment variables in your shell before running the tests. This step was covered in [Run](./how-to-run.md#run) section. Then you need to execute the following command in the `webapp` directory:

```shell
go test -v
```
`-v` option enables verbose output, showing all tests that are run and their results.
