import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import App from "./App";

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
});
