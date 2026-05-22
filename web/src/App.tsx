import { lazy, Suspense } from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "./components/layout/AppShell";

const Home = lazy(() => import("./pages/Home"));
const CliReference = lazy(() => import("./pages/CliReference"));
const McpConfig = lazy(() => import("./pages/McpConfig"));
const Settings = lazy(() => import("./pages/Settings"));

// Generators
const UuidPage = lazy(() => import("./pages/tools/Uuid"));
const FakerPage = lazy(() => import("./pages/tools/Faker"));
const IdPage = lazy(() => import("./pages/tools/Id"));
const TotpPage = lazy(() => import("./pages/tools/Totp"));
const CryptoPage = lazy(() => import("./pages/tools/Crypto"));

// Formatters
const JsonPage = lazy(() => import("./pages/tools/Json"));
const YamlPage = lazy(() => import("./pages/tools/Yaml"));
const CsvPage = lazy(() => import("./pages/tools/Csv"));
const SqlPage = lazy(() => import("./pages/tools/Sql"));
const CodePage = lazy(() => import("./pages/tools/Code"));
const MarkdownPage = lazy(() => import("./pages/tools/Markdown"));

// Converters
const EncodingPage = lazy(() => import("./pages/tools/Encoding"));
const DataPage = lazy(() => import("./pages/tools/Data"));
const ColorPage = lazy(() => import("./pages/tools/Color"));
const TimePage = lazy(() => import("./pages/tools/Time"));
const TimezonePage = lazy(() => import("./pages/tools/Timezone"));
const MathPage = lazy(() => import("./pages/tools/Math"));

// Analyzers
const DiffPage = lazy(() => import("./pages/tools/Diff"));
const RegexPage = lazy(() => import("./pages/tools/Regex"));
const CronPage = lazy(() => import("./pages/tools/Cron"));
const JwtPage = lazy(() => import("./pages/tools/Jwt"));
const StringPage = lazy(() => import("./pages/tools/StringTools"));
const UrlPage = lazy(() => import("./pages/tools/Url"));
const HeadersPage = lazy(() => import("./pages/tools/Headers"));
const DnsPage = lazy(() => import("./pages/tools/Dns"));
const HttpPage = lazy(() => import("./pages/tools/Http"));
const IpPage = lazy(() => import("./pages/tools/Ip"));

// DevOps
const GitPage = lazy(() => import("./pages/tools/Git"));
const DockerfilePage = lazy(() => import("./pages/tools/DockerfileTool"));
const EnvPage = lazy(() => import("./pages/tools/Env"));
const K8sPage = lazy(() => import("./pages/tools/K8s"));

function PageFallback() {
  return (
    <div className="px-md md:px-lg py-md md:py-lg">
      <div className="max-w-md skeleton h-8 rounded mb-3" />
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
        <div className="skeleton h-64 rounded" />
        <div className="skeleton h-64 rounded" />
      </div>
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AppShell>
        <Suspense fallback={<PageFallback />}>
          <Routes>
            <Route path="/" element={<Home />} />

            {/* Tool routes */}
            <Route path="/tools/uuid" element={<UuidPage />} />
            <Route path="/tools/faker" element={<FakerPage />} />
            <Route path="/tools/id" element={<IdPage />} />
            <Route path="/tools/totp" element={<TotpPage />} />
            <Route path="/tools/crypto" element={<CryptoPage />} />

            <Route path="/tools/json" element={<JsonPage />} />
            <Route path="/tools/yaml" element={<YamlPage />} />
            <Route path="/tools/csv" element={<CsvPage />} />
            <Route path="/tools/sql" element={<SqlPage />} />
            <Route path="/tools/code" element={<CodePage />} />
            <Route path="/tools/markdown" element={<MarkdownPage />} />

            <Route path="/tools/encoding" element={<EncodingPage />} />
            <Route path="/tools/data" element={<DataPage />} />
            <Route path="/tools/color" element={<ColorPage />} />
            <Route path="/tools/time" element={<TimePage />} />
            <Route path="/tools/timezone" element={<TimezonePage />} />
            <Route path="/tools/math" element={<MathPage />} />

            <Route path="/tools/diff" element={<DiffPage />} />
            <Route path="/tools/regex" element={<RegexPage />} />
            <Route path="/tools/cron" element={<CronPage />} />
            <Route path="/tools/jwt" element={<JwtPage />} />
            <Route path="/tools/string" element={<StringPage />} />
            <Route path="/tools/url" element={<UrlPage />} />
            <Route path="/tools/headers" element={<HeadersPage />} />
            <Route path="/tools/dns" element={<DnsPage />} />
            <Route path="/tools/http" element={<HttpPage />} />
            <Route path="/tools/ip" element={<IpPage />} />

            <Route path="/tools/git" element={<GitPage />} />
            <Route path="/tools/dockerfile" element={<DockerfilePage />} />
            <Route path="/tools/env" element={<EnvPage />} />
            <Route path="/tools/k8s" element={<K8sPage />} />

            {/* Special pages */}
            <Route path="/cli" element={<CliReference />} />
            <Route path="/mcp" element={<McpConfig />} />
            <Route path="/settings" element={<Settings />} />

            {/* Legacy redirects */}
            <Route path="/uuid" element={<Navigate to="/tools/uuid" replace />} />
            <Route path="/json" element={<Navigate to="/tools/json" replace />} />
            <Route path="/diff" element={<Navigate to="/tools/diff" replace />} />
            <Route path="/regex" element={<Navigate to="/tools/regex" replace />} />
            <Route path="/cron" element={<Navigate to="/tools/cron" replace />} />
            <Route path="/jwt" element={<Navigate to="/tools/jwt" replace />} />
            <Route path="/tz" element={<Navigate to="/tools/timezone" replace />} />
            <Route path="/faker" element={<Navigate to="/tools/faker" replace />} />

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Suspense>
      </AppShell>
    </BrowserRouter>
  );
}
