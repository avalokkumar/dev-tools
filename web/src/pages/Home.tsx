import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { ArrowRight, Search, Sparkles, Star } from "lucide-react";
import {
  CATEGORIES,
  TOOLS,
  TOOLS_BY_CATEGORY,
  TOTAL_OPS,
  type Category,
} from "../lib/catalog";
import { cn } from "../lib/cn";
import { PillTabs } from "../components/ui/Tabs";

const FAVOURITES_KEY = "devforge-favourites";
const DEFAULT_FAVS = ["uuid", "json", "jwt", "diff"];

function loadFavs(): string[] {
  try {
    const raw = localStorage.getItem(FAVOURITES_KEY);
    if (!raw) return DEFAULT_FAVS;
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed) ? parsed : DEFAULT_FAVS;
  } catch {
    return DEFAULT_FAVS;
  }
}

export default function Home() {
  const [filter, setFilter] = useState<Category | "All">("All");
  const [query, setQuery] = useState("");
  const [favs, setFavs] = useState<string[]>(loadFavs);

  useEffect(() => {
    try {
      localStorage.setItem(FAVOURITES_KEY, JSON.stringify(favs));
    } catch {
      /* ignore */
    }
  }, [favs]);

  function toggleFav(slug: string) {
    setFavs((f) => (f.includes(slug) ? f.filter((s) => s !== slug) : [...f, slug]));
  }

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    let list = filter === "All" ? TOOLS : TOOLS_BY_CATEGORY[filter];
    if (q) {
      list = list.filter(
        (t) =>
          t.name.toLowerCase().includes(q) ||
          t.tagline.toLowerCase().includes(q) ||
          t.category.toLowerCase().includes(q),
      );
    }
    return list;
  }, [filter, query]);

  const favTools = favs.map((s) => TOOLS.find((t) => t.slug === s)).filter(Boolean) as typeof TOOLS;

  return (
    <div className="px-md md:px-lg py-md md:py-xl flex flex-col gap-xl max-w-[1400px] mx-auto w-full">
      {/* Hero */}
      <section className="text-center max-w-3xl mx-auto w-full mt-md md:mt-lg">
        <h2 className="font-display-brand text-display-brand text-on-surface mb-1">
          {TOTAL_OPS} operations. {TOOLS.length} tools. <span className="text-secondary-container">One forge.</span>
        </h2>
        <p className="font-body-md text-body-md text-on-surface-variant mb-md">
          A local-first developer toolkit available as CLI, Web, and MCP server.
        </p>
        <div className="relative w-full max-w-2xl mx-auto group">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-secondary-container h-6 w-6" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="What do you need to build today? Try 'JSON', 'JWT', or 'UUID'…"
            className="w-full bg-surface-container-lowest border border-outline/20 rounded-lg font-body-lg text-body-lg text-on-surface pl-14 pr-24 py-4 shadow-sm focus:outline-none focus:ring-2 focus:ring-secondary-container focus:border-secondary-container transition-all"
          />
          <div className="absolute right-3 top-1/2 -translate-y-1/2 hidden md:flex gap-1 pointer-events-none">
            <kbd className="bg-surface-container border border-outline/10 text-on-surface-variant font-data-label text-data-label px-2 py-1 rounded">⌘</kbd>
            <kbd className="bg-surface-container border border-outline/10 text-on-surface-variant font-data-label text-data-label px-2 py-1 rounded">K</kbd>
          </div>
        </div>
      </section>

      {/* Pinned tools */}
      {favTools.length > 0 && (
        <section>
          <div className="flex items-center justify-between mb-md">
            <div className="flex items-center gap-2">
              <Star className="h-5 w-5 text-secondary-container" />
              <h3 className="font-headline-page text-headline-page text-on-surface">Pinned Tools</h3>
            </div>
          </div>
          <div className="flex gap-md overflow-x-auto pb-2 hide-scrollbar snap-x">
            {favTools.map((t) => {
              const Icon = t.icon;
              return (
                <Link
                  key={t.slug}
                  to={`/tools/${t.slug}`}
                  className="bg-surface-container-lowest border border-outline/10 rounded-lg p-md shadow-sm min-w-[220px] snap-start hover:border-secondary-container transition-colors group"
                >
                  <div className="flex items-center gap-sm mb-2">
                    <Icon className="h-5 w-5 text-on-surface-variant group-hover:text-secondary-container transition-colors" />
                    <span className="font-body-md text-body-md font-medium text-on-surface">
                      {t.name}
                    </span>
                  </div>
                  <p className="font-body-sm text-body-sm text-on-surface-variant line-clamp-1">
                    {t.tagline}
                  </p>
                </Link>
              );
            })}
          </div>
        </section>
      )}

      {/* Tool grid */}
      <section>
        <div className="flex items-center justify-between mb-md gap-md flex-wrap">
          <div className="flex items-center gap-2">
            <Sparkles className="h-5 w-5 text-on-surface-variant" />
            <h3 className="font-headline-page text-headline-page text-on-surface">All Tools</h3>
          </div>
          <PillTabs
            tabs={[{ id: "All", label: "All Tools" }, ...CATEGORIES.map((c) => ({ id: c.key, label: c.key }))]}
            active={filter}
            onChange={(id) => setFilter(id as typeof filter)}
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-md">
          {filtered.map((t) => {
            const Icon = t.icon;
            const isFav = favs.includes(t.slug);
            return (
              <Link
                key={t.slug}
                to={`/tools/${t.slug}`}
                className="bg-surface-container-lowest border border-outline/10 rounded-lg shadow-sm hover:border-secondary-container transition-colors group flex flex-col h-full overflow-hidden relative"
              >
                <div className="p-4 flex-1">
                  <div className="flex items-start justify-between mb-md">
                    <div className="bg-surface-container-low p-2 rounded text-on-surface-variant group-hover:text-secondary-container group-hover:bg-secondary-container/10 transition-colors">
                      <Icon className="h-6 w-6" />
                    </div>
                    <div className="flex items-center gap-1">
                      <button
                        onClick={(e) => {
                          e.preventDefault();
                          toggleFav(t.slug);
                        }}
                        aria-label={isFav ? "Unpin" : "Pin"}
                        className="p-1 text-on-surface-variant hover:text-secondary-container transition-colors"
                      >
                        <Star
                          className={cn(
                            "h-4 w-4 transition-colors",
                            isFav && "fill-secondary-container text-secondary-container",
                          )}
                        />
                      </button>
                      <span className="bg-surface-container border border-outline/10 text-on-surface-variant font-data-label text-data-label px-2 py-1 rounded uppercase">
                        {t.category}
                      </span>
                    </div>
                  </div>
                  <h4 className="font-body-md text-body-md font-medium text-on-surface mb-1">
                    {t.name}
                  </h4>
                  <p className="font-body-sm text-body-sm text-on-surface-variant">{t.tagline}</p>
                </div>
                <div className="px-4 py-3 border-t border-outline/10 bg-surface-bright flex justify-between items-center">
                  <span className="font-data-label text-data-label text-on-surface-variant">
                    {t.ops.length} {t.ops.length === 1 ? "op" : "ops"}
                  </span>
                  <ArrowRight className="h-4 w-4 text-on-surface-variant group-hover:text-secondary-container transition-colors" />
                </div>
              </Link>
            );
          })}
        </div>

        {filtered.length === 0 && (
          <div className="text-center py-xl text-on-surface-variant">
            No tools match "{query}".
          </div>
        )}
      </section>
    </div>
  );
}
