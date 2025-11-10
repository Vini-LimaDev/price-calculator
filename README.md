# price-calculator

Calculadora de Impostos em **Go** — projeto didático para praticar `structs`, ponteiros, leitura/escrita de arquivos e regras de imposto.

## Resumo
Projeto simples em Go que recebe uma entrada (JSON ou CLI), aplica regras de impostos/aliquotas por item, agrega resultados e gera saída em `saida.json` e um relatório `.txt`. Ideal para aprender organização de código, tratamento de erros e uso eficiente de memória com ponteiros.

## Funcionalidades
- Estruturas (`struct`) para representar Produto, Taxa e Resultado  
- Passagem de structs por ponteiro para evitar cópias desnecessárias  
- Leitura/Escrita de arquivos (`.json`, `.txt`)  
- Tratamento de erros consistente (retorno `error`)  
- Execução via binário ou `go run`

## Exemplo (resumo)
```go
func main() {
    produtos := LoadFromFile("entrada.json")
    for i := range produtos {
        CalcTaxes(&produtos[i]) // ponteiro evita cópia
    }
    SaveResults("saida.json", produtos)
}
