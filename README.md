# 🧮 Calculadora de Preço  

Aplicação em **Go (Golang)** com **frontend em React + TypeScript + Vite + TailwindCSS**, desenvolvida para **calcular impostos e taxas sobre valores**.  
O projeto foi criado com fins **didáticos**, explorando desde `structs`, **ponteiros**, **leitura e escrita de arquivos** até a integração entre backend e frontend via API local.

---

## 📘 Visão Geral

O sistema lê valores de entrada — digitados manualmente ou importados de um arquivo `.txt` — aplica **taxas/aliquotas** sobre cada valor, e exibe os resultados em uma **tabela organizada**.  
Os resultados também são exportados para um arquivo `.json`, que pode ser baixado pelo usuário.

Cada linha da tabela mostra:
- 💰 Valor **antes** da taxa  
- 💹 Valor **após** a taxa  
- ⚖️ Diferença entre ambos  

---

## ⚙️ Funcionalidades

### 🧠 Backend (Go)
- Leitura de arquivos `.txt` com valores base  
- Aplicação de taxas configuradas pelo usuário  
- Escrita de resultados em `.json`  
- Estruturas (`structs`) para representar preços, taxas e resultados  
- Passagem por ponteiro (`*T`) para otimizar o uso de memória  
- Modularização em pacotes (`cmd`, `conversao`, `leitura_arq`, `precos`)  
- Tratamento de erros com retorno `error`  

### 💻 Frontend (Vite + React + TypeScript)
- Interface moderna e responsiva com **TailwindCSS**  
- Campos para digitar valores e taxas manualmente  
- Upload de arquivo `.txt` contendo valores  
- Botão para executar o cálculo via API  
- Exibição de resultados em tabela dinâmica  
- Opção para **baixar o resultado em JSON**
---

## 🚀 Execução do Projeto
- Backend (Go)
  - Por padrão roda em: _http://localhost:8080_
```bash
cd cmd/server
go run main.go
```

- Frontend(Vite)
```bash
cd price-calculator-frontend
npm install
npm run dev
```
---
## 🧩 Estrutura do Projeto

```bash
📦 price-calculator
 ┣ 📂 cmd
 ┃ ┣ 📂 cli                      
 ┃ ┗ 📂 server                   # Servidor HTTP principal (API Go)
 ┃   ┗ 📜 main.go                # Ponto de entrada do backend
 ┣ 📂 conversao
 ┃ ┗ 📜 converter.go             # Regras de cálculo e aplicação de taxas
 ┣ 📂 leitura_arq
 ┃ ┗ 📜 filemanager.go           # Funções de leitura/escrita de arquivos
 ┣ 📂 precos
 ┃ ┗ 📜 precos.go                # Definição de structs (Preço, Taxa, Resultado)
 ┣ 📂 price-calculator-frontend
 ┃ ┣ 📂 public                   # Recursos públicos do front
 ┃ ┣ 📂 src
 ┃ ┃ ┣ 📂 assets                 # Imagens e ícones
 ┃ ┃ ┣ 📜 App.tsx                # Componente principal do app
 ┃ ┃ ┣ 📜 App.css                # Estilos específicos do app
 ┃ ┃ ┣ 📜 main.tsx               # Ponto de entrada React/Vite
 ┃ ┃ ┣ 📜 index.css              # Estilos globais + Tailwind base
 ┃ ┣ 📜 .env                     # Variáveis de ambiente (URL do backend)
 ┣ 📜 precos.txt                 # Arquivo de entrada de exemplo
 ┣ 📜 resultado.json             # Exemplo de resultado gerado
 ┣ 📜 go.mod                     # Dependências Go
