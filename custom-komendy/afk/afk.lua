{{/*
	This command allows users to set an AFK message with optional duration.
	Usage:

	-afk <message>
	-afk <message> -d <duration>

	Recommended trigger: Command trigger with trigger `-afk`
*/}}
{{ $time := 0 }}
{{ if .CmdArgs }}
	{{ $args := .CmdArgs }}
	{{ $message := "" }}
	{{ $duration := 0 }}
	{{ $skip := false }}
	{{- range $i, $v := $args -}}
		{{ if and (gt (len $v) 1) (not $skip) }}
			{{ if and (eq $v "-d") (gt (len $args) (add $i 1)) }}
				{{ $duration = index $args (add $i 1) }}
				{{ $skip = true }}
			{{ else }}
				{{ $message = joinStr " " $message $v }}
				{{ $skip = false }}
			{{ end }}
		{{ else if not $skip }}
			{{ $skip = false }}
			{{ $message = joinStr " " $message $v }}
		{{ else if $skip }}
			{{ $skip = false }}
		{{ end }}
	{{- end -}}
	{{ $parsedDur := 0 }}
	{{ with and $duration (toDuration $duration) }} {{ $parsedDur = . }} {{ end }}
	{{ if $parsedDur }}
		{{ dbSetExpire .User.ID "afk" $message (div $parsedDur 1000000000) }}
	{{ else }} {{ dbSet .User.ID "afk" $message }} {{ end }}
	{{ .User.Mention }}, Ustawiłem twój AFK na `{{ $message }}`.
{{ else }}
	{{ if dbGet .User.ID "afk" }}
		{{ dbDel .User.ID "afk" }}
		{{ .User.Mention }}, twój AFK został wyłączony.
	{{ else }}
		Wpisz wiadomość używając `r/afk <reason>`, z długością `r/afk -d <duration> <reason>` albo usuń AFK używając `r/afk` (albo nie wpisałeś powodu, albo użyłeś r/afk kiedy nie byłeś AFK)
	{{ end }}
{{ end }}