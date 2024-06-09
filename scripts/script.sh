#!/usr/bin/env bash

topologias=("topologia_ciclo_3" 
"topologia_grid3x3" 
"topologia_linha" 
"topologia_tres_triangulos" 
"topologia_triangulo")

echo "Escolha a topologia:"

declare -i i

for ((i=0;i<${#topologias[@]};i++)); do
    echo "[$i] ${topologias[$i]}"
done

declare -i index

read -p "" index

declare topologia

topologia=${topologias[$index]}

echo "Você escolheu $topologia!"

declare -i n

n=$(cat ../txts/$topologia/len.txt)

echo "Uma topologia de $n nós!" 

declare -i PORT

HOST=$(ip address | grep -oE "\b192.168.[0-9]{1,3}.[0-9]{1,3}\b" | head -n 1)

for ((i=1; i<=2;i++)); do
    PORT=($i+5000)
    echo "Criando nó $i"
    gnome-terminal -- bash -c "../src/UP2P '$HOST:$PORT' '../txts/$topologia/$i.txt' '../txts/$topologia/pares_$i.txt'" bash
done

echo ''
echo '\O/ MASTER DOS NÓS, ESCOLHA UMA OPÇÃO \O/'
echo ''

declare -i opcao

while [[ 1=1 ]]; do
    echo ''
    echo "[0]: VISUALIZAR TOPOLOGIA"
    echo "[1]: ENCERRAR TODOS NÓS"
    echo ''
    read -p "" opcao
    if [[ $opcao == 0 ]]; then
        cat "../txts/$topologia/topologia.txt"
    elif [[ $opcao == 1 ]]; then
        pkill -f UP2P
        break
    else
        echo 'ESCOLHA UMA OPÇÃO ENTRE 0 E 1'
    fi
done
        