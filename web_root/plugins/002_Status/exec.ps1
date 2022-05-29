
function WriteStats() {
    Write-Output (Get-WmiObject -Class win32_processor -ErrorAction Stop | Measure-Object -Property LoadPercentage -Average | Select-Object Average).Average
    $ComputerMemory = Get-WmiObject -Class win32_operatingsystem -ErrorAction Stop
    Write-Output -NoEnumerate $ComputerMemory.TotalVisibleMemorySize
    Write-Output -NoEnumerate $ComputerMemory.FreePhysicalMemory
}

function WriteProcs {
    Get-Process | Select-Object -Property ProcessName,Id,CPU | ConvertTo-Json
}

$opts = $args[0].Split(",")

switch ( $opts[0] ) {
    "stats" {WriteStats}
    "procs" {WriteProcs}
}

[Environment]::Exit(0)