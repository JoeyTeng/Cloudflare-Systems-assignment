# Readme

## New knowledge

It was long since I last use Go (Golang), and I had no prior experience of using it as a pure HTTP server. Through this journey, though I encounter a confusing behaviour (and very little documented) that a `WriteHeader` function call to `http.ResponseWriter` would eat-up the `http.SetCookie` call. This took me much time to try-and-error with the help of `cURL`.

Besides, this is the first few times I wrote JavaScript / TypeScript. Though I had to learn and try new APIs like `fetch`, the whole journey is quite smooth and the only big trick to me is the `JSON.stringify` with `Map` after ES6.

## Most Difficult Part

I would say in the System Assignment, requirement "JWT can be verified with the public key". In the task note, it only states `The body of the response should simply be your RSA public key in plain text.`, but has no requirement to the format of the public key. It took me some time to figure out that byte stream / encoded hex stream is not desired, but the string stream of a PEM format.

## Extra Credit

Due to the time constraint I did not tried much; indeed, I had not even finish the task "Verify endpoint provides meaningful status code when token is invalid". However, I tried this "README.txt" endpoint.
