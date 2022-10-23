---
To: {{ .To }}
{{- if .Cc }}
Cc: {{ range .Cc }} {{ . }}, {{ end }}
{{ end }}
Subject: {{ .Subject }}
---

{{ .Body }}
