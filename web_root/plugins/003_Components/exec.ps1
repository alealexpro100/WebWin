$ErrorActionPreference = "Stop"
# See https://stackoverflow.com/questions/9948517 for getting exitcode of normal exe's

function GetFeatures() {
    Get-WindowsFeature |  Where-Object {$_.Depth -eq 1 } | Select-Object -Property DisplayName,Installed,Description,DependsOn,Name | ConvertTo-Json
}

function InstallFeature($name, $AllSub, $Manage) {
    if ($AllSub -eq "1") {
        $AllSub = $true
    } else {
        $AllSub = $false
    }
    if ($Manage -eq "1") {
        $Manage = $true
    } else {
        $Manage = $false
    }
    Install-WindowsFeature -Name $name -IncludeAllSubFeature:$AllSub -IncludeManagementTools:$Manage
}

function UninstallFeature($name) {
    Uninstall-WindowsFeature -Name $name
}

$opts = $args[0].Split(",")

# Normally you should only add or change functions.
switch ( $opts[0] ) {
    "get" {GetFeatures}
    "install" {InstallFeature $opts[1] $opts[2] $opts[3]}
    "uninstall" {UninstallFeature $opts[1]}
}

[Environment]::Exit(0)