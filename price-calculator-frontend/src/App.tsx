import React, { useMemo, useRef, useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";

/* ========= Helpers ========= */
function parsePrices(text: string): number[] {
  if (!text) return [];
  const tokens = text
    .split(/[\n,;\s]+/)
    .map((t) => t.trim())
    .filter(Boolean);
  const nums = tokens.map((t) => Number(t.replace(",", ".")));
  const bad = nums.findIndex((n) => Number.isNaN(n));
  if (bad !== -1) throw new Error(`Valor inválido: "${tokens[bad]}"`);
  return nums;
}

function formatCurrency(n: number) {
  return new Intl.NumberFormat("pt-BR", { 
    style: "currency", 
    currency: "BRL" 
  }).format(n);
}

function downloadJSON(filename: string, data: unknown) {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url; 
  a.download = filename;
  document.body.appendChild(a); 
  a.click(); 
  a.remove();
  URL.revokeObjectURL(url);
}

/* ========= Tipos / API ========= */
interface ResultadoGo {
  taxa: number;
  precos: number[];
  precos_incluindo_taxas: string[];
}

const API_BASE = import.meta.env.VITE_API_BASE as string | undefined;

async function apiCalcular(precos: number[], taxas: number[]): Promise<ResultadoGo[]> {
  if (!API_BASE) return taxas.map((t) => calcularLocal(precos, t));
  const r = await fetch(`${API_BASE}/calcular`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ precos, taxas }),
  });
  if (!r.ok) throw new Error("Falha ao calcular no backend.");
  return (await r.json()) as ResultadoGo[];
}

function calcularLocal(precos: number[], taxa: number): ResultadoGo {
  const c = precos.map((p) => Number((p * (1 + taxa)).toFixed(2)));
  return { taxa, precos, precos_incluindo_taxas: c.map((x) => x.toFixed(2)) };
}

