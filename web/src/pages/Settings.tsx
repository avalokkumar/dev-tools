import { useEffect, useState } from "react";
import { Settings as SettingsIcon, Sun, Moon, Monitor, Trash2 } from "lucide-react";
import { Card, CardHeader, Badge } from "../components/ui/Card";
import { Button } from "../components/ui/Button";
import { TOTAL_OPS, TOOLS, findTool } from "../lib/catalog";
import { cn } from "../lib/cn";
import { Link } from "react-router-dom";

type Theme = "light" | "dark" | "system";

export default function SettingsPage() {
  const [theme, setTheme] = useState<Theme>(() => {
    try {
      return (localStorage.getItem("devforge-theme") as Theme) || "system";
    } catch {
      return "system";
    }
  });
  const [favourites, setFavourites] = useState<string[]>(() => {
    try {
      const raw = localStorage.getItem("devforge-favourites");
      return raw ? JSON.parse(raw) : [];
    } catch {
      return [];
    }
  });

  useEffect(() => {
    try {
      localStorage.setItem("devforge-theme", theme);
    } catch {
      /* ignore */
    }
    if (theme === "dark") {
      document.documentElement.classList.add("dark");
    } else if (theme === "light") {
      document.documentElement.classList.remove("dark");
    } else {
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      document.documentElement.classList.toggle("dark", prefersDark);
    }
  }, [theme]);

  function clearFavs() {
    setFavourites([]);
    try {
      localStorage.setItem("devforge-favourites", "[]");
    } catch {
      /* ignore */
    }
  }

  const favTools = favourites.map((s) => findTool(s)).filter(Boolean);

  return (
    <div className="px-md md:px-lg py-md md:py-lg max-w-[1100px] mx-auto w-full flex flex-col gap-lg">
      <header className="flex items-start gap-md">
        <div className="h-12 w-12 rounded-md bg-primary text-on-primary flex items-center justify-center shrink-0">
          <SettingsIcon className="h-6 w-6" />
        </div>
        <div>
          <h1 className="font-display-brand text-display-brand text-on-surface leading-none">
            Settings
          </h1>
          <p className="font-body-md text-body-md text-on-surface-variant mt-1">
            Local preferences. Stored in your browser.
          </p>
        </div>
      </header>

      <Card>
        <CardHeader title="Appearance" />
        <div className="flex flex-col gap-md">
          <div>
            <label className="font-data-label text-data-label uppercase text-on-surface-variant mb-2 block">
              Theme
            </label>
            <div className="grid grid-cols-3 gap-2 max-w-md">
              {[
                { id: "light" as Theme, label: "Light", Icon: Sun },
                { id: "dark" as Theme, label: "Dark", Icon: Moon },
                { id: "system" as Theme, label: "System", Icon: Monitor },
              ].map((t) => (
                <button
                  key={t.id}
                  onClick={() => setTheme(t.id)}
                  className={cn(
                    "flex flex-col items-center gap-2 p-3 rounded border transition-colors",
                    theme === t.id
                      ? "border-secondary-container bg-secondary-container/10 text-on-surface"
                      : "border-outline/20 hover:border-outline/40 text-on-surface-variant",
                  )}
                >
                  <t.Icon className="h-5 w-5" />
                  <span className="font-body-sm text-body-sm">{t.label}</span>
                </button>
              ))}
            </div>
          </div>
        </div>
      </Card>

      <Card>
        <CardHeader
          title="Pinned tools"
          trailing={
            favTools.length > 0 && (
              <Button variant="outline" size="sm" onClick={clearFavs}>
                <Trash2 className="h-4 w-4" /> Clear all
              </Button>
            )
          }
        />
        {favTools.length === 0 ? (
          <div className="font-body-sm text-body-sm text-on-surface-variant">
            No pinned tools yet — visit the home dashboard and click the star icon on any tool card.
          </div>
        ) : (
          <ul className="grid grid-cols-1 md:grid-cols-2 gap-2">
            {favTools.map((tool) => {
              const t = tool!;
              const Icon = t.icon;
              return (
                <li key={t.slug}>
                  <Link
                    to={`/tools/${t.slug}`}
                    className="flex items-center gap-2 px-3 py-2 bg-surface-container-low rounded hover:bg-surface-container transition-colors"
                  >
                    <Icon className="h-4 w-4 text-on-surface-variant" />
                    <span className="font-body-sm text-body-sm text-on-surface flex-1">{t.name}</span>
                    <Badge>{t.category}</Badge>
                  </Link>
                </li>
              );
            })}
          </ul>
        )}
      </Card>

      <Card>
        <CardHeader title="Telemetry" />
        <p className="font-body-sm text-body-sm text-on-surface-variant">
          Anonymous usage telemetry is opt-in and stored locally at
          <code className="font-code-block bg-surface-container-low px-1 mx-1 rounded">~/.devforge/events.jsonl</code>.
          No data is uploaded. Toggle via the <code className="font-code-block bg-surface-container-low px-1 mx-1 rounded">DEVFORGE_TELEMETRY</code> env var or the
          {" "}
          <code className="font-code-block bg-surface-container-low px-1 mx-1 rounded">devforge telemetry</code> CLI command.
        </p>
      </Card>

      <Card>
        <CardHeader title="About" />
        <dl className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <Item label="Version" value="0.1.0-dev" />
          <Item label="Tools" value={String(TOOLS.length)} />
          <Item label="Operations" value={String(TOTAL_OPS)} />
          <Item label="Surfaces" value="CLI · Web · MCP" />
        </dl>
      </Card>
    </div>
  );
}

function Item({ label, value }: { label: string; value: string }) {
  return (
    <div className="px-3 py-2 bg-surface-container-low rounded">
      <dt className="font-data-label text-data-label uppercase text-on-surface-variant">{label}</dt>
      <dd className="font-body-md text-body-md text-on-surface">{value}</dd>
    </div>
  );
}
