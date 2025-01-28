# Настройка доступа к домашнему серверу через удалённый VPS

Данный проект позволяет настроить доступ к домашнему серверу через удалённый VPS с использованием VPN. Домашний сервер будет подключён к VPS через OpenVPN, а весь входящий трафик на VPS будет перенаправляться на домашний сервер.

## Шаги настройки

### 1. Настройка OpenVPN
#### Установка OpenVPN на VPS
1. Установите OpenVPN:
   sudo apt update
   sudo apt install openvpn
2. Настройте OpenVPN-сервер. 
    Для этого вам понадобится .ovpn-файл, 
    содержащий настройки соединения. Убедитесь, 
    что сервер настроен на использование статического публичного IP-адреса.
    
    
    <h2>Настройка OpenVPN на домашнем сервере</h2>
Установите OpenVPN:

sudo apt update
sudo apt install openvpn

Скопируйте конфигурационный файл (например, home-server.ovpn) на домашний сервер в директорию /home/<username>/.

{client
dev tun
proto udp
remote <VPS_PUBLIC_IP> 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
cipher AES-256-CBC
verb 3

<ca>
-----BEGIN CERTIFICATE-----
(Содержимое файла ca.crt)
-----END CERTIFICATE-----
</ca>
<cert>
-----BEGIN CERTIFICATE-----
(Содержимое файла home-server.crt)
-----END CERTIFICATE-----
</cert>
<key>
-----BEGIN PRIVATE KEY-----
(Содержимое файла home-server.key)
-----END PRIVATE KEY-----
</key>}

Полная настройка OpenVPN для подключения к домашнему серверу через VPS
1.  Настройка VPS
Установка необходимых пакетов На вашем VPS выполните следующие команды:

sudo apt update
sudo apt install openvpn easy-rsa -y
Настройка инфраструктуры PKI (создание сертификатов и ключей)

Создайте папку для EasyRSA и переместитесь в неё:

make-cadir ~/openvpn-ca
cd ~/openvpn-ca
Инициализируйте инфраструктуру PKI:

./easyrsa init-pki
Затем создайте центр сертификации (CA):

./easyrsa build-ca nopass
После этого будет создан файл ca.crt (сертификат центра сертификации).

Создание серверного сертификата и ключей

Генерация ключа и сертификата для сервера:

./easyrsa gen-req server nopass
./easyrsa sign-req server server
Подтвердите подпись, введя yes. Эти команды создадут ключ и сертификат сервера.

Создание клиентского сертификата и ключа

Создайте сертификат для домашнего сервера (клиента):

./easyrsa gen-req home-server nopass
./easyrsa sign-req client home-server
Подтвердите подпись, введя yes. Теперь у вас есть сертификаты клиента.

Создание файла dh.pem и CRL Сгенерируйте параметры Диффи-Хеллмана (DH):

./easyrsa gen-dh
Создайте список отозванных сертификатов (CRL):

./easyrsa gen-crl
Копирование сертификатов в папку OpenVPN Переместите все необходимые файлы в директорию OpenVPN:

sudo cp pki/ca.crt /etc/openvpn/
sudo cp pki/issued/server.crt /etc/openvpn/
sudo cp pki/private/server.key /etc/openvpn/
sudo cp pki/dh.pem /etc/openvpn/
sudo cp pki/crl.pem /etc/openvpn/

2. Создание конфигурационного файла для сервера Создайте файл /etc/openvpn/server.conf:

sudo nano /etc/openvpn/server.conf
Вставьте в него следующие настройки:

ini
Copy
Edit
port 1194
proto udp
dev tun
ca ca.crt
cert server.crt
key server.key
dh dh.pem
crl-verify crl.pem
server 10.8.0.0 255.255.255.0
ifconfig-pool-persist ipp.txt
push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS 8.8.8.8"
keepalive 10 120
cipher AES-256-CBC
user nobody
group nogroup
persist-key
persist-tun
status openvpn-status.log
log /var/log/openvpn.log
verb 3

<b>Проверьте подключение:</b>

sudo openvpn --config /home/<username>/home-server.ovpn
Если подключение успешно, сервер начнёт видеть домашний сервер.

<h2>Добавление OpenVPN в автозагрузку на домашнем сервере</h2>
Для автоматического запуска OpenVPN при включении сервера:
Создайте файл службы systemd:

sudo nano /etc/systemd/system/openvpn-home.service

Вставьте следующий код:

[Unit]
Description=OpenVPN Client for Home Server
After=network.target

[Service]
ExecStart=/usr/sbin/openvpn --config /home/<username>/home-server.ovpn
Restart=always
User=root

[Install]
WantedBy=multi-user.target

<b>Сохраните файл и выполните команды:</b>

sudo systemctl daemon-reload
sudo systemctl enable openvpn-home.service
sudo systemctl start openvpn-home.service

Проверьте статус:

sudo systemctl status openvpn-home.service

<h2>Настройка проброса портов на VPS</h2>
Включение IP Forwarding
На VPS включите пересылку IP-пакетов:

sudo sysctl -w net.ipv4.ip_forward=1

Для сохранения изменений отредактируйте файл /etc/sysctl.conf:

sudo nano /etc/sysctl.conf

Добавьте или измените строку:

net.ipv4.ip_forward=1
Примените изменения:

sudo sysctl -p

Настройка NAT (проброс портов)
Для перенаправления всего трафика (кроме порта 22) с VPS на домашний сервер:

Установите правило в iptables:

sudo iptables -t nat -A PREROUTING -p tcp ! --dport 22 -j DNAT --to-destination <локальный_IP_ПК>
sudo iptables -t nat -A POSTROUTING -j MASQUERADE

Сохраните правила:

sudo apt install iptables-persistent -y
sudo netfilter-persistent save

Теперь любой входящий трафик (кроме порта 22) на VPS будет перенаправляться на домашний сервер.