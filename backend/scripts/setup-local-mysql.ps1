# Run once in an elevated PowerShell window. Resets the local MySQL80 root
# password, creates an isolated qare database/user, and writes backend/.env.
# The SQLite database is not touched by this script.
$ErrorActionPreference = 'Stop'

$server = 'C:\Program Files\MySQL\MySQL Server 8.0\bin\mysqld.exe'
$client = 'C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe'
$admin = 'C:\Program Files\MySQL\MySQL Server 8.0\bin\mysqladmin.exe'
$config = 'C:\ProgramData\MySQL\MySQL Server 8.0\my.ini'
$backend = Split-Path -Parent $PSScriptRoot
$tempBase = Join-Path $env:TEMP ('qare-mysql-' + [guid]::NewGuid().ToString('N'))
$initFile = $tempBase + '.sql'
$rootOptionFile = $tempBase + '-root.cnf'
$appOptionFile = $tempBase + '-app.cnf'
$serverLog = $tempBase + '.log'
$port = 3307

function Write-PrivateFile([string]$path, [string]$content) {
    [System.IO.File]::WriteAllText($path, $content, (New-Object System.Text.UTF8Encoding($false)))
}

function Wait-ForPort([bool]$open) {
    for ($attempt = 0; $attempt -lt 60; $attempt++) {
        $socket = New-Object System.Net.Sockets.TcpClient
        $connected = $false
        try {
            $result = $socket.BeginConnect('127.0.0.1', $port, $null, $null)
            if ($result.AsyncWaitHandle.WaitOne(250)) {
                $socket.EndConnect($result)
                $connected = $true
            }
        } catch {
            $connected = $false
        } finally {
            $socket.Close()
        }
        if ($connected -eq $open) { return }
        Start-Sleep -Milliseconds 500
    }
    throw "MySQL port $port did not reach the expected state. Check $serverLog"
}

function Invoke-MySQLTool([string]$program, [string[]]$arguments) {
    $result = Start-Process -FilePath $program -ArgumentList $arguments -PassThru -Wait -WindowStyle Hidden -RedirectStandardOutput ($tempBase + '.out') -RedirectStandardError ($tempBase + '.err')
    return $result.ExitCode
}

$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = New-Object Security.Principal.WindowsPrincipal($identity)
if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw 'Open PowerShell as Administrator to run this script.'
}
foreach ($file in @($server, $client, $admin, $config)) {
    if (-not (Test-Path -LiteralPath $file)) { throw "Missing MySQL file: $file" }
}

do {
    $secure = Read-Host 'Choose a NEW MySQL root password (12+ characters: letters, digits, ! @ # % . _ -)' -AsSecureString
    $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    try {
        $rootPassword = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
    } finally {
        [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
    }
    if ($rootPassword -cnotmatch '^[A-Za-z0-9!@#%._-]{12,}$') {
        Write-Host 'Invalid password. Please use at least 12 characters with only letters, digits, ! @ # % . _ -' -ForegroundColor Yellow
    }
} while ($rootPassword -cnotmatch '^[A-Za-z0-9!@#%._-]{12,}$')
$random = New-Object byte[] 24
$rng = [Security.Cryptography.RandomNumberGenerator]::Create()
try { $rng.GetBytes($random) } finally { $rng.Dispose() }
$appPassword = -join ($random | ForEach-Object { $_.ToString('x2') })
$appUser = 'qare_app_' + $appPassword.Substring(0, 8)

$sql = @"
ALTER USER 'root'@'localhost' IDENTIFIED BY '$rootPassword';
CREATE DATABASE IF NOT EXISTS qare CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER '$appUser'@'127.0.0.1' IDENTIFIED BY '$appPassword';
GRANT ALL PRIVILEGES ON qare.* TO '$appUser'@'127.0.0.1';
"@
Write-PrivateFile $initFile $sql
Write-PrivateFile $rootOptionFile "[client]`r`nuser=root`r`npassword=`"$rootPassword`"`r`nhost=127.0.0.1`r`nport=$port`r`n"
Write-PrivateFile $appOptionFile "[client]`r`nuser=$appUser`r`npassword=`"$appPassword`"`r`nhost=127.0.0.1`r`nport=$port`r`n"

$temporaryServer = $null
try {
    Stop-Service MySQL80
    (Get-Service MySQL80).WaitForStatus('Stopped', [TimeSpan]::FromSeconds(30))
    Wait-ForPort $false

    $arguments = @("--defaults-file=`"$config`"", "--init-file=`"$initFile`"", '--console')
    $temporaryServer = Start-Process -FilePath $server -ArgumentList $arguments -PassThru -WindowStyle Hidden -RedirectStandardOutput $serverLog -RedirectStandardError ($serverLog + '.err')
    Wait-ForPort $true

    $connected = $false
    for ($attempt = 0; $attempt -lt 40; $attempt++) {
        $clientArgs = @("--defaults-extra-file=`"$appOptionFile`"", '--batch', '--skip-column-names', 'qare', '-e', '"SELECT 1"')
        if ((Invoke-MySQLTool $client $clientArgs) -eq 0) { $connected = $true; break }
        if ($temporaryServer.HasExited) { break }
        Start-Sleep -Milliseconds 500
    }
    if (-not $connected) { throw "MySQL did not apply the reset. Check $serverLog and $serverLog.err" }

    if ((Invoke-MySQLTool $admin @("--defaults-extra-file=`"$rootOptionFile`"", 'shutdown')) -ne 0) {
        throw 'Could not stop the temporary MySQL server cleanly.'
    }
    Wait-ForPort $false
    Remove-Item -LiteralPath $initFile -Force
    Start-Service MySQL80
    (Get-Service MySQL80).WaitForStatus('Running', [TimeSpan]::FromSeconds(30))
    Wait-ForPort $true
    if ((Invoke-MySQLTool $client $clientArgs) -ne 0) { throw 'The qare account could not connect after the service restart.' }

    Write-PrivateFile (Join-Path $backend '.env') "PORT=8080`r`nMYSQL_DSN=$($appUser):$appPassword@tcp(127.0.0.1:$port)/qare`r`nFRONTEND_ORIGIN=http://localhost:3000`r`n"
    Write-Host 'Success: MySQL80 is running on port 3307; qare database/user created; backend/.env written.'
} finally {
    foreach ($path in @($initFile, $rootOptionFile, $appOptionFile)) {
        Remove-Item -LiteralPath $path -Force -ErrorAction SilentlyContinue
    }
    if ((Get-Service MySQL80 -ErrorAction SilentlyContinue).Status -eq 'Stopped') {
        Start-Service MySQL80 -ErrorAction SilentlyContinue
    }
}
