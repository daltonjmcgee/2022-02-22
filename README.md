# 2022-02-22

This is a continuation of [this week's project](https://github.com/daltonjmcgee/2022-01-17).

I did a little more work on it by building the portfolio for [Executive Producer Jen McKenzie](https://jenmckenzie.com). This app is currently running that website as well via an nginx proxy.

The plan is two-fold:
- to make the directory structure better, more extensible and installable using go modules.
- to add an optional Admin Panel to update the `noSQL.json` file via the web. My biggest concern with that is authentication. I'd LOVE to get away with not using any SQL database, but I don't know how secure that is.

### Technologies

- [Golang](https://go.dev/)

### Resources and details

This project comes out of my interest in building a super-light HTTP Server that works with GET methods and uses no dependencies. I want it to have some easy to use feature that allow you to build a website with HTML/CSS and JS when you don't need to use a database of any sort.

### Features
- URL -> `public/[filename].[html]` mapping. If you go to `website.com/hello` you'll be served the file from `./public/hello.html`. This works with subdirectories as well, e.g. `website.com/pages/hello` will serve `./public/pages/hello.html`.
- Static files being served from the /static/ folder.
- Dynamic filenames that map to a key/value pair in a noSQL JSON file. e.g. [id].html will look for an "id" tag in all of the entries in the JSON "database". [title].html will look for the 'title' tag in all entries. First entry found is what is returned. In theory IDs should be a UUID so there is never a possible conflict.
