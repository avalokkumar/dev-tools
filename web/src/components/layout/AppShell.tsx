import { useState, type ReactNode } from "react";
import { Sidebar } from "./Sidebar";
import { Header } from "./Header";

interface AppShellProps {
  children: ReactNode;
}

export function AppShell({ children }: AppShellProps) {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="min-h-screen bg-background text-on-background flex">
      <Sidebar isOpen={menuOpen} onClose={() => setMenuOpen(false)} />
      <div className="flex-1 flex flex-col xl:ml-sidebar-width min-w-0">
        <Header onMenuToggle={() => setMenuOpen((v) => !v)} />
        <main className="flex-1 min-w-0">{children}</main>
        <footer className="border-t border-outline/10 px-md md:px-lg py-md text-on-surface-variant font-body-sm text-body-sm flex items-center justify-between gap-2">
          <span>
            DevForge — local-first developer toolkit. CLI · Web · MCP.
          </span>
          <span className="font-data-label text-data-label uppercase">v0.1</span>
        </footer>
      </div>
    </div>
  );
}
