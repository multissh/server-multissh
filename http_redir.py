from sanic import Sanic, exceptions, response

app = Sanic('http_redir')
app.static('/.well-known', './.well-known', name='ssl')

@app.exception(exceptions.NotFound, exceptions.MethodNotSupported)
def redirect_everything_else(request, exception):
    server, path = request.server_name, request.path
    if server and path.startswith('/'):
        return response.redirect(f'https://{server}{path}', status=308)
    return response.text("Bad Request. Please use HTTPS!", status=400)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=False, access_log=False, auto_reload=True)
