$ErrorActionPreference = "Stop"
# See https://stackoverflow.com/questions/9948517 for getting exitcode of normal exe's

function GetFeatures() {
    Get-WindowsFeature |  Where-Object {$_.Depth -eq 1 } | Select-Object -Property Name,DisplayName,Installed,DependsOn,Depth | ConvertTo-Json
}

function InstallFeature($name, $AllSub, $Manage) {
    $options = {}
    if ($AllSub -eq 1) {
        $options += "-IncludeAllSubFeature"
    }
    if ($Manage -eq 1) {
        $options += "-IncludeManagementTools"
    }
    Install-WindowsFeature -Name Web-Server $options
}

function UninstallFeature($name) {
    Uninstall-WindowsFeature -Name $name
}

$opts = $args[0].Split(",")

# Normally you should only add or change functions.
switch ( $opts[0] ) {
    "get" {GetFeatures}
    "install" {InstallFeature}
    "uninstall" {UninstallFeature}
}

[Environment]::Exit(0)