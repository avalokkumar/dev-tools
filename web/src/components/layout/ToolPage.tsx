import { type ReactNode } from "react";
import { Link } from "react-router-dom";
import { ChevronRight, Home } from "lucide-react";
import type { ToolMeta } from "../../lib/catalog";
import { Badge } from "../ui/Card";

interface ToolPageProps {
  tool: ToolMeta;
  children: ReactNode;
  trailing?: ReactNode;
}

export function ToolPage({ tool, children, trailing }: ToolPageProps) {
  const Icon = tool.icon;
  return (
    <div className="px-md md:px-lg py-md md:py-lg flex flex-col gap-lg max-w-[1400px] mx-auto w-full">
      {/* Breadcrumbs */}
      <nav aria-label="Breadcrumb" className="flex items-center gap-1.5 text-on-surface-variant font-body-sm text-body-sm">
        <Link to="/" className="inline-flex items-center gap-1 hover:text-on-surface transition-colors">
          <Home className="h-3.5 w-3.5" />
          <span>Home</span>
        </Link>
        <ChevronRight className="h-3.5 w-3.5 opacity-60" />
        <span className="text-on-surface-variant">{tool.category}</span>
        <ChevronRight className="h-3.5 w-3.5 opacity-60" />
        <span className="text-on-surface font-medium">{tool.name}</span>
      </nav>

      {/* Page header */}
      <header className="flex flex-wrap items-start justify-between gap-md">
        <div className="flex items-start gap-md min-w-0">
          <div className="shrink-0 h-12 w-12 rounded-md bg-primary text-on-primary flex items-center justify-center">
            <Icon className="h-6 w-6" />
          </div>
          <div className="min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="font-display-brand text-display-brand text-on-surface leading-none">
                {tool.name}
              </h1>
              <Badge>{tool.category}</Badge>
            </div>
            <p className="font-body-md text-body-md text-on-surface-variant mt-1">{tool.tagline}</p>
          </div>
        </div>
        {trailing && <div className="flex items-center gap-2">{trailing}</div>}
      </header>

      <div>{children}</div>
    </div>
  );
}

interface SplitPaneProps {
  left: ReactNode;
  right: ReactNode;
}
export function SplitPane({ left, right }: SplitPaneProps) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-md">
      <div className="min-w-0">{left}</div>
      <div className="min-w-0">{right}</div>
    </div>
  );
}
