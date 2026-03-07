export function formatNumber(value: number) {
  return new Intl.NumberFormat().format(value ?? 0);
}

export function formatPercent(value: number) {
  return `${((value ?? 0) * 100).toFixed(1)}%`;
}

export function formatDate(value?: string | null) {
  if (!value) {
    return "-";
  }
  return new Date(value).toLocaleString();
}
