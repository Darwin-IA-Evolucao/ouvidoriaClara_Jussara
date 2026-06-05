import { createRequire } from "module";
import { writeFileSync } from "fs";
import { resolve, dirname } from "path";
import { fileURLToPath } from "url";

const require = createRequire(import.meta.url);
const XLSX = require("xlsx");

const __dirname = dirname(fileURLToPath(import.meta.url));
const input = resolve(__dirname, "CEP_Sorocaba_regiao_preenchida.xlsx");
const output = resolve(__dirname, "enderecos_insert.sql");

function sqlEscape(value) {
  return String(value ?? "").replace(/'/g, "''");
}

function toSqlValue(value) {
  return `'${sqlEscape(value)}'`;
}

const wb = XLSX.readFile(input);
const ws = wb.Sheets[wb.SheetNames[0]];
const rows = XLSX.utils.sheet_to_json(ws);

const chunkSize = 500;
const lines = [
  "-- Gerado por convert_enderecos.mjs",
  `-- Total de registros: ${rows.length}`,
  "",
];

const validRows = rows.filter((row) => row.LOGRADOURO);

for (let i = 0; i < validRows.length; i += chunkSize) {
  const chunk = validRows.slice(i, i + chunkSize);
  const values = chunk
    .map((row) => {
      const logradouro = row.LOGRADOURO ?? "";
      const bairro = row.BAIRRO ?? "";
      const regiao = row["Região"] ?? "";
      return `(${toSqlValue(logradouro)}, ${toSqlValue(bairro)}, ${toSqlValue(regiao)})`;
    })
    .join(",\n");

  lines.push("INSERT INTO enderecos (logradouro, bairro, regiao) VALUES");
  lines.push(values + ";");
  lines.push("");
}

writeFileSync(output, lines.join("\n"), "utf8");
console.log(`Gerado: ${output} (${validRows.length} registros)`);
