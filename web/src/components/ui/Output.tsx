import { useState, type ReactNode } from "react";
import { Check, Copy, Download, Inbox, AlertTriangle, Info, XCircle } from "lucide-react";
import { cn } from "../../lib/cn";
import type { Diagnostic } from "../../lib/api";

interface CopyButtonProps {
  text: string;
  label?: string;
  className?: string;
  size?: "sm" | "md";
}
export function CopyButton({ text, label = "Copy", className, size = "sm" }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);
  async function copy() {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* clipboard unavailable */
    }
  }
  return (
    <button
      type="button"
      onClick={copy}
      className={cn(
        "inline-flex items-center gap-1.5 rounded-sm border border-outline/20 bg-surface-container-lowest text-on-surface-variant hover:text-on-surface hover:border-outline/40 transition-colors",
        size === "sm" ? "px-2 py-1 text-xs" : "px-3 py-1.5 text-body-sm",
        className,
      )}
      aria-label={copied ? "Copied" : label}
    >
      {copied ? (
        <Check className="h-3.5 w-3.5 text-tertiary-container" />
      ) : (
        <Copy className="h-3.5 w-3.5" />
      )}
      <span className="font-data-label text-data-label uppercase">
        {copied ? "Copied" : label}
      </span>
    </button>
  );
}

interface DownloadButtonProps {
  text: string;
  filename: string;
  mime?: string;
  label?: string;
  className?: string;
}
export function DownloadButton({ text, filename, mime = "text/plain", label = "Download", className }: DownloadButtonProps) {
  function download() {
    const blob = new Blob([text], { type: mime });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  }
  return (
    <button
      type="button"
      onClick={download}
      className={cn(
        "inline-flex items-center gap-1.5 rounded-sm border border-outline/20 bg-surface-container-lowest text-on-surface-variant hover:text-on-surface hover:border-outline/40 transition-colors px-2 py-1 text-xs",
        className,
      )}
    >
      <Download className="h-3.5 w-3.5" />
      <span className="font-data-label text-data-label uppercase">{label}</span>
    </button>
  );
}

interface CodeBlockProps {
  code: string;
  language?: string;
  className?: string;
  showCopy?: boolean;
  download?: { filename: string; mime?: string };
  maxHeight?: string;
  trailing?: ReactNode;
}
export function CodeBlock({
  code,
  language,
  className,
  showCopy = true,
  download,
  maxHeight = "60vh",
  trailing,
}: CodeBlockProps) {
  return (
    <div
      className={cn(
        "relative bg-primary text-on-primary rounded border border-outline/10 overflow-hidden",
        className,
      )}
    >
      <div className="flex items-center justify-between px-3 py-2 border-b border-outline/10 bg-primary-container">
        <span className="font-data-label text-data-label uppercase text-on-primary/60">
          {language || "output"}
        </span>
        <div className="flex items-center gap-1">
          {trailing}
          {download && (
            <DownloadButton
              text={code}
              filename={download.filename}
              mime={download.mime}
              className="!bg-transparent !border-on-primary/20 !text-on-primary/70 hover:!text-on-primary"
            />
          )}
          {showCopy && (
            <CopyButton
              text={code}
              className="!bg-transparent !border-on-primary/20 !text-on-primary/70 hover:!text-on-primary"
            />
          )}
        </div>
      </div>
      <pre
        className="thin-scrollbar overflow-auto p-3 font-code-block text-code-block whitespace-pre"
        style={{ maxHeight }}
      >
        <code>{code}</code>
      </pre>
    </div>
  );
}

interface DiagnosticsProps {
  items?: Diagnostic[];
  className?: string;
}
export function Diagnostics({ items, className }: DiagnosticsProps) {
  if (!items || items.length === 0) return null;
  const order: Record<number, number> = { 2: 0, 1: 1, 0: 2 };
  const sorted = [...items].sort((a, b) => order[a.severity] - order[b.severity]);
  return (
    <ul className={cn("flex flex-col gap-2", className)}>
      {sorted.map((d, i) => {
        const cfg =
          d.severity === 2
            ? { Icon: XCircle, border: "border-l-error", bg: "bg-error/5", text: "text-error" }
            : d.severity === 1
              ? {
                  Icon: AlertTriangle,
                  border: "border-l-secondary-container",
                  bg: "bg-secondary-container/10",
                  text: "text-on-secondary-container",
                }
              : {
                  Icon: Info,
                  border: "border-l-sky-aqua",
                  bg: "bg-sky-aqua/10",
                  text: "text-primary",
                };
        return (
          <li
            key={i}
            className={cn(
              "flex items-start gap-2 border-l-2 pl-3 py-2 pr-2 rounded-sm",
              cfg.border,
              cfg.bg,
            )}
          >
            <cfg.Icon className={cn("h-4 w-4 mt-0.5 shrink-0", cfg.text)} />
            <div className="flex-1 min-w-0">
              <div className={cn("font-body-sm text-body-sm", cfg.text)}>{d.message}</div>
              {d.code && (
                <div className="font-data-label text-data-label text-on-surface-variant mt-0.5">
                  {d.code}
                </div>
              )}
            </div>
          </li>
        );
      })}
    </ul>
  );
}

interface EmptyStateProps {
  title: string;
  description?: string;
  icon?: ReactNode;
  action?: ReactNode;
  className?: string;
}
export function EmptyState({ title, description, icon, action, className }: EmptyStateProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center text-center px-6 py-12 rounded-lg border border-dashed border-outline/30 bg-surface-container-low/30",
        className,
      )}
    >
      <div className="text-on-surface-variant/60 mb-3">{icon ?? <Inbox className="h-8 w-8" />}</div>
      <p className="font-body-md text-body-md text-on-surface font-medium">{title}</p>
      {description && (
        <p className="font-body-sm text-body-sm text-on-surface-variant mt-1 max-w-md">
          {description}
        </p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

interface ErrorBannerProps {
  error: unknown;
  onDismiss?: () => void;
}
export function ErrorBanner({ error, onDismiss }: ErrorBannerProps) {
  const message =
    error instanceof Error ? error.message : typeof error === "string" ? error : String(error);
  return (
    <div className="flex items-start gap-2 rounded border-l-2 border-l-error bg-error/5 px-3 py-2">
      <XCircle className="h-4 w-4 mt-0.5 text-error shrink-0" />
      <div className="flex-1 font-body-sm text-body-sm text-error break-all">{message}</div>
      {onDismiss && (
        <button onClick={onDismiss} className="text-error/70 hover:text-error">
          ×
        </button>
      )}
    </div>
  );
}
