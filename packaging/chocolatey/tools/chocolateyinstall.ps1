$ErrorActionPreference = 'Stop'

$packageArgs = @{
    packageName    = $env:chocolateyPackageName
    unzipLocation  = Split-Path -Parent $MyInvocation.MyCommand.Definition
    url64bit       = 'https://github.com/socialbu/socialbu-cli/releases/download/v__VERSION__/socialbu___VERSION___windows_amd64.zip'
    checksum64     = '__CHECKSUM64__'
    checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
