# server-multissh

![Website](https://img.shields.io/website?url=https%3A%2F%2Fmultissh.github.io&up_message=online&down_message=offline&logo=googlechrome&label=demo%20website)
![Go Version](https://img.shields.io/badge/go-1.18-blue)
![License](https://img.shields.io/github/license/multissh/server-multissh)

This is the backend server for [multissh.github.io](https://multissh.github.io).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Build](#build)

## Installation

Follow these steps to install and set up the project on your local machine:

1. **Clone the repository**

   Use the following command to clone this repository:
   ```sh
   git clone https://github.com/multissh/server-multissh.git
   cd server-multissh
   ```

2. **Set up the API key**

   The project uses an API key for authentication. Replace the `Api_key` variable in `handler.go` with your actual API key. Rebuilt the ssh-client.
    ```go
    const Api_key = 'dTAu1iOvOfxQ63BZsYQpDqvyHMjeD8itjZ7GTs'
    ```
3. **Set up the SSL Certificate**

    Run the following command to generate SSL
    ```sh
    # example if ssh-client-amd64 in /root
    # change yourdomain.com
    apt-get remove certbot
    python3 -m pip install certbot
    certbot certonly -d yourdomain.com -d www.yourdomain.com --webroot-path /root
    openssl crl2pkcs7 -nocrl -certfile /etc/letsencrypt/live/yourdomain.com/fullchain.pem | openssl pkcs7 -print_certs -out /root/cert.crt
    openssl pkey -in /etc/letsencrypt/live/yourdomain.com/privkey.pem -out /root/private.key
    ```

    Run the following command to auto-update SSL `At 01:01 on day-of-month 1`
    ```sh
    # example if ssh-client-amd64 in /root
    # change yourdomain.com
    echo 'certbot renew --dry-run --webroot-path /root' > renew.sh
    echo 'openssl crl2pkcs7 -nocrl -certfile /etc/letsencrypt/live/yourdomain.com/fullchain.pem | openssl pkcs7 -print_certs -out /root/cert.crt' >> renew.sh
    echo 'openssl pkey -in /etc/letsencrypt/live/yourdomain.com/privkey.pem -out /root/private.key' >> renew.sh
    echo 'systemctl restart multissh-server' >> renew.sh
    chmod +x renew.sh
    echo '1 1 1 * * root /root/renew.sh > /dev/null' >> /etc/crontab
    ```

4. **Run the server**

    choose for your os environment
    ```sh
    ssh-client.exe       # windows amd64
    ./ssh-client-arm64   # linux arm64
    ./ssh-client-amd64   # linux amd64
    ```

## Usage

Alternatively, you can use `screen` or `systemctl` to manage the server process.
    
**Using screen :**

```sh
screen -Sdm /root/ssh-client-amd64
```

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
ExecStart=/root/ssh-client-amd64
User=yourusername
Restart=always

[Install]
WantedBy=multi-user.target
```
Save and close the file, then start the service:
```sh
sudo systemctl daemon-reload
sudo systemctl enable multissh-server
sudo systemctl start multissh-server
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
