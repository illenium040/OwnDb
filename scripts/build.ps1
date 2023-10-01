$projectSoucrePath = $pwd
$binName = "owndb.exe"

$env:GOOS = "windows";
$env:GOARCH = "amd64";
go build -o "$projectSoucrePath/bin/$binName" "./src"

Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH

echo "Build successfully!"