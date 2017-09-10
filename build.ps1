$dirty = &{ git diff-index --quiet HEAD; if (!$?) { echo "-dirty" } }
$rev = (git rev-parse --short HEAD)
$zipfile = "email_report-$rev$dirty.zip"

rm -fo *.zip
go clean .
go install
go build -ldflags="-w -s" .
zip a $zipfile email_report.exe
