-- Example 1: Parse template from string
local tmpl = template.parse([[
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
</head>
<body>
    <h1>{{.title}}</h1>
    <p>{{.content}}</p>
    <ul>
    {{range .items}}
        <li>{{.}}</li>
    {{end}}
    </ul>
</body>
</html>
]])

local data = {
    title = "My Page",
    content = "Welcome to my page!",
    items = {"Item 1", "Item 2", "Item 3"}
}

local html = tmpl(data)
print(html)

-- Example 2: Parse template from file
local file_tmpl = template.parse_file("examples/templates/page.html")
local file_html = file_tmpl({
    title = "From File",
    content = "This content comes from a file template",
    items = {"File Item 1", "File Item 2"}
})
print(file_html)

-- Example 3: Parse multiple template files
local multi_tmpl = template.parse_files("examples/templates/layout.html", "examples/templates/content.html")
local multi_html = multi_tmpl({
    title = "Multiple Templates",
    content = "This uses multiple template files",
    items = {"Multi Item 1", "Multi Item 2"}
})
print(multi_html)

-- Example 4: Parse templates using glob pattern
local glob_tmpl = template.parse_glob("examples/templates/*.html")
local glob_html = glob_tmpl({
    title = "Glob Templates",
    content = "This uses glob pattern to find templates",
    items = {"Glob Item 1", "Glob Item 2"}
})
print(glob_html) 