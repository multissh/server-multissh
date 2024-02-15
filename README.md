# server-multissh

![Website](https://img.shields.io/website?url=https%3A%2F%2Fmultissh.github.io&up_message=online&down_message=offline&logo=googlechrome&label=demo%20website)
![Python Version](https://img.shields.io/badge/python-3.9-blue)
![Go Version](https://img.shields.io/badge/go-1.17-blue)
![License](https://img.shields.io/github/license/multissh/server-multissh)
![Sanic](https://img.shields.io/badge/Sanic-23.6.0-blue)
![ParallelSSH](https://img.shields.io/badge/ParallelSSH-2.5.0-blue)
![Gevent](https://img.shields.io/badge/Gevent-21.8.0-blue)
![HTTPX](https://img.shields.io/badge/HTTPX-0.25.2-blue)

This is the backend server for [multissh.github.io](https://multissh.github.io).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Ssl](#ssl-certificate-generation)
- [Build](#build)

## Installation

Follow these steps to install and set up the project on your local machine:

1. **Clone the repository**

   Use the following command to clone this repository:
   ```sh
   git clone https://github.com/multissh/server-multissh.git
   cd server-multissh
   ```

2. **Install dependencies**
    ```sh
    python3 -m pip install sanic==23.6.0 parallel-ssh==2.5.0 gevent==21.8.0 httpx==0.25.2
    ```
    or
    ```sh
    python -m pip install sanic==23.6.0 parallel-ssh==2.5.0 gevent==21.8.0 httpx==0.25.2
    ```

3. **Set up the API key**
    The project uses an API key for authentication. Replace the `api_key` variable in `app.py` with your actual API key.
    example :
    ```py
    api_key = 'dTAu1iOvOfxQ63BZsYQpDqvyHMjeD8itjZ7GTs'
    ```
4. **Run the server**
    This project consists of two Python scripts that need to be run: `app.py` and `http_redir.py`.
    ```bash
    python3 app.py
    python3 http_redir.py
    ```
    choose for your os environment
    ```
    ssh-client.exe
    ./ssh-client-arm64
    ./ssh-client-amd64
    ```

## SSL Certificate Generation
This project requires an SSL certificate and key. Follow these steps to generate them:

1. **Generate a private key**

   Run the following command to generate a private key:
   ```sh
    openssl genrsa -out private.key 2048
    openssl req -new -key private.key -out cert.csr
    openssl x509 -req -days 365 -in cert.csr -signkey private.key -out cert.crt
   ```

    After generating the certificate and key, place them in the root directory of the project and ensure they are named cert.crt and private.key respectively.
    Please note that this will generate a self-signed certificate. It's recommended to use a certificate from a trusted certificate authority for production environments.

## Usage

Alternatively, you can use `screen` or `systemctl` to manage the server process.
    
**Using screen :**

```sh
screen -S server-session python3 app.py
screen -S server-session python3 http_redir.py
```

To detach from the screen session, press Ctrl+A followed by D.

**Using systemctl :**

Create a new service file (replace yourusername with your actual username):
```sh
sudo nano /etc/systemd/system/multissh-server.service
```
Add the following content to the file:
```
[Unit]
Description=MultiSSH Server

[Service]
ExecStart=/usr/bin/python3 /home/yourusername/server-multissh/app.py
User=yourusername
Restart=always

[Install]
WantedBy=multi-user.target
```
Save and close the file, then start the service:
```sh
sudo systemctl start multissh-server
```
To enable the service to start on boot:
```sh
sudo systemctl enable multissh-server
```

## Build
This project includes build scripts for both Windows and Unix-based systems. Depending on your operating system, you can use either of these scripts to build the project:
- **Windows**

  If you're using a Windows system, you can use the `build.bat` script. Open a command prompt in the project directory and run the following command:

  ```cmd
  build.bat
  ```
- **Unix-based systems (Linux, macOS)***

    If you're using a Unix-based system like Linux or macOS, you can use the build.sh script. Open a terminal in the project directory and run the following command:

    ```sh
    ./build.sh
    ```
