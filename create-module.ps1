# Check if moduleName argument is provided
if (-not $args[0]) {
    Write-Host "Usage: .\script.ps1 moduleName"
    exit 1
}

$moduleName = $args[0]

# Define the base directory
$baseDir = "internal/$moduleName"

# Create the directory structure
$dirs = @(
    "$baseDir/delivery/http",
    "$baseDir/dto",
    "$baseDir/interfaces",
    "$baseDir/repository",
    "$baseDir/usecase"
)

foreach ($dir in $dirs) {
    if (-not (Test-Path -Path $dir)) {
        New-Item -ItemType Directory -Path $dir | Out-Null
    }
}

# Create the files
$files = @(
    "$baseDir/delivery/http/handlers.go",
    "$baseDir/delivery/http/routers.go",
    "$baseDir/interfaces/delivery.go",
    "$baseDir/interfaces/repository.go",
    "$baseDir/interfaces/usecase.go",
    "$baseDir/repository/repository.go",
    "$baseDir/usecase/usecase.go"
)

foreach ($file in $files) {
    if (-not (Test-Path -Path $file)) {
        New-Item -ItemType File -Path $file | Out-Null
    }
}

Write-Host "Folder and file structure created under internal/$moduleName"
