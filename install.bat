@echo off
setlocal enabledelayedexpansion

rem ============================================================
rem  Install / uninstall every command in this repository.
rem    install.bat            ... install all commands
rem    install.bat uninstall  ... uninstall all commands
rem
rem  NOTE: Keep this file ASCII-only. cmd.exe misreads batch files
rem        that contain UTF-8 Japanese text, so comments and
rem        messages are written in English on purpose.
rem ============================================================

rem Move to the folder that holds this file (the repository root),
rem so the script behaves the same no matter where it is run from.
cd /d "%~dp0"

set "RC=0"
set /a OK=0, NG=0
set "NOPATH="
set "BIN="

rem The first argument selects the action (default: install).
set "MODE=%~1"
if not defined MODE set "MODE=install"
if /i "%MODE%"=="install" goto :check_go
if /i "%MODE%"=="uninstall" goto :check_go
goto :usage

:check_go
rem Go is required both to build and to locate the install folder.
where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go was not found. Install it from https://go.dev/dl/
    set "RC=1"
    goto :end
)

rem Treat every folder that has a go.mod as one command.
rem A newly added command folder is picked up without editing this file.
for /d %%D in (*) do if exist "%%D\go.mod" call :%MODE%_one "%%D"

rem Show the summary.
echo.
if /i "%MODE%"=="install" (
    echo Install result: %OK% succeeded, %NG% failed
    if defined BIN echo Install folder: !BIN!
) else (
    echo Uninstall result: %OK% removed, %NG% failed
)
if defined NOPATH (
    echo.
    echo [WARN] The install folder is not on PATH, so the commands cannot be run by name.
    echo        Add the install folder above to your user "Path" environment variable,
    echo        then open a new Command Prompt.
)
if %NG% gtr 0 set "RC=1"
goto :end

:usage
echo Usage: %~nx0 [install ^| uninstall]
echo   install    install all commands. This is the default.
echo   uninstall  uninstall all commands.
set "RC=1"
goto :end

rem ------------------------------------------------------------
rem  Install one command.
rem    %1 = folder name of the command
rem ------------------------------------------------------------
:install_one
pushd "%~1"
call :get_target
go install .
set "ERR=%ERRORLEVEL%"
popd
if not "%ERR%"=="0" (
    echo [NG] %~1
    set /a NG+=1
    exit /b 0
)
echo [OK] %~1
set /a OK+=1
if not defined TARGET exit /b 0

rem Check that typing the command name really runs the one just installed.
rem   not found       -> the install folder is not on PATH
rem   found elsewhere -> a same-named command comes earlier on PATH
set "FIRST="
for /f "delims=" %%W in ('where "%NAME%" 2^>nul') do if not defined FIRST set "FIRST=%%~fW"
if not defined FIRST set "NOPATH=1" & exit /b 0
if /i not "%FIRST%"=="%TARGET%" echo      [WARN] another "%NAME%" comes first on PATH: !FIRST!
exit /b 0

rem ------------------------------------------------------------
rem  Uninstall one command.
rem    %1 = folder name of the command
rem ------------------------------------------------------------
:uninstall_one
pushd "%~1"
call :get_target
popd
if not defined TARGET (
    echo [NG] %~1 : could not determine the install path
    set /a NG+=1
    exit /b 0
)
if not exist "%TARGET%" (
    echo [--] %~1 : not installed
    exit /b 0
)
del "%TARGET%"
if exist "%TARGET%" (
    echo [NG] %~1 : could not delete !TARGET!
    set /a NG+=1
    exit /b 0
)
echo [OK] %~1 : removed
set /a OK+=1
exit /b 0

rem ------------------------------------------------------------
rem  Ask Go where "go install" puts the executable of the command
rem  in the current folder.
rem    TARGET = full path of the executable
rem    NAME   = command name without extension
rem    BIN    = install folder
rem ------------------------------------------------------------
:get_target
set "TARGET="
set "NAME="
for /f "delims=" %%T in ('go list -f "{{.Target}}" . 2^>nul') do set "TARGET=%%~fT" & set "NAME=%%~nT" & set "BIN=%%~dpT"
if defined TARGET set "BIN=%BIN:~0,-1%"
exit /b 0

:end
rem When started by double-click, keep the window open so the result can be read.
set "CMDLINE=!cmdcmdline!"
if not "!CMDLINE:%~nx0=!"=="!CMDLINE!" pause
endlocal & exit /b %RC%
