# Crompressor WASM

> Motor de compressão [Crompressor](https://github.com/MrJc01/crompressor) compilado para WebAssembly — roda no navegador.

---

## O que é?

O **Crompressor WASM** permite executar o motor de compressão CROM diretamente no navegador,
sem servidor backend. Ideal para demos, playgrounds e integração com aplicações web.

### Funcionalidades

- `cromPack()` — Comprimir dados no navegador
- `cromUnpack()` — Descomprimir dados no navegador
- `cromInfo()` — Informações do motor

---

## Build

```bash
make build    # Compila para WASM
make serve    # Inicia servidor local em :8080
make clean    # Remove artefatos
```

**Pré-requisitos:** Go 1.22+ com suporte a `GOOS=js GOARCH=wasm`.

---

## Arquitetura

```
crompressor-wasm/
├── cmd/wasm/          ← Entrypoint WASM (main.go)
├── www/               ← Demo HTML + wasm_exec.js
│   ├── index.html
│   ├── crompressor.wasm  (gerado)
│   └── wasm_exec.js      (gerado)
├── go.mod
├── Makefile
└── README.md
```

---

## Ecossistema

| Repositório | Papel |
|-------------|-------|
| [crompressor](https://github.com/MrJc01/crompressor) | Motor core |
| [crompressor-wasm](https://github.com/MrJc01/crompressor-wasm) | WebAssembly (este repo) |
| [crompressor-gui](https://github.com/MrJc01/crompressor-gui) | Interface gráfica nativa |

---

## Licença

MIT — veja [LICENSE](LICENSE).
