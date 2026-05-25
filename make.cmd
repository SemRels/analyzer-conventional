@REM SPDX-License-Identifier: Apache-2.0
@REM SPDX-FileCopyrightText: 2026 The go-semrel Authors
@echo off
setlocal

if "%PLUGIN_NAME%"=="" set "PLUGIN_NAME=analyzer-conventional"
if "%DIST_DIR%"=="" set "DIST_DIR=dist"

if "%~1"=="" goto usage
if /I "%~1"=="build" goto build
if /I "%~1"=="test" goto test
if /I "%~1"=="lint" goto lint
if /I "%~1"=="coverage" goto coverage
if /I "%~1"=="build-all-platforms" goto buildall
if /I "%~1"=="release" goto release
if /I "%~1"=="clean" goto clean

goto usage

:build
if not exist bin mkdir bin
call go build -o "bin\%PLUGIN_NAME%.exe" .\cmd\plugin
exit /b %ERRORLEVEL%

:test
call go test -v ./...
exit /b %ERRORLEVEL%

:lint
call golangci-lint run
exit /b %ERRORLEVEL%

:coverage
call go test -cover ./...
exit /b %ERRORLEVEL%

:buildall
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"
call :buildtarget linux amd64 "%DIST_DIR%\%PLUGIN_NAME%-linux-amd64"
if errorlevel 1 exit /b %ERRORLEVEL%
call :buildtarget linux arm64 "%DIST_DIR%\%PLUGIN_NAME%-linux-arm64"
if errorlevel 1 exit /b %ERRORLEVEL%
call :buildtarget darwin amd64 "%DIST_DIR%\%PLUGIN_NAME%-darwin-amd64"
if errorlevel 1 exit /b %ERRORLEVEL%
call :buildtarget darwin arm64 "%DIST_DIR%\%PLUGIN_NAME%-darwin-arm64"
if errorlevel 1 exit /b %ERRORLEVEL%
call :buildtarget windows amd64 "%DIST_DIR%\%PLUGIN_NAME%-windows-amd64.exe"
if errorlevel 1 exit /b %ERRORLEVEL%
call :buildtarget windows arm64 "%DIST_DIR%\%PLUGIN_NAME%-windows-arm64.exe"
exit /b %ERRORLEVEL%

:buildtarget
set "TARGET_GOOS=%~1"
set "TARGET_GOARCH=%~2"
set "TARGET_OUTPUT=%~3"
set "GOOS=%TARGET_GOOS%"
set "GOARCH=%TARGET_GOARCH%"
call go build -o "%TARGET_OUTPUT%" .\cmd\plugin
set "RESULT=%ERRORLEVEL%"
set "GOOS="
set "GOARCH="
exit /b %RESULT%

:release
echo Build all release artifacts locally with "make build-all-platforms" and push a v*.*.* tag to trigger .github/workflows/release.yml
exit /b 0

:clean
call go clean
exit /b %ERRORLEVEL%

:usage
echo Usage: make ^<build^|test^|lint^|coverage^|build-all-platforms^|release^|clean^>
exit /b 1
