<html>
<head>
 <title> isula repo {{.website}}</title>
 <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<pre>
	{{range $key, $value := .files}}
	<a href="{{$value}}">{{$key}}</a>
	{{end}}
</pre>
</html>
