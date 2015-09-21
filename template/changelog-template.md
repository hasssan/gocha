## Changelog for {{.AppName}} {{.Version}}

### Chores
{{range $scope, $elmts := .Chores}}
    {{if $scope == "none"}}
    foo
    {{end}}
### {{$scope}}
    {{range $elmts}}
    - {{.Subject}}
    {{end}}
{{end}}