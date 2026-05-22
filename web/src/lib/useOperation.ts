import { useCallback, useRef, useState } from "react";

export interface OperationState<O> {
  data: O | null;
  error: Error | null;
  loading: boolean;
}

export function useOperation<I, O>(fn: (input: I) => Promise<O>) {
  const [state, setState] = useState<OperationState<O>>({
    data: null,
    error: null,
    loading: false,
  });
  const abortRef = useRef<AbortController | null>(null);

  const run = useCallback(
    async (input: I) => {
      abortRef.current?.abort();
      const ctrl = new AbortController();
      abortRef.current = ctrl;
      setState((s) => ({ ...s, loading: true, error: null }));
      try {
        const data = await fn(input);
        if (!ctrl.signal.aborted) {
          setState({ data, error: null, loading: false });
        }
        return data;
      } catch (e) {
        const err = e instanceof Error ? e : new Error(String(e));
        if (!ctrl.signal.aborted) {
          setState({ data: null, error: err, loading: false });
        }
        throw err;
      }
    },
    [fn],
  );

  const reset = useCallback(() => {
    abortRef.current?.abort();
    setState({ data: null, error: null, loading: false });
  }, []);

  return { ...state, run, reset };
}
