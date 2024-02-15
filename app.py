#!/usr/bin/env python3
# python3 -m pip install sanic==23.6.0
# python3 -m pip install parallel-ssh==2.5.0
# python3 -m pip install gevent==21.8.0
# python3 -m pip install httpx==0.25.2

import sanic, json, httpx, asyncio
from sanic.exceptions import NotFound
from pssh.clients import ParallelSSHClient

api_key = 'dTAu1iOvOfxQ63BZsYQpDqvyHMjeD8itjZ7GTs'
app = sanic.Sanic(__name__)

class CustomString(str):
    def __eq__(self, other):
        return self is other
    def __hash__(self):
        return id(self)

def json_cfg(s):
    cfg = {}
    data = s.split('\n')
    for item in data:
        if item:
            row = item.split('||')
            user = row[0].split(':')[0]
            ppas = row[0].split(':')[1]
            ip = row[1].split(':')[0]
            pport = row[1].split(':')[1]
            data = {CustomString(ip): {'user': user, 'password': ppas, 'port': int(pport)}}
            cfg.update(data)
    return cfg

async def ws_receiver(msg_json, ws):
    if msg_json != 'pong':
        try:
            msg = json.loads(msg_json)
            if 'key' not in msg or msg.get('key') != api_key:
                await ws.send("red;;;error api_key not valid!")
            else:
                cfg = json_cfg(msg.get('cfg'))
                cmd = msg.get('cmd')
                hosts = cfg.keys()
                client = ParallelSSHClient(hosts, host_config=cfg, num_retries=1, timeout=60)
                if 'cus_cmd' in cmd:
                    split_cmd = cmd.split('\n')
                    cus_cmd = split_cmd[0].replace('cus_cmd = ', '').split(', ')
                    cmd = '\n'.join(split_cmd[1:])
                    host_list = [cmd.replace('cus_cmd', cus_cmd[i]) for i in range(len(hosts))]
                else:
                    host_list = [cmd for _ in range(len(hosts))]
                host_args = tuple(host_list)
                output = client.run_command("%s", host_args=host_args, stop_on_errors=False, read_timeout=60)
                for i, host_output in enumerate(output):
                    await ws.send(f"blue;;;{host_output.host}:22~# {host_list[i]}")
                    try:
                        for line in host_output.stdout:
                            await ws.send(f"green;;;{line}")
                        for line in host_output.stderr:
                            await ws.send(f"yellow;;;{line.replace('bash: ', '')}")
                    except Exception:
                        await ws.send(f"red;;;{host_output.exception.__class__.__name__}")
                        pass
        except:
            await ws.send("red;;;error processing data!")
        await ws.send("blue;;;All job Done!")

async def get_data(url, headers):
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(url, headers=headers)
            return response.text
    except:
        return None

@app.on_response
async def add_request_id_header(request, response):
    response.headers['Access-Control-Allow-Origin'] = '*'

@app.exception(NotFound)
async def ignore_404s(request, exception):
    return sanic.response.text("404 Not Found!", status=404)

@app.websocket('/run')
async def start(request, ws):
    async for msg in ws:
        try:
            await ws.send('ping')
            await ws_receiver(msg, ws)
        except asyncio.CancelledError:
            pass

@app.get('/server')
async def api_data_server(request):
    if 'key' not in request.args or request.args.get('key') != api_key:
        return sanic.response.text("404 Not Found!", status=404)
    headers = {'Authorization': f"token {request.args.get('token')}"}
    data = await get_data(request.args.get('url'), headers)
    if not data:
        return sanic.response.text("404 Not Found!", status=404)
    return sanic.response.text(data)

@app.get('/snippets')
async def api_data_snippets(request):
    if 'key' not in request.args or request.args.get('key') != api_key:
        return sanic.response.text("404 Not Found!", status=404)
    headers = {'Authorization': f"token {request.args.get('token')}"}
    data = await get_data(request.args.get('url'), headers)
    if not data:
        return sanic.response.text("404 Not Found!", status=404)
    return sanic.response.text(data)

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=443, debug=False, access_log=False, auto_reload=True, ssl={
        'cert': './cert.crt',
        'key': './private.key',
        'names': []
    })
