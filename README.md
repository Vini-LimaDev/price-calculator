# price-calculator

Calculadora de Impostos em **Go** — projeto didático para praticar `structs`, ponteiros, leitura/escrita de arquivos, importação de funções de outros arquivos.

## Resumo
Projeto simples em Go que recebe uma entrada dos valores apartir do arquivo `precos.txt`, aplica regras de impostos/aliquotas por item, agrega resultados e gera saída em `resultado_x.json`. Ideal para aprender organização de código, tratamento de erros e uso eficiente de memória com ponteiros.

## Funcionalidades
- Estruturas (`struct`) para representar Preços, Taxas e Resultado (Preços com a taxa inclusa)  
- Passagem de structs por ponteiro para evitar cópias desnecessárias  
- Leitura/Escrita de arquivos (`.json`, `.txt`)  
- Tratamento de erros consistente (retorno `error`)  
- Execução via `go run`

