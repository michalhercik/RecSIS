[CmdletBinding()]
param (
    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string]$MigrationsDir,

    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string]$Container
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

# configuration
$Database = $env:POSTGRES_DB
$User     = $env:POSTGRES_OWNER

if (-not (Test-Path $MigrationsDir)) {
    throw "Migration directory not found: $MigrationsDir"
}

$files = @(Get-ChildItem `
     -Path $MigrationsDir `
     -Filter "*.sql" `
     | Sort-Object Name)

if ($files.Count -eq 0) {
    Write-Host "No migration files found."
    exit 0
}

Write-Host "Running migrations in a single transaction..."

# build one SQL stream
$sql = New-Object System.Text.StringBuilder

$sql.AppendLine("BEGIN;") | Out-Null
$sql.AppendLine("\set ON_ERROR_STOP on") | Out-Null

foreach ($file in $files) {
    Write-Host "Including migration: $($file.Name)"
    $sql.AppendLine("-- $($file.Name)") | Out-Null
    $sql.AppendLine((Get-Content $file.FullName -Raw)) | Out-Null
}

$sql.AppendLine("COMMIT;") | Out-Null

# execute once
$sql.ToString() |
    docker exec -i $Container `
        psql `
        -U $User `
        -d $Database

if ($LASTEXITCODE -ne 0) {
    throw "Migration transaction failed. All changes rolled back."
}

Write-Host "All migrations committed successfully."
