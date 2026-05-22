import { useState } from "react";
import { Link, NavLink, useLocation } from "react-router-dom";
import {
  ChevronDown,
  ChevronRight,
  Settings,
  BookOpen,
  Terminal,
  Bot,
  Wrench,
  AlignLeft,
  Repeat2,
  ScanSearch,
  Cog,
  type LucideIcon,
} from "lucide-react";
import { CATEGORIES, TOOLS_BY_CATEGORY, type Category } from "../../lib/catalog";
import { cn } from "../../lib/cn";

const CATEGORY_ICONS: Record<Category, LucideIcon> = {
  Generators: Wrench,
  Formatters: AlignLeft,
  Converters: Repeat2,
  Analyzers: ScanSearch,
  DevOps: Cog,
};

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({ isOpen, onClose }: SidebarProps) {
  const location = useLocation();
  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>(() => {
    // open the group containing the active route
    const init: Record<string, boolean> = {};
    CATEGORIES.forEach((c) => {
      init[c.key] = TOOLS_BY_CATEGORY[c.key].some((t) =>
        location.pathname.startsWith(`/tools/${t.slug}`),
      );
    });
    // default: open Generators when nothing active
    if (!Object.values(init).some(Boolean)) init.Generators = true;
    return init;
  });

  function toggle(key: string) {
    setOpenGroups((g) => ({ ...g, [key]: !g[key] }));
  }

  return (
    <>
      {isOpen && (
        <button
          aria-label="Close menu"
          className="xl:hidden fixed inset-0 z-30 bg-primary/40 backdrop-blur-sm"
          onClick={onClose}
        />
      )}
      <nav
        aria-label="Primary"
        className={cn(
          "fixed left-0 top-0 z-40 h-full w-sidebar-width bg-primary text-on-primary border-r border-on-primary/10 flex flex-col py-md transition-transform xl:translate-x-0",
          isOpen ? "translate-x-0" : "-translate-x-full xl:translate-x-0",
        )}
      >
        {/* Brand */}
        <Link to="/" onClick={onClose} className="px-lg pb-md pt-sm shrink-0 group">
          <h1 className="font-display-brand text-display-brand text-on-primary leading-none group-hover:text-secondary-container transition-colors">
            DevForge
          </h1>
          <p className="font-body-sm text-body-sm text-on-primary/70 mt-1">Developer Workspace</p>
        </Link>

        {/* Tools */}
        <div className="flex-1 flex flex-col gap-1 px-sm overflow-y-auto thin-scrollbar pb-md">
          <p className="font-data-label text-data-label text-on-primary/50 uppercase px-sm pt-sm pb-1">
            Tools
          </p>

          {CATEGORIES.map((c) => {
            const Icon = CATEGORY_ICONS[c.key];
            const tools = TOOLS_BY_CATEGORY[c.key];
            const open = !!openGroups[c.key];
            const groupActive = tools.some((t) =>
              location.pathname.startsWith(`/tools/${t.slug}`),
            );
            return (
              <div key={c.key} className="flex flex-col">
                <button
                  onClick={() => toggle(c.key)}
                  className={cn(
                    "flex items-center gap-2 px-sm py-2 rounded-md text-left transition-colors",
                    groupActive
                      ? "text-on-primary bg-primary-container"
                      : "text-on-primary/70 hover:text-on-primary hover:bg-primary-container/50",
                  )}
                >
                  <Icon className="h-5 w-5 shrink-0" />
                  <span className="flex-1 font-body-md text-body-md">{c.key}</span>
                  {open ? (
                    <ChevronDown className="h-4 w-4 opacity-60" />
                  ) : (
                    <ChevronRight className="h-4 w-4 opacity-60" />
                  )}
                </button>
                {open && (
                  <ul className="ml-3 border-l border-on-primary/10 pl-2 mt-0.5 mb-1 flex flex-col">
                    {tools.map((t) => {
                      const ToolIcon = t.icon;
                      return (
                        <li key={t.slug}>
                          <NavLink
                            to={`/tools/${t.slug}`}
                            onClick={onClose}
                            className={({ isActive }) =>
                              cn(
                                "relative flex items-center gap-2 px-sm py-1.5 rounded text-body-sm font-body-sm transition-colors",
                                isActive
                                  ? "bg-secondary-container text-on-secondary"
                                  : "text-on-primary/70 hover:text-on-primary hover:bg-primary-container/50",
                              )
                            }
                          >
                            <ToolIcon className="h-4 w-4 shrink-0 opacity-80" />
                            <span className="truncate">{t.name}</span>
                          </NavLink>
                        </li>
                      );
                    })}
                  </ul>
                )}
              </div>
            );
          })}
        </div>

        {/* Footer nav */}
        <div className="mt-auto px-sm flex flex-col gap-1 pt-md border-t border-on-primary/10">
          <FooterLink to="/cli" icon={Terminal} label="CLI Reference" onClick={onClose} />
          <FooterLink to="/mcp" icon={Bot} label="MCP Config" onClick={onClose} />
          <FooterLink to="/docs" icon={BookOpen} label="Documentation" onClick={onClose} external />
          <FooterLink to="/settings" icon={Settings} label="Settings" onClick={onClose} />
        </div>
      </nav>
    </>
  );
}

function FooterLink({
  to,
  icon: Icon,
  label,
  onClick,
  external,
}: {
  to: string;
  icon: LucideIcon;
  label: string;
  onClick?: () => void;
  external?: boolean;
}) {
  if (external) {
    return (
      <a
        href="https://github.com"
        target="_blank"
        rel="noreferrer"
        className="flex items-center gap-2 px-sm py-2 rounded-md text-on-primary/70 hover:text-on-primary hover:bg-primary-container/50 transition-colors"
      >
        <Icon className="h-5 w-5" />
        <span className="font-body-md text-body-md">{label}</span>
      </a>
    );
  }
  return (
    <NavLink
      to={to}
      onClick={onClick}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-2 px-sm py-2 rounded-md transition-colors",
          isActive
            ? "bg-primary-container text-on-primary"
            : "text-on-primary/70 hover:text-on-primary hover:bg-primary-container/50",
        )
      }
    >
      <Icon className="h-5 w-5" />
      <span className="font-body-md text-body-md">{label}</span>
    </NavLink>
  );
}
