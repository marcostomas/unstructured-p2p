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

for ((i=0; i<n;i++)); do
    PORT=($i+1)
    echo "Criando nó $i"
    gnome-terminal -- bash -c "../src/UP2P 'localhost:500$PORT' '../txts/$topologia/vizinhos[$i].txt' '../txts/$topologia/lista_chave_valor[$i].txt'" bash
done