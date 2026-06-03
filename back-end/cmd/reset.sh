#!/bin/bash

# Verifica se foi passado o nome do app como argumento
if [ -z "$1" ]; then
  echo "Uso: ./update-go-app.sh nome_do_app"
  exit 1
fi

APP_NAME=$1

echo "Iniciando atualização do serviço '$APP_NAME'..."

# Parar o serviço systemd
echo "Parando serviço: $APP_NAME"
sudo systemctl stop "$APP_NAME"

# Desabilitar o serviço para evitar reinício automático durante a atualização
echo "Desabilitando serviço temporariamente"
sudo systemctl disable "$APP_NAME"

# Matar qualquer processo restante com o nome do app
echo "Matando processos em execução com nome '$APP_NAME'"
sudo pkill "$APP_NAME" 2>/dev/null

# Compilar o novo binário
echo "Compilando nova versão com go build..."
go build -o "$APP_NAME" main.go

# Mover o binário para /usr/local/bin
echo "Movendo binário para /usr/local/bin"
sudo mv "$APP_NAME" /usr/local/bin

# Recarregar os serviços do systemd
echo "Recarregando systemd"
sudo systemctl daemon-reexec
sudo systemctl daemon-reload

# Habilitar novamente o serviço para iniciar no boot
echo "Reabilitando serviço"
sudo systemctl enable "$APP_NAME"

# Iniciar o serviço
echo "Iniciando serviço"
sudo systemctl start "$APP_NAME"

# Mostrar status
echo "Status do serviço:"
sudo systemctl status "$APP_NAME" --no-pager