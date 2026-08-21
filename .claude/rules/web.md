---
paths:
  - "web/src/**"
  - "web/*.ts"
  - "web/*.json"
  - "web/.oxlintrc.json"
---

# Web stack

- React 19 + TypeScript, bundled with Vite 8 (rolldown-vite) and `@vitejs/plugin-react`. Routing is `react-router-dom` (`BrowserRouter` in `main.tsx`, `Routes`/`Route` in `App.tsx`).
- shadcn/ui is installed (`--template vite -b base`, Base UI primitives, Nova preset): `components.json`, `src/components/ui/`, `src/lib/utils.ts`. The `@` path alias (`./src`) is configured in both `vite.config.ts` and the `tsconfig.*.json` files' `paths` (no `baseUrl` — it's deprecated under this TS version's `moduleResolution: "bundler"`, and `paths` resolves relative to the tsconfig without it).
- Tailwind CSS v4 via `@tailwindcss/vite`, imported in `src/index.css`.
- `src/index.css` carries two token systems that both need to stay in sync: the site's own tokens (`--text`, `--text-h`, `--bg`, `--code-bg`) and shadcn's tokens (`--background`, `--foreground`, `--border`, `--accent`, etc., which shadcn normally toggles via a `.dark` class). Since this project has no theme-toggle mechanism, both sets are switched together by the same plain `@media (prefers-color-scheme: dark)` block — the `.dark` class's values are duplicated into it so shadcn components also follow the OS theme automatically. If a manual toggle is ever added, that duplication should be reconsidered.
- Linting is via Oxlint (`web/.oxlintrc.json`), not ESLint. Type-aware lint rules are not enabled by default; see `web/README.md` for how to add them via `oxlint-tsgolint` if needed.
- `tsconfig.json` is a project-references root pointing at `tsconfig.app.json` (app code) and `tsconfig.node.json` (Vite config).
