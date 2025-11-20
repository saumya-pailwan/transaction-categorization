"use client";
export default function TransactionTable({ transactions }) {
  return (
    <table className="min-w-full border border-gray-300">
      <thead>
        <tr className="bg-gray-100">
          <th className="p-2 border">Date</th>
          <th className="p-2 border">Merchant</th>
          <th className="p-2 border">Amount</th>
          <th className="p-2 border">Category</th>
          <th className="p-2 border">Confidence</th>
        </tr>
      </thead>
      <tbody>
        {transactions.map((t) => (
          <tr key={t.id}>
            <td className="p-2 border">{t.date}</td>
            <td className="p-2 border">{t.merchant_name}</td>
            <td className="p-2 border">{t.amount}</td>
            <td className="p-2 border">{t.category_id || "—"}</td>
            <td className="p-2 border">{t.confidence ?? "—"}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}