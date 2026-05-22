import { forwardRef, type InputHTMLAttributes, type SelectHTMLAttributes, type TextareaHTMLAttributes, type ReactNode } from "react";
import { cn } from "../../lib/cn";

interface FieldShellProps {
  label?: string;
  hint?: string;
  error?: string;
  htmlFor?: string;
  trailing?: ReactNode;
  className?: string;
  children: ReactNode;
}

export function Field({ label, hint, error, htmlFor, trailing, className, children }: FieldShellProps) {
  return (
    <div className={cn("flex flex-col gap-1.5", className)}>
      {label && (
        <div className="flex items-center justify-between">
          <label htmlFor={htmlFor} className="font-data-label text-data-label uppercase text-on-surface-variant">
            {label}
          </label>
          {trailing}
        </div>
      )}
      {children}
      {error ? (
        <p className="font-body-sm text-body-sm text-error">{error}</p>
      ) : hint ? (
        <p className="font-body-sm text-body-sm text-on-surface-variant">{hint}</p>
      ) : null}
    </div>
  );
}

const baseInput =
  "w-full bg-surface-container-lowest border border-outline/20 rounded font-body-md text-body-md text-on-surface " +
  "px-3 py-2 transition-all placeholder:text-on-surface-variant/50 " +
  "focus:outline-none focus:ring-2 focus:ring-sky-aqua/40 focus:border-sky-aqua " +
  "disabled:opacity-60 disabled:cursor-not-allowed";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  hint?: string;
  error?: string;
}
export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, hint, error, className, id, ...props }, ref) => (
    <Field label={label} hint={hint} error={error} htmlFor={id}>
      <input ref={ref} id={id} className={cn(baseInput, className)} {...props} />
    </Field>
  ),
);
Input.displayName = "Input";

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  hint?: string;
  error?: string;
  mono?: boolean;
}
export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ label, hint, error, mono = true, className, id, rows = 8, ...props }, ref) => (
    <Field label={label} hint={hint} error={error} htmlFor={id}>
      <textarea
        ref={ref}
        id={id}
        rows={rows}
        spellCheck={false}
        className={cn(
          baseInput,
          "leading-relaxed resize-y",
          mono && "font-code-block text-code-block",
          className,
        )}
        {...props}
      />
    </Field>
  ),
);
Textarea.displayName = "Textarea";

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  hint?: string;
  error?: string;
  options: { label: string; value: string }[];
}
export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, hint, error, options, className, id, ...props }, ref) => (
    <Field label={label} hint={hint} error={error} htmlFor={id}>
      <select
        ref={ref}
        id={id}
        className={cn(baseInput, "appearance-none bg-no-repeat pr-9", className)}
        style={{
          backgroundImage:
            "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2345464c' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E\")",
          backgroundPosition: "right 0.75rem center",
          backgroundSize: "12px 12px",
        }}
        {...props}
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    </Field>
  ),
);
Select.displayName = "Select";

interface ToggleProps {
  checked: boolean;
  onChange: (v: boolean) => void;
  label: string;
  hint?: string;
  disabled?: boolean;
}
export function Toggle({ checked, onChange, label, hint, disabled }: ToggleProps) {
  return (
    <label
      className={cn(
        "flex items-start gap-3 cursor-pointer select-none",
        disabled && "opacity-50 cursor-not-allowed",
      )}
    >
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        disabled={disabled}
        onClick={() => !disabled && onChange(!checked)}
        className={cn(
          "relative mt-0.5 inline-flex h-5 w-9 flex-shrink-0 items-center rounded-full transition-colors",
          checked ? "bg-secondary-container" : "bg-outline/30",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-aqua",
        )}
      >
        <span
          className={cn(
            "inline-block h-4 w-4 transform rounded-full bg-white shadow-sm transition-transform",
            checked ? "translate-x-4" : "translate-x-0.5",
          )}
        />
      </button>
      <div className="flex-1 min-w-0">
        <div className="font-body-sm text-body-sm text-on-surface">{label}</div>
        {hint && <div className="font-body-sm text-body-sm text-on-surface-variant">{hint}</div>}
      </div>
    </label>
  );
}

interface NumberInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "type" | "onChange"> {
  label?: string;
  hint?: string;
  value: number;
  onChange: (n: number) => void;
  min?: number;
  max?: number;
  step?: number;
}
export function NumberInput({ label, hint, value, onChange, min, max, step = 1, ...props }: NumberInputProps) {
  return (
    <Field label={label} hint={hint}>
      <input
        type="number"
        value={value}
        min={min}
        max={max}
        step={step}
        onChange={(e) => onChange(Number(e.target.value))}
        className={cn(baseInput)}
        {...props}
      />
    </Field>
  );
}
