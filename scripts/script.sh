#!/usr/bin/env bash

topologias=("topologia_ciclo_3" "topologia_grid3x3" "topologia_linha" "topologia_tres_triangulos" "topologia_triangulo")

echo "Escolha a topologia:"

for ((i=1;i<=${#topologias[@]};i++)); do
    echo ${topologias[i]}
done

declare topologia

read -p "" topologia

declare valido

valido=${0}

for elemento in "${topologias[@]}"; do
    if [[ $topologia == $elemento ]]; then
        valido=${1}
        break
    fi
done

if $valido; then
    echo "Opção válida."
else
    echo "Opção inválida."
    exit 1
fi

declare -i n

n=$(cat ../txts/$topologia/len.txt)

echo $n

for ((i=0; i<1000;i++)); do
    gnome-terminal -- bash -c "../src/UP2P 'localhost:500$i' '../txts/$topologia/vizinhos[0].txt' '../txts/$topologia/lista_chave_valor[0].txt'" bash
done