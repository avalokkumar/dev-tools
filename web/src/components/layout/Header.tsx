import { useEffect, useState, useRef, useMemo } from "react";
import { useNavigate, useLocation, Link } from "react-router-dom";
import { Menu, Search, Settings, Sun, Moon, Command } from "lucide-react";
import { TOOLS } from "../../lib/catalog";
import { cn } from "../../lib/cn";

interface HeaderProps {
  onMenuToggle: () => void;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const [paletteOpen, setPaletteOpen] = useState(false);

  // Open palette on Cmd/Ctrl+K
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setPaletteOpen((v) => !v);
      } else if (e.key === "Escape") {
        setPaletteOpen(false);
      }
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  return (
    <>
      <header
        className={cn(
          "sticky top-0 z-20 h-header-height shrink-0",
          "bg-surface/90 backdrop-blur border-b border-outline/10",
          "flex items-center justify-between gap-md px-md md:px-lg",
        )}
      >
        <div className="flex items-center gap-md flex-1 min-w-0">
          <button
            onClick={onMenuToggle}
            aria-label="Toggle menu"
            className="xl:hidden text-on-surface hover:text-secondary-container transition-colors"
          >
            <Menu className="h-6 w-6" />
          </button>
          <button
            onClick={() => setPaletteOpen(true)}
            className="hidden md:flex items-center gap-2 w-full max-w-md bg-surface-container-low border border-outline/20 rounded px-3 py-2 text-on-surface-variant hover:border-outline/40 transition-colors text-left"
          >
            <Search className="h-4 w-4" />
            <span className="font-body-sm text-body-sm flex-1">Search 75+ tools…</span>
            <kbd className="font-data-label text-data-label bg-surface-container border border-outline/10 rounded px-1.5 py-0.5">
              ⌘K
            </kbd>
          </button>
          <button
            onClick={() => setPaletteOpen(true)}
            aria-label="Search"
            className="md:hidden text-on-surface-variant hover:text-on-surface"
          >
            <Search className="h-5 w-5" />
          </button>
        </div>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <Link
            to="/settings"
            aria-label="Settings"
            className="p-2 text-on-surface-variant hover:text-secondary-container hover:bg-surface-container-low rounded transition-colors"
          >
            <Settings className="h-5 w-5" />
          </Link>
        </div>
      </header>
      {paletteOpen && <CommandPalette onClose={() => setPaletteOpen(false)} />}
    </>
  );
}

function ThemeToggle() {
  const [dark, setDark] = useState(() =>
    typeof window !== "undefined" && document.documentElement.classList.contains("dark"),
  );
  function toggle() {
    const next = !dark;
    setDark(next);
    document.documentElement.classList.toggle("dark", next);
    try {
      localStorage.setItem("devforge-theme", next ? "dark" : "light");
    } catch {
      /* ignore */
    }
  }
  return (
    <button
      onClick={toggle}
      aria-label="Toggle theme"
      className="p-2 text-on-surface-variant hover:text-secondary-container hover:bg-surface-container-low rounded transition-colors"
    >
      {dark ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
    </button>
  );
}

interface CommandPaletteProps {
  onClose: () => void;
}
function CommandPalette({ onClose }: CommandPaletteProps) {
  const [query, setQuery] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);
  const navigate = useNavigate();
  const location = useLocation();
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  // Close on route change — but not on the initial mount (otherwise the
  // palette would open and immediately close itself).
  const initialPath = useRef(location.pathname);
  useEffect(() => {
    if (location.pathname !== initialPath.current) onClose();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname]);

  const results = useMemo(() => {
    const q = query.trim().toLowerCase();
    const items = TOOLS.map((t) => ({
      slug: t.slug,
      name: t.name,
      tagline: t.tagline,
      category: t.category,
      icon: t.icon,
    }));
    if (!q) return items.slice(0, 10);
    return items
      .filter(
        (t) =>
          t.name.toLowerCase().includes(q) ||
          t.tagline.toLowerCase().includes(q) ||
          t.category.toLowerCase().includes(q),
      )
      .slice(0, 12);
  }, [query]);

  function onKey(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActiveIndex((i) => Math.min(results.length - 1, i + 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActiveIndex((i) => Math.max(0, i - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      const sel = results[activeIndex];
      if (sel) navigate(`/tools/${sel.slug}`);
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 bg-primary/40 backdrop-blur-sm flex items-start justify-center pt-[10vh] px-4 animate-fade-in"
      onClick={onClose}
    >
      <div
        className="w-full max-w-xl bg-surface-container-lowest rounded-lg shadow-sm border border-outline/20 overflow-hidden"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center gap-3 px-4 py-3 border-b border-outline/10">
          <Command className="h-5 w-5 text-on-surface-variant" />
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              setActiveIndex(0);
            }}
            onKeyDown={onKey}
            placeholder="Search tools by name, category, or what you want to do…"
            className="flex-1 bg-transparent text-body-md font-body-md text-on-surface placeholder:text-on-surface-variant/60 outline-none"
          />
          <kbd className="font-data-label text-data-label bg-surface-container border border-outline/10 rounded px-1.5 py-0.5 text-on-surface-variant">
            esc
          </kbd>
        </div>
        <ul className="max-h-[60vh] overflow-y-auto thin-scrollbar py-1">
          {results.length === 0 ? (
            <li className="px-4 py-6 text-center text-on-surface-variant font-body-sm text-body-sm">
              No results for "{query}"
            </li>
          ) : (
            results.map((r, i) => {
              const Icon = r.icon;
              const isActive = i === activeIndex;
              return (
                <li key={r.slug}>
                  <button
                    onMouseEnter={() => setActiveIndex(i)}
                    onClick={() => navigate(`/tools/${r.slug}`)}
                    className={cn(
                      "w-full flex items-center gap-3 px-4 py-2.5 text-left transition-colors",
                      isActive ? "bg-surface-container-low" : "hover:bg-surface-container-low/60",
                    )}
                  >
                    <Icon className="h-5 w-5 text-on-surface-variant shrink-0" />
                    <div className="flex-1 min-w-0">
                      <div className="font-body-md text-body-md text-on-surface truncate">
                        {r.name}
                      </div>
                      <div className="font-body-sm text-body-sm text-on-surface-variant truncate">
                        {r.tagline}
                      </div>
                    </div>
                    <span className="font-data-label text-data-label uppercase text-on-surface-variant shrink-0">
                      {r.category}
                    </span>
                  </button>
                </li>
              );
            })
          )}
        </ul>
      </div>
    </div>
  );
}
