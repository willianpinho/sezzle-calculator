import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import App from "./App";
import { parseFiniteNumber } from "./number";

const fetchMock = vi.fn();

beforeEach(() => {
  fetchMock.mockReset();
  vi.stubGlobal("fetch", fetchMock);
});

function response(body: unknown, ok = true) {
  return Promise.resolve({ ok, json: () => Promise.resolve(body) });
}

function enterNumbers(a: string, b = "") {
  fireEvent.change(screen.getByLabelText("First number"), {
    target: { value: a },
  });
  if (b)
    fireEvent.change(screen.getByLabelText("Second number"), {
      target: { value: b },
    });
}

function deferred<T>() {
  let resolve: (value: T) => void;
  let reject: (reason?: unknown) => void;
  const promise = new Promise<T>((resolvePromise, rejectPromise) => {
    resolve = resolvePromise;
    reject = rejectPromise;
  });
  return { promise, resolve: resolve!, reject: reject! };
}

describe("App", () => {
  it("submits a binary calculation and displays the server result", async () => {
    fetchMock.mockReturnValue(response({ result: 5 }));
    render(<App />);
    enterNumbers("2", "3");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));

    await screen.findByText("5");
    expect(fetchMock).toHaveBeenCalledWith(
      "/api/calculate",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ operation: "add", a: 2, b: 3 }),
      }),
    );
  });

  it("omits the second number for square root", async () => {
    fetchMock.mockReturnValue(response({ result: 3 }));
    render(<App />);
    fireEvent.change(screen.getByLabelText("Operation"), {
      target: { value: "sqrt" },
    });
    expect(screen.queryByLabelText("Second number")).not.toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "9" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));

    await screen.findByText("3");
    expect(JSON.parse(fetchMock.mock.calls[0][1].body)).toEqual({
      operation: "sqrt",
      a: 9,
    });
  });

  it("shows API and local validation errors", async () => {
    fetchMock.mockReturnValue(
      response({ error: { message: "Cannot divide by zero." } }, false),
    );
    render(<App />);
    fireEvent.change(screen.getByLabelText("Operation"), {
      target: { value: "divide" },
    });
    enterNumbers("4", "0");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Cannot divide by zero.",
    );

    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    expect(screen.getByRole("alert")).toHaveTextContent("Enter a valid number");
  });

  it("clears a completed result when an operand changes", async () => {
    fetchMock.mockReturnValue(response({ result: 5 }));
    render(<App />);
    enterNumbers("2", "3");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    await screen.findByText("5");

    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "7" },
    });
    await waitFor(() =>
      expect(screen.queryByText("5")).not.toBeInTheDocument(),
    );
    expect(
      screen.getByText("Enter numbers, then calculate."),
    ).toBeInTheDocument();
  });

  it("parses finite values and rejects empty or non-finite values", () => {
    expect(parseFiniteNumber("2.5")).toBe(2.5);
    expect(parseFiniteNumber(" ")).toBeNull();
    expect(parseFiniteNumber("Infinity")).toBeNull();
  });

  it("shows a safe fallback for a non-Error network failure", async () => {
    fetchMock.mockRejectedValue("offline");
    render(<App />);
    enterNumbers("2", "3");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "A network error occurred. Please try again.",
    );
  });

  it("keeps an aborted request silent", async () => {
    fetchMock.mockRejectedValue(new DOMException("Aborted", "AbortError"));
    render(<App />);
    enterNumbers("2", "3");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));

    await waitFor(() =>
      expect(screen.queryByRole("alert")).not.toBeInTheDocument(),
    );
  });

  it("ignores a late response or error after operands change", async () => {
    const lateSuccess = deferred<{
      ok: boolean;
      json: () => Promise<{ result: number }>;
    }>();
    const lateFailure = deferred<never>();
    fetchMock
      .mockReturnValueOnce(lateSuccess.promise)
      .mockReturnValueOnce(lateFailure.promise);
    render(<App />);
    enterNumbers("2", "3");
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "4" },
    });
    lateSuccess.resolve({
      ok: true,
      json: () => Promise.resolve({ result: 5 }),
    });

    await waitFor(() =>
      expect(screen.queryByText("5")).not.toBeInTheDocument(),
    );
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "6" },
    });
    lateFailure.reject(new Error("late failure"));

    await waitFor(() =>
      expect(screen.queryByRole("alert")).not.toBeInTheDocument(),
    );
  });

  it("updates percentage guidance and validates a missing second number", () => {
    render(<App />);
    fireEvent.change(screen.getByLabelText("Operation"), {
      target: { value: "percentage" },
    });
    expect(
      screen.getByText(
        "Calculates the first number as a percentage of the second.",
      ),
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText("Number to find the percentage of"),
    ).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText("First number"), {
      target: { value: "25" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Calculate" }));
    expect(screen.getByRole("alert")).toHaveTextContent("Enter a valid number");
  });
});
