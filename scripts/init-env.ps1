param([string]$path)

get-content $path | foreach {
     $name, $value = $_.split("=")
     set-content env:\$name $value
}