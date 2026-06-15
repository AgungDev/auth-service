$ErrorActionPreference = 'Stop'

$swagCmd = Get-Command swag -ErrorAction SilentlyContinue
if (-not $swagCmd) {
    Write-Host 'Swag CLI not found. Installing github.com/swaggo/swag/cmd/swag@latest...'
    go install github.com/swaggo/swag/cmd/swag@latest
    $gopath = go env GOPATH
    $swagPath = Join-Path $gopath 'bin\swag.exe'
    if (-not (Test-Path $swagPath)) {
        Write-Error "swag was installed, but $swagPath was not found. Ensure GOPATH\bin is in your PATH and reopen PowerShell."
    }
    & $swagPath init -g cmd/server/main.go
} else {
    swag init -g cmd/server/main.go
}
