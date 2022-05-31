$ErrorActionPreference = "Stop"
# See https://stackoverflow.com/questions/9948517 for getting exitcode of normal exe's

function InstallSSHD() {
    Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
    Start-Service sshd
    Set-Service -Name sshd -StartupType 'Automatic'
    if (!(Get-NetFirewallRule -Name "OpenSSH-Server-In-TCP" -ErrorAction SilentlyContinue | Select-Object Name, Enabled)) {
        New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
    }
}

$opts = $args[0].Split(",")

# Normally you should only add or change functions.
switch ( $opts[0] ) {
    "sshd" {InstallSSHD}
}

[Environment]::Exit(0)