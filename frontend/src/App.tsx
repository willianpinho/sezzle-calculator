import { type FormEvent, useRef, useState } from "react";
import { calculate } from "./api";
import {
  operations,
  operationNeedsSecondNumber,
  secondNumberLabel,
  type Operation,
} from "./calculator";
import { parseFiniteNumber } from "./number";

const initialOperation: Operation = "add";

export default function App() {
  const [operation, setOperation] = useState<Operation>(initialOperation);
  const [firstNumber, setFirstNumber] = useState("");
  const [secondNumber, setSecondNumber] = useState("");
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const controllerRef = useRef<AbortController | null>(null);
  const requestId = useRef(0);
  const needsSecondNumber = operationNeedsSecondNumber(operation);

  function resetResponse() {
    requestId.current += 1;
    controllerRef.current?.abort();
    controllerRef.current = null;
    setIsLoading(false);
    setResult(null);
    setError("");
  }

  function updateOperation(value: Operation) {
    resetResponse();
    setOperation(value);
  }

  function updateFirstNumber(value: string) {
    resetResponse();
    setFirstNumber(value);
  }

  function updateSecondNumber(value: string) {
    resetResponse();
    setSecondNumber(value);
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const a = parseFiniteNumber(firstNumber);
    const b = parseFiniteNumber(secondNumber);

    if (a === null || (needsSecondNumber && b === null)) {
      setResult(null);
      setError("Enter a valid number in every required field.");
      return;
    }

    controllerRef.current?.abort();
    const controller = new AbortController();
    controllerRef.current = controller;
    const currentRequest = ++requestId.current;
    setIsLoading(true);
    setResult(null);
    setError("");

    try {
      const response = await calculate(
        needsSecondNumber ? { operation, a, b: b as number } : { operation, a },
        controller.signal,
      );
      if (currentRequest === requestId.current) setResult(response);
    } catch (reason) {
      if (
        currentRequest === requestId.current &&
        !(reason instanceof DOMException && reason.name === "AbortError")
      ) {
        setError(
          reason instanceof Error
            ? reason.message
            : "A network error occurred. Please try again.",
        );
      }
    } finally {
      if (currentRequest === requestId.current) setIsLoading(false);
    }
  }

  return (
    <main className="page-shell">
      <section className="calculator" aria-labelledby="calculator-title">
        <header>
          <p className="eyebrow">SEZZLE TOOLKIT</p>
          <h1 id="calculator-title">Calculator</h1>
          <p className="intro">
            Quick, clear calculations for the numbers in front of you.
          </p>
        </header>

        <form onSubmit={handleSubmit} noValidate>
          <div className="field">
            <label htmlFor="operation">Operation</label>
            <select
              id="operation"
              value={operation}
              onChange={(event) =>
                updateOperation(event.target.value as Operation)
              }
            >
              {operations.map((item) => (
                <option key={item.value} value={item.value}>
                  {item.label}
                </option>
              ))}
            </select>
          </div>

          <div className="numbers">
            <div className="field">
              <label htmlFor="first-number">First number</label>
              <input
                id="first-number"
                type="number"
                inputMode="decimal"
                step="any"
                value={firstNumber}
                onChange={(event) => updateFirstNumber(event.target.value)}
              />
            </div>
            {needsSecondNumber && (
              <div className="field">
                <label htmlFor="second-number">
                  {secondNumberLabel(operation)}
                </label>
                <input
                  id="second-number"
                  type="number"
                  inputMode="decimal"
                  step="any"
                  value={secondNumber}
                  onChange={(event) => updateSecondNumber(event.target.value)}
                />
              </div>
            )}
          </div>
          {operation === "percentage" && (
            <p className="hint">
              Calculates the first number as a percentage of the second.
            </p>
          )}
          <button type="submit" disabled={isLoading}>
            {isLoading ? "Calculating…" : "Calculate"}
          </button>
        </form>

        <output className="result" aria-live="polite" aria-atomic="true">
          {isLoading && "Calculating…"}
          {result !== null && (
            <>
              {" "}
              <span>Result</span>
              <strong>{result}</strong>
            </>
          )}
          {error && (
            <span className="error" role="alert">
              {error}
            </span>
          )}
          {!isLoading && result === null && !error && (
            <span>Enter numbers, then calculate.</span>
          )}
        </output>
      </section>
    </main>
  );
}
