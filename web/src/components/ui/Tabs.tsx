import { type ReactNode } from "react";
import { cn } from "../../lib/cn";

interface Tab {
  id: string;
  label: string;
  icon?: ReactNode;
}

interface TabsProps {
  tabs: Tab[];
  active: string;
  onChange: (id: string) => void;
  className?: string;
}

export function Tabs({ tabs, active, onChange, className }: TabsProps) {
  return (
    <div
      role="tablist"
      className={cn(
        "flex items-center gap-1 border-b border-outline/10 overflow-x-auto hide-scrollbar",
        className,
      )}
    >
      {tabs.map((t) => {
        const isActive = t.id === active;
        return (
          <button
            key={t.id}
            role="tab"
            aria-selected={isActive}
            onClick={() => onChange(t.id)}
            className={cn(
              "relative flex items-center gap-2 px-4 py-2.5 font-body-sm text-body-sm whitespace-nowrap transition-colors",
              "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-aqua rounded-t",
              isActive
                ? "text-on-surface font-semibold"
                : "text-on-surface-variant hover:text-on-surface",
            )}
          >
            {t.icon}
            <span>{t.label}</span>
            {isActive && (
              <span className="absolute inset-x-3 -bottom-px h-0.5 bg-secondary-container rounded-full" />
            )}
          </button>
        );
      })}
    </div>
  );
}

interface PillTabsProps extends TabsProps {}

export function PillTabs({ tabs, active, onChange, className }: PillTabsProps) {
  return (
    <div className={cn("flex items-center gap-2 flex-wrap", className)}>
      {tabs.map((t) => {
        const isActive = t.id === active;
        return (
          <button
            key={t.id}
            onClick={() => onChange(t.id)}
            className={cn(
              "inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full font-body-sm text-body-sm transition-colors border",
              isActive
                ? "bg-primary text-on-primary border-primary"
                : "bg-surface-container border-outline/10 text-on-surface hover:border-outline/30",
            )}
          >
            {t.icon}
            <span>{t.label}</span>
          </button>
        );
      })}
    </div>
  );
}
