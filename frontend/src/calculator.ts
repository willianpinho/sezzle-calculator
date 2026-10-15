export const operations = [
  { value: "add", label: "Add (+)", needsSecondNumber: true },
  { value: "subtract", label: "Subtract (−)", needsSecondNumber: true },
  { value: "multiply", label: "Multiply (×)", needsSecondNumber: true },
  { value: "divide", label: "Divide (÷)", needsSecondNumber: true },
  { value: "power", label: "Power (xʸ)", needsSecondNumber: true },
  { value: "sqrt", label: "Square root (√)", needsSecondNumber: false },
  { value: "percentage", label: "Percentage (% of)", needsSecondNumber: true },
] as const;

export type Operation = (typeof operations)[number]["value"];

export function operationNeedsSecondNumber(operation: Operation): boolean {
  return (
    operations.find((item) => item.value === operation)?.needsSecondNumber ??
    true
  );
}

export function secondNumberLabel(operation: Operation): string {
  return operation === "percentage"
    ? "Number to find the percentage of"
    : "Second number";
}