/* ========= App ========= */
export default function App() {
  const [precosTexto, setPrecosTexto] = useState<string>("");
  const [taxasTexto, setTaxasTexto] = useState<string>("");
  const [precosArquivo, setPrecosArquivo] = useState<number[] | null>(null);
  const [resultados, setResultados] = useState<ResultadoGo[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState<string | null>(null);
  const [erro, setErro] = useState<string | null>(null);

  const fileRef = useRef<HTMLInputElement | null>(null);

  const precosAtuais = useMemo<number[]>(() => {
    if (precosArquivo && precosArquivo.length) return precosArquivo;
    try {
      return parsePrices(precosTexto);
    } catch {
      return [];
    }
  }, [precosArquivo, precosTexto]);

  const handleUploadClick = () => fileRef.current?.click();

  const onFileChange: React.ChangeEventHandler<HTMLInputElement> = async (e) => {
    const f = e.target.files?.[0];
    if (!f) return;
    try {
      setErro(null);
      const text = await f.text();
      const nums = parsePrices(text);
      setPrecosArquivo(nums);
      setResultados(null);
      // SEMPRE atualiza o campo de texto com os valores do arquivo
      setPrecosTexto(nums.join('\n'));
      setMsg(`${nums.length} preços carregados do arquivo.`);
    } catch (err: any) {
      setErro(err?.message || "Falha ao ler arquivo.");
    } finally {
      e.currentTarget.value = "";
    }
  };

  const parseTaxas = (): number[] => {
    const arr = taxasTexto
      .split(/[,;\s]+/)
      .map((t) => t.trim())
      .filter(Boolean)
      .map((t) => Number(t.replace(",", ".")))
      .map((t) => t / 100);
    if (arr.some((n) => Number.isNaN(n) || !Number.isFinite(n) || n < 0)) {
      throw new Error("Há taxas inválidas. Use números (ex.: 10, 7.5, 15)");
    }
    return arr;
  };

  const handleCalcular = async () => {
    try {
      setErro(null); 
      setMsg(null);
      const precos = precosAtuais;
      if (!precos.length) throw new Error("Nenhum preço informado.");
      const taxas = parseTaxas();
      setLoading(true);
      const res = await apiCalcular(precos, taxas);
      setResultados(res);
      setMsg(`Calculado ${res.length} cenário(s) de taxas para ${precos.length} produto(s).`);
    } catch (err: any) {
      setErro(err?.message || "Erro ao calcular.");
    } finally {
      setLoading(false);
    }
  };

  const limparArquivo = () => {
    setPrecosArquivo(null);
    setResultados(null);
    // SEMPRE limpa o campo de texto
    setPrecosTexto("");
    // Limpa o valor do input file para permitir re-upload do mesmo arquivo
    if (fileRef.current) {
      fileRef.current.value = "";
    }
    setMsg("Arquivo removido. Campo de valores limpo.");
  };

  // Auto-hide para mensagens
  useEffect(() => {
    if (msg) {
      const timer = setTimeout(() => {
        setMsg(null);
      }, 4000); // 4 segundos
      return () => clearTimeout(timer);
    }
  }, [msg]);

  useEffect(() => {
    if (erro) {
      const timer = setTimeout(() => {
        setErro(null);
      }, 6000); // 6 segundos (erro fica mais tempo)
      return () => clearTimeout(timer);
    }
  }, [erro]);

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        {/* Título */}
        <h1 className="text-4xl font-bold text-center mb-8 text-gray-800">
          💸 Calculadora de Preços
        </h1>

        {/* Campos de entrada lado a lado */}
        <div className="flex flex-col md:flex-row gap-4 mb-4">
          <div className="flex-1">
            <textarea
              value={precosTexto}
              onChange={(e) => setPrecosTexto(e.target.value)}
              placeholder="Digitar valores"
              rows={precosTexto.split('\n').length || 3}
              disabled={!!precosArquivo}
              className={`w-full min-h-[80px] max-h-[200px] px-4 py-3 border-2 border-gray-300 rounded-lg resize-y focus:outline-none focus:border-blue-500 flex items-start ${
                precosArquivo ? 'bg-gray-100 cursor-not-allowed text-gray-500' : ''
              }`}
              style={{ 
                verticalAlign: 'top',
                lineHeight: '1.5'
              }}
            />
            <p className="text-xs text-gray-500 mt-1">
              {precosAtuais.length} valor(es) detectado(s)
              {precosArquivo && " (valores do arquivo)"}
            </p>
          </div>
          <div className="flex-1">
            <input
              value={taxasTexto}
              onChange={(e) => setTaxasTexto(e.target.value)}
              placeholder="Digitar taxas (ex: 10, 7, 15)"
              className="w-full h-[80px] px-4 py-3 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 flex items-center"
            />
            <p className="text-xs text-gray-500 mt-1">
              Taxas em percentual (%)
            </p>
          </div>
        </div>

        {/* Botões */}
        <div className="flex flex-col sm:flex-row gap-4 mb-8 items-center justify-center">
          <input 
            ref={fileRef} 
            type="file" 
            accept=".txt" 
            onChange={onFileChange} 
            className="hidden" 
          />
          <button
            onClick={handleUploadClick}
            className="px-4 py-3 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
          >
            📁 Selecionar arquivo
          </button>
          <button
            onClick={handleCalcular}
            disabled={loading || !precosAtuais.length}
            className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400 transition-colors"
          >
            {loading ? "🔄 Calculando..." : "🧮 Calcular"}
          </button>
          {precosArquivo && (
            <button
              onClick={limparArquivo}
              className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors text-sm"
            >
              ❌ Limpar arquivo
            </button>
          )}
        </div>

        {/* Status Messages com animação */}
        <AnimatePresence>
          {erro && (
            <motion.div
              initial={{ opacity: 0, y: -20, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              transition={{ duration: 0.3, ease: "easeOut" }}
              className="mb-4 p-3 bg-red-100 border border-red-300 rounded-lg text-red-700 text-sm text-center"
            >
              ⚠️ {erro}
            </motion.div>
          )}
        </AnimatePresence>

        <AnimatePresence>
          {!erro && msg && (
            <motion.div
              initial={{ opacity: 0, y: -20, scale: 0.95 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -10, scale: 0.95 }}
              transition={{ duration: 0.3, ease: "easeOut" }}
              className="mb-4 p-3 bg-blue-100 border border-blue-300 rounded-lg text-blue-700 text-sm text-center"
            >
              ✅ {msg}
            </motion.div>
          )}
        </AnimatePresence>

        {/* Tabela para organizar os valores pós taxa */}
        <div className="bg-white border-2 border-gray-300 rounded-lg overflow-hidden">
          <div className="bg-gray-50 px-4 py-3 border-b border-gray-300">
            <h2 className="text-lg font-semibold text-gray-800 text-center">
              📊 Tabela para organizar os valores pós taxa
            </h2>
          </div>

          {!resultados || resultados.length === 0 ? (
            <div className="h-64 flex items-center justify-center text-gray-500">
              <div className="text-center">
                <div className="text-4xl mb-2">📈</div>
                <p className="font-medium">Nenhum resultado ainda</p>
                <p className="text-sm">Preencha os dados acima e clique em Calcular</p>
              </div>
            </div>
          ) : (
            <div className="space-y-6 p-6">
              {resultados.map((resultado, index) => (
                <TabelaResultado key={index} resultado={resultado} />
              ))}
            </div>
          )}
        </div>

        {/* Dica API */}
        <motion.section 
          initial={{ opacity: 0 }} 
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="mt-8 text-center text-xs text-gray-500"
        >
          {API_BASE ? (
            <p>
              🔗 Usando backend Go em <span className="font-mono text-gray-700 bg-gray-200 px-2 py-1 rounded">{API_BASE}</span>
            </p>
          ) : (
            <p>
              💻 Cálculo executado localmente no navegador. Configure <span className="font-mono bg-gray-200 px-2 py-1 rounded text-gray-700">VITE_API_BASE</span> para usar a API Go.
            </p>
          )}
        </motion.section>
      </div>
    </div>
  );
}

/* ========= Componente de Tabela de Resultado ========= */
function TabelaResultado({ resultado }: { resultado: ResultadoGo }) {
  const { taxa, precos, precos_incluindo_taxas } = resultado;
  const taxaPercent = (taxa * 100).toFixed(1);

  return (
    <div className="border border-gray-200 rounded-lg overflow-hidden">
      <div className="bg-blue-50 px-4 py-3 border-b border-gray-200 flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-gray-800">Taxa aplicada: {taxaPercent}%</h3>
          <p className="text-xs text-gray-600">{precos.length} produto(s)</p>
        </div>
        <button
          onClick={() => downloadJSON(`resultado_${Math.round(taxa * 100)}.json`, resultado)}
          className="px-3 py-1 bg-green-500 hover:bg-green-600 text-white rounded text-sm font-medium transition-colors"
        >
          📁 Download JSON
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left text-gray-700 font-medium border-r border-gray-200">
                Preço Original
              </th>
              <th className="px-4 py-2 text-left text-gray-700 font-medium border-r border-gray-200">
                Com Taxa ({taxaPercent}%)
              </th>
              <th className="px-4 py-2 text-left text-gray-700 font-medium">
                Diferença
              </th>
            </tr>
          </thead>
          <tbody>
            {precos.map((preco, i) => {
              const precoComTaxa = Number(precos_incluindo_taxas[i]);
              const diferenca = precoComTaxa - preco;
              return (
                <tr key={i} className="border-t border-gray-200 hover:bg-gray-50">
                  <td className="px-4 py-3 font-mono text-gray-800 border-r border-gray-200">
                    {formatCurrency(preco)}
                  </td>
                  <td className="px-4 py-3 font-mono text-green-600 border-r border-gray-200">
                    {formatCurrency(precoComTaxa)}
                  </td>
                  <td className="px-4 py-3 font-mono text-blue-600">
                    +{formatCurrency(diferenca)}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
