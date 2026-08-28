[CmdletBinding()]
param(
    [string]$Version = $env:SOCIALBU_VERSION,
    [string]$InstallDir = $env:SOCIALBU_INSTALL_DIR,
    [switch]$NoPathUpdate
)

$ErrorActionPreference = "Stop"
$repository = "usamaejaz/socialbu-cli"

if ([string]::IsNullOrWhiteSpace($Version)) {
    $Version = "latest"
}
if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    $localAppData = [Environment]::GetFolderPath("LocalApplicationData")
    $InstallDir = Join-Path $localAppData "Programs\SocialBu"
}

$architecture = switch ([Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()) {
    "X64" { "amd64" }
    "Arm64" { "arm64" }
    default { throw "Unsupported architecture: $($_)" }
}

$asset = "socialbu_windows_$architecture.exe"
if ($Version -eq "latest") {
    $downloadBase = "https://github.com/$repository/releases/latest/download"
} else {
    $tag = if ($Version.StartsWith("v")) { $Version } else { "v$Version" }
    $downloadBase = "https://github.com/$repository/releases/download/$tag"
}

$tempDir = Join-Path ([IO.Path]::GetTempPath()) ("socialbu-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir | Out-Null

try {
    $binaryPath = Join-Path $tempDir $asset
    $checksumsPath = Join-Path $tempDir "checksums.txt"
    if ([string]::IsNullOrWhiteSpace($env:GITHUB_TOKEN)) {
        Invoke-WebRequest -Uri "$downloadBase/$asset" -OutFile $binaryPath
        Invoke-WebRequest -Uri "$downloadBase/checksums.txt" -OutFile $checksumsPath
    } else {
        $apiHeaders = @{
            "Accept" = "application/vnd.github+json"
            "Authorization" = "Bearer $($env:GITHUB_TOKEN)"
            "X-GitHub-Api-Version" = "2022-11-28"
        }
        $releaseApi = if ($Version -eq "latest") {
            "https://api.github.com/repos/$repository/releases/latest"
        } else {
            "https://api.github.com/repos/$repository/releases/tags/$tag"
        }
        $release = Invoke-RestMethod -Uri $releaseApi -Headers $apiHeaders

        foreach ($download in @(
            @{ Name = $asset; Path = $binaryPath },
            @{ Name = "checksums.txt"; Path = $checksumsPath }
        )) {
            $releaseAsset = $release.assets | Where-Object { $_.name -eq $download.Name } | Select-Object -First 1
            if ($null -eq $releaseAsset) {
                throw "Release asset not found: $($download.Name)"
            }
            $downloadHeaders = @{
                "Accept" = "application/octet-stream"
                "Authorization" = "Bearer $($env:GITHUB_TOKEN)"
                "X-GitHub-Api-Version" = "2022-11-28"
            }
            Invoke-WebRequest -Uri $releaseAsset.url -Headers $downloadHeaders -OutFile $download.Path
        }
    }

    $assetPattern = [regex]::Escape($asset)
    $checksumMatch = Select-String -Path $checksumsPath -Pattern "^([0-9a-fA-F]{64})\s+$assetPattern$" | Select-Object -First 1
    if ($null -eq $checksumMatch) {
        throw "No checksum found for $asset"
    }

    $expectedChecksum = $checksumMatch.Matches[0].Groups[1].Value.ToLowerInvariant()
    $actualChecksum = (Get-FileHash -Algorithm SHA256 -LiteralPath $binaryPath).Hash.ToLowerInvariant()
    if ($actualChecksum -ne $expectedChecksum) {
        throw "Checksum verification failed for $asset"
    }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $targetPath = Join-Path $InstallDir "socialbu.exe"
    Copy-Item -Force -LiteralPath $binaryPath -Destination $targetPath

    if (-not $NoPathUpdate) {
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $pathEntries = @($userPath -split ";" | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
        if ($pathEntries -notcontains $InstallDir) {
            $newUserPath = (@($pathEntries) + $InstallDir) -join ";"
            [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
            Write-Output "Added $InstallDir to your user PATH."
        }
        if (($env:Path -split ";") -notcontains $InstallDir) {
            $env:Path = "$env:Path;$InstallDir"
        }
    }

    Write-Output "Installed socialbu to $targetPath"
    & $targetPath version
} finally {
    Remove-Item -Recurse -Force -LiteralPath $tempDir -ErrorAction SilentlyContinue
}
