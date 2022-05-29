
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

switch ( $opts[0] ) {
    "get" {GetFeatures}
    "install" {InstallFeature}
    "uninstall" {UninstallFeature}
}

[Environment]::Exit(0)