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

for ((i=0; i<2;i++)); do
    PORT=($i+1)
    echo "Criando nó $i"
    gnome-terminal -- bash -c "../src/UP2P '$HOST:500$PORT' '../txts/$topologia/vizinhos[$i].txt' '../txts/$topologia/lista_chave_valor[$i].txt'" bash
done