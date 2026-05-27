@echo off
setlocal

set "BINARY_NAME=roamer.exe"

for /f "usebackq delims=" %%i in (`go env GOBIN`) do set "GO_BIN=%%i"

if not "%GO_BIN%"=="" goto target_resolved

for /f "usebackq delims=" %%i in (`go env GOPATH`) do set "GO_PATH=%%i"
if "%GO_PATH%"=="" (
    echo Neither GOBIN nor GOPATH is configured.
    exit /b 1
)

:target_resolved
if "%GO_BIN%"=="" set "GO_BIN=%GO_PATH%\bin"
set "TARGET=%GO_BIN%\%BINARY_NAME%"
set "BUILD_TARGET=%TEMP%\%BINARY_NAME%"

echo Installing Roamer
echo Target: %TARGET%

if not exist "%GO_BIN%" mkdir "%GO_BIN%"

go build -trimpath -o "%BUILD_TARGET%" .\cmd
if errorlevel 1 (
    echo Build failed.
    exit /b 1
)

copy /Y "%BUILD_TARGET%" "%TARGET%" >nul
if errorlevel 1 (
    echo Install failed.
    echo Close any running Roamer instance and try again.
    exit /b 1
)

echo %PATH% | find /I "%GO_BIN%" >nul
if errorlevel 1 (
    echo.
    echo Warning: %GO_BIN% is not on PATH.
    echo Add it to PATH to launch roamer from Windows search or a terminal.
) else (
    echo Install directory is on PATH.
)

echo Done.
