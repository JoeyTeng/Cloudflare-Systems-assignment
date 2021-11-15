const POST_KEY_PREFIX = 'POST.'
const POST_NAMESPACE = 'kv'

addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})

/**
 * Respond with hello worker text
 * @param {Request} request
 */
async function handleRequest(request) {
    const END_POINT = new Map([
        ['/post', post]
    ])
    const url = new URL(request.url)

    if (END_POINT.has(url.pathname)) {
        return END_POINT.get(url.pathname)(request)
    }

    return new Response('Endpoint Not Found', {
        headers: { 'content-type': 'text/plain' },
        status: 404,
    })
}

async function post(request) {
    async function withGet() {
        let result = await kv.list({ 'prefix': POST_KEY_PREFIX })
        let postKeys = []
        do {
            postKeys.push(result.keys.map(keyInfo => keyInfo.name))
            result = await kv.list({ 'cursor': result.cursor })
        } while (!result.list_complete)

        postKeys = postKeys.flatMap(keys => {
            if (typeof keys === 'string') {
                return [keys]
            }
            return keys
        })
        const posts = (await Promise.all(postKeys.map(key => kv.get(key)))).filter(post => post != null)

        return new Response(JSON.stringify(posts, null, 2), {
            headers: {
                "content-type": "application/json;charset=UTF-8"
            },
            status: 200,
            statusText: 'success',
        })
    }

    async function withPost() {
        const { headers } = request
        const contentType = headers.get("content-type") || ""
        if (!contentType.includes("application/json")) {
            return new Response('Only accept JSON format. Please check request content type (' + request.contentType + ').', {
                status: 400,
                statusText: 'Invalid body: only accept JSON.',
            })
        }

        var body = new Map()
        const rawBodyJSON = await request.json()
        try {
            body = new Map(Object.entries(rawBodyJSON))
        } catch (e) {
            return new Response('Invalid JSON passed ' + e, {
                status: 400,
                statusText: 'Invalid JSON body.',
            })
        }

        let required_fields = ['title', 'username', 'content']
        required_fields = required_fields.filter(key => !body.has(key))
        if (required_fields.length > 0) {
            const missing_fields = required_fields.toString()
            return new Response(missing_fields + ' required', {
                status: 400,
                statusText: 'missing field(s) ' + missing_fields,
            })
        }
        const time = Date.now().toString()
        const uuid = await crypto.randomUUID()
            // key: POST.<TIMESTAMP>.<uuid>
        const key = POST_KEY_PREFIX + time + '.' + uuid

        const value = JSON.stringify(rawBodyJSON)
        await kv.put(key, value)

        const metadata = {
            'time': time,
            'uuid': uuid,
        }
        return new Response(JSON.stringify(metadata), {
            headers: {
                "content-type": "application/json;charset=UTF-8"
            },
            status: 200,
            statusText: 'success',
        })

    }

    if (request.method === 'GET') {
        return withGet()
    }
    if (request.method === 'POST') {
        return withPost()
    }

    return new Response('Post created successfully', {
        status: 201,
        statusText: 'success'
    })
}