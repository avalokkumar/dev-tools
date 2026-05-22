import type { ReactNode } from "react";
import { cn } from "../../lib/cn";

interface CardProps {
  children: ReactNode;
  className?: string;
  padded?: boolean;
}
export function Card({ children, className, padded = true }: CardProps) {
  return (
    <div
      className={cn(
        "bg-surface-container-lowest border border-outline/10 rounded-lg shadow-sm",
        padded && "p-md",
        className,
      )}
    >
      {children}
    </div>
  );
}

interface CardHeaderProps {
  title: string;
  subtitle?: string;
  trailing?: ReactNode;
  icon?: ReactNode;
}
export function CardHeader({ title, subtitle, trailing, icon }: CardHeaderProps) {
  return (
    <div className="flex items-start justify-between gap-md pb-md border-b border-outline/10 mb-md">
      <div className="flex items-start gap-sm min-w-0">
        {icon && <div className="text-secondary-container shrink-0 mt-0.5">{icon}</div>}
        <div className="min-w-0">
          <h3 className="font-headline-page text-headline-page text-on-surface truncate">{title}</h3>
          {subtitle && (
            <p className="font-body-sm text-body-sm text-on-surface-variant mt-1">{subtitle}</p>
          )}
        </div>
      </div>
      {trailing && <div className="flex items-center gap-2 shrink-0">{trailing}</div>}
    </div>
  );
}

interface SectionProps {
  title?: string;
  description?: string;
  trailing?: ReactNode;
  children: ReactNode;
  className?: string;
}
export function Section({ title, description, trailing, children, className }: SectionProps) {
  return (
    <section className={cn("mb-lg", className)}>
      {(title || trailing) && (
        <div className="flex items-end justify-between gap-md mb-md">
          <div>
            {title && (
              <h3 className="font-headline-page text-headline-page text-on-surface">{title}</h3>
            )}
            {description && (
              <p className="font-body-sm text-body-sm text-on-surface-variant mt-1">{description}</p>
            )}
          </div>
          {trailing}
        </div>
      )}
      {children}
    </section>
  );
}

interface BadgeProps {
  children: ReactNode;
  tone?: "neutral" | "success" | "warning" | "info" | "error";
  className?: string;
}
export function Badge({ children, tone = "neutral", className }: BadgeProps) {
  const tones: Record<string, string> = {
    neutral: "bg-surface-container border border-outline/10 text-on-surface-variant",
    success: "bg-tertiary-fixed text-on-tertiary-fixed",
    warning: "bg-error-container text-on-error-container",
    info: "bg-sky-aqua/15 text-primary",
    error: "bg-error text-on-error",
  };
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 px-2 py-0.5 rounded font-data-label text-data-label uppercase",
        tones[tone],
        className,
      )}
    >
      {children}
    </span>
  );
}
