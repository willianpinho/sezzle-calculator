import { describe, expect, it } from "vitest";
import {
  operationNeedsSecondNumber,
  operations,
  secondNumberLabel,
  type Operation,
} from "./calculator";

describe("calculator operation helpers", () => {
  it("identifies unary and binary operations", () => {
    for (const operation of operations) {
      expect(operationNeedsSecondNumber(operation.value)).toBe(
        operation.value !== "sqrt",
      );
    }
  });

  it("uses safe defaults and operation-specific labels", () => {
    expect(operationNeedsSecondNumber("unknown" as Operation)).toBe(true);
    expect(secondNumberLabel("percentage")).toBe(
      "Number to find the percentage of",
    );
    expect(secondNumberLabel("multiply")).toBe("Second number");
  });
});
