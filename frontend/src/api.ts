import type { Operation } from "./calculator";

type CalculationRequest = { operation: Operation; a: number; b?: number };
type ApiPayload = { error?: unknown; result?: unknown };

export const REQUEST_TIMEOUT_MS = 10_000;

function asPayload(value: unknown): ApiPayload {
  return value !== null && typeof value === "object"
    ? (value as ApiPayload)
    : {};
}

export async function calculate(
  request: CalculationRequest,
  signal: AbortSignal,
): Promise<number> {
  const timeoutController = new AbortController();
  const timeout = setTimeout(
    () => timeoutController.abort(),
    REQUEST_TIMEOUT_MS,
  );
  try {
    const response = await fetch("/api/calculate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(request),
      signal: AbortSignal.any([signal, timeoutController.signal]),
    });
    const payload = asPayload(await response.json().catch(() => null));

    if (!response.ok) {
      const message =
        payload.error !== null &&
        typeof payload.error === "object" &&
        "message" in payload.error &&
        typeof payload.error.message === "string"
          ? payload.error.message
          : null;
      throw new Error(message ?? "The calculation could not be completed.");
    }
    if (
      typeof payload.result !== "number" ||
      !Number.isFinite(payload.result)
    ) {
      throw new Error("The calculator returned an invalid result.");
    }
    return payload.result;
  } catch (reason) {
    if (timeoutController.signal.aborted) {
      throw new Error("The calculation timed out. Please try again.", {
        cause: reason,
      });
    }
    throw reason;
  } finally {
    clearTimeout(timeout);
  }
}
