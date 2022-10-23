---
To: {{ .To }}
{{- if .Cc }}
Cc: {{ .Cc }}
{{- end }}
Subject: {{ .Subject }}
---

{{ .Body }}
