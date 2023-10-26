export function formatCurrency(currency: string | number): string {
  return currency
    .toString()
    .replace(
      /(\d+?)(\d{3})?(\d{3})?([.,]\d+)?$/,
      (_, $1, $2, $3, $4) => `R$ ${$1}${$2 ? '.' + $2 : ''}${$3 ? '.' + $3 : ''}${$4 ? $4.replace('.', ',') : ',00'}`
    );
}