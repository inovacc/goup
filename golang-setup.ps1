# install golang and all the tools needed to work with it
# gotests, gomodifytags, impl, goplay, dlv, staticcheck, gopls

# Check if running with PowerShell
if ($PSVersionTable.PSEdition -ne 'Core' -and $PSVersionTable.PSEdition -ne 'Desktop') {
    Write-Error "This script must be run with PowerShell or PowerShell Core."
    exit 1
}

# Function to check internet connection
function Test-InternetConnection {
    $null = Test-Connection -ComputerName "www.google.com" -Count 1 -Quiet
    return $?
}

# Function to install Go
function Install-Go {
    param (
        [string]$version
    )

    $downloadUrl = "https://golang.org/dl/$version.windows-amd64.msi"
    $msiFile = Join-Path $env:TEMP "go.msi"

    Invoke-WebRequest -UseBasicParsing -Uri $downloadUrl -OutFile $msiFile
    Start-Process msiexec -Wait -ArgumentList "/i $msiFile /quiet ADDLOCAL=ALL"
    Remove-Item $msiFile
}

# Function to install Go tools
function Install-GoTools {
    $tools = @(
        "github.com/cweill/gotests",
        "github.com/fatih/gomodifytags",
        "github.com/josharian/impl",
        "github.com/haya14busa/goplay/cmd/goplay",
        "github.com/go-delve/delve/cmd/dlv",
        "honnef.co/go/tools/cmd/staticcheck",
        "golang.org/x/tools/gopls"
    )

    foreach ($tool in $tools) {
        go install $tool@latest
    }
}

# Get latest Go version from website
$versionUrl = "https://golang.org/dl/?mode=json"
$latestVersion = (Invoke-WebRequest -UseBasicParsing -Uri $versionUrl | ConvertFrom-Json | Select-Object -First 1).version

# Check for internet connection
if (Test-InternetConnection) {
    Write-Output "Internet connection available."

    # Check if Go is already installed
    if (Get-Command go -ErrorAction SilentlyContinue) {
        $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'

        if ($installedVersion -eq "go$latestVersion") {
            Write-Output "Go $installedVersion is already installed."
        } else {
            Write-Output "Installing Go $latestVersion..."
            Install-Go -version $latestVersion

            if (Get-Command go -ErrorAction SilentlyContinue) {
                $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'
                Write-Output "Go $installedVersion has been installed."
            } else {
                Write-Error "Failed to install Go."
            }
        }
    } else {
        Write-Output "Installing Go $latestVersion..."
        Install-Go -version $latestVersion

        if (Get-Command go -ErrorAction SilentlyContinue) {
            $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'
            Write-Output "Go $installedVersion has been installed."
        } else {
            Write-Error "Failed to install Go."
        }
    }

    # Set environment variables
    $goPath = Join-Path $env:USERPROFILE "go"
    [System.Environment]::SetEnvironmentVariable("GOPATH", $goPath, "Machine")
    [System.Environment]::SetEnvironmentVariable("Path", "$env:Path;$($goPath)\bin", "Machine")

    # Install Go tools
    Write-Output "Installing Go tools..."
    Install-GoTools

    Write-Host "All tools successfully installed. You are ready to Go. :)"
} else {
    Write-Warning "Internet connection not available."

    # Set proxy
    $proxyUrl = "http://10.0.0.24:8080" # Replace with your proxy URL
    $proxy = New-Object System.Net.WebProxy($proxyUrl,$true)
    [System.Net.WebRequest]::DefaultWebProxy = $proxy

    # Check if proxy works
    if (Test-InternetConnection) {
        Write-Output "Proxy connection established."

        # Retry installation with proxy
        if (Get-Command go -ErrorAction SilentlyContinue) {
            $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'

            if ($installedVersion -eq "go$latestVersion") {
                Write-Output "Go $installedVersion is already installed."
            } else {
                Write-Output "Installing Go $latestVersion with proxy..."
                Install-Go -version $latestVersion

                if (Get-Command go -ErrorAction SilentlyContinue) {
                    $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'
                    Write-Output "Go $installedVersion has been installed."
                } else {
                    Write-Error "Failed to install Go."
                }
            }
        } else {
            Write-Output "Installing Go $latestVersion with proxy..."
            Install-Go -version $latestVersion

            if (Get-Command go -ErrorAction SilentlyContinue) {
                $installedVersion = go version | Select-String -Pattern 'go[0-9]+\.[0-9]+(\.[0-9]+)?'
                Write-Output "Go $installedVersion has been installed."
            } else {
                Write-Error "Failed to install Go."
            }
        }

        # Set environment variables
        $goPath = Join-Path $env:USERPROFILE "go"
        [System.Environment]::SetEnvironmentVariable("GOPATH", $goPath, "Machine")
        [System.Environment]::SetEnvironmentVariable("Path", "$env:Path;$($goPath)\bin", "Machine")

        # Install Go tools
        Write-Output "Installing Go tools with proxy..."
        Install-GoTools

        Write-Host "All tools successfully installed. You are ready to Go. :)"
    } else {
        Write-Error "Failed to establish proxy connection."
    }
}