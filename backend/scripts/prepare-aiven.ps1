$ErrorActionPreference = 'Stop'

$root = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$cache = Join-Path $root '.cache'
$certificateTarget = Join-Path $root 'backend/aiven-ca.pem'

Write-Host 'Aiven MySQL > Overview > Connection information'
$hostName = (Read-Host 'Host').Trim()
$portText = (Read-Host 'Port').Trim()
$database = (Read-Host 'Database [defaultdb]').Trim()
$user = (Read-Host 'User [avnadmin]').Trim()
$password = Read-Host 'Password (hidden)' -AsSecureString
$certificateSource = (Read-Host 'Full path to downloaded CA Certificate (ca.pem)').Trim('"', "'", ' ')

if (-not $database) { $database = 'defaultdb' }
if (-not $user) { $user = 'avnadmin' }
if (-not $hostName -or $hostName -match '[\s/:@]') { throw 'Invalid host.' }
$portNumber = 0
if (-not [int]::TryParse($portText, [ref]$portNumber) -or $portNumber -lt 1 -or $portNumber -gt 65535) { throw 'Invalid port.' }
if (-not $database -or $database -match '[\s/?@]') { throw 'Invalid database name.' }
if (-not $user -or $user -match '[\s/:@]') { throw 'Invalid user.' }
if ($password.Length -eq 0) { throw 'Password cannot be empty.' }
if (-not (Test-Path -LiteralPath $certificateSource -PathType Leaf)) { throw 'CA certificate file not found.' }
$certificate = Get-Content -LiteralPath $certificateSource -Raw
if ($certificate -notmatch '-----BEGIN CERTIFICATE-----') { throw 'Invalid CA certificate file.' }

New-Item -ItemType Directory -Force -Path $cache | Out-Null
Copy-Item -LiteralPath $certificateSource -Destination $certificateTarget -Force
$connection = @{
    host = $hostName
    port = $portNumber
    database = $database
    user = $user
    encrypted_password = ($password | ConvertFrom-SecureString)
}
$json = $connection | ConvertTo-Json -Compress
[System.IO.File]::WriteAllText((Join-Path $cache 'aiven-connection.json'), $json, [System.Text.UTF8Encoding]::new($false))
Write-Host 'Saved local Aiven connection information. No password was printed.'
