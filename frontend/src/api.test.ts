import { describe, expect, it, vi } from "vitest";
import { calculate, REQUEST_TIMEOUT_MS } from "./api";

describe("calculate", () => {
  it("rejects a malformed successful response", async () => {
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockResolvedValue({ ok: true, json: () => Promise.resolve(null) }),
    );

    await expect(
      calculate({ operation: "add", a: 2, b: 3 }, new AbortController().signal),
    ).rejects.toThrow("invalid result");
  });

  it("uses a safe message for malformed error payloads", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        json: () => Promise.resolve({ error: null }),
      }),
    );

    await expect(
      calculate({ operation: "add", a: 2, b: 3 }, new AbortController().signal),
    ).rejects.toThrow("could not be completed");
  });

  it("reports a timeout without leaving a pending request", async () => {
    vi.useFakeTimers();
    vi.stubGlobal(
      "fetch",
      vi.fn(
        (_url, options: RequestInit) =>
          new Promise((_, reject) => {
            (options.signal as AbortSignal).addEventListener("abort", () =>
              reject(new DOMException("Aborted", "AbortError")),
            );
          }),
      ),
    );
    const request = calculate(
      { operation: "sqrt", a: 9 },
      new AbortController().signal,
    );
    const assertion = expect(request).rejects.toThrow("timed out");
    await vi.advanceTimersByTimeAsync(REQUEST_TIMEOUT_MS);

    await assertion;
    vi.useRealTimers();
  });

  it("keeps the timeout active while the response body is pending", async () => {
    vi.useFakeTimers();
    vi.stubGlobal(
      "fetch",
      vi.fn((_url, options: RequestInit) =>
        Promise.resolve({
          ok: true,
          json: () =>
            new Promise((_, reject) => {
              (options.signal as AbortSignal).addEventListener("abort", () =>
                reject(new DOMException("Aborted", "AbortError")),
              );
            }),
        }),
      ),
    );
    const request = calculate(
      { operation: "sqrt", a: 9 },
      new AbortController().signal,
    );
    const assertion = expect(request).rejects.toThrow("timed out");
    await vi.advanceTimersByTimeAsync(REQUEST_TIMEOUT_MS);

    await assertion;
    vi.useRealTimers();
  });
});
